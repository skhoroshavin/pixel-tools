package atlas

type rootJson struct {
	Meta   meta    `json:"meta"`
	Frames []frame `json:"frames"`
}

type meta struct {
	Image  string `json:"image"`
	Format string `json:"format"`
	Size   size   `json:"size"`
	Scale  string `json:"scale"`
}

type frame struct {
	Name          string `json:"filename"`
	Frame         rect   `json:"frame"`
	Scale9Borders *rect  `json:"scale9Borders,omitempty"`

	Trimmed          bool  `json:"trimmed,omitempty"`
	SpriteSourceSize *rect `json:"spriteSourceSize,omitempty"`
	SourceSize       *size `json:"sourceSize,omitempty"`

	Data map[string]any `json:"data,omitempty"`
}

type size struct {
	W int `json:"w"`
	H int `json:"h"`
}

type rect struct {
	X int `json:"x"`
	Y int `json:"y"`
	W int `json:"w"`
	H int `json:"h"`
}
