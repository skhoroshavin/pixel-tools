package tsx

import (
	"image"
	"path/filepath"

	"pixel-tools/pkg/imgutil"
)

type Image struct {
	Source string `xml:"source,attr"`
	Width  int    `xml:"width,attr"`
	Height int    `xml:"height,attr"`

	Data image.Image `xml:"-"`
}

func (i *Image) PostLoad(basePath string) {
	i.Data = imgutil.GetImage(filepath.Join(basePath, i.Source))
}
