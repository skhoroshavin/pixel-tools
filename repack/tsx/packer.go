package tsx

import (
	"image"
	"log"
)

func NewPacker(srcTilesets []*Tileset) *Packer {
	if len(srcTilesets) < 1 {
		log.Fatalf("At least one source tileset is required for repacking")
	}

	minTileWidth := srcTilesets[0].TileWidth
	minTileHeight := srcTilesets[0].TileHeight
	for _, ts := range srcTilesets[1:] {
		if ts.TileWidth < minTileWidth {
			minTileWidth = ts.TileWidth
		}
		if ts.TileHeight < minTileHeight {
			minTileHeight = ts.TileHeight
		}
	}

	return &Packer{
		srcTilesets:    srcTilesets,
		repackedTileID: make(map[GlobalTileID]GlobalTileID),
		tileWidth:      minTileWidth,
		tileHeight:     minTileHeight,
	}
}

type Packer struct {
	srcTilesets    []*Tileset
	tilesToRepack  []*Tile
	repackedTileID map[GlobalTileID]GlobalTileID
	tileWidth      int
	tileHeight     int
}

func (p *Packer) UseTileID(tileID GlobalTileID) GlobalTileID {
	if repackedTileID, ok := p.repackedTileID[tileID]; ok {
		return repackedTileID
	}

	var tile *Tile
	for _, ts := range p.srcTilesets {
		tile = ts.Tile(tileID)
		if tile != nil {
			break
		}
	}
	if tile == nil {
		log.Fatalf("Failed to find source tile %d", tileID)
	}

	p.tilesToRepack = append(p.tilesToRepack, tile)
	repackedTileID := GlobalTileID(len(p.tilesToRepack))
	p.repackedTileID[tileID] = repackedTileID

	for _, anim := range tile.Animation {
		animTileID := tile.Tileset.GlobalTileID(anim.TileID)
		if _, ok := p.repackedTileID[animTileID]; ok {
			continue
		}
		animTile := tile.Tileset.Tile(animTileID)
		if animTile == nil {
			log.Fatalf("Failed to find source animation tile %d", animTileID)
		}
		p.tilesToRepack = append(p.tilesToRepack, animTile)
		p.repackedTileID[animTileID] = GlobalTileID(len(p.tilesToRepack))
	}

	return repackedTileID
}

func (p *Packer) BuildNewTileset(name string) *Tileset {
	tileCount := len(p.tilesToRepack)
	imageColumns := nextPowerOfTwoSquare(tileCount)
	imageWidth := p.tileWidth * imageColumns
	imageHeight := p.tileHeight * imageColumns

	repackedTileset := &Tileset{
		FirstGID:   1,
		Name:       name,
		TileWidth:  p.tileWidth,
		TileHeight: p.tileHeight,
		TileCount:  tileCount,
		Columns:    imageColumns,
		Image: Image{
			Source: name + ".png",
			Width:  imageWidth,
			Height: imageHeight,
			Data:   image.NewRGBA(image.Rect(0, 0, imageWidth, imageHeight)),
		},
	}

	for i, tile := range p.tilesToRepack {
		if tile.Tileset.TileWidth != p.tileWidth || tile.Tileset.TileHeight != p.tileHeight {
			log.Fatalf("Tileset %s has different tile size than base tile size, please use object layer instead", tile.Tileset.Name)
		}

		repackedTile := &Tile{
			ID:          LocalTileID(i),
			Type:        tile.Type,
			ObjectGroup: tile.ObjectGroup,
			Properties:  tile.Properties,
			Tileset:     repackedTileset,
		}
		for _, frame := range tile.Animation {
			animTileID := tile.Tileset.GlobalTileID(frame.TileID)
			repackedAnimTileID := p.repackedTileID[animTileID]
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

func (p *Packer) RepackedTileID(id GlobalTileID) GlobalTileID {
	return p.repackedTileID[id]
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
