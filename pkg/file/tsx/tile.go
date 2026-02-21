package tsx

import (
	"image"
)

type Tile struct {
	ID          LocalTileID `xml:"id,attr"`
	Type        string      `xml:"type,attr"`
	Animation   []Frame     `xml:"animation>frame"`
	Properties  []Property  `xml:"properties>property"`
	ObjectGroup ObjectGroup `xml:"objectgroup"`

	Image  *Image `xml:"image"`
	X      int    `xml:"x,attr"`
	Y      int    `xml:"y,attr"`
	Width  int    `xml:"width,attr"`
	Height int    `xml:"height,attr"`

	Tileset *Tileset `xml:"-"`
	data    image.Image
}

func (t *Tile) PostLoad(basePath string, tileset *Tileset) {
	t.Tileset = tileset
	if t.Image != nil {
		t.Image.PostLoad(basePath)
	}
}

func (t *Tile) Data() image.Image {
	if t.data == nil {
		if t.Image == nil {
			w := t.Tileset.TileWidth
			h := t.Tileset.TileHeight
			srcX := (int(t.ID) % t.Tileset.Columns) * w
			srcY := (int(t.ID) / t.Tileset.Columns) * h
			t.data = t.Tileset.Image.Data.SubImage(image.Rect(srcX, srcY, srcX+w, srcY+h))
		} else {
			t.data = t.Image.Data
		}
	}
	return t.data
}

type Frame struct {
	TileID   LocalTileID `xml:"tileid,attr"`
	Duration int         `xml:"duration,attr"` // milliseconds
}

type GlobalTileID uint32
type LocalTileID uint32

func (t GlobalTileID) WithoutFlags() GlobalTileID {
	return t & 0x0fffffff
}

func (t GlobalTileID) Flags() uint32 {
	return uint32(t) >> 28
}

func (t GlobalTileID) WithFlags(flags uint32) GlobalTileID {
	return GlobalTileID(uint32(t) | (flags << 28))
}
