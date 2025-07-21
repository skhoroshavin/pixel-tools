package tmj

import "repack/tsx"

type Property struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Value string `json:"value"`
}

func convertProperty(prop tsx.Property) Property {
	return Property{
		Name:  prop.Name,
		Type:  prop.Type,
		Value: prop.Value,
	}
}

func convertProperties(props []tsx.Property) []Property {
	var res []Property
	for _, p := range props {
		res = append(res, convertProperty(p))
	}
	return res
}
