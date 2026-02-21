package main

import (
	"fmt"
	"image"
	"os"
	"path"
	"path/filepath"
	"pixel-tools/cmd/fontpack/config"
	"pixel-tools/pkg/atlas"
	"pixel-tools/pkg/file/bmfont"
	"pixel-tools/pkg/file/png"
	"pixel-tools/pkg/imgutil"
)

func main() {
	if len(os.Args) != 3 {
		println("Usage: fontpack <font-config> <output-dir>")
		os.Exit(1)
	}

	fontConfigFile := os.Args[1]
	outputDir := os.Args[2]

	fmt.Println("Input config:", fontConfigFile)
	fmt.Println("Output directory:", outputDir)

	fa := NewFontAtlas()
	for _, font := range config.Read(fontConfigFile) {
		fa.AddFont(font)
	}

	fa.Build()
	outputImage := filepath.Join(outputDir, "fonts.png")
	fa.SaveImage(outputImage)
	fa.SaveBMFonts(outputDir)
}

func NewFontAtlas() *FontAtlas {
	return &FontAtlas{
		atlas: atlas.New(),
		fonts: make([]config.Font, 0),
	}
}

type FontAtlas struct {
	atlas *atlas.Atlas

	fonts   []config.Font
	bmfonts []*bmfont.BMFont
}

func (fa *FontAtlas) AddFont(font config.Font) {
	fa.fonts = append(fa.fonts, font)
	bmf := bmfont.New(font.Name, font.Size, "")
	fa.bmfonts = append(fa.bmfonts, bmf)
	bmf.AddChar(32, 0, 0, 0, 0, 0, 0, font.SpaceWidth+font.LetterSpacing)

	img := png.Read(filepath.Join(filepath.Dir(os.Args[1]), font.Name+".png"))
	minTop := font.Size
	maxBottom := 0

	for y, str := range font.Letters {
		x := 0
		for _, chr := range str {
			glyphImg := img.SubImage(image.Rect(x*font.Size, y*font.Size, (x+1)*font.Size, (y+1)*font.Size))
			rightMargin := imgutil.GetRightMargin(glyphImg)
			topMargin := imgutil.GetTopMargin(glyphImg)
			bottomMargin := imgutil.GetBottomMargin(glyphImg)

			w := font.Size - rightMargin
			h := font.Size - bottomMargin - topMargin

			if w > 0 && h > 0 {
				if topMargin < minTop {
					minTop = topMargin
				}
				if font.Size-bottomMargin > maxBottom {
					maxBottom = font.Size - bottomMargin
				}

				actualGlyph := img.SubImage(image.Rect(x*font.Size, y*font.Size+topMargin, (x+1)*font.Size-rightMargin, (y+1)*font.Size-bottomMargin))
				spriteName := fmt.Sprintf("%s_%d", font.Name, chr)
				fa.atlas.AddSprite(spriteName, actualGlyph, nil, nil)
				bmf.AddChar(chr, 0, 0, w, h, 0, topMargin, w+font.LetterSpacing)
			}
			x++
		}
	}

	for i := range bmf.Chars.Char {
		if bmf.Chars.Char[i].ID != 32 {
			bmf.Chars.Char[i].YOffset -= minTop
			bmf.Chars.Char[i].YOffset += font.TopOffset
		} else {
			bmf.Chars.Char[i].YOffset = -minTop + font.TopOffset
		}
	}
	bmf.Common.LineHeight = maxBottom - minTop + font.LineSpacing
}

func (fa *FontAtlas) Build() {
	fa.atlas.Pack()

	for _, bmf := range fa.bmfonts {
		bmf.Common.Base = 0
		for i := range bmf.Chars.Char {
			char := &bmf.Chars.Char[i]
			if char.ID == 32 {
				continue
			}
			spriteName := fmt.Sprintf("%s_%d", bmf.Info.Face, char.ID)
			sprite := fa.atlas.GetSprite(spriteName)
			if sprite != nil {
				char.X = sprite.X
				char.Y = sprite.Y
			}
		}
		bmf.Chars.Count = len(bmf.Chars.Char)
	}
}

func (fa *FontAtlas) SaveImage(filePath string) {
	fa.atlas.SaveImage(filePath)
	for i := range fa.bmfonts {
		fa.bmfonts[i].Pages.Page[0].File = path.Base(filePath)
	}
}

func (fa *FontAtlas) SaveBMFonts(outputDir string) {
	for _, bf := range fa.bmfonts {
		bf.Save(filepath.Join(outputDir, bf.Info.Face+".bmfont"))
	}
}
