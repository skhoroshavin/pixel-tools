package tmx

import (
	"log"
	"repack/tsx"
	"strconv"
	"strings"
)

type Layer struct {
	Name       string         `xml:"name,attr"`
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

func (l *Layer) Decode() {
	if l.Data.Encoding != "csv" {
		log.Fatalf("Unsupported encoding: %s", l.Data.Encoding)
	}

	lines := strings.Split(strings.TrimSpace(l.Data.Value), ",")
	l.Data.Decoded = make([]tsx.GlobalTileID, len(lines))
	for i, entry := range lines {
		entry = strings.TrimSpace(entry)
		if entry == "" {
			continue
		}
		num, err := strconv.ParseInt(entry, 10, 32)
		if err != nil {
			log.Fatalf("Failed to parse CSV entry: %v", err)
		}
		l.Data.Decoded[i] = tsx.GlobalTileID(num)
	}
}
