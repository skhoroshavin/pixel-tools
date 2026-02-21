package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"pixel-tools/cmd/tilepack/builder"
	"pixel-tools/pkg/file/tmx"
	"strings"
)

func main() {
	if len(os.Args) != 3 {
		println("Usage: tilepack <input-dir> <output-dir>")
		os.Exit(1)
	}

	inputDir := os.Args[1]
	outputDir := os.Args[2]

	fmt.Println("Input directory:", inputDir)
	fmt.Println("Output directory:", outputDir)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	entries, err := os.ReadDir(inputDir)
	if err != nil {
		log.Fatalf("Failed to read input directory: %v", err)
	}

	b := builder.New()
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".tmx") {
			continue
		}

		inputTmx := filepath.Join(inputDir, entry.Name())
		baseName := strings.TrimSuffix(entry.Name(), ".tmx")

		fmt.Println("Processing", entry.Name())
		src := tmx.Load(inputTmx)
		b.AddTilemap(baseName, src)
	}

	fmt.Println("Saving to", outputDir)
	b.Save(outputDir)
}
