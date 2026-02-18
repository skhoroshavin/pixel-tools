package imgutil

import (
	"image"
	"image/png"
	"log"
	"os"
	"sync"
)

// ReadImage loads a PNG image from the given path.
// It logs a fatal error if the file cannot be opened or decoded.
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
func GetImage(path string) image.Image {
	muImgCache.Lock()
	defer muImgCache.Unlock()

	if img, ok := imgCache[path]; ok {
		return img
	}

	img := ReadImage(path)
	imgCache[path] = img
	return img
}

// GetRightMargin returns the number of empty columns on the right side
// of the specified region within the image.
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

// GetTopMargin returns the number of empty rows on the top side
// of the specified region within the image.
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

// GetBottomMargin returns the number of empty rows on the bottom side
// of the specified region within the image.
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

// NumberUniqueColors returns the count of unique colors in the image.
func NumberUniqueColors(img image.Image) int {
	colors := make(map[uint32]struct{})
	bounds := img.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			color := (r << 24) | (g << 16) | (b << 8) | a
			colors[color] = struct{}{}
		}
	}
	return len(colors)
}

var (
	muImgCache sync.Mutex
	imgCache   = make(map[string]image.Image)
)
