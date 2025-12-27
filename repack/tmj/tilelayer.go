package tmj

import (
	"repack/tmx"
	"repack/tsx"
)

type TileLayer struct {
	ID         int                `json:"id"`
	Name       string             `json:"name"`
	Type       string             `json:"type"` // "tilelayer"
	Class      string             `json:"class"`
	X          int                `json:"x"`
	Y          int                `json:"y"`
	Width      int                `json:"width"`
	Height     int                `json:"height"`
	Opacity    float64            `json:"opacity"`
	Visible    bool               `json:"visible"`
	Data       []tsx.GlobalTileID `json:"data"`
	Properties []Property         `json:"properties,omitempty"`
}

func convertTileLayer(layer *tmx.TileLayer, layerID int) TileLayer {
	props := convertProperties(layer.Properties)
	// So far Phaser doesn't load layer.class, so this is a workaround
	if layer.Class != "" {
		props = append(props, Property{
			Name:  "class",
			Type:  "string",
			Value: layer.Class,
		})
	}

	return TileLayer{
		ID:         layerID,
		Name:       layer.Name,
		Type:       "tilelayer",
		Class:      layer.Class,
		X:          0,
		Y:          0,
		Width:      layer.Width,
		Height:     layer.Height,
		Opacity:    1,
		Visible:    true,
		Data:       layer.Data.Decoded,
		Properties: props,
	}
}
