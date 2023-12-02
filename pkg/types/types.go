package wlptypes

type FilesToFileOperationConfig struct {
	SourceFilesPathes []string
	TargetFilePath    string
}

type FileToFileOperationConfig struct {
	SourceFilePath string
	TargetFilePath string
}
type DirToFileOperationConfig struct {
	SourceDirPath  string
	TargetFilePath string
}
type FileToDirOperationConfig struct {
	SourceFilePath string
	TargetDirPath  string
}
type DirToDirOperationConfig struct {
	SourceDirPath string
	TargetDirPath string
}

type FileToFileOperation func(c *FileToFileOperationConfig) error
type FilesToFileOperation func(c *FilesToFileOperationConfig) error
type DirToFileOperation func(c *DirToFileOperationConfig) error
type FileToDirOperation func(c *FileToDirOperationConfig) error
type DirToDirOperation func(c *DirToDirOperationConfig) error
