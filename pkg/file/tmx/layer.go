package tmx

import (
	"encoding/xml"
	"log"
	"pixel-tools/pkg/file/tsx"
	"strconv"
	"strings"
)

type Layer struct {
	XMLName xml.Name

	// Common attributes
	Name       string         `xml:"name,attr"`
	Properties []tsx.Property `xml:"properties>property"`

	// Tile layer attributes
	Class  string    `xml:"class,attr"`
	Width  int       `xml:"width,attr"`
	Height int       `xml:"height,attr"`
	Data   LayerData `xml:"data"`

	// Object group attributes
	ID      int          `xml:"id,attr"`
	Objects []tsx.Object `xml:"object"`
}

func (l *Layer) IsTileLayer() bool {
	return l.XMLName.Local == "layer"
}

func (l *Layer) IsObjectGroup() bool {
	return l.XMLName.Local == "objectgroup"
}

func (l *Layer) AsTileLayer() *TileLayer {
	return &TileLayer{
		Name:       l.Name,
		Class:      l.Class,
		Width:      l.Width,
		Height:     l.Height,
		Data:       l.Data,
		Properties: l.Properties,
	}
}

func (l *Layer) AsObjectGroup() tsx.ObjectGroup {
	return tsx.ObjectGroup{
		ID:         l.ID,
		Name:       l.Name,
		Class:      l.Class,
		Objects:    l.Objects,
		Properties: l.Properties,
	}
}

type TileLayer struct {
	Name       string         `xml:"name,attr"`
	Class      string         `xml:"class,attr"`
	Width      int            `xml:"width,attr"`
	Height     int            `xml:"height,attr"`
	Data       LayerData      `xml:"data"`
	Properties []tsx.Property `xml:"properties>property"`
}

type LayerData struct {
	Encoding string             `xml:"encoding,attr"`
	Value    string             `xml:",chardata"`
	Decoded  []tsx.GlobalTileID `xml:"-"`
}

func (l *LayerData) Decode() {
	if l.Encoding != "csv" {
		log.Fatalf("Unsupported encoding: %s", l.Encoding)
	}

	lines := strings.Split(strings.TrimSpace(l.Value), ",")
	l.Decoded = make([]tsx.GlobalTileID, len(lines))
	for i, entry := range lines {
		entry = strings.TrimSpace(entry)
		if entry == "" {
			continue
		}
		num, err := strconv.ParseUint(entry, 10, 32)
		if err != nil {
			log.Fatalf("Failed to parse CSV entry: %v", err)
		}
		l.Decoded[i] = tsx.GlobalTileID(num)
	}
}
