package tmj

import (
	"encoding/json"
	"log"
	"os"
	"repack/tmx"
)

type Map struct {
	Version          string    `json:"version"`
	CompressionLevel int       `json:"compressionlevel"`
	Height           int       `json:"height"`
	Infinite         bool      `json:"infinite"`
	Layers           []any     `json:"layers"` // tile + object layers
	TileWidth        int       `json:"tilewidth"`
	TileHeight       int       `json:"tileheight"`
	Width            int       `json:"width"`
	Orientation      string    `json:"orientation"`
	Tilesets         []Tileset `json:"tilesets"`
}

func (m *Map) Save(path string) {
	raw, err := json.Marshal(m)
	if err != nil {
		log.Fatalf("Failed to encode map: %v", err)
	}
	err = os.WriteFile(path, raw, 0644)
	if err != nil {
		log.Fatalf("Failed to write map: %v", err)
	}
}

func ConvertFromTMX(src *tmx.Map, name string) Map {
	res := Map{
		Version:          src.Version,
		CompressionLevel: -1,
		Height:           src.Height,
		Width:            src.Width,
		TileWidth:        src.TileWidth,
		TileHeight:       src.TileHeight,
		Orientation:      src.Orientation,
		Infinite:         false,
	}

	layerID := 1

	// Convert tile layers
	for _, layer := range src.Layers {
		jsonLayer := TileLayer{
			Data:    layer.Data.Decoded,
			Height:  layer.Height,
			ID:      layerID,
			Name:    layer.Name,
			Opacity: 1,
			Type:    "tilelayer",
			Visible: true,
			Width:   layer.Width,
			X:       0,
			Y:       0,
		}
		res.Layers = append(res.Layers, jsonLayer)
		layerID++
	}

	// Convert object layers
	for _, og := range src.ObjectGroups {
		var objects []Object
		for _, obj := range og.Objects {
			props := make(map[string]string)
			for _, p := range obj.Properties {
				props[p.Name] = p.Value
			}
			objects = append(objects, Object{
				ID:         obj.ID,
				Name:       obj.Name,
				Type:       obj.Type,
				X:          obj.X,
				Y:          obj.Y,
				Width:      obj.Width,
				Height:     obj.Height,
				Rotation:   obj.Rotation,
				Properties: props,
			})
		}

		res.Layers = append(res.Layers, ObjectLayer{
			ID:      layerID,
			Name:    og.Name,
			Objects: objects,
			Opacity: 1,
			Type:    "objectgroup",
			Visible: true,
			X:       0,
			Y:       0,
		})
		layerID++
	}

	// Convert tilesets
	for _, ts := range src.Tilesets {
		tileset := Tileset{
			FirstGID:    ts.FirstGID,
			Name:        ts.Name,
			TileWidth:   ts.TileWidth,
			TileHeight:  ts.TileHeight,
			TileCount:   ts.TileCount,
			Columns:     ts.Columns,
			Image:       ts.Image.Source,
			ImageWidth:  ts.Image.Width,
			ImageHeight: ts.Image.Height,
		}

		for _, tile := range ts.Tiles {
			var anim []Frame
			for _, f := range tile.Animation {
				anim = append(anim, Frame{
					TileID:   f.TileID,
					Duration: f.Duration,
				})
			}

			var props []Property
			for _, p := range tile.Properties {
				props = append(props, Property{
					Name:  p.Name,
					Type:  "string", //p.Type
					Value: p.Value,
				})
			}
			tileset.Tiles = append(tileset.Tiles, Tile{
				ID:         tile.ID,
				Animation:  anim,
				Properties: props,
			})
		}

		res.Tilesets = append(res.Tilesets, tileset)
	}
	return res
}
