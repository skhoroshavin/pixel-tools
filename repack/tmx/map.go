package tmx

import (
	"repack/atlas"
	"repack/tsx"
)

type Map struct {
	Version      string            `xml:"version,attr"`
	Orientation  string            `xml:"orientation,attr"`
	RenderOrder  string            `xml:"renderorder,attr"`
	Width        int               `xml:"width,attr"`
	Height       int               `xml:"height,attr"`
	TileWidth    int               `xml:"tilewidth,attr"`
	TileHeight   int               `xml:"tileheight,attr"`
	Tilesets     []*tsx.Tileset    `xml:"tileset"`
	Layers       []*Layer          `xml:"layer"`
	ObjectGroups []tsx.ObjectGroup `xml:"objectgroup"`
	Properties   []tsx.Property    `xml:"properties>property"`

	atlas *atlas.Atlas
}

func (m *Map) Repack(name string) {
	tilePacker := tsx.NewPacker(m.Tilesets)
	for _, layer := range m.Layers {
		for i, tileID := range layer.Data.Decoded {
			if tileID != 0 {
				layer.Data.Decoded[i] = tilePacker.UseTileID(tileID)
			}
		}
	}
	tileset := tilePacker.BuildNewTileset(name)
	m.Tilesets = []*tsx.Tileset{tileset}
	m.atlas.UseTileset(tileset)

	for _, objectGroup := range m.ObjectGroups {
		for i, obj := range objectGroup.Objects {
			if obj.GID != 0 {
				objectGroup.Objects[i].GID = m.atlas.UseTile(obj.GID)
			}
		}
	}
	m.atlas.Pack()
}

func (m *Map) SaveAtlas(baseName string) {
	m.atlas.Save(baseName)
}
