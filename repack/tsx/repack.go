package tsx

import (
	"image"
	"log"
)

func NewRepacker(srcTilesets []*Tileset) *Repacker {
	if len(srcTilesets) < 1 {
		log.Fatalf("At least one source tileset is required for repacking")
	}

	for _, ts := range srcTilesets[1:] {
		if ts.TileWidth != srcTilesets[0].TileWidth ||
			ts.TileHeight != srcTilesets[0].TileHeight {
			log.Fatalf("All source tilesets must have the same tile size")
		}
	}

	return &Repacker{
		srcTilesets:    srcTilesets,
		repackedTileID: make(map[GlobalTileID]GlobalTileID),
	}
}

type Repacker struct {
	srcTilesets    []*Tileset
	tilesToRepack  []*Tile
	repackedTileID map[GlobalTileID]GlobalTileID
}

func (r *Repacker) UseTileID(tileID GlobalTileID) {
	if _, ok := r.repackedTileID[tileID]; ok {
		return
	}

	var tile *Tile
	for _, ts := range r.srcTilesets {
		tile = ts.Tile(tileID)
		if tile != nil {
			break
		}
	}
	if tile == nil {
		log.Fatalf("Failed to find source tile %d", tileID)
	}

	r.tilesToRepack = append(r.tilesToRepack, tile)
	r.repackedTileID[tileID] = GlobalTileID(len(r.tilesToRepack))

	for _, anim := range tile.Animation {
		animTileID := tile.Tileset.GlobalTileID(anim.TileID)
		if _, ok := r.repackedTileID[animTileID]; ok {
			continue
		}
		animTile := tile.Tileset.Tile(animTileID)
		if animTile == nil {
			log.Fatalf("Failed to find source animation tile %d", animTileID)
		}
		r.tilesToRepack = append(r.tilesToRepack, animTile)
		r.repackedTileID[animTileID] = GlobalTileID(len(r.tilesToRepack))
	}
}

func (r *Repacker) BuildNewTileset(name string) *Tileset {
	tileset := r.srcTilesets[0]

	tileCount := len(r.tilesToRepack)
	imageColumns := nextPowerOfTwoSquare(tileCount)
	imageWidth := tileset.TileWidth * imageColumns
	imageHeight := tileset.TileHeight * imageColumns

	repackedTileset := &Tileset{
		FirstGID:   1,
		Name:       name,
		TileWidth:  tileset.TileWidth,
		TileHeight: tileset.TileHeight,
		TileCount:  tileCount,
		Columns:    imageColumns,
		Image: Image{
			Source: name + ".png",
			Width:  imageWidth,
			Height: imageHeight,
			Data:   image.NewRGBA(image.Rect(0, 0, imageWidth, imageHeight)),
		},
	}

	for i, tile := range r.tilesToRepack {
		repackedTile := &Tile{
			ID:          LocalTileID(i),
			Type:        tile.Type,
			ObjectGroup: tile.ObjectGroup,
			Properties:  tile.Properties,
			Tileset:     repackedTileset,
		}
		for _, frame := range tile.Animation {
			animTileID := tile.Tileset.GlobalTileID(frame.TileID)
			repackedAnimTileID := r.repackedTileID[animTileID]
			repackedTile.Animation = append(repackedTile.Animation, Frame{
				TileID:   LocalTileID(repackedAnimTileID - 1),
				Duration: frame.Duration,
			})
		}
		repackedTile.CopyImageFrom(tile)
		repackedTileset.Tiles = append(repackedTileset.Tiles, repackedTile)
	}

	return repackedTileset
}

func (r *Repacker) RepackedTileID(id GlobalTileID) GlobalTileID {
	return r.repackedTileID[id]
}

func nextPowerOfTwoSquare(n int) int {
	res := 1
	for {
		if res*res >= n {
			break
		}
		res *= 2
	}
	return res
}
