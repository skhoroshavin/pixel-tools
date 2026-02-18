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

func GetRightMargin(img image.Image, x, y, width, height int) int {
	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {
			_, _, _, a := img.At(x+width-1-i, y+j).RGBA()
			if a != 0 {
				return i
			}
		}
	}
	return width
}

func GetTopMargin(img image.Image, x, y, width, height int) int {
	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			_, _, _, a := img.At(x+j, y+i).RGBA()
			if a != 0 {
				return i
			}
		}
	}
	return height
}

func GetBottomMargin(img image.Image, x, y, width, height int) int {
	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			_, _, _, a := img.At(x+j, y+height-1-i).RGBA()
			if a != 0 {
				return i
			}
		}
	}
	return height
}
