package imgutil

import "image"

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
