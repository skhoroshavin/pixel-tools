package util

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

type Config struct {
	Reference ReferenceConfig `yaml:"reference"`
	Recolor   RecolorConfig   `yaml:"recolor"`
}
type ReferenceConfig struct {
	BaseLUT         string `yaml:"baseLUT"`
	Folder          string `yaml:"folder"`
	OriginalSuffix  string `yaml:"originalSuffix"`
	RecoloredSuffix string `yaml:"recoloredSuffix"`
	ResultingLUT    string `yaml:"resultingLUT"`
}

type RecolorConfig struct {
	Folder          string `yaml:"folder"`
	OriginalSuffix  string `yaml:"originalSuffix"`
	RecoloredSuffix string `yaml:"recoloredSuffix"`
}

func LoadConfig(path string) *Config {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		log.Fatalf("Failed to parse config file: %v", err)
	}

	ref := cfg.Reference
	recolor := cfg.Recolor

	if ref.BaseLUT == "" && ref.Folder == "" {
		log.Fatalf("Either 'baseLUT' or 'folder' must be specified in reference config")
	}

	if ref.Folder != "" {
		if ref.OriginalSuffix == "" {
			log.Fatalf("reference.originalSuffix must be specified when reference.folder is set")
		}
		if ref.RecoloredSuffix == "" {
			log.Fatalf("reference.recoloredSuffix must be specified when reference.folder is set")
		}
	}

	if recolor.Folder != "" {
		if recolor.OriginalSuffix == "" {
			log.Fatalf("recolor.originalSuffix must be specified when recolor.folder is set")
		}
		if recolor.RecoloredSuffix == "" {
			log.Fatalf("recolor.recoloredSuffix must be specified when recolor.folder is set")
		}
	}

	return &cfg
}
