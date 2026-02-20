# pixel-tools
This is a set of utilities useful for pixel art

## [Repack](repack)

A CLI tool for repacking and optimizing [Tiled](https://www.mapeditor.org)
map files (*.tmx) and their associated tilesets. The tool analyzes tile usage
in maps and creates optimized tilesets containing only the tiles that are
actually used in the maps.

### Features

- Analyzes TMX map files to identify used tiles
- Creates optimized tilesets by removing unused tiles
- Preserves tile properties and object layers
- Supports PNG image format for tilesets
- Handles both embedded and external tilesets
- Packs differently-sized object layer tiles into separate atlas
- Outputs JSON tilemaps and sprite atlases that integrate seamlessly with phaser.io engine

### Usage

```bash
repack <source-folder> <destination-folder>
```

This command will:
1. Recursively scan all `.tmx` files in the source folder
2. Generate optimized output in the destination folder:
    - JSON tilemaps and atlas mappings (compatible with Phaser.io engine)
    - Compressed PNG files with only used tiles and sprites

Example input and output structures:
```
source/
├── episode1/
│ ├── level1.tmx
│ └── level2.tmx
├── episode2/
│ └── level1.tmx
└── tilesets/
  ├── terrain.png
  └── trees.png

destination/
├── episode1/
│ ├── level1.tmj
│ ├── level1.png
│ ├── level1.atlas
│ ├── level2.tmj
│ ├── level2.png
│ └── level2.atlas
└── episode2/
  ├── level1.tmj
  ├── level1.png
  └── level1.atlas
```

## [Recolor](recolor)

A simple CLI tool for color conversion between images using a color lookup
table (LUT). It has two main functions:

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

## Terms of use

### Code

All tools are free to download and use for any purpose, including commercial
projects. Full source code is available on GitHub under GPL3 license, which
means you're also free to modify and redistribute these tools - the main
restriction is that if (and only if) you're redistributing modified versions,
you must also provide source code of your modifications to your customers.
See the file [LICENSE](https://github.com/skhoroshavin/phaser-pixui/blob/main/LICENSE)
for the full license text.

### Example art assets

The example project included in this repository uses pixel art assets created
by Gabriel Lima aka [tiopalada](https://tiopalada.itch.io), with small modifications
by me. Big thanks to him for creating beautiful pixel art and putting it into the
public domain. His work is licensed under [CC BY 4.0](https://creativecommons.org/licenses/by/4.0/),
which means it can be freely shared and used, even commercially, but attribution
to original author is required, as well as indication whether it was modified.
