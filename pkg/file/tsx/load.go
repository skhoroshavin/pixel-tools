package tsx

import (
	"encoding/xml"
	"log"
	"os"
	"path/filepath"
)

func Load(filePath string) *Tileset {
	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Failed to read tileset %s: %v", filePath, err)
	}

	var res Tileset
	if err = xml.Unmarshal(data, &res); err != nil {
		log.Fatalf("Failed to parse tileset %s: %v", filePath, err)
	}
	res.PostLoad(filepath.Dir(filePath))

	return &res
}
