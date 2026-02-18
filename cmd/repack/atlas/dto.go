package atlas

import (
	"pixel-tools/cmd/repack/tsx"
)

type RootJson struct {
	Meta   Meta    `json:"meta"`
	Frames []Frame `json:"frames"`
}

type Meta struct {
	Image  string `json:"image"`
	Format string `json:"format"`
	Size   Size   `json:"size"`
	Scale  string `json:"scale"`
}

type Size struct {
	W int `json:"w"`
	H int `json:"h"`
}

type Frame struct {
	Filename string `json:"filename"`
	Frame    Rect   `json:"frame"`

	Animation   []AnimationFrame `json:"animation,omitempty"`
	ObjectGroup *ObjectGroup     `json:"objectgroup,omitempty"`
	Properties  []Property       `json:"properties,omitempty"`
}

type Rect struct {
	X int `json:"x"`
	Y int `json:"y"`
	W int `json:"w"`
	H int `json:"h"`
}

type AnimationFrame struct {
	TileID   tsx.LocalTileID `json:"tileid"`
	Duration int             `json:"duration"`
}

type ObjectGroup struct {
	Name       string     `json:"name"`
	Type       string     `json:"type"` // "objectgroup"
	Opacity    float64    `json:"opacity"`
	Visible    bool       `json:"visible"`
	X          int        `json:"x"`
	Y          int        `json:"y"`
	Objects    []Object   `json:"objects"`
	Properties []Property `json:"properties"`
}

type Object struct {
	Name       string           `json:"name"`
	Type       string           `json:"type"`
	X          float64          `json:"x"`
	Y          float64          `json:"y"`
	Width      float64          `json:"width"`
	Height     float64          `json:"height"`
	Rotation   float64          `json:"rotation"`
	GID        tsx.GlobalTileID `json:"gid,omitempty"`
	Polygon    []Point          `json:"polygon,omitempty"`
	Properties []Property       `json:"properties,omitempty"`
}

type Point struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type Property struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Value string `json:"value"`
}
