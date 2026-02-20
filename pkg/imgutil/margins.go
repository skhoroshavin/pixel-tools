package imgutil

import (
	"image"
)

// GetLeftMargin returns the number of empty columns on the left side
// of the image.
func GetLeftMargin(img image.Image) int {
	bounds := img.Bounds()
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			_, _, _, a := img.At(x, y).RGBA()
			if a != 0 {
				return x - bounds.Min.X
			}
		}
	}
	return bounds.Max.X - bounds.Min.X
}

// GetRightMargin returns the number of empty columns on the right side
// of the image.
func GetRightMargin(img image.Image) int {
	bounds := img.Bounds()
	for x := bounds.Max.X - 1; x >= bounds.Min.X; x-- {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			_, _, _, a := img.At(x, y).RGBA()
			if a != 0 {
				return bounds.Max.X - 1 - x
			}
		}
	}
	return bounds.Max.X - bounds.Min.X
}

// GetTopMargin returns the number of empty rows on the top side
// of the image.
func GetTopMargin(img image.Image) int {
	bounds := img.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			_, _, _, a := img.At(x, y).RGBA()
			if a != 0 {
				return y - bounds.Min.Y
			}
		}
	}
	return bounds.Max.Y - bounds.Min.Y
}

// GetBottomMargin returns the number of empty rows on the bottom side
// of the image.
func GetBottomMargin(img image.Image) int {
	bounds := img.Bounds()
	for y := bounds.Max.Y - 1; y >= bounds.Min.Y; y-- {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			_, _, _, a := img.At(x, y).RGBA()
			if a != 0 {
				return bounds.Max.Y - 1 - y
			}
		}
	}
	return bounds.Max.Y - bounds.Min.Y
}
