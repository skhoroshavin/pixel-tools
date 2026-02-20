package tmj

import (
	"encoding/json"
	"log"
	"os"

	"pixel-tools/cmd/tilepack/tmx"
)

type Map struct {
	Version          string    `json:"version"`
	CompressionLevel int       `json:"compressionlevel"`
	Height           int       `json:"height"`
	Infinite         bool      `json:"infinite"`
	TileWidth        int       `json:"tilewidth"`
	TileHeight       int       `json:"tileheight"`
	Width            int       `json:"width"`
	Orientation      string    `json:"orientation"`
	Layers           []any     `json:"layers"` // tile + object layers
	Tilesets         []Tileset `json:"tilesets"`
}

func (m *Map) Save(path string) {
	raw, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		log.Fatalf("Failed to encode map: %v", err)
	}
	err = os.WriteFile(path, raw, 0644)
	if err != nil {
		log.Fatalf("Failed to write map: %v", err)
	}
}

func ConvertFromTMX(src *tmx.Map) Map {
	res := Map{
		Version:          src.Version,
		CompressionLevel: -1,
		Height:           src.Height,
		Width:            src.Width,
		TileWidth:        src.TileWidth,
		TileHeight:       src.TileHeight,
		Orientation:      src.Orientation,
		Infinite:         false,
	}

	layerID := 1

	// Convert layers
	for _, layer := range src.Layers {
		switch {
		case layer.IsTileLayer():
			res.Layers = append(res.Layers, convertTileLayer(layer.AsTileLayer(), layerID))
		case layer.IsObjectGroup():
			res.Layers = append(res.Layers, convertObjectGroup(layer.AsObjectGroup()))
		}
		layerID++
	}

	// Convert tilesets
	for _, ts := range src.Tilesets {
		res.Tilesets = append(res.Tilesets, convertTileset(ts))
	}
	return res
}
