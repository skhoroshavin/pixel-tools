# pixel-tools

A set of CLI tools for building pixel art asset pipelines for games,
with first-class support for the [Phaser](https://phaser.io/) game engine.

## Overview

- **[Atlaspack](#atlaspack)** — Packs multiple PNG images into a single optimized texture atlas with support for nineslice and spritesheets
- **[Fontpack](#fontpack)** — Creates packed texture atlases and BMFont descriptor files from pixel art font sheets
- **[Tilepack](#tilepack)** — Repacks and optimizes Tiled map files (TMX) and their tilesets, outputting Phaser-compatible JSON tilemaps
- **[Recolor](#recolor)** — Converts image colors using a color lookup table (LUT) built from example image pairs

All tools are written in Go and distributed as standalone binaries.
They are also available as an npm package with an included Vite plugin
for automatic asset processing during development and production builds.

## Installation

### Manually from release archive

1.  Download the latest release archive for your platform (Windows, Linux, or macOS) 
    from [itch.io](https://skhoroshavin.itch.io/pixel-tools) or the
    [GitHub](https://github.com/skhoroshavin/pixel-tools/releases) releases page.
2.  Extract the archive to a folder of your choice.
3.  Add the folder containing the binaries to your system's `PATH` to use them from any directory.

### Using npm

This package is also available on npm, making it easy to integrate into web game development workflows.

```bash
npm install --save-dev pixel-tools
```

You can then use the tools directly via `npx`:

```bash
npx atlaspack atlas.yaml output/game
npx fontpack fonts.yaml output/
npx tilepack tilemaps/ output/
```

Or add commands to the `package.json` scripts section to create a reusable asset processing script.

```json
{
  "scripts": {
    "build:assets": "atlaspack assets/atlas.yaml src/assets/game && fontpack assets/fonts.yaml src/assets/fonts"
  }
}
```

The npm package also includes a Vite plugin that automatically runs pixel-tools
on each build, both during development (with hot-reload when source assets change)
and for production builds. For a usage example, see the
[phaser-example](phaser-example) folder.

## Atlaspack

A CLI tool for creating packed texture atlases from multiple PNG images. It supports
individual sprites, nineslice definitions, and spritesheets with named animations.

### Features

- Packs multiple PNG images into a single optimized texture atlas
- Trims empty space around sprites to minimize atlas size
- Supports nineslice metadata for UI elements
- Supports extracting sprites from spritesheet PNGs:
  - Individual sprites can be extracted just by specifying sprite width and height
  - Named sprite sequences can be defined with automatic deduplication of frames referenced multiple times
  - Spritesheet definitions can be imported from external YAML files
- Generates JSON `.atlas` file compatible with the Phaser game engine

### Usage

```bash
atlaspack <config.yaml> <output-base>
```

The tool will process the configuration and output:
- `<output-base>.png`: The packed texture atlas.
- `<output-base>.atlas`: A JSON descriptor file containing sprite coordinates and nineslice data.

### Configuration Example

```yaml
# Name of the sprite in the atlas.
- name: hero
  # Source image path relative to the config file
  image: character.png
  # Optional spritesheet definition
  spritesheet:
    # Width and height of each sprite in the sheet
    sprite_width: 48
    sprite_height: 48
    # Optional mapping of sprite indices to names or animation sequences
    sprite_names:
      idle_front: 0
      walk_front: [0, 1, 2, 0, 4, 5]
# Another sprite with 9-slice borders.
- name: frame
  image: gui_9Slices.png
  nineslice:
    x: 8
    y: 8
    w: 32
    h: 32
```

For a more detailed example, please see the [example folder](cmd/atlaspack/example), which
contains a complete configuration with source images and [expected output](cmd/atlaspack/example/expected_output).

## Fontpack

A CLI tool for creating packed texture atlases and [BMFont](https://www.angelcode.com/products/bmfont/)
character descriptor files from pixel art font sheets. It automates glyph extraction,
atlas packing, and coordinate mapping.

### Features

- Trims empty space around glyphs to minimize atlas size
- Automatically calculates line height and offsets for consistent baseline alignment
- Packs multiple fonts into a single texture atlas
- Generates XML-based `.bmfont` files, which can be loaded directly by the Phaser game engine

### Usage

```bash
fontpack <font-config.yaml> <output-dir>
```

The tool expects a YAML configuration file and an output directory. Source PNG
images are resolved relative to the config file (named after the font's `name` property). It outputs:
- `fonts.png`: A packed texture atlas containing all glyphs.
- `<font_name>.bmfont`: Individual XML descriptor files for each font.

### Configuration Example

```yaml
# Name of the font. In this case awesome.png will be searched in the same
# folder, where this config.yaml is located, and awesome.bmfont will be generated
# in the output directory
- name: awesome
  # Font cell width and height in the source image in pixels
  size: 16
  # Space between lines in pixels. Line height will be calculated automatically
  # as a height of the biggest character plus line_spacing
  line_spacing: 1
  # Space between letters in pixels
  letter_spacing: 1
  # Width of the space character in pixels
  space_width: 4
  # List of characters to include in the font. Positions of the letters
  # in this list should correspond to positions of letters in the source image
  letters:
    - "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
    - "abcdefghijklmnopqrstuvwxyz"
    - "0123456789!\"#$%&'()*+,-./:"
# Name of another font to include in the atlas
- name: another
  size: 16
  line_spacing: 1
# ...and so on
```

For a more detailed example, please see the [example folder](cmd/fontpack/example), which
contains both an example font config with source font images, and an [output](cmd/fontpack/example/expected_output)
expected to be generated by the tool.

## Tilepack

A CLI tool for repacking and optimizing [Tiled](https://www.mapeditor.org)
map files (*.tmx) and their associated tilesets. It analyzes tile usage
across maps and creates optimized tilesets containing only the tiles that are
actually used.

### Features

- Analyzes TMX map files to identify used tiles
- Creates an optimized tileset by removing unused tiles
- Packs all tiles into a single texture atlas
- Preserves tile properties and object layers
- Supports PNG image format for tilesets
- Handles both embedded and external tilesets
- Outputs JSON tilemaps and a sprite atlas that integrates seamlessly with Phaser engine

### Usage

```bash
tilepack <source-dir> <destination-dir>
```

This command will:
1. Scan all `.tmx` files in the source directory and all associated tilesets
2. Generate optimized output in the destination folder:
    - JSON tilemaps (compatible with Phaser engine)
    - `tileset.png` and `tileset.atlas` containing all used tiles and sprites

Example input and output structures:
```
source/
├── level1.tmx
└── level2.tmx

somewhere-referred-from-tmx-files/
├── terrain.tsx
├── terrain.png
├── trees.tsx
└── trees.png

destination/
├── level1.tmj
├── level2.tmj
├── tileset.png
└── tileset.atlas
```

For a more detailed example, please see the [example folder](cmd/tilepack/example), which
contains sample TMX map files with tilesets, and the [expected output](cmd/tilepack/example/expected_output)
generated by the tool.

## Recolor

Note: this tool has pretty different interface from others, and is planned to be significantly reworked in the future.

A CLI tool for batch color conversion of images using a color lookup
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

All tools are free to download, use, modify, and redistribute for any purpose,
including commercial projects. Full source code is available on GitHub under the
[MIT license](LICENSE).

### Example art assets

Example assets included in this repository were created by Gabriel Lima aka 
[tiopalada](https://tiopalada.itch.io), with small modifications by me. Big thanks to him for
creating beautiful pixel art and putting it into the public domain.
His work is licensed partially under [CC0 1.0 Universal](https://creativecommons.org/publicdomain/zero/1.0/),
and partially under [CC BY 4.0](https://creativecommons.org/licenses/by/4.0/),
for details please see LICENSE files in corresponding example folders.
