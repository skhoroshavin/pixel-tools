# Tilepack Example

This is a set of example assets to test the tile pack tool, containing
a collection of [Tiled](https://www.mapeditor.org) maps and their associated
tilesets. If you run `tilepack` in this directory with the following command:
```bash
tilepack ./tilemaps ./output
```
it will create an `output` directory with JSON tilemaps for each
TMX file, an optimized `tileset.png` atlas containing all used tiles and
sprites, and a `tileset.atlas` file with the atlas mappings.

Expected output can be found in the `expected_output` directory.
