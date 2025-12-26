package tmx

import (
	"path/filepath"
	"repack/atlas"
	"repack/tsx"
	"repack/util"
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
}

func (m *Map) Repack(name string) {
	objPacker := atlas.New(m.Tilesets)
	for _, objectGroup := range m.ObjectGroups {
		for i, obj := range objectGroup.Objects {
			if obj.GID != 0 {
				objectGroup.Objects[i].GID = objPacker.UseTileID(obj.GID)
			}
		}
	}

	tilePacker := tsx.NewPacker(m.Tilesets)
	for _, layer := range m.Layers {
		for i, tileID := range layer.Data.Decoded {
			if tileID != 0 {
				layer.Data.Decoded[i] = tilePacker.UseTileID(tileID)
			}
		}
	}

	m.Tilesets = []*tsx.Tileset{tilePacker.BuildNewTileset(name)}
}

func (m *Map) SaveImages(path string) {
	for _, tileset := range m.Tilesets {
		util.WriteImage(tileset.Image.Data, filepath.Join(path, tileset.Image.Source))
	}
}
