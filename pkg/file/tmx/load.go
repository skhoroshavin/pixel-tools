package tmx

import (
	"encoding/xml"
	"log"
	"os"
	"path/filepath"
	"pixel-tools/pkg/file/tsx"
)

func Load(filePath string) *Map {
	bytes, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}

	var res Map
	err = xml.Unmarshal(bytes, &res)
	if err != nil {
		log.Fatalf("Failed to parse file: %v", err)
	}

	for _, layer := range res.Layers {
		if layer.IsTileLayer() {
			layer.Data.Decode()
		}
	}

	for i, tileset := range res.Tilesets {
		if tileset.Source == "" {
			tileset.PostLoad(filepath.Dir(filePath))
		} else {
			external := tsx.Load(filepath.Join(filepath.Dir(filePath), tileset.Source))
			external.FirstGID = tileset.FirstGID
			res.Tilesets[i] = external
		}
	}

	return &res
}
