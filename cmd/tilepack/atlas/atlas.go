package atlas

import (
	"fmt"
	"image"
	"log"
	"path"
	"pixel-tools/cmd/tilepack/tmj"

	atlaspkg "pixel-tools/pkg/atlas"

	"pixel-tools/cmd/tilepack/tsx"
)

func New(name string, tilesets []*tsx.Tileset) *Atlas {
	return &Atlas{
		atlas:           atlaspkg.New(name),
		name:            name,
		tilesets:        tilesets,
		repackedTiles:   make(map[tsx.GlobalTileID]tsx.GlobalTileID),
		repackedSprites: make(map[tsx.GlobalTileID]tsx.GlobalTileID),
	}
}

type Atlas struct {
	atlas    *atlaspkg.Atlas
	name     string
	tilesets []*tsx.Tileset

	repackedTiles   map[tsx.GlobalTileID]tsx.GlobalTileID
	repackedSprites map[tsx.GlobalTileID]tsx.GlobalTileID
	tiles           []*Sprite
	sprites         []*Sprite
}

type Sprite struct {
	ID   string
	Tile *tsx.Tile
}

func (a *Atlas) UseTile(tileID tsx.GlobalTileID) tsx.GlobalTileID {
	if repackedTileID, ok := a.repackedTiles[tileID]; ok {
		return repackedTileID
	}

	tile := a.findTile(tileID)
	if len(a.tiles) > 0 {
		firstTile := a.tiles[0].Tile
		if tile.Tileset.TileWidth != firstTile.Tileset.TileWidth || tile.Tileset.TileHeight != firstTile.Tileset.TileHeight {
			log.Fatalf("All tiles must have the same size, please use object layer instead")
		}
	}

	repackedTileID := tsx.GlobalTileID(len(a.repackedTiles) + 1)
	a.repackedTiles[tileID] = repackedTileID

	w := tile.Tileset.TileWidth
	h := tile.Tileset.TileHeight
	srcX := (int(tile.ID) % tile.Tileset.Columns) * w
	srcY := (int(tile.ID) / tile.Tileset.Columns) * h

	subImg := tile.Tileset.Image.Data.SubImage(image.Rect(srcX, srcY, srcX+w, srcY+h))

	a.atlas.AddTile(subImg)
	a.tiles = append(a.tiles, &Sprite{
		Tile: tile,
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

	var subImg image.Image
	if tile.Image == nil {
		w := tile.Tileset.TileWidth
		h := tile.Tileset.TileHeight
		srcX := (int(tile.ID) % tile.Tileset.Columns) * w
		srcY := (int(tile.ID) / tile.Tileset.Columns) * h
		subImg = tile.Tileset.Image.Data.SubImage(image.Rect(srcX, srcY, srcX+w, srcY+h))
	} else {
		subImg = tile.Image.Data.SubImage(image.Rect(tile.X, tile.Y, tile.X+tile.Width, tile.Y+tile.Height))
	}

	data := convertTileData(tile)

	a.atlas.AddSprite(fmt.Sprintf("%d", repackedTileID), subImg, data)
	a.sprites = append(a.sprites, &Sprite{
		ID:   fmt.Sprintf("%d", repackedTileID),
		Tile: tile,
	})

	for _, anim := range tile.Animation {
		a.UseSprite(tile.Tileset.GlobalTileID(anim.TileID))
	}
	return repackedTileID
}

func (a *Atlas) Pack() *tsx.Tileset {
	a.atlas.Pack()
	return a.buildTileset()
}

func (a *Atlas) Save(baseName string) {
	a.atlas.SaveImage(baseName+".png")
	a.atlas.SaveJSON(baseName+".atlas", path.Base(baseName)+".png")
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

func (a *Atlas) buildTileset() *tsx.Tileset {
	tileCount := len(a.tiles)
	tileWidth := 16
	tileHeight := 16
	if tileCount != 0 {
		tileWidth = a.tiles[0].Tile.Tileset.TileWidth
		tileHeight = a.tiles[0].Tile.Tileset.TileHeight
	}

	size := a.atlas.Size()
	tileColumns := size / tileWidth
	tileset := &tsx.Tileset{
		FirstGID:   1,
		Name:       a.name,
		TileWidth:  tileWidth,
		TileHeight: tileHeight,
		TileCount:  tileCount,
		Columns:    tileColumns,
		Image: &tsx.Image{
			Source: a.name + ".png",
			Width:  size,
			Height: size,
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

func convertTileData(tile *tsx.Tile) map[string]any {
	res := make(map[string]any)

	animation := tmj.ConvertAnimation(tile.Animation)
	if len(animation) > 0 {
		res["animation"] = animation
	}

	objectGroup := tmj.ConvertOptionalObjectGroup(tile.ObjectGroup)
	if objectGroup != nil {
		res["objectgroup"] = objectGroup
	}

	properties := tmj.ConvertProperties(tile.Properties)
	if len(properties) > 0 {
		res["properties"] = properties
	}

	if len(res) == 0 {
		return nil
	}

	return res
}
