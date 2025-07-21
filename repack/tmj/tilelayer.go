package tmj

import (
	"repack/tmx"
	"repack/tsx"
)

type TileLayer struct {
	ID      int                `json:"id"`
	Name    string             `json:"name"`
	Type    string             `json:"type"` // "tilelayer"
	X       int                `json:"x"`
	Y       int                `json:"y"`
	Width   int                `json:"width"`
	Height  int                `json:"height"`
	Opacity float64            `json:"opacity"`
	Visible bool               `json:"visible"`
	Data    []tsx.GlobalTileID `json:"data"`
}

func convertTileLayer(layer *tmx.Layer, layerID int) TileLayer {
	return TileLayer{
		ID:      layerID,
		Name:    layer.Name,
		Type:    "tilelayer",
		X:       0,
		Y:       0,
		Width:   layer.Width,
		Height:  layer.Height,
		Opacity: 1,
		Visible: true,
		Data:    layer.Data.Decoded,
	}
}
