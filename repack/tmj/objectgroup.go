package tmj

import "repack/tsx"

type ObjectGroup struct {
	Name    string   `json:"name"`
	Type    string   `json:"type"` // "objectgroup"
	Opacity float64  `json:"opacity"`
	Visible bool     `json:"visible"`
	X       int      `json:"x"`
	Y       int      `json:"y"`
	Objects []Object `json:"objects"`
}

func convertObjectGroup(og tsx.ObjectGroup) ObjectGroup {
	return ObjectGroup{
		Name:    og.Name,
		Type:    "objectgroup",
		Objects: convertObjects(og.Objects),
		Opacity: 1,
		Visible: true,
		X:       0,
		Y:       0,
	}
}

func convertObjectGroupOptional(og tsx.ObjectGroup) *ObjectGroup {
	if len(og.Objects) == 0 {
		return nil
	}

	res := convertObjectGroup(og)
	return &res
}
