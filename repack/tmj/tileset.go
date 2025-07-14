package tmj

import "repack/tsx"

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

type Tile struct {
	ID         tsx.LocalTileID `json:"id"`
	Animation  []Frame         `json:"animation,omitempty"`
	Properties []Property      `json:"properties,omitempty"`
}

type Frame struct {
	TileID   tsx.LocalTileID `json:"tileid"`
	Duration int             `json:"duration"`
}

type Property struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Value string `json:"value"`
}
