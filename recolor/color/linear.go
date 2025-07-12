package color

import "math"

type Linear struct {
	r, g, b, a float64
}

func (l Linear) Add(c Linear) Linear {
	l.r += c.r
	l.g += c.g
	l.b += c.b
	l.a += c.a
	return l
}

func (l Linear) MulC(v float64) Linear {
	l.r *= v
	l.g *= v
	l.b *= v
	l.a *= v
	return l
}

func LinearFromSRGB(c SRGB) Linear {
	return Linear{
		r: toLinear(c.R),
		g: toLinear(c.G),
		b: toLinear(c.B),
		a: float64(c.A),
	}
}

func InterpolateLinear(a, b Linear, t float64) Linear {
	return Linear{
		r: a.r*(1-t) + b.r*t,
		g: a.g*(1-t) + b.g*t,
		b: a.b*(1-t) + b.b*t,
		a: a.a*(1-t) + b.a*t,
	}
}

func DistanceLinear(a, b Linear) float64 {
	dr := a.r - b.r
	dg := a.g - b.g
	db := a.b - b.b
	da := 1000 * (a.a - b.a) // To make sure we don't mix different alphas
	return math.Sqrt(dr*dr + dg*dg + db*db + da*da)
}
