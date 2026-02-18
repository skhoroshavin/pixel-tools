package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"pixel-tools/cmd/repack/tmj"
	"pixel-tools/cmd/repack/tmx"
)

func main() {
	if len(os.Args) != 3 {
		println("Usage: repack <input-dir> <output-dir>")
		os.Exit(1)
	}

	inputDir := os.Args[1]
	outputDir := os.Args[2]

	fmt.Println("Input directory:", inputDir)
	fmt.Println("Output directory:", outputDir)

	err := filepath.Walk(inputDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || !strings.HasSuffix(path, ".tmx") {
			return nil
		}

		inputTmx := path
		fullName := strings.TrimSuffix(strings.TrimPrefix(inputTmx, inputDir), ".tmx")
		shortName := filepath.Base(fullName)

		outputTmj := filepath.Join(outputDir, fullName+".tmj")
		outputAtlas := strings.TrimSuffix(outputTmj, ".tmj")
		tmjDir := filepath.Dir(outputTmj)

		fmt.Println("Processing", strings.TrimPrefix(inputTmx, inputDir))

		err = os.MkdirAll(tmjDir, 0775)
		if err != nil {
			log.Fatalf("Failed to create output directory: %v", err)
		}

		src := tmx.Load(inputTmx)
		src.Repack(shortName)
		dst := tmj.ConvertFromTMX(src)

		src.SaveAtlas(outputAtlas)
		dst.Save(outputTmj)

		return nil
	})
	if err != nil {
		log.Fatalf("Failed to walk input directory: %v", err)
	}
}
