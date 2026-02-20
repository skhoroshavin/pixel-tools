package atlas

type RootJson struct {
	Meta   Meta    `json:"meta"`
	Frames []Frame `json:"frames"`
}

type Meta struct {
	Image  string `json:"image"`
	Format string `json:"format"`
	Size   Size   `json:"size"`
	Scale  string `json:"scale"`
}

type Size struct {
	W int `json:"w"`
	H int `json:"h"`
}

type Frame struct {
	Filename string         `json:"filename"`
	Frame    Rect           `json:"frame"`
	Data     map[string]any `json:"data"`
}

type Rect struct {
	X int `json:"x"`
	Y int `json:"y"`
	W int `json:"w"`
	H int `json:"h"`
}
