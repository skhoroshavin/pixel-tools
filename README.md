# pixel-tools
This is a set of utilities useful for pixel art

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


## [Recolor](recolor)

A simple CLI tool for color conversion between images using a color lookup table (LUT).
It has two main functions:

1. Building a LUT by analyzing pairs of source and target images to learn color mappings
2. Applying the LUT to convert new images from source to target color style

The tool was originally developed to convert [Time Fantasy](https://itch.io/c/201945/time-fantasy-rpg-asset-packs)
tilesets to match the SNES-style colors of [Time Elements](https://itch.io/c/3379349/time-elements-snes-style-game-assets)
assets (both created by [finalbossblues](https://finalbossblues.itch.io/)).
However, it can be used for any image color conversion where you have example
pairs showing the desired color mapping.

A pre-built [LUT file](recolor/tf_to_e.lut) is included for Time Fantasy to
Time Elements conversion. Please see the example [config file](recolor/recolor-config.yaml)
for details on all available options.
