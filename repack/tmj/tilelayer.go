package tmj

import "repack/tsx"

type TileLayer struct {
	Data    []tsx.GlobalTileID `json:"data"`
	Height  int                `json:"height"`
	ID      int                `json:"id"`
	Name    string             `json:"name"`
	Opacity float64            `json:"opacity"`
	Type    string             `json:"type"` // "tilelayer"
	Visible bool               `json:"visible"`
	Width   int                `json:"width"`
	X       int                `json:"x"`
	Y       int                `json:"y"`
}
