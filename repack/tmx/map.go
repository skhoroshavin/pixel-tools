package tmx

import (
	"path/filepath"
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
	repacker := tsx.NewRepacker(m.Tilesets)
	for _, layer := range m.Layers {
		for _, tileID := range layer.Data.Decoded {
			if tileID == 0 {
				continue
			}
			repacker.UseTileID(tileID)
		}
	}

	m.Tilesets = []*tsx.Tileset{repacker.BuildNewTileset(name)}
	for i, layer := range m.Layers {
		newData := make([]tsx.GlobalTileID, len(layer.Data.Decoded))
		for j, tileID := range layer.Data.Decoded {
			newData[j] = repacker.RepackedTileID(tileID)
		}
		m.Layers[i].Data.Decoded = newData
	}
}

func (m *Map) SaveImages(path string) {
	for _, tileset := range m.Tilesets {
		util.WriteImage(tileset.Image.Data, filepath.Join(path, tileset.Image.Source))
	}
}
