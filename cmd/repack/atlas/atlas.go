package atlas

import (
	"encoding/json"
	"fmt"
	"image"
	"image/draw"
	"log"
	"os"
	"path"
	"sort"

	"pixel-tools/pkg/imgutil"
	"pixel-tools/cmd/repack/tsx"
)

func New(name string, tilesets []*tsx.Tileset) *Atlas {
	const initialSize = 16

	return &Atlas{
		name:            name,
		tilesets:        tilesets,
		repackedTiles:   make(map[tsx.GlobalTileID]tsx.GlobalTileID),
		repackedSprites: make(map[tsx.GlobalTileID]tsx.GlobalTileID),
		size:            initialSize,
		skyline:         make([]int, initialSize),
	}
}

type Atlas struct {
	name     string
	tilesets []*tsx.Tileset

	repackedTiles   map[tsx.GlobalTileID]tsx.GlobalTileID
	repackedSprites map[tsx.GlobalTileID]tsx.GlobalTileID
	tiles           []*Sprite
	sprites         []*Sprite

	size    int
	skyline []int
}

type Sprite struct {
	ID    string
	Tile  *tsx.Tile
	Frame Rect
}

func (a *Atlas) UseTile(tileID tsx.GlobalTileID) tsx.GlobalTileID {
	if repackedTileID, ok := a.repackedTiles[tileID]; ok {
		return repackedTileID
	}

	tile := a.findTile(tileID)
	if len(a.tiles) > 0 && !tile.Tileset.HasSameTileSize(a.tiles[0].Tile.Tileset) {
		log.Fatalf("All tiles must have the same size, please use object layer instead")
	}

	repackedTileID := tsx.GlobalTileID(len(a.repackedTiles) + 1)
	a.repackedTiles[tileID] = repackedTileID
	a.tiles = append(a.tiles, &Sprite{
		Tile: tile,
		Frame: Rect{
			W: tile.Tileset.TileWidth,
			H: tile.Tileset.TileHeight,
		},
	})

	for _, anim := range tile.Animation {
		a.UseTile(tile.Tileset.GlobalTileID(anim.TileID))
	}
	return repackedTileID
}

func (a *Atlas) UseSprite(tileID tsx.GlobalTileID) tsx.GlobalTileID {
	if repackedTileID, ok := a.repackedSprites[tileID]; ok {
		return repackedTileID
	}

	tile := a.findTile(tileID)
	repackedTileID := tsx.GlobalTileID(1000000 + len(a.repackedSprites))
	a.repackedSprites[tileID] = repackedTileID

	var frame Rect
	if tile.Image == nil {
		frame = Rect{
			W: tile.Tileset.TileWidth,
			H: tile.Tileset.TileHeight,
		}
	} else {
		frame = Rect{
			W: tile.Width,
			H: tile.Height,
		}
	}

	a.sprites = append(a.sprites, &Sprite{
		ID:    fmt.Sprintf("%d", repackedTileID),
		Tile:  tile,
		Frame: frame,
	})

	for _, anim := range tile.Animation {
		a.UseSprite(tile.Tileset.GlobalTileID(anim.TileID))
	}
	return repackedTileID
}

func (a *Atlas) Pack() *tsx.Tileset {
	// Sort sprites by size
	sort.Slice(a.sprites, func(i, j int) bool {
		iw := a.sprites[i].Frame.W
		ih := a.sprites[i].Frame.H
		jw := a.sprites[j].Frame.W
		jh := a.sprites[j].Frame.H
		if iw*ih > jw*jh {
			return true
		}
		if iw*ih < jw*jh {
			return false
		}
		return iw > ih
	})

	// Repack until the size becomes stable
	lastSize := a.size
	for {
		a.setSkyline(0)
		for _, tile := range a.tiles {
			a.packSprite(tile)
		}

		a.setSkyline(a.skyline[0])
		for _, sprite := range a.sprites {
			a.packSprite(sprite)
		}

		if lastSize != a.size {
			lastSize = a.size
		} else {
			return a.buildTileset()
		}
	}
}

func (a *Atlas) Save(baseName string) {
	img := image.NewRGBA(image.Rect(0, 0, a.size, a.size))
	data := RootJson{
		Meta: Meta{
			Image:  path.Base(baseName) + ".png",
			Format: "RGBA8888",
			Size: Size{
				W: a.size,
				H: a.size,
			},
			Scale: "1",
		},
	}

	for _, tile := range a.tiles {
		dstX := tile.Frame.X
		dstY := tile.Frame.Y
		w := tile.Frame.W
		h := tile.Frame.H

		srcTile := tile.Tile
		srcColumns := srcTile.Tileset.Columns
		srcX := (int(srcTile.ID) % srcColumns) * w
		srcY := (int(srcTile.ID) / srcColumns) * h

		draw.Draw(img, image.Rect(dstX, dstY, dstX+w, dstY+h),
			srcTile.Tileset.Image.Data, image.Pt(srcX, srcY), draw.Src)
	}

	for _, sprite := range a.sprites {
		dstX := sprite.Frame.X
		dstY := sprite.Frame.Y
		w := sprite.Frame.W
		h := sprite.Frame.H

		srcTile := sprite.Tile
		srcColumns := srcTile.Tileset.Columns
		var srcX, srcY int
		var srcImage image.Image
		if srcTile.Image == nil {
			srcX = (int(srcTile.ID) % srcColumns) * w
			srcY = (int(srcTile.ID) / srcColumns) * h
			srcImage = srcTile.Tileset.Image.Data
		} else {
			srcX = srcTile.X
			srcY = srcTile.Y
			srcImage = srcTile.Image.Data
		}

		draw.Draw(img, image.Rect(dstX, dstY, dstX+w, dstY+h),
			srcImage, image.Pt(srcX, srcY), draw.Src)

		data.Frames = append(data.Frames, Frame{
			Filename: sprite.ID,
			Frame:    sprite.Frame,

			Animation:   convertAnimation(srcTile.Animation),
			ObjectGroup: convertObjectGroup(srcTile.ObjectGroup),
			Properties:  convertProperties(srcTile.Properties),
		})
	}

	imgutil.WriteImage(img, baseName+".png")

	atlasJSON, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Fatalf("Failed to encode atlas JSON: %v", err)
	}
	err = os.WriteFile(baseName+".atlas", atlasJSON, 0644)
	if err != nil {
		log.Fatalf("Failed to write atlas JSON: %v", err)
	}
}

func (a *Atlas) setSkyline(level int) {
	a.skyline = make([]int, a.size)
	for i := range a.skyline {
		a.skyline[i] = level
	}
}

func (a *Atlas) findTile(id tsx.GlobalTileID) *tsx.Tile {
	for _, ts := range a.tilesets {
		if tile := ts.Tile(id); tile != nil {
			return tile
		}
	}
	log.Fatalf("Failed to find source tile %d", id)
	return nil
}

func (a *Atlas) packSprite(sprite *Sprite) {
	x := a.findBestPosition(sprite)
	sprite.Frame.X = x
	sprite.Frame.Y = a.skyline[x]
	a.insertFrame(x, sprite.Frame)
}

func (a *Atlas) findBestPosition(sprite *Sprite) int {
	bestX := -1
	for x, y := range a.skyline {
		// If we don't fit into the frame, then skip
		if !a.isFitting(x, sprite.Frame) {
			continue
		}
		// If we already have candidate best point - compare with it and discard if we are not better
		if bestX != -1 && a.skyline[bestX] <= y {
			continue
		}
		bestX = x
	}
	// If we found the best point - return it
	if bestX >= 0 {
		return bestX
	}
	// Otherwise double the atlas size and try again
	a.skyline = append(a.skyline, make([]int, a.size)...)
	a.size *= 2
	return a.findBestPosition(sprite)
}

func (a *Atlas) isFitting(x int, frame Rect) bool {
	y := a.skyline[x]
	// Return false if we don't fit to frame
	if (x+frame.W > a.size) || (y+frame.H > a.size) {
		return false
	}
	// Check whether we overlap with the next skyline points
	for nx := x + 1; nx < len(a.skyline); nx++ {
		// Point is outside the sprite - we can stop iterating, safely assuming that sprite fits
		if nx >= x+frame.W {
			return true
		}
		// Point is inside the sprite - we can stop iterating, safely assuming that sprite doesn't fit
		if a.skyline[nx] > y {
			return false
		}
	}
	return true
}

func (a *Atlas) insertFrame(x int, frame Rect) {
	cy := a.skyline[x] + frame.H
	for cx := x; cx < x+frame.W; cx++ {
		a.skyline[cx] = cy
	}
}

func (a *Atlas) buildTileset() *tsx.Tileset {
	tileCount := len(a.tiles)
	tileWidth := 16
	tileHeight := 16
	if tileCount != 0 {
		tileWidth = a.tiles[0].Tile.Tileset.TileWidth
		tileHeight = a.tiles[0].Tile.Tileset.TileHeight
	}

	tileColumns := a.size / tileWidth
	tileset := &tsx.Tileset{
		FirstGID:   1,
		Name:       a.name,
		TileWidth:  tileWidth,
		TileHeight: tileHeight,
		TileCount:  tileCount,
		Columns:    tileColumns,
		Image: &tsx.Image{
			Source: a.name + ".png",
			Width:  a.size,
			Height: a.size,
		},
	}

	for i, sprite := range a.tiles {
		tile := sprite.Tile
		repackedTile := &tsx.Tile{
			ID:          tsx.LocalTileID(i),
			Type:        tile.Type,
			ObjectGroup: tile.ObjectGroup,
			Properties:  tile.Properties,
			Tileset:     tileset,
		}
		for _, frame := range sprite.Tile.Animation {
			animTileID := tile.Tileset.GlobalTileID(frame.TileID)
			repackedAnimTileID := a.repackedTiles[animTileID]
			repackedTile.Animation = append(repackedTile.Animation, tsx.Frame{
				TileID:   tsx.LocalTileID(repackedAnimTileID - 1),
				Duration: frame.Duration,
			})
		}
		tileset.Tiles = append(tileset.Tiles, repackedTile)
	}

	return tileset
}
