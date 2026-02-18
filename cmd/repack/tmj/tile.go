package tmj

import (
	"pixel-tools/cmd/repack/tsx"
)

type Tile struct {
	ID          tsx.LocalTileID `json:"id"`
	Type        string          `json:"type,omitempty"`
	Animation   []Frame         `json:"animation,omitempty"`
	ObjectGroup *ObjectGroup    `json:"objectgroup,omitempty"`
	Properties  []Property      `json:"properties,omitempty"`
}

type Frame struct {
	TileID   tsx.LocalTileID `json:"tileid"`
	Duration int             `json:"duration"`
}

func convertTile(t *tsx.Tile) Tile {
	return Tile{
		ID:          t.ID,
		Type:        t.Type,
		ObjectGroup: convertOptionalObjectGroup(t.ObjectGroup),
		Animation:   convertAnimation(t.Animation),
		Properties:  convertProperties(t.Properties),
	}
}

func convertTiles(t []*tsx.Tile) []Tile {
	res := make([]Tile, len(t))
	for i, tile := range t {
		res[i] = convertTile(tile)
	}
	return res
}

func convertAnimation(anim []tsx.Frame) []Frame {
	res := make([]Frame, len(anim))
	for i, f := range anim {
		res[i] = Frame{
			TileID:   f.TileID,
			Duration: f.Duration,
		}
	}
	return res
}
