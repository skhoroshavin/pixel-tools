package atlas

import (
	"encoding/json"
	"fmt"
	"image"
	"image/draw"
	"log"
	"os"
	"path"
	"repack/tsx"
	"repack/util"
	"sort"
)

func New(tilesets []*tsx.Tileset) *Atlas {
	return &Atlas{
		tilesets:       tilesets,
		repackedTileID: make(map[tsx.GlobalTileID]tsx.GlobalTileID),
		size:           16,
		skyline:        []Point{{X: 0, Y: 0}},
	}
}

type Atlas struct {
	tilesets       []*tsx.Tileset
	repackedTileID map[tsx.GlobalTileID]tsx.GlobalTileID
	sprites        []*Sprite
	size           int
	skyline        []Point
}

type Sprite struct {
	ID      tsx.GlobalTileID
	SrcTile *tsx.Tile
	Frame   Rect
}

type Point struct {
	X, Y int
}

func (a *Atlas) UseTileID(tileID tsx.GlobalTileID) tsx.GlobalTileID {
	if repackedTileID, ok := a.repackedTileID[tileID]; ok {
		return repackedTileID
	}

	var tile *tsx.Tile
	for _, ts := range a.tilesets {
		tile = ts.Tile(tileID)
		if tile != nil {
			break
		}
	}
	if tile == nil {
		log.Fatalf("Failed to find source tile %d", tileID)
	}

	repackedTileID := tsx.GlobalTileID(1000000 + len(a.repackedTileID))
	a.repackedTileID[tileID] = repackedTileID
	a.sprites = append(a.sprites, &Sprite{
		ID:      repackedTileID,
		SrcTile: tile,
		Frame: Rect{
			W: tile.Tileset.TileWidth,
			H: tile.Tileset.TileHeight,
		},
	})

	for _, anim := range tile.Animation {
		a.UseTileID(tile.Tileset.GlobalTileID(anim.TileID))
	}

	return repackedTileID
}

func (a *Atlas) Pack() {
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

	for _, sprite := range a.sprites {
		a.packSprite(sprite)
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

	for _, sprite := range a.sprites {
		srcTile := sprite.SrcTile
		srcColumns := srcTile.Tileset.Columns

		w := sprite.Frame.W
		h := sprite.Frame.H

		srcX := (int(srcTile.ID) % srcColumns) * w
		srcY := (int(srcTile.ID) / srcColumns) * h
		dstX := sprite.Frame.X
		dstY := sprite.Frame.Y

		draw.Draw(img, image.Rect(dstX, dstY, dstX+w, dstY+h),
			srcTile.Tileset.Image.Data, image.Pt(srcX, srcY), draw.Src)

		data.Frames = append(data.Frames, Frame{
			Filename: fmt.Sprintf("%d", sprite.ID),
			Frame:    sprite.Frame,
		})
	}

	util.WriteImage(img, baseName+".png")

	atlasJSON, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Fatalf("Failed to encode atlas JSON: %v", err)
	}
	err = os.WriteFile(baseName+".json", atlasJSON, 0644)
	if err != nil {
		log.Fatalf("Failed to write atlas JSON: %v", err)
	}
}

func (a *Atlas) packSprite(sprite *Sprite) {
	idx := a.findBestPosition(sprite)
	pt := a.skyline[idx]
	sprite.Frame.X = pt.X
	sprite.Frame.Y = pt.Y
	a.insertFrame(idx, sprite.Frame)
}

func (a *Atlas) findBestPosition(sprite *Sprite) int {
	bestIdx := -1
	var bestPoint Point
	for i, pt := range a.skyline {
		// If we don't fit into frame - skip
		if !a.isFitting(i, sprite.Frame) {
			continue
		}
		// If we already have candidate best point - compare with it and discard if we are not better
		if bestIdx != -1 && bestPoint.Y <= pt.Y {
			continue
		}
		bestIdx = i
		bestPoint = pt
	}
	// If we found best point - return it
	if bestIdx >= 0 {
		return bestIdx
	}
	// Otherwise double atlas size and try again
	a.size *= 2
	return a.findBestPosition(sprite)
}

func (a *Atlas) isFitting(idx int, frame Rect) bool {
	base := a.skyline[idx]
	// Return false if we don't fit to frame
	if (base.X+frame.W > a.size) || (base.Y+frame.H > a.size) {
		return false
	}
	// Check whether we overlap with the next skyline points
	for i := idx + 1; i < len(a.skyline); i++ {
		pt := a.skyline[i]
		// Point is outside the sprite - we can stop iterating, safely assuming that sprite fits
		if pt.X >= base.X+frame.W {
			return true
		}
		// Point is inside the sprite - we can stop iterating, safely assuming that sprite doesn't fit
		if pt.Y > base.Y {
			return false
		}
	}
	return true
}

func (a *Atlas) insertFrame(idx int, frame Rect) {
	pt := a.skyline[idx]
	if idx == len(a.skyline)-1 {
		a.skyline = append(a.skyline, Point{X: pt.X + frame.W, Y: 0})
	} else {
		curY := a.skyline[idx].Y
		nextIdx := idx + 1
		nextPt := a.skyline[nextIdx]
		for {
			if nextPt.X >= pt.X+frame.W {
				a.insertPoint(nextIdx, Point{X: pt.X + frame.W, Y: curY})
				break
			}

			a.removePoint(nextIdx)
			curY = nextPt.Y
			if nextIdx < len(a.skyline) {
				nextPt = a.skyline[nextIdx]
			} else {
				nextPt = Point{X: a.size, Y: 0}
			}
		}
	}
	a.skyline[idx] = Point{X: pt.X, Y: pt.Y + frame.H}

	for i := 0; i < len(a.skyline)-1; i++ {
		cur := a.skyline[i]
		next := a.skyline[i+1]
		if cur.Y == next.Y {
			a.removePoint(i + 1)
			i--
		}
	}
}

func (a *Atlas) insertPoint(idx int, pt Point) {
	a.skyline = append(a.skyline, Point{})
	for i := len(a.skyline) - 1; i > idx; i-- {
		a.skyline[i] = a.skyline[i-1]
	}
	a.skyline[idx] = pt
}

func (a *Atlas) removePoint(idx int) {
	for i := idx + 1; i < len(a.skyline); i++ {
		a.skyline[i-1] = a.skyline[i]
	}
	a.skyline = a.skyline[:len(a.skyline)-1]
}
