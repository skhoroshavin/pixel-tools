package tmj

import (
	"log"
	"pixel-tools/pkg/file/tsx"
	"strconv"
	"strings"
)

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
		Properties: ConvertProperties(obj.Properties),
	}
}

func convertObjects(objs []tsx.Object) []Object {
	var res []Object
	for _, obj := range objs {
		res = append(res, convertObject(obj))
	}
	return res
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
