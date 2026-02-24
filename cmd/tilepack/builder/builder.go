package builder

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"pixel-tools/pkg/atlas"
	"pixel-tools/pkg/file/tmj"
	"pixel-tools/pkg/file/tmx"
	"pixel-tools/pkg/file/tsx"
)

func New(padding int) *Builder {
	return &Builder{
		tileset:  newTileRegistry(padding),
		tilemaps: make(map[string]*mapContext),
	}
}

type Builder struct {
	tileset  *tileRegistry
	tilemaps map[string]*mapContext
}

type tileKey struct {
	tileset string
	id      tsx.LocalTileID
}

func (b *Builder) AddTilemap(name string, tilemap *tmx.Map) {
	b.tilemaps[name] = newMapContext(b.tileset, tilemap)
}

func (b *Builder) Save(outputDir string) {
	err := os.MkdirAll(outputDir, 0775)
	if err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	b.tileset.atlas.Pack()
	b.tileset.atlas.Save(filepath.Join(outputDir, "tileset"))

	for name, tilemap := range b.tilemaps {
		tilemap.save(outputDir, name)
	}
}

func convertTileData(tile *tsx.Tile) map[string]any {
	res := make(map[string]any)

	if animation := tmj.ConvertAnimation(tile.Animation); len(animation) > 0 {
		res["animation"] = animation
	}

	if objectGroup := tmj.ConvertOptionalObjectGroup(tile.ObjectGroup); objectGroup != nil {
		res["objectgroup"] = objectGroup
	}

	if properties := tmj.ConvertProperties(tile.Properties); len(properties) > 0 {
		res["properties"] = properties
	}

	if len(res) == 0 {
		return nil
	}

	return res
}

type tileRegistry struct {
	atlas *atlas.Atlas
	tiles []*tsx.Tile

	repackedTiles   map[tileKey]tsx.GlobalTileID
	repackedSprites map[tileKey]tsx.GlobalTileID
}

func newTileRegistry(padding int) *tileRegistry {
	return &tileRegistry{
		atlas: atlas.New(atlas.Config{Padding: padding}),

		repackedTiles:   make(map[tileKey]tsx.GlobalTileID),
		repackedSprites: make(map[tileKey]tsx.GlobalTileID),
	}
}

func (r *tileRegistry) addTile(tile *tsx.Tile) tsx.GlobalTileID {
	if len(r.tiles) > 0 {
		firstTile := r.tiles[0]
		if tile.Tileset.TileWidth != firstTile.Tileset.TileWidth || tile.Tileset.TileHeight != firstTile.Tileset.TileHeight {
			log.Fatalf("All tiles must have the same size, please use object layer instead")
		}
	}

	repackedTileID, created := r.add(r.repackedTiles, tile, 1)
	if !created {
		return repackedTileID
	}

	r.atlas.AddTile(tile.Data())
	r.tiles = append(r.tiles, tile)
	return repackedTileID
}

func (r *tileRegistry) addSprite(tile *tsx.Tile) tsx.GlobalTileID {
	repackedTileID, created := r.add(r.repackedSprites, tile, 1000000)
	if !created {
		return repackedTileID
	}

	data := convertTileData(tile)
	r.atlas.AddSprite(atlas.SpriteConfig{
		Name:  fmt.Sprintf("%d", repackedTileID),
		Image: tile.Data(),
		Data:  data,
	})
	return repackedTileID
}

func (r *tileRegistry) add(reg map[tileKey]tsx.GlobalTileID, tile *tsx.Tile, offset int) (gid tsx.GlobalTileID, created bool) {
	key := tileKey{tile.Tileset.Name, tile.ID}
	if repackedTileID, ok := reg[key]; ok {
		return repackedTileID, false
	}

	repackedTileID := tsx.GlobalTileID(len(reg) + offset)
	reg[key] = repackedTileID
	return repackedTileID, true
}

func (r *tileRegistry) buildTileset(usedTiles map[tsx.GlobalTileID]bool) *tsx.Tileset {
	tileCount := len(r.tiles)
	tileWidth := 16
	tileHeight := 16
	if tileCount != 0 {
		tileWidth = r.tiles[0].Tileset.TileWidth
		tileHeight = r.tiles[0].Tileset.TileHeight
	}

	width := r.atlas.Width()
	tileColumns := width / tileWidth
	tileset := &tsx.Tileset{
		FirstGID:   1,
		Name:       "tileset",
		TileWidth:  tileWidth,
		TileHeight: tileHeight,
		TileCount:  tileCount,
		Columns:    tileColumns,
		Image: &tsx.Image{
			Source: "tileset.png",
			Width:  width,
			Height: r.atlas.Height(),
		},
	}

	for i, tile := range r.tiles {
		if !usedTiles[tileset.GlobalTileID(tsx.LocalTileID(i))] {
			continue
		}

		repackedTile := &tsx.Tile{
			ID:          tsx.LocalTileID(i),
			Type:        tile.Type,
			ObjectGroup: tile.ObjectGroup,
			Properties:  tile.Properties,
			Tileset:     tileset,
		}
		for _, frame := range tile.Animation {
			animTile := tile.Tileset.Tile(tile.Tileset.GlobalTileID(frame.TileID))
			animKey := tileKey{animTile.Tileset.Name, animTile.ID}
			repackedAnimTileID := r.repackedTiles[animKey]
			repackedTile.Animation = append(repackedTile.Animation, tsx.Frame{
				TileID:   tsx.LocalTileID(repackedAnimTileID - 1),
				Duration: frame.Duration,
			})
		}
		tileset.Tiles = append(tileset.Tiles, repackedTile)
	}

	return tileset
}

type mapContext struct {
	tileset   *tileRegistry
	tilemap   *tmx.Map
	usedTiles map[tsx.GlobalTileID]bool
}

func newMapContext(tileset *tileRegistry, tilemap *tmx.Map) *mapContext {
	ctx := &mapContext{
		tileset:   tileset,
		tilemap:   tilemap,
		usedTiles: make(map[tsx.GlobalTileID]bool),
	}

	for _, layer := range tilemap.Layers {
		switch {
		case layer.IsTileLayer():
			for i, tileID := range layer.Data.Decoded {
				if tileID != 0 {
					repackedTileID := ctx.useTile(tileID.WithoutFlags())
					layer.Data.Decoded[i] = repackedTileID.WithFlags(tileID.Flags())
				}
			}
		case layer.IsObjectGroup():
			for i, obj := range layer.Objects {
				if obj.GID != 0 {
					repackedTileID := ctx.useSprite(obj.GID.WithoutFlags())
					layer.Objects[i].GID = repackedTileID.WithFlags(obj.GID.Flags())
				}
			}
		}
	}
	return ctx
}

func (m *mapContext) useTile(tileID tsx.GlobalTileID) tsx.GlobalTileID {
	tile := m.findTile(tileID)
	repackedTileID := m.tileset.addTile(tile)
	m.usedTiles[repackedTileID] = true

	for _, anim := range tile.Animation {
		m.useTile(tile.Tileset.GlobalTileID(anim.TileID))
	}
	return repackedTileID
}

func (m *mapContext) useSprite(tileID tsx.GlobalTileID) tsx.GlobalTileID {
	tile := m.findTile(tileID)
	repackedTileID := m.tileset.addSprite(tile)
	m.usedTiles[repackedTileID] = true

	for _, anim := range tile.Animation {
		m.useSprite(tile.Tileset.GlobalTileID(anim.TileID))
	}
	return repackedTileID
}

func (m *mapContext) save(outputDir string, exportName string) {
	outputTmj := path.Join(outputDir, exportName+".tmj")

	tileset := m.tileset.buildTileset(m.usedTiles)
	m.tilemap.Tilesets = []*tsx.Tileset{tileset}
	dst := tmj.ConvertFromTMX(m.tilemap)
	dst.Save(outputTmj)
}

func (m *mapContext) findTile(id tsx.GlobalTileID) *tsx.Tile {
	for _, ts := range m.tilemap.Tilesets {
		if tile := ts.Tile(id); tile != nil {
			return tile
		}
	}
	log.Fatalf("Failed to find source tile %d", id)
	return nil
}
