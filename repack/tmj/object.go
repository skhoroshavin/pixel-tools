package tmj

import "repack/tsx"

type Object struct {
	Name       string     `json:"name"`
	Type       string     `json:"type"`
	X          float64    `json:"x"`
	Y          float64    `json:"y"`
	Width      float64    `json:"width"`
	Height     float64    `json:"height"`
	Rotation   float64    `json:"rotation"`
	Properties []Property `json:"properties,omitempty"`
}

func convertObject(obj tsx.Object) Object {
	return Object{
		Name:       obj.Name,
		Type:       obj.Type,
		X:          obj.X,
		Y:          obj.Y,
		Width:      obj.Width,
		Height:     obj.Height,
		Rotation:   obj.Rotation,
		Properties: convertProperties(obj.Properties),
	}
}

func convertObjects(objs []tsx.Object) []Object {
	var res []Object
	for _, obj := range objs {
		res = append(res, convertObject(obj))
	}
	return res
}
