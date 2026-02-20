package tmj

import "pixel-tools/cmd/tilepack/tsx"

type Property struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Value string `json:"value"`
}

func ConvertProperty(prop tsx.Property) Property {
	return Property{
		Name:  prop.Name,
		Type:  prop.Type,
		Value: prop.Value,
	}
}

func ConvertProperties(props []tsx.Property) []Property {
	var res []Property
	for _, p := range props {
		res = append(res, ConvertProperty(p))
	}
	return res
}
