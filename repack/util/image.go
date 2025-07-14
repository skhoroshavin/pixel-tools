package util

import (
	"image"
	"image/png"
	"log"
	"os"
)

func ReadImage(path string) image.Image {
	f, err := os.Open(path)
	if err != nil {
		log.Fatalf("Failed to open image: %v", err)
	}
	defer f.Close()
	img, err := png.Decode(f)
	if err != nil {
		log.Fatalf("Failed to decode image: %v", err)
	}
	return img
}

func WriteImage(img image.Image, path string) {
	f, err := os.Create(path)
	if err != nil {
		log.Fatalf("Failed to open image: %v", err)
	}
	defer f.Close()
	err = png.Encode(f, img)
	if err != nil {
		log.Fatalf("Failed to encode image: %v", err)
	}
}
