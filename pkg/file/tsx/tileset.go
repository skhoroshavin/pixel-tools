package tsx

import (
	"log"
)

type Tileset struct {
	FirstGID   GlobalTileID `xml:"firstgid,attr"`
	Source     string       `xml:"source,attr"` // only used if external
	Name       string       `xml:"name,attr"`
	TileWidth  int          `xml:"tilewidth,attr"`
	TileHeight int          `xml:"tileheight,attr"`
	TileCount  int          `xml:"tilecount,attr"`
	Columns    int          `xml:"columns,attr"`
	Image      *Image       `xml:"image"`
	Tiles      []*Tile      `xml:"tile"`
}

func (ts *Tileset) HasSameTileSize(other *Tileset) bool {
	return ts.TileWidth == other.TileWidth && ts.TileHeight == other.TileHeight
}

func (ts *Tileset) PostLoad(basePath string) {
	if ts.Image != nil {
		ts.Image.PostLoad(basePath)
	}

	for _, tile := range ts.Tiles {
		tile.PostLoad(basePath, ts)
	}
}

func (ts *Tileset) HasTile(id GlobalTileID) bool {
	return id >= ts.FirstGID && id < ts.FirstGID+GlobalTileID(ts.TileCount)
}

func (ts *Tileset) LocalTileID(id GlobalTileID) LocalTileID {
	if !ts.HasTile(id) {
		log.Fatalf("Cannot find tile %d in tileset %s", id, ts.Name)
	}
	return LocalTileID(id - ts.FirstGID)
}

func (ts *Tileset) GlobalTileID(id LocalTileID) GlobalTileID {
	return GlobalTileID(id) + ts.FirstGID
}

func (ts *Tileset) Tile(id GlobalTileID) *Tile {
	if !ts.HasTile(id) {
		return nil
	}

	localID := ts.LocalTileID(id)
	for _, tile := range ts.Tiles {
		if tile.ID == localID {
			return tile
		}
	}

	tile := &Tile{
		ID:      localID,
		Tileset: ts,
	}
	ts.Tiles = append(ts.Tiles, tile)
	return tile
}
