package main

import (
	"flag"
	"fmt"
	"image"
	"log"
	"maps"
	"os"
	"path"
	"path/filepath"
	"pixel-tools/cmd/atlaspack/config"
	"pixel-tools/pkg/atlas"
	"pixel-tools/pkg/file/png"
	"slices"
)

func main() {
	var padding int
	flag.IntVar(&padding, "padding", 0, "padding between sprites in the atlas")
	flag.Usage = func() {
		_, _ = fmt.Fprintln(os.Stderr, "Usage: atlaspack [options] <config-file> <output-base>")
		flag.PrintDefaults()
	}
	flag.Parse()

	if flag.NArg() != 2 {
		flag.Usage()
		os.Exit(1)
	}
	configFile := flag.Arg(0)
	outputBase := flag.Arg(1)

	fmt.Println("Input config:", configFile)
	fmt.Println("Output base:", outputBase)
	if padding > 0 {
		fmt.Println("Padding:", padding)
	}
	if err := os.MkdirAll(path.Dir(outputBase), 0755); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	b := NewBuilder(filepath.Dir(configFile), padding)

	for _, frame := range config.Read(configFile) {
		b.AddFrame(&frame)
	}

	b.Save(outputBase)
}

func NewBuilder(inputDir string, padding int) *Builder {
	return &Builder{
		atlas:    atlas.New(atlas.Config{Padding: padding}),
		inputDir: inputDir,
	}
}

type Builder struct {
	atlas    *atlas.Atlas
	inputDir string
}

func (b *Builder) AddFrame(frame *config.Frame) {
	img := b.loadImage(frame)

	if frame.Spritesheet == nil {
		b.atlas.AddSprite(atlas.SpriteConfig{
			Name:      frame.Name,
			Image:     img,
			NineSlice: frame.NineSlice,
		})
		return
	}

	ss := NewSpritesheet(img, frame.Spritesheet)

	for _, sprite := range ss.UniqueSprites() {
		name := frame.Name + sprite.Name
		b.atlas.AddSprite(atlas.SpriteConfig{
			Name:      name,
			Image:     sprite.Image,
			NineSlice: frame.NineSlice,
		})
	}

	for _, ref := range ss.SpriteRefs() {
		name := frame.Name + ref.Name
		source := frame.Name + ss.UniqueSprites()[ref.SourceIndex].Name
		if name != source {
			b.atlas.AddSpriteRef(name, source, nil)
		}
	}
}

func (b *Builder) Save(outputBase string) {
	b.atlas.Pack()
	b.atlas.Save(outputBase)
}

func (b *Builder) loadImage(frame *config.Frame) *image.RGBA {
	imagePath := frame.Image
	if imagePath == "" {
		imagePath = frame.Name + ".png"
	}
	return png.Read(filepath.Join(b.inputDir, imagePath))
}

func NewSpritesheet(img *image.RGBA, cfg *config.Spritesheet) *Spritesheet {
	bounds := img.Bounds()
	cols := bounds.Dx() / cfg.SpriteWidth
	rows := bounds.Dy() / cfg.SpriteHeight

	res := &Spritesheet{
		img:           img,
		rows:          rows,
		cols:          cols,
		width:         cfg.SpriteWidth,
		height:        cfg.SpriteHeight,
		usedPositions: make(map[int]int),
	}
	res.build(cfg)
	return res
}

type Spritesheet struct {
	img *image.RGBA

	rows, cols    int
	width, height int

	usedPositions map[int]int // sprite index in source sheet -> position in unique sprites slice
	uniqueSprites []Sprite
	spriteRefs    []SpriteRef
}

type Sprite struct {
	Name  string
	Image image.Image
}

type SpriteRef struct {
	Name        string
	SourceIndex int
}

func (ss *Spritesheet) UniqueSprites() []Sprite { return ss.uniqueSprites }

func (ss *Spritesheet) SpriteRefs() []SpriteRef { return ss.spriteRefs }

func (ss *Spritesheet) build(cfg *config.Spritesheet) {
	if cfg.SpriteNames != nil {
		ss.buildNamed(cfg.SpriteNames)
		return
	}

	start := cfg.StartSprite
	count := cfg.SpriteCount

	maxCount := ss.cols * ss.rows
	if count == 0 || (start+count) > maxCount {
		count = maxCount - start
	}

	for i := range count {
		ss.addSprite(fmt.Sprintf("%02d", i), start+i)
	}
}

func (ss *Spritesheet) buildNamed(names map[string]any) {
	for _, name := range slices.Sorted(maps.Keys(names)) {
		val := names[name]
		name = "_" + name
		switch v := val.(type) {
		case int:
			ss.addSprite(name, v)
		case []any:
			for i, item := range v {
				if idx, ok := item.(int); ok {
					ss.addSprite(fmt.Sprintf("%s%02d", name, i), idx)
				}
			}
		}
	}
}

func (ss *Spritesheet) addSprite(name string, pos int) {
	ref, isRef := ss.usedPositions[pos]
	if !isRef {
		ref = len(ss.uniqueSprites)
		ss.usedPositions[pos] = ref
	}

	ss.spriteRefs = append(ss.spriteRefs, SpriteRef{Name: name, SourceIndex: ref})
	if isRef {
		return
	}

	col := pos % ss.cols
	row := pos / ss.cols
	rect := image.Rect(col*ss.width, row*ss.height, (col+1)*ss.width, (row+1)*ss.height)
	img := ss.img.SubImage(rect)
	ss.uniqueSprites = append(ss.uniqueSprites, Sprite{Name: name, Image: img})
}
