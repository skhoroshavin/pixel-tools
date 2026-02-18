package lut

import (
	"pixel-tools/cmd/recolor/color"
)

type ColorDistribution map[color.SRGB]float64

func (c ColorDistribution) Value() color.SRGB {
	sumValues := color.Linear{}
	sumWeights := 0.0
	for c, w := range c {
		sumValues = sumValues.Add(color.LinearFromSRGB(c).MulC(w))
		sumWeights += w
	}
	return color.SRGBFromLinear(sumValues.MulC(1.0 / sumWeights))
}
