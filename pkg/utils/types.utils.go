package utils

type FilesToFileOperationConfig struct {
	BinaryPath        string
	SourceFilesPathes []string
	TargetFilePath    string
}

type FileToFileOperationConfig struct {
	BinaryPath     string
	SourceFilePath string
	TargetFilePath string
}
type DirToFileOperationConfig struct {
	BinaryPath     string
	SourceDirPath  string
	TargetFilePath string
}
type FileToDirOperationConfig struct {
	BinaryPath     string
	SourceFilePath string
	TargetDirPath  string
}
type DirToDirOperationConfig struct {
	BinaryPath    string
	SourceDirPath string
	TargetDirPath string
}

type FileToFileOperation func(c *FileToFileOperationConfig) bool
type FilesToFileOperation func(c *FilesToFileOperationConfig) bool
type DirToFileOperation func(c *DirToFileOperationConfig) bool
type FileToDirOperation func(c *FileToDirOperationConfig) bool
type DirToDirOperation func(c *DirToDirOperationConfig) bool
