package main

import (
	"fmt"
	"image"
	"os"
	"pixel-tools/pkg/fileutil"
	"pixel-tools/pkg/imgutil"

	"pixel-tools/cmd/fontpack/util"
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

	fontConfig := util.ReadFontConfig(fontConfigFile)
	for _, font := range fontConfig.Fonts {
		img := fileutil.ReadImage(font.Name + ".png")

		bmfont := util.BMFont{
			Info: util.FontInfo{
				Face: font.Name,
				Size: font.Size,
			},
			Common: util.Common{
				Pages: 1,
			},
			Pages: util.Pages{
				Page: []util.Page{
					{
						File: font.Name + ".png",
					},
				},
			},
		}

		bmfont.Chars.Count = 1
		bmfont.Chars.Char = append(bmfont.Chars.Char, util.Char{
			ID:       32,
			Letter:   " ",
			X:        0,
			Y:        0,
			Width:    0,
			Height:   0,
			XAdvance: font.SpaceWidth + font.LetterSpacing,
		})

		minTopMargin := font.Size
		minBottomMargin := font.Size
		for y, str := range font.Letters {
			x := 0
			for _, chr := range str {
				bmfont.Chars.Count++
				glyph := img.SubImage(image.Rect(x*font.Size, y*font.Size, (x+1)*font.Size, (y+1)*font.Size))
				rightMargin := imgutil.GetRightMargin(glyph)
				topMargin := imgutil.GetTopMargin(glyph)
				bottomMargin := imgutil.GetBottomMargin(glyph)

				bmfont.Chars.Char = append(bmfont.Chars.Char, util.Char{
					ID:       int(chr),
					X:        x * font.Size,
					Y:        y*font.Size + topMargin,
					Width:    font.Size - rightMargin,
					Height:   font.Size - bottomMargin - topMargin,
					XOffset:  0,
					YOffset:  topMargin,
					XAdvance: font.Size - rightMargin + font.LetterSpacing,
					Letter:   fmt.Sprintf("%c", chr),
				})

				minTopMargin = min(minTopMargin, topMargin)
				minBottomMargin = min(minBottomMargin, bottomMargin)
				x++
			}
		}
		bmfont.Common.LineHeight = font.Size + font.LineSpacing - minTopMargin - minBottomMargin
		for i := range bmfont.Chars.Char[1:] {
			bmfont.Chars.Char[i].YOffset -= minTopMargin
		}

		util.SaveBMFont(bmfont, outputDir+"/"+font.Name+".xml")
	}
}
