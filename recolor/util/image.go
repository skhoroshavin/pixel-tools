package util

import (
	"image"
	"image/png"
	"log"
	"os"
	"recolor/color"
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

func NumberUniqueColors(img image.Image) int {
	colors := make(map[color.SRGB]struct{})
	bounds := img.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			colors[color.SRGBFromImg(img.At(x, y))] = struct{}{}
		}
	}
	return len(colors)
}
