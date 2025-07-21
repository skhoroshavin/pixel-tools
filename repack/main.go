package main

import (
	"os"
	"path/filepath"
	"repack/tmj"
	"repack/tmx"
	"strings"
)

func main() {
	if len(os.Args) != 3 {
		println("Usage: repack <input.tmx> <output-dir>")
		os.Exit(1)
	}

	inputTmx := os.Args[1]
	outputDir := os.Args[2]

	outputName := strings.TrimSuffix(filepath.Base(inputTmx), ".tmx")

	src := tmx.Load(inputTmx)
	src.Repack(outputName)
	dst := tmj.ConvertFromTMX(src)

	src.SaveImages(outputDir)
	dst.Save(filepath.Join(outputDir, outputName+".tmj"))
}
