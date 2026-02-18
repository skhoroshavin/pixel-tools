package util

import (
	"encoding/json"
	"log"
	"os"
)

func ReadFontConfig(path string) (res FontConfig) {
	rawConfig, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Failed to read font config: %v", err)
	}

	err = json.Unmarshal(rawConfig, &res)
	if err != nil {
		log.Fatalf("Failed to parse font config: %v", err)
	}

	return res
}

type FontConfig struct {
	Fonts []Font `json:"fonts"`
}

type Font struct {
	Name          string   `json:"name"`
	Size          int      `json:"size"`
	TopOffset     int      `json:"top_offset"`
	LineSpacing   int      `json:"line_spacing"`
	LetterSpacing int      `json:"letter_spacing"`
	SpaceWidth    int      `json:"space_width"`
	Letters       []string `json:"letters"`
}
