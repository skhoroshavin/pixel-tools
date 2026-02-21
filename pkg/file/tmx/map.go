package tmx

import "pixel-tools/pkg/file/tsx"

type Map struct {
	Version     string         `xml:"version,attr"`
	Orientation string         `xml:"orientation,attr"`
	RenderOrder string         `xml:"renderorder,attr"`
	Width       int            `xml:"width,attr"`
	Height      int            `xml:"height,attr"`
	TileWidth   int            `xml:"tilewidth,attr"`
	TileHeight  int            `xml:"tileheight,attr"`
	Tilesets    []*tsx.Tileset `xml:"tileset"`
	Layers      []*Layer       `xml:",any"`
	Properties  []tsx.Property `xml:"properties>property"`
}
