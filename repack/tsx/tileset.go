package tsx

import (
	"image"
	"image/draw"
	"log"
)

type Tileset struct {
	FirstGID   GlobalTileID `xml:"firstgid,attr"`
	Source     string       `xml:"source,attr"` // only used if external
	Name       string       `xml:"name,attr"`
	TileWidth  int          `xml:"tilewidth,attr"`
	TileHeight int          `xml:"tileheight,attr"`
	TileCount  int          `xml:"tilecount,attr"`
	Columns    int          `xml:"columns,attr"`
	Image      Image        `xml:"image"`
	Tiles      []*Tile      `xml:"tile"`
}

type Tile struct {
	ID         LocalTileID `xml:"id,attr"`
	Animation  []Frame     `xml:"animation>frame"`
	Properties []Property  `xml:"properties>property"`

	Tileset *Tileset `xml:"-"`
}

type Frame struct {
	TileID   LocalTileID `xml:"tileid,attr"`
	Duration int         `xml:"duration,attr"` // milliseconds
}

type GlobalTileID int
type LocalTileID int

func (ts *Tileset) PostLoad(basePath string) {
	ts.Image.PostLoad(basePath)
	for _, tile := range ts.Tiles {
		tile.Tileset = ts
	}
}

func (ts *Tileset) HasTile(id GlobalTileID) bool {
	return id >= ts.FirstGID && id < ts.FirstGID+GlobalTileID(ts.TileCount)
}

func (ts *Tileset) LocalTileID(id GlobalTileID) LocalTileID {
	if !ts.HasTile(id) {
		log.Fatalf("Cannot find tile %d in tileset %s", id, ts.Name)
	}
	return LocalTileID(id - ts.FirstGID)
}

func (ts *Tileset) GlobalTileID(id LocalTileID) GlobalTileID {
	return GlobalTileID(id) + ts.FirstGID
}

func (ts *Tileset) Tile(id GlobalTileID) *Tile {
	if !ts.HasTile(id) {
		return nil
	}

	localID := ts.LocalTileID(id)
	for _, tile := range ts.Tiles {
		if tile.ID == localID {
			return tile
		}
	}

	tile := &Tile{
		ID:      localID,
		Tileset: ts,
	}
	ts.Tiles = append(ts.Tiles, tile)
	return tile
}

func (t *Tile) CopyImageFrom(srcTile *Tile) {
	tileWidth := t.Tileset.TileWidth
	tileHeight := t.Tileset.TileHeight
	srcColumns := srcTile.Tileset.Columns
	dstColumns := t.Tileset.Columns

	srcX := (int(srcTile.ID) % srcColumns) * tileWidth
	srcY := (int(srcTile.ID) / srcColumns) * tileHeight
	dstX := (int(t.ID) % dstColumns) * tileWidth
	dstY := (int(t.ID) / dstColumns) * tileHeight

	dstImage, ok := t.Tileset.Image.Data.(draw.Image)
	if !ok {
		log.Fatalf("Cannot draw to tileset image")
	}

	draw.Draw(dstImage, image.Rect(dstX, dstY, dstX+tileWidth, dstY+tileHeight),
		srcTile.Tileset.Image.Data, image.Pt(srcX, srcY), draw.Src)
}
