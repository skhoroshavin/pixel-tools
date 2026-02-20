package atlas

import (
	"encoding/json"
	"image"
	"image/draw"
	"log"
	"os"
	"path"
	"pixel-tools/pkg/fileutil"
	"sort"
)

func New(name string) *Atlas {
	const initialSize = 16

	return &Atlas{
		name:    name,
		size:    initialSize,
		skyline: make([]int, initialSize),
	}
}

type Atlas struct {
	name string

	tiles   []frame
	sprites []frame

	size    int
	skyline []int
}

type frame struct {
	name string
	img  image.Image
	x, y int
	w, h int
	data map[string]any
}

func (a *Atlas) Size() int { return a.size }

func (a *Atlas) AddTile(tile image.Image) {
	if len(a.tiles) > 0 && !tile.Bounds().Size().Eq(a.tiles[0].img.Bounds().Size()) {
		log.Fatalf("All tiles must have the same size")
	}
	a.tiles = append(a.tiles, frame{
		img: tile,
		w:   tile.Bounds().Size().X,
		h:   tile.Bounds().Size().Y,
	})
}

func (a *Atlas) AddSprite(name string, sprite image.Image, data map[string]any) {
	a.sprites = append(a.sprites, frame{
		name: name,
		img:  sprite,
		w:    sprite.Bounds().Size().X,
		h:    sprite.Bounds().Size().Y,
		data: data,
	})
}

func (a *Atlas) Pack() {
	// Sort sprites by size
	sort.Slice(a.sprites, func(i, j int) bool {
		iw := a.sprites[i].w
		ih := a.sprites[i].h
		jw := a.sprites[j].w
		jh := a.sprites[j].h
		if iw*ih > jw*jh {
			return true
		}
		if iw*ih < jw*jh {
			return false
		}
		return iw > ih
	})

	// Repack until the size becomes stable
	lastSize := a.size
	for {
		a.setSkyline(0)
		for i := range a.tiles {
			a.packSprite(&a.tiles[i])
		}

		a.setSkyline(a.skyline[0])
		for i := range a.sprites {
			a.packSprite(&a.sprites[i])
		}

		if lastSize != a.size {
			lastSize = a.size
		} else {
			return
		}
	}
}

func (a *Atlas) Save(baseName string) {
	img := image.NewRGBA(image.Rect(0, 0, a.size, a.size))
	data := RootJson{
		Meta: Meta{
			Image:  path.Base(baseName) + ".png",
			Format: "RGBA8888",
			Size: Size{
				W: a.size,
				H: a.size,
			},
			Scale: "1",
		},
	}

	for _, tile := range a.tiles {
		a.drawFrame(img, &tile)
	}

	for _, sprite := range a.sprites {
		a.drawFrame(img, &sprite)

		data.Frames = append(data.Frames, Frame{
			Filename: sprite.name,
			Frame:    Rect{X: sprite.x, Y: sprite.y, W: sprite.w, H: sprite.h},
			Data:     sprite.data,
		})
	}

	fileutil.WriteImage(img, baseName+".png")

	atlasJSON, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Fatalf("Failed to encode atlas JSON: %v", err)
	}
	err = os.WriteFile(baseName+".atlas", atlasJSON, 0644)
	if err != nil {
		log.Fatalf("Failed to write atlas JSON: %v", err)
	}
}

func (a *Atlas) setSkyline(level int) {
	a.skyline = make([]int, a.size)
	for i := range a.skyline {
		a.skyline[i] = level
	}
}

func (a *Atlas) packSprite(f *frame) {
	x := a.findBestPosition(f)
	f.x = x
	f.y = a.skyline[x]
	a.insertFrame(x, f.w, f.h)
}

func (a *Atlas) findBestPosition(f *frame) int {
	bestX := -1
	for x, y := range a.skyline {
		// If we don't fit into the frame, then skip
		if !a.isFitting(x, f.w, f.h) {
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
	a.skyline = append(a.skyline, make([]int, a.size)...)
	a.size *= 2
	return a.findBestPosition(f)
}

func (a *Atlas) isFitting(x int, w int, h int) bool {
	y := a.skyline[x]
	// Return false if we don't fit to frame
	if (x+w > a.size) || (y+h > a.size) {
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

func (a *Atlas) drawFrame(dstImage *image.RGBA, f *frame) {
	draw.Draw(dstImage, image.Rect(f.x, f.y, f.x+f.w, f.y+f.h),
		f.img, f.img.Bounds().Min, draw.Src)
}
