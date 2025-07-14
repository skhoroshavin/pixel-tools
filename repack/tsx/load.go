package tsx

import (
	"encoding/xml"
	"log"
	"os"
	"path/filepath"
)

func Load(path string) *Tileset {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Failed to read tileset %s: %v", path, err)
	}

	var res Tileset
	if err = xml.Unmarshal(data, &res); err != nil {
		log.Fatalf("Failed to parse tileset %s: %v", path, err)
	}
	res.PostLoad(filepath.Dir(path))

	return &res
}
