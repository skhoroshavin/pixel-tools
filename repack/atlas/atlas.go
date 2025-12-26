package atlas

import "repack/tsx"

func New(tilesets []*tsx.Tileset) *Atlas {
	return &Atlas{
		tilesets: tilesets,
	}
}

type Atlas struct {
	tilesets []*tsx.Tileset
}

func (a *Atlas) UseTileID(id tsx.GlobalTileID) tsx.GlobalTileID {
	return id
}
