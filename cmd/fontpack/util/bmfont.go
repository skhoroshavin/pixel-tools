package util

import (
	"encoding/xml"
	"log"
	"os"
	"strings"
)

func SaveBMFont(data BMFont, path string) {
	xmlData, err := xml.MarshalIndent(&data, "", "  ")
	if err != nil {
		log.Fatalf("Failed to encode BMFont: %+v", err)
	}
	rawData := `<?xml version="1.0" encoding="UTF-8"?>` + "\n" + string(xmlData)
	rawData = strings.ReplaceAll(rawData, "></info>", "/>")
	rawData = strings.ReplaceAll(rawData, "></common>", "/>")
	rawData = strings.ReplaceAll(rawData, "></page>", "/>")
	rawData = strings.ReplaceAll(rawData, "></char>", "/>")

	err = os.WriteFile(path, []byte(rawData), 0644)
	if err != nil {
		log.Fatalf("Failed to write BMFont: %+v", err)
	}
}

type BMFont struct {
	XMLName xml.Name `xml:"font"`

	Info   FontInfo `xml:"info"`
	Common Common   `xml:"common"`
	Pages  Pages    `xml:"pages"`
	Chars  Chars    `xml:"chars"`
}

type FontInfo struct {
	Face string `xml:"face,attr"`
	Size int    `xml:"size,attr"`
}

type Common struct {
	LineHeight int `xml:"lineHeight,attr"`
	Base       int `xml:"base,attr"`
	Pages      int `xml:"pages,attr"`
}

type Pages struct {
	Page []Page `xml:"page"`
}

type Page struct {
	ID   int    `xml:"id,attr"`
	File string `xml:"file,attr"`
}

type Chars struct {
	Count int    `xml:"count,attr"`
	Char  []Char `xml:"char"`
}

type Char struct {
	ID       int    `xml:"id,attr"`
	Letter   string `xml:"letter,attr"`
	X        int    `xml:"x,attr"`
	Y        int    `xml:"y,attr"`
	Width    int    `xml:"width,attr"`
	Height   int    `xml:"height,attr"`
	XOffset  int    `xml:"xoffset,attr"`
	YOffset  int    `xml:"yoffset,attr"`
	XAdvance int    `xml:"xadvance,attr"`
}
