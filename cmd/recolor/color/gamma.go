package color

import "math"

const gamma = 2.2
const invGamma = 1.0 / gamma

func toLinear(c uint8) float64 {
	return math.Pow(float64(c), gamma)
}

func toGamma(c float64) uint8 {
	return uint8(math.Pow(c, invGamma))
}
