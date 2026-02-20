package png

import (
	"image"
	"image/png"
	"log"
	"os"
)

// Read loads a PNG image from the given path and returns it as *image.RGBA.
// It logs a fatal error if the file cannot be opened or decoded.
func Read(filePath string) *image.RGBA {
	f, err := os.Open(filePath)
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

// Write saves an image to the given path as PNG.
// It logs a fatal error if the file cannot be created or encoded.
func Write(img image.Image, filePath string) {
	f, err := os.Create(filePath)
	if err != nil {
		log.Fatalf("Failed to open image: %v", err)
	}
	defer f.Close()
	err = png.Encode(f, img)
	if err != nil {
		log.Fatalf("Failed to encode image: %v", err)
	}
}
