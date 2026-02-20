package tsx

import (
	"image"
	"path/filepath"
	"pixel-tools/pkg/fileutil"
)

type Image struct {
	Source string `xml:"source,attr"`
	Width  int    `xml:"width,attr"`
	Height int    `xml:"height,attr"`

	Data *image.RGBA `xml:"-"`
}

func (i *Image) PostLoad(basePath string) {
	i.Data = fileutil.GetImage(filepath.Join(basePath, i.Source))
}
