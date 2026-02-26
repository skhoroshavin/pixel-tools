export interface AssetsConfig {
    source_path: string
    destination_path: string
    padding?: number
    copy?: PathConfig[]
    fonts?: PathConfig[]
    atlases?: PathConfig[]
    tilemaps?: PathConfig[]
}

export interface PathConfig {
    source: string
    target?: string
}
