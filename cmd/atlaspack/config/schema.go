package config

import (
	"log"
	"os"
	"path/filepath"
	"pixel-tools/pkg/atlas"

	"gopkg.in/yaml.v3"
)

func Read(filePath string) []Frame {
	rawConfig, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Failed to read atlas config: %v", err)
	}

	var res []Frame
	err = yaml.Unmarshal(rawConfig, &res)
	if err != nil {
		log.Fatalf("Failed to parse atlas config: %v", err)
	}

	configDir := filepath.Dir(filePath)
	for i := range res {
		if res[i].Spritesheet != nil && res[i].Spritesheet.File != "" {
			externalPath := filepath.Join(configDir, res[i].Spritesheet.File)
			res[i].Spritesheet = LoadSpritesheet(externalPath)
		}
	}

	return res
}

func LoadSpritesheet(filePath string) *Spritesheet {
	rawData, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Failed to read spritesheet config: %v", err)
	}

	var ss Spritesheet
	err = yaml.Unmarshal(rawData, &ss)
	if err != nil {
		log.Fatalf("Failed to parse spritesheet config: %v", err)
	}

	return &ss
}

type Frame struct {
	Name        string           `yaml:"name"`
	Image       string           `yaml:"image,omitempty"`
	Spritesheet *Spritesheet     `yaml:"spritesheet,omitempty"`
	NineSlice   *atlas.NineSlice `yaml:"nineslice,omitempty"`
}

type Spritesheet struct {
	File         string         `yaml:"file,omitempty"`
	SpriteWidth  int            `yaml:"sprite_width,omitempty"`
	SpriteHeight int            `yaml:"sprite_height,omitempty"`
	StartSprite  int            `yaml:"start_sprite,omitempty"`
	SpriteCount  int            `yaml:"sprite_count,omitempty"`
	SpriteNames  map[string]any `yaml:"sprite_names,omitempty"`
}
