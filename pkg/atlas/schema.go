package atlas

type rootJson struct {
	Meta   meta    `json:"meta"`
	Frames []frame `json:"frames"`
}

type meta struct {
	Image  string `json:"image"`
	Format string `json:"format"`
	Size   Size   `json:"size"`
	Scale  string `json:"scale"`
}

type frame struct {
	Name          string `json:"filename"`
	Frame         Rect   `json:"frame"`
	Scale9Borders *Rect  `json:"scale9Borders,omitempty"`

	Trimmed          bool  `json:"trimmed,omitempty"`
	SpriteSourceSize *Rect `json:"spriteSourceSize,omitempty"`
	SourceSize       *Size `json:"sourceSize,omitempty"`

	Data map[string]any `json:"data,omitempty"`
}

type Size struct {
	W int `json:"w"`
	H int `json:"h"`
}

type Rect struct {
	X int `json:"x"`
	Y int `json:"y"`
	W int `json:"w"`
	H int `json:"h"`
}
