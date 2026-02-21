package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"pixel-tools/cmd/tilepack/atlas"
	"pixel-tools/pkg/file/tmj"
	"pixel-tools/pkg/file/tmx"
	"pixel-tools/pkg/file/tsx"
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
		a := atlas.New(shortName, src.Tilesets)
		repack(src, a)
		dst := tmj.ConvertFromTMX(src)

		a.Save(outputAtlas)
		dst.Save(outputTmj)

		return nil
	})
	if err != nil {
		log.Fatalf("Failed to walk input directory: %v", err)
	}
}

func repack(m *tmx.Map, a *atlas.Atlas) {
	for _, layer := range m.Layers {
		switch {
		case layer.IsTileLayer():
			for i, tileID := range layer.Data.Decoded {
				if tileID != 0 {
					layer.Data.Decoded[i] = a.UseTile(tileID.WithoutFlags()).WithFlags(tileID.Flags())
				}
			}
		case layer.IsObjectGroup():
			for i, obj := range layer.Objects {
				if obj.GID != 0 {
					layer.Objects[i].GID = a.UseSprite(obj.GID.WithoutFlags()).WithFlags(obj.GID.Flags())
				}
			}
		}
	}

	m.Tilesets = []*tsx.Tileset{a.Pack()}
}
