# Phaser + pixel-tools example

This project demonstrates how to integrate `pixel-tools` into a Phaser game development workflow using the Vite plugin.
Also, a [live demo](https://pixel-tools-phaser-example.pages.dev) is available.

## Overview

This example shows:
- Using the `pixel-tools` Vite plugin for both development and production.
- Automatic generation of texture atlases (`atlaspack`), bitmap fonts (`fontpack`), and optimized tilemaps (`tilepack`).
- Loading and using these generated assets in a Phaser 4 scene.
- Hot-reload support: when you modify source assets (like `sprites.yaml` or Tiled maps), the plugin automatically regenerates the packed assets and triggers a browser refresh.

## Project Structure

- `assets/`: Contains source assets (YAML configurations, raw PNGs, TMX maps).
- `public/packed_assets/`: The destination for optimized assets. These files are added to `.gitignore` as they are built from the source assets.
- `src/main.ts`: A Phaser scene that loads and uses the generated assets.
- `vite/`:
    - `assets.mjs`: Actual configuration for the asset pipeline.
    - `config.dev.mjs`: Vite config for development (uses `processAssetsDev`).
    - `config.prod.mjs`: Vite config for production (uses `processAssetsProd`).

## Getting Started

First you need to build the `pixel-tools` project, as this example depends on it,
and then initialize the example project:

```bash
# In the root of pixel-tools project
npm install
npm run build
cd phaser-example
npm install
```

Then, to start the development server with hot-reload:

```bash
npm run dev
```

Or to build the project for production:

```bash
npm run build
```
