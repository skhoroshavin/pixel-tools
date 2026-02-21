package tmj

import "pixel-tools/pkg/file/tsx"

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

func convertObjectGroup(og tsx.ObjectGroup) ObjectGroup {
	props := ConvertProperties(og.Properties)

	// So far Phaser doesn't load layer.class, so this is a workaround
	if og.Class != "" {
		props = append(props, Property{
			Name:  "class",
			Type:  "string",
			Value: og.Class,
		})
	}

	return ObjectGroup{
		Name:       og.Name,
		Type:       "objectgroup",
		Objects:    convertObjects(og.Objects),
		Properties: props,
		Opacity:    1,
		Visible:    true,
		X:          0,
		Y:          0,
	}
}

func ConvertOptionalObjectGroup(og tsx.ObjectGroup) *ObjectGroup {
	if len(og.Objects) == 0 {
		return nil
	}

	res := convertObjectGroup(og)
	return &res
}
