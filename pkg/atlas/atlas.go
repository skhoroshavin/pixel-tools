package atlas

import (
	"encoding/json"
	"image"
	"image/draw"
	"log"
	"os"
	"pixel-tools/pkg/file/png"
	"sort"
)

func New(name string) *Atlas {
	const initialSize = 16

	return &Atlas{
		name:        name,
		width:       initialSize,
		height:      initialSize,
		spriteIndex: make(map[string]int),
		skyline:     make([]int, initialSize),
	}
}

type Atlas struct {
	name string

	tiles   []Frame
	sprites []Frame

	spriteIndex map[string]int

	width   int
	height  int
	skyline []int
}

type Frame struct {
	Name  string
	Image image.Image
	X, Y  int
	W, H  int
	Data  map[string]any
}

func (a *Atlas) Width() int { return a.width }

func (a *Atlas) Height() int { return a.height }

func (a *Atlas) GetSprite(name string) *Frame {
	idx, ok := a.spriteIndex[name]
	if !ok {
		return nil
	}
	return &a.sprites[idx]
}

func (a *Atlas) AddTile(tile image.Image) {
	if len(a.tiles) > 0 && !tile.Bounds().Size().Eq(a.tiles[0].Image.Bounds().Size()) {
		log.Fatalf("All tiles must have the same size")
	}
	a.tiles = append(a.tiles, Frame{
		Image: tile,
		W:     tile.Bounds().Size().X,
		H:     tile.Bounds().Size().Y,
	})
}

func (a *Atlas) AddSprite(name string, sprite image.Image, data map[string]any) {
	a.spriteIndex[name] = len(a.sprites)
	a.sprites = append(a.sprites, Frame{
		Name:  name,
		Image: sprite,
		W:     sprite.Bounds().Size().X,
		H:     sprite.Bounds().Size().Y,
		Data:  data,
	})
}

func (a *Atlas) Pack() {
	// Sort sprites by size
	sort.Slice(a.sprites, func(i, j int) bool {
		iw := a.sprites[i].W
		ih := a.sprites[i].H
		jw := a.sprites[j].W
		jh := a.sprites[j].H
		if iw*ih > jw*jh {
			return true
		}
		if iw*ih < jw*jh {
			return false
		}
		return iw > ih
	})

	// Rebuild index after sorting
	for i, s := range a.sprites {
		a.spriteIndex[s.Name] = i
	}

	// Repack until the size becomes stable
	lastWidth := a.width
	for {
		a.setSkyline(0)
		for i := range a.tiles {
			a.packSprite(&a.tiles[i])
		}

		a.setSkyline(a.skyline[0])
		for i := range a.sprites {
			a.packSprite(&a.sprites[i])
		}

		if lastWidth == a.width {
			return
		}
		lastWidth = a.width
	}
}

func (a *Atlas) SaveImage(filePath string) {
	img := image.NewRGBA(image.Rect(0, 0, a.width, a.height))

	for _, tile := range a.tiles {
		a.drawFrame(img, &tile)
	}

	for _, sprite := range a.sprites {
		a.drawFrame(img, &sprite)
	}

	png.Write(img, filePath)
}

func (a *Atlas) SaveJSON(filePath string, imagePath string) {
	data := rootJson{
		Meta: meta{
			Image:  imagePath,
			Format: "RGBA8888",
			Size: size{
				W: a.width,
				H: a.height,
			},
			Scale: "1",
		},
	}

	for _, s := range a.sprites {
		data.Frames = append(data.Frames, frame{
			Filename: s.Name,
			Frame: rect{
				X: s.X,
				Y: s.Y,
				W: s.W,
				H: s.H,
			},
			Data: s.Data,
		})
	}

	atlasJSON, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Fatalf("Failed to encode atlas JSON: %v", err)
	}
	err = os.WriteFile(filePath, atlasJSON, 0644)
	if err != nil {
		log.Fatalf("Failed to write atlas JSON: %v", err)
	}
}

func (a *Atlas) setSkyline(level int) {
	a.skyline = make([]int, a.width)
	for i := range a.skyline {
		a.skyline[i] = level
	}
}

func (a *Atlas) packSprite(f *Frame) {
	x := a.findBestPosition(f)
	f.X = x
	f.Y = a.skyline[x]
	a.insertFrame(x, f.W, f.H)
}

func (a *Atlas) findBestPosition(f *Frame) int {
	bestX := -1
	for x, y := range a.skyline {
		// If we don't fit into the frame, then skip
		if !a.isFitting(x, f.W, f.H) {
			continue
		}
		// If we already have a candidate best point - compare with it and discard if we are not better
		if bestX != -1 && a.skyline[bestX] <= y {
			continue
		}
		bestX = x
	}
	// If we found the best point - return it
	if bestX >= 0 {
		return bestX
	}
	// Otherwise double the atlas size and try again
	if a.width == a.height {
		a.skyline = append(a.skyline, make([]int, a.width)...)
		a.width *= 2
	} else {
		a.height *= 2
	}
	return a.findBestPosition(f)
}

func (a *Atlas) isFitting(x int, w int, h int) bool {
	y := a.skyline[x]
	// Return false if we don't fit to frame
	if (x+w > a.width) || (y+h > a.height) {
		return false
	}
	// Check whether we overlap with the next skyline points
	for nx := x + 1; nx < len(a.skyline); nx++ {
		// Point is outside the sprite - we can stop iterating, safely assuming that sprite fits
		if nx >= x+w {
			return true
		}
		// Point is inside the sprite - we can stop iterating, safely assuming that sprite doesn't fit
		if a.skyline[nx] > y {
			return false
		}
	}
	return true
}

func (a *Atlas) insertFrame(x int, w int, h int) {
	cy := a.skyline[x] + h
	for cx := x; cx < x+w; cx++ {
		a.skyline[cx] = cy
	}
}

func (a *Atlas) drawFrame(dstImage *image.RGBA, f *Frame) {
	draw.Draw(dstImage, image.Rect(f.X, f.Y, f.X+f.W, f.Y+f.H),
		f.Image, f.Image.Bounds().Min, draw.Src)
}
