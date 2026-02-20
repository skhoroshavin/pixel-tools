package png

import (
	"image"
	"sync"
)

// Get loads an image from the given path with caching.
// Subsequent calls for the same path return the cached image.
func Get(path string) *image.RGBA {
	muImgCache.Lock()
	defer muImgCache.Unlock()

	if img, ok := imgCache[path]; ok {
		return img
	}

	img := Read(path)
	imgCache[path] = img
	return img
}

var (
	muImgCache sync.Mutex
	imgCache   = make(map[string]*image.RGBA)
)
