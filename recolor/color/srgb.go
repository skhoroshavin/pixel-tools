package color

import (
	"fmt"
	"image/color"
	"strconv"
)

type SRGB color.RGBA

func (s SRGB) RGBA() (r, g, b, a uint32) {
	return color.RGBA(s).RGBA()
}

func (s SRGB) String() string {
	return fmt.Sprintf("#%02X%02X%02X%02X", s.R, s.G, s.B, s.A)
}

func (s SRGB) MarshalText() ([]byte, error) {
	return []byte(s.String()), nil
}

func (s *SRGB) UnmarshalText(b []byte) error {
	if len(b) != 9 {
		return fmt.Errorf("invalid color string: %s", string(b))
	}
	hex := string(b[1:])
	ri, err := strconv.ParseUint(hex[0:2], 16, 8)
	if err != nil {
		return err
	}
	gi, err := strconv.ParseUint(hex[2:4], 16, 8)
	if err != nil {
		return err
	}
	bi, err := strconv.ParseUint(hex[4:6], 16, 8)
	if err != nil {
		return err
	}
	ai, err := strconv.ParseUint(hex[6:8], 16, 8)
	if err != nil {
		return err
	}
	s.R, s.G, s.B, s.A = uint8(ri), uint8(gi), uint8(bi), uint8(ai)
	return nil
}

func SRGBFromImg(c color.Color) SRGB {
	return SRGB(color.RGBAModel.Convert(c).(color.RGBA))
}

func SRGBFromLinear(c Linear) SRGB {
	return SRGB{
		R: toGamma(c.r),
		G: toGamma(c.g),
		B: toGamma(c.b),
		A: uint8(c.a),
	}
}

func Interpolate(a, b SRGB, t float64) SRGB {
	return SRGBFromLinear(
		InterpolateLinear(LinearFromSRGB(a), LinearFromSRGB(b), t),
	)
}

func Distance(a, b SRGB) float64 {
	return DistanceLinear(LinearFromSRGB(a), LinearFromSRGB(b))
}
