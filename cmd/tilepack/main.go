package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"pixel-tools/cmd/tilepack/builder"
	"pixel-tools/pkg/file/tmx"
	"strings"
)

func main() {
	var padding int
	flag.IntVar(&padding, "padding", 0, "padding between sprites in the atlas")
	flag.Usage = func() {
		_, _ = fmt.Fprintln(os.Stderr, "Usage: tilepack [options] <input-dir> <output-dir>")
		flag.PrintDefaults()
	}
	flag.Parse()

	if flag.NArg() != 2 {
		flag.Usage()
		os.Exit(1)
	}

	inputDir := flag.Arg(0)
	outputDir := flag.Arg(1)

	fmt.Println("Input directory:", inputDir)
	fmt.Println("Output directory:", outputDir)
	if padding > 0 {
		fmt.Println("Padding:", padding)
	}
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	entries, err := os.ReadDir(inputDir)
	if err != nil {
		log.Fatalf("Failed to read input directory: %v", err)
	}

	b := builder.New(padding)
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
