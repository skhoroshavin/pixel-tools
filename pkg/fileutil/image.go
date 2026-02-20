package fileutil

import (
	"image"
	"image/png"
	"log"
	"os"
	"sync"
)

// ReadImage loads a PNG image from the given path and returns it as *image.RGBA.
// It logs a fatal error if the file cannot be opened or decoded.
func ReadImage(path string) *image.RGBA {
	f, err := os.Open(path)
	if err != nil {
		log.Fatalf("Failed to open image: %v", err)
	}
	defer f.Close()

	img, err := png.Decode(f)
	if err != nil {
		log.Fatalf("Failed to decode image: %v", err)
	}

	rgba, ok := img.(*image.RGBA)
	if ok {
		return rgba
	}

	bounds := img.Bounds()
	rgba = image.NewRGBA(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			rgba.Set(x, y, img.At(x, y))
		}
	}
	return rgba
}

// WriteImage saves an image to the given path as PNG.
// It logs a fatal error if the file cannot be created or encoded.
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

// GetImage loads an image from the given path with caching.
// Subsequent calls for the same path return the cached image.
func GetImage(path string) *image.RGBA {
	muImgCache.Lock()
	defer muImgCache.Unlock()

	if img, ok := imgCache[path]; ok {
		return img
	}

	img := ReadImage(path)
	imgCache[path] = img
	return img
}

var (
	muImgCache sync.Mutex
	imgCache   = make(map[string]*image.RGBA)
)
