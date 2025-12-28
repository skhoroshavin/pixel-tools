package atlas

import (
	"log"
	"repack/tsx"
	"strconv"
	"strings"
)

func convertAnimation(anim []tsx.Frame) []AnimationFrame {
	res := make([]AnimationFrame, len(anim))
	for i, f := range anim {
		res[i] = AnimationFrame{
			TileID:   f.TileID,
			Duration: f.Duration,
		}
	}
	return res
}

func convertObjectGroup(og tsx.ObjectGroup) *ObjectGroup {
	if len(og.Objects) == 0 {
		return nil
	}
	props := convertProperties(og.Properties)

	// So far Phaser doesn't load layer.class, so this is a workaround
	if og.Class != "" {
		props = append(props, Property{
			Name:  "class",
			Type:  "string",
			Value: og.Class,
		})
	}

	return &ObjectGroup{
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

func convertObjects(objs []tsx.Object) []Object {
	var res []Object
	for _, obj := range objs {
		res = append(res, convertObject(obj))
	}
	return res
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
		GID:        obj.GID,
		Polygon:    convertPolygon(obj.Polygon),
		Properties: convertProperties(obj.Properties),
	}
}

func convertPolygon(polygon tsx.Polygon) []Point {
	if polygon.Points == "" {
		return nil
	}

	points := strings.Split(polygon.Points, " ")
	res := make([]Point, len(points))
	for i, point := range points {
		if point == "" {
			continue
		}
		coords := strings.Split(point, ",")
		if len(coords) != 2 {
			log.Fatalf("Invalid polygon point: %s", point)
		}
		x, err := strconv.ParseFloat(coords[0], 32)
		if err != nil {
			log.Fatalf("Failed to parse polygon point x: %v", err)
		}
		y, err := strconv.ParseFloat(coords[1], 32)
		if err != nil {
			log.Fatalf("Failed to parse polygon point y: %v", err)
		}
		res[i] = Point{X: x, Y: y}
	}
	return res
}

func convertProperties(props []tsx.Property) []Property {
	var res []Property
	for _, p := range props {
		res = append(res, convertProperty(p))
	}
	return res
}

func convertProperty(prop tsx.Property) Property {
	return Property{
		Name:  prop.Name,
		Type:  prop.Type,
		Value: prop.Value,
	}
}
