
/**
 * Configuration for the asset processing pipeline.
 * @type {import('pixel-tools/config').AssetsConfig}
 */
export const assetsConfig = {
    /** Path to the source assets directory, relative to the project root. */
    source_path: 'assets',

    /** Path where the processed assets will be stored, relative to the project root. */
    destination_path: 'public/packed_assets',

    /** Padding to use when packing textures. */
    padding: 0,

    /**
     * List of font configurations to process using fontpack.
     * Each entry specifies a source YAML file.
     */
    fonts: [{
        source: 'fonts.yaml',
    }],

    /**
     * List of texture atlas configurations to process using atlaspack.
     * Each entry specifies a source YAML file.
     *
     * In a bigger project, this is likely to be looking something more like:
     * [
     *     { source: 'ui.yaml' },
     *     { source: 'level1/sprites.yaml' },
     *     { source: 'level2/sprites.yaml' },
     *     { source: 'level3/sprites.yaml' },
     * ]
     */
    atlases: [
        { source: 'sprites.yaml' },
    ],

    /**
     * List of tilemap configurations to process using tilepack.
     * Each entry specifies a source directory containing TMX files.
     *
     * In a bigger project, this is likely to be looking something like:
     * [
     *     { source: 'level1' },
     *     { source: 'level2' },
     *     { source: 'level3' },
     * ]
     */
    tilemaps: [
        { source: '' }
    ]
};
