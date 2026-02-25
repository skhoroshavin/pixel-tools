export interface AssetsConfig {
    source_path: string
    destination_path: string
    padding?: number
    fonts?: PathConfig[]
    atlases?: PathConfig[]
    tilemaps?: PathConfig[]
}

export interface PathConfig {
    source: string
    target?: string
}
