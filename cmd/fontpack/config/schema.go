package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

func Read(filePath string) []Font {
	rawConfig, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Failed to read font config: %v", err)
	}

	var res []Font
	err = yaml.Unmarshal(rawConfig, &res)
	if err != nil {
		log.Fatalf("Failed to parse font config: %v", err)
	}

	return res
}

type Font struct {
	Name          string   `yaml:"name"`
	Size          int      `yaml:"size"`
	TopOffset     int      `yaml:"top_offset"`
	LineSpacing   int      `yaml:"line_spacing"`
	LetterSpacing int      `yaml:"letter_spacing"`
	SpaceWidth    int      `yaml:"space_width"`
	Letters       []string `yaml:"letters"`
}
