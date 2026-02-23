package atlas

import (
	"encoding/json"
	"image"
	"image/draw"
	"log"
	"os"
	"path"
	"pixel-tools/pkg/file/png"
	"pixel-tools/pkg/imgutil"
	"sort"
)

func New() *Atlas {
	const initialSize = 16

	return &Atlas{
		width:       initialSize,
		height:      initialSize,
		spriteIndex: make(map[string]int),
		skyline:     make([]int, initialSize),
	}
}

type Atlas struct {
	tiles     []Frame
	sprites   []Frame
	frameRefs []FrameRef

	spriteIndex map[string]int

	width   int
	height  int
	skyline []int
}

type Frame struct {
	frame
	Image image.Image
}

type NineSlice = rect

type FrameRef struct {
	Name        string
	SourceFrame string
	Data        map[string]any
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
		frame: frame{
			Frame: rect{
				W: tile.Bounds().Size().X,
				H: tile.Bounds().Size().Y,
			},
		},
		Image: tile,
	})
}

func (a *Atlas) AddSprite(name string, sprite image.Image, nineSlice *NineSlice, data map[string]any) *Frame {
	a.spriteIndex[name] = len(a.sprites)
	if nineSlice != nil {
		return a.addUntrimmed(name, sprite, nineSlice, data)
	}

	left := imgutil.GetLeftMargin(sprite)
	right := imgutil.GetRightMargin(sprite)
	top := imgutil.GetTopMargin(sprite)
	bottom := imgutil.GetBottomMargin(sprite)
	if left == 0 && right == 0 && top == 0 && bottom == 0 {
		return a.addUntrimmed(name, sprite, nineSlice, data)
	}

	originalW := sprite.Bounds().Dx()
	originalH := sprite.Bounds().Dy()
	trimmedW := originalW - left - right
	trimmedH := originalH - top - bottom

	res := Frame{
		frame: frame{
			Name: name,
			Frame: rect{
				W: trimmedW,
				H: trimmedH,
			},
			Trimmed: true,
			SpriteSourceSize: &rect{
				X: left,
				Y: top,
				W: trimmedW,
				H: trimmedH,
			},
			SourceSize: &size{
				W: originalW,
				H: originalH,
			},
			Data: data,
		},
		Image: sprite,
	}
	a.sprites = append(a.sprites, res)
	return &res
}

func (a *Atlas) AddSpriteRef(name string, sourceName string, data map[string]any) {
	a.frameRefs = append(a.frameRefs, FrameRef{
		Name:        name,
		SourceFrame: sourceName,
		Data:        data,
	})
}

func (a *Atlas) Pack() {
	// Sort sprites by size
	sort.Slice(a.sprites, func(i, j int) bool {
		iw := a.sprites[i].Frame.W
		ih := a.sprites[i].Frame.H
		jw := a.sprites[j].Frame.W
		jh := a.sprites[j].Frame.H
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

func (a *Atlas) Save(baseName string) {
	a.SaveImage(baseName + ".png")
	a.SaveJSON(baseName+".atlas", path.Base(baseName)+".png")
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
		data.Frames = append(data.Frames, s.frame)
	}

	for _, ref := range a.frameRefs {
		sourceFrame := a.GetSprite(ref.SourceFrame)
		if sourceFrame == nil {
			log.Fatalf("Frame reference %q refers to non-existent sprite %q", ref.Name, ref.SourceFrame)
		}
		refFrame := sourceFrame.frame
		refFrame.Name = ref.Name
		refFrame.Data = ref.Data
		data.Frames = append(data.Frames, refFrame)
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

func (a *Atlas) addUntrimmed(name string, sprite image.Image, nineSlice *NineSlice, data map[string]any) *Frame {
	res := Frame{
		frame: frame{
			Name: name,
			Frame: rect{
				W: sprite.Bounds().Dx(),
				H: sprite.Bounds().Dy(),
			},
			Scale9Borders: nineSlice,
			Data:          data,
		},
		Image: sprite,
	}
	a.sprites = append(a.sprites, res)
	return &res
}

func (a *Atlas) setSkyline(level int) {
	a.skyline = make([]int, a.width)
	for i := range a.skyline {
		a.skyline[i] = level
	}
}

func (a *Atlas) packSprite(f *Frame) {
	x := a.findBestPosition(f)
	f.Frame.X = x
	f.Frame.Y = a.skyline[x]
	a.insertFrame(x, f.Frame.W, f.Frame.H)
}

func (a *Atlas) findBestPosition(f *Frame) int {
	bestX := -1
	for x, y := range a.skyline {
		// If we don't fit into the frame, then skip
		if !a.isFitting(x, f.Frame.W, f.Frame.H) {
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
	if f.Frame.W == 0 || f.Frame.H == 0 {
		return
	}

	sourcePoint := f.Image.Bounds().Min
	if f.SpriteSourceSize != nil {
		sourcePoint.X += f.SpriteSourceSize.X
		sourcePoint.Y += f.SpriteSourceSize.Y
	}
	draw.Draw(dstImage, image.Rect(f.Frame.X, f.Frame.Y, f.Frame.X+f.Frame.W, f.Frame.Y+f.Frame.H),
		f.Image, sourcePoint, draw.Src)
}
