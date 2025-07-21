package tsx

type ObjectGroup struct {
	ID         int        `xml:"id,attr"`
	Name       string     `xml:"name,attr"`
	Objects    []Object   `xml:"object"`
	Properties []Property `xml:"properties>property"`
}

type Object struct {
	ID         int        `xml:"id,attr"`
	Name       string     `xml:"name,attr"`
	Type       string     `xml:"type,attr"`
	X          float64    `xml:"x,attr"`
	Y          float64    `xml:"y,attr"`
	Width      float64    `xml:"width,attr"`
	Height     float64    `xml:"height,attr"`
	Rotation   float64    `xml:"rotation,attr"`
	Properties []Property `xml:"properties>property"`
}
