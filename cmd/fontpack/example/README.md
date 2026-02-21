# Font Pack Example

This is a set of example assets to test the font pack tool, containing
single `fonts.yaml` configuration for a tool, and three example `png`
font sheets. If you run `fontpack` in this directory with following command:
```bash
fontpack fonts.yaml ./output
```
it will create `output` directory with a single `fonts.png` atlas file,
containing all the glyphs from the three fonts, and a `.bmfont` file with
symbol descriptions for each of the original PNG font.

Expected output can be found in `expected_output` directory.
