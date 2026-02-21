# Atlaspack Example

This is a set of example assets to test the atlas pack tool, containing
the `atlas.yaml` configuration and a set of PNG sprites. If you run 
`atlaspack` in this directory with the following command:
```bash
atlaspack atlas.yaml ./output/game
```
it will create an `output` directory with a `game.png` atlas file,
containing all the sprites from the input images, and a `game.atlas` file with
metadata descriptions for each sprite in the atlas.

Expected output can be found in the `expected_output` directory.
