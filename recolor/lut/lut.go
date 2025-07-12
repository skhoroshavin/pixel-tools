package lut

import (
	"encoding/json"
	"image"
	"log"
	"math"
	"os"
	"recolor/color"
)

func New() LUT {
	return LUT{
		lut:   make(map[color.SRGB]ColorDistribution),
		cache: make(map[color.SRGB]color.SRGB),
	}
}

type LUT struct {
	lut   map[color.SRGB]ColorDistribution
	cache map[color.SRGB]color.SRGB
}

func (l *LUT) Save(path string) {
	data, err := json.MarshalIndent(l.lut, "", "  ")
	if err != nil {
		log.Fatalf("Failed to encode LUT: %v", err)
	}

	f, err := os.Create(path)
	if err != nil {
		log.Fatalf("Failed to open LUT for writing: %v", err)
	}
	defer f.Close()

	_, err = f.Write(data)
	if err != nil {
		log.Fatalf("Failed to write LUT: %v", err)
	}
}

func (l *LUT) Load(path string) {
	l.invalidateCache()

	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Failed to read LUT: %v", err)
	}

	err = json.Unmarshal(data, &l.lut)
	if err != nil {
		log.Fatalf("Failed to decode LUT: %v", err)
	}
}

func (l *LUT) AddColorMapping(src, dst color.SRGB) {
	l.invalidateCache()

	if outputEntry, found := l.lut[src]; found {
		outputEntry[dst] += 1.0
		return
	}

	l.lut[src] = ColorDistribution{dst: 1}
}

func (l *LUT) AddImageMapping(src, dst image.Image) {
	bounds := src.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			srcColor := color.SRGBFromImg(src.At(x, y))
			dstColor := color.SRGBFromImg(dst.At(x, y))
			l.AddColorMapping(srcColor, dstColor)
		}
	}
}

func (l *LUT) LookupColor(c color.SRGB) color.SRGB {
	res, found := l.cache[c]
	if found {
		return res
	}

	out, found := l.lut[c]
	if found {
		res = out.Value()
		l.cache[c] = res
		return res
	}

	var minDistValue, secondMinDistValue ColorDistribution
	minDist := math.MaxFloat64
	secondMinDist := math.MaxFloat64
	for in, out := range l.lut {
		d := color.Distance(c, in)
		if d > secondMinDist {
			continue
		}
		if d < minDist {
			secondMinDist = minDist
			secondMinDistValue = minDistValue

			minDist = d
			minDistValue = out
		} else {
			secondMinDist = d
			secondMinDistValue = out
		}
	}

	t := minDist / (minDist + secondMinDist)
	res = color.Interpolate(minDistValue.Value(), secondMinDistValue.Value(), t)
	l.cache[c] = res
	return res
}

func (l *LUT) ApplyToColor(c color.SRGB) color.SRGB {
	return l.LookupColor(c)
}

func (l *LUT) ApplyToImage(img image.Image) image.Image {
	bounds := img.Bounds()
	res := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			res.Set(x, y, l.ApplyToColor(color.SRGBFromImg(img.At(x, y))))
		}
	}
	return res
}

func (l *LUT) invalidateCache() {
	if len(l.cache) > 0 {
		l.cache = make(map[color.SRGB]color.SRGB)
	}
}
