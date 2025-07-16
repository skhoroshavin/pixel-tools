# pixel-tools
This is a set of utilities useful for pixel art

## [Recolor](recolor)

A simple CLI tool that can build color lookup table (LUT) based on set of
original and their corresponding recolored images, and then apply it to another
set of images in original style to recolor them. This was originally developed
to convert tilesets from [Time Fantasy](https://itch.io/c/201945/time-fantasy-rpg-asset-packs)
asset packs created by [finalbossblues](https://finalbossblues.itch.io/) to
color style of [Time Elements](https://itch.io/c/3379349/time-elements-snes-style-game-assets)
assets from same author, but can be used for any other similar recoloring
needs. For Time Fantasy to Time Elements recoloring a premade [LUT](recolor/tf_to_e.lut)
is supplied. For information on available options please check example [config](recolor/recolor-config.yaml)
file.

## [Repack](repack)

A CLI tool for repacking and optimizing Tiled map files (*.tmx) and their
associated tilesets. The tool analyzes tile usage in maps and creates optimized
tilesets containing only the tiles that are actually used in the maps.

### Features:

- Analyzes TMX map files to identify used tiles
- Creates optimized tilesets by removing unused tiles
- Preserves tile properties and object layers
- Supports PNG image format for tilesets
- Handles both embedded and external tilesets
- All tilesets used in a map should have the same tile size
