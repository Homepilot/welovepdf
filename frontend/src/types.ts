export enum FileType {
    PDF = 'PDF',
    IMAGE = 'IMAGE'
}

export enum PageName {
    HOME = 'HOME',
    MERGE = 'MERGE',
    CONVERT_IMG = 'CONVERT_IMG',
    COMPRESS = 'COMPRESS',
    RESIZE = 'RESIZE',
}

export enum CompressionMode {
    OPTIMIZE = "Optimisation",
    COMPRESS = "Compression", 
    EXTREME = "Compression extreme"
}

export type FileInfo = {
    name: string;
    id: string;
}