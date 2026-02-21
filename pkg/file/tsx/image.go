package tsx

import (
	"image"
	"path/filepath"
	"pixel-tools/pkg/file/png"
)

type Image struct {
	Source string `xml:"source,attr"`
	Width  int    `xml:"width,attr"`
	Height int    `xml:"height,attr"`

	Data *image.RGBA `xml:"-"`
}

func (i *Image) PostLoad(basePath string) {
	i.Data = png.Get(filepath.Join(basePath, i.Source))
}
