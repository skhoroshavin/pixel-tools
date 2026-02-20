package tmj

import "pixel-tools/cmd/tilepack/tsx"

type Tileset struct {
	FirstGID    tsx.GlobalTileID `json:"firstgid"`
	Name        string           `json:"name"`
	TileWidth   int              `json:"tilewidth"`
	TileHeight  int              `json:"tileheight"`
	TileCount   int              `json:"tilecount"`
	Columns     int              `json:"columns"`
	Image       string           `json:"image"`
	ImageWidth  int              `json:"imagewidth"`
	ImageHeight int              `json:"imageheight"`
	Tiles       []Tile           `json:"tiles,omitempty"`
}

func convertTileset(ts *tsx.Tileset) Tileset {
	return Tileset{
		FirstGID:    ts.FirstGID,
		Name:        ts.Name,
		TileWidth:   ts.TileWidth,
		TileHeight:  ts.TileHeight,
		TileCount:   ts.TileCount,
		Columns:     ts.Columns,
		Image:       ts.Image.Source,
		ImageWidth:  ts.Image.Width,
		ImageHeight: ts.Image.Height,
		Tiles:       convertTiles(ts.Tiles),
	}
}
