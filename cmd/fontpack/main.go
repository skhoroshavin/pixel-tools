package main

import (
	"flag"
	"fmt"
	"image"
	"log"
	"os"
	"path/filepath"
	"pixel-tools/cmd/fontpack/config"
	"pixel-tools/pkg/atlas"
	"pixel-tools/pkg/file/bmfont"
	"pixel-tools/pkg/file/png"
)

func main() {
	var padding int
	flag.IntVar(&padding, "padding", 0, "padding between sprites in the atlas")
	flag.Usage = func() {
		_, _ = fmt.Fprintln(os.Stderr, "Usage: fontpack [options] <font-config> <output-dir>")
		flag.PrintDefaults()
	}
	flag.Parse()

	if flag.NArg() != 2 {
		flag.Usage()
		os.Exit(1)
	}

	fontConfigFile := flag.Arg(0)
	outputDir := flag.Arg(1)

	fmt.Println("Input config:", fontConfigFile)
	fmt.Println("Output directory:", outputDir)
	if padding > 0 {
		fmt.Println("Padding:", padding)
	}
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	fa := NewFontAtlas(filepath.Dir(fontConfigFile), padding)
	for _, font := range config.Read(fontConfigFile) {
		fa.AddFont(font)
	}

	fa.Build()
	outputImage := filepath.Join(outputDir, "fonts.png")
	fa.SaveImage(outputImage)
	fa.SaveBMFonts(outputDir)
}

func NewFontAtlas(inputDir string, padding int) *FontAtlas {
	return &FontAtlas{
		atlas:    atlas.New(atlas.Config{Padding: padding}),
		inputDir: inputDir,
		fonts:    make([]config.Font, 0),
	}
}

type FontAtlas struct {
	atlas    *atlas.Atlas
	inputDir string

	fonts      []config.Font
	bmfonts    []*bmfont.BMFont
	minYOffset []int
}

func (fa *FontAtlas) AddFont(font config.Font) {
	img := png.Read(filepath.Join(fa.inputDir, font.Name+".png"))

	minYOffset := font.Size
	maxBottom := 0

	for y, str := range font.Letters {
		x := 0
		for _, chr := range str {
			glyphName := fmt.Sprintf("%s_%d", font.Name, chr)
			glyphImg := img.SubImage(image.Rect(x*font.Size, y*font.Size, (x+1)*font.Size, (y+1)*font.Size))
			glyph := fa.atlas.AddSprite(atlas.SpriteConfig{Name: glyphName, Image: glyphImg, NineSlice: nil, Data: nil})
			if glyph.Trimmed {
				minYOffset = min(minYOffset, glyph.SpriteSourceSize.Y)
				maxBottom = max(maxBottom, glyph.Frame.H+glyph.SpriteSourceSize.Y)
			} else {
				minYOffset = min(minYOffset, 0)
				maxBottom = max(maxBottom, font.Size)
			}
			x++
		}
	}

	bmf := bmfont.New(font.Name, font.Size, "")
	bmf.Common.LineHeight = maxBottom - minYOffset + font.LineSpacing
	bmf.AddChar(32, atlas.Rect{}, 0, 0, font.SpaceWidth+font.LetterSpacing)

	fa.fonts = append(fa.fonts, font)
	fa.bmfonts = append(fa.bmfonts, bmf)
	fa.minYOffset = append(fa.minYOffset, minYOffset)
}

func (fa *FontAtlas) Build() {
	fa.atlas.Pack()

	for i, font := range fa.fonts {
		bmf := fa.bmfonts[i]
		minYOffset := fa.minYOffset[i]

		for _, str := range font.Letters {
			for _, chr := range str {
				glyphName := fmt.Sprintf("%s_%d", font.Name, chr)
				glyph := fa.atlas.GetSprite(glyphName)
				if glyph == nil {
					continue
				}

				if !glyph.Trimmed {
					bmf.AddChar(chr, glyph.Frame, 0, 0,
						font.Size+font.LetterSpacing)
					continue
				}

				bmf.AddChar(chr, glyph.Frame,
					glyph.SpriteSourceSize.X,
					glyph.SpriteSourceSize.Y-minYOffset,
					glyph.SpriteSourceSize.X+glyph.Frame.W+font.LetterSpacing)
			}
		}
	}
}

func (fa *FontAtlas) SaveImage(filePath string) {
	fa.atlas.SaveImage(filePath)
	for _, bmf := range fa.bmfonts {
		bmf.Pages.Page[0].File = filepath.Base(filePath)
	}
}

func (fa *FontAtlas) SaveBMFonts(outputDir string) {
	for _, bf := range fa.bmfonts {
		bf.Save(filepath.Join(outputDir, bf.Info.Face+".bmfont"))
	}
}
