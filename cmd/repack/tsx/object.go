package tsx

type ObjectGroup struct {
	ID         int        `xml:"id,attr"`
	Name       string     `xml:"name,attr"`
	Class      string     `xml:"class,attr"`
	Objects    []Object   `xml:"object"`
	Properties []Property `xml:"properties>property"`
}

type Object struct {
	ID         int          `xml:"id,attr"`
	Name       string       `xml:"name,attr"`
	Type       string       `xml:"type,attr"`
	X          float64      `xml:"x,attr"`
	Y          float64      `xml:"y,attr"`
	Width      float64      `xml:"width,attr"`
	Height     float64      `xml:"height,attr"`
	Rotation   float64      `xml:"rotation,attr"`
	GID        GlobalTileID `xml:"gid,attr"`
	Polygon    Polygon      `xml:"polygon"`
	Properties []Property   `xml:"properties>property"`
}

type Polygon struct {
	Points string `xml:"points,attr"`
}
