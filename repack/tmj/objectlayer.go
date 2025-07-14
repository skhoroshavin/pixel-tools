package tmj

type ObjectLayer struct {
	ID      int      `json:"id"`
	Name    string   `json:"name"`
	Objects []Object `json:"objects"`
	Opacity float64  `json:"opacity"`
	Type    string   `json:"type"` // "objectgroup"
	Visible bool     `json:"visible"`
	X       int      `json:"x"`
	Y       int      `json:"y"`
}

type Object struct {
	ID         int               `json:"id"`
	Name       string            `json:"name"`
	Type       string            `json:"type"`
	X          float64           `json:"x"`
	Y          float64           `json:"y"`
	Width      float64           `json:"width"`
	Height     float64           `json:"height"`
	Rotation   float64           `json:"rotation"`
	Properties map[string]string `json:"properties,omitempty"`
}
