package tmx

import (
	"encoding/xml"
	"log"
	"os"
	"path/filepath"
	"repack/atlas"
	"repack/tsx"
)

func Load(path string) *Map {
	bytes, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}

	var res Map
	err = xml.Unmarshal(bytes, &res)
	if err != nil {
		log.Fatalf("Failed to parse file: %v", err)
	}

	for _, layer := range res.Layers {
		layer.Decode()
	}

	for i, tileset := range res.Tilesets {
		if tileset.Source == "" {
			tileset.PostLoad(filepath.Dir(path))
		} else {
			external := tsx.Load(filepath.Join(filepath.Dir(path), tileset.Source))
			external.FirstGID = tileset.FirstGID
			res.Tilesets[i] = external
		}
	}

	res.atlas = atlas.New(res.Tilesets)
	return &res
}
