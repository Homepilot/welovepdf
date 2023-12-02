package commands

import (
	"log/slog"
	"os"
	"path"
	wlptypes "welovepdf/pkg/types"
)

func buildCompressAllFilesInDir(
	tempDirPath string,
	compressSinglePageFile func(targetImageQuality int, config *wlptypes.FileToFileOperationConfig) error,
) func(imageQuality int, compressionConfig *wlptypes.DirToDirOperationConfig) error {
	return func(quality int, config *wlptypes.DirToDirOperationConfig) error {
		filesToCompress, err := os.ReadDir(config.SourceDirPath)
		if err != nil {
			slog.Error("CompressAllFilesInDir : Error reading files in dir : %s", slog.String("sourcedir", config.SourceDirPath))
			return err
		}

		for _, file := range filesToCompress {
			compressionErr := compressSinglePageFile(quality, &wlptypes.FileToFileOperationConfig{
				SourceFilePath: path.Join(config.SourceDirPath, file.Name()),
				TargetFilePath: path.Join(config.TargetDirPath, file.Name()),
			})
			if compressionErr != nil {
				return compressionErr
			}
		}
		return nil
	}
}
