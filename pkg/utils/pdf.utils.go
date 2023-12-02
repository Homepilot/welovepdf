package utils

import (
	"log/slog"
	"os"
	"path"
	wlptypes "welovepdf/pkg/types"
)

func CompressSinglePageFile(tempDirPath string, targetImageQuality int, viewJpegScriptPath string, compressionConfig *wlptypes.FileToFileOperationConfig) error {
	tempFilePath := GetNewTempFilePath(tempDirPath, "jpg")
	defer os.Remove(tempFilePath)

	err := convertToLowQualityJpeg(targetImageQuality, &wlptypes.FileToFileOperationConfig{
		SourceFilePath: compressionConfig.SourceFilePath,
		TargetFilePath: tempFilePath,
	})

	if err != nil {
		slog.Error("CompressSinglePageFile : Error converting file to JPG: %s", slog.String("filepath", tempFilePath))
		return err
	}

	return convertJpegToPdf(viewJpegScriptPath, &wlptypes.FileToFileOperationConfig{
		TargetFilePath: compressionConfig.TargetFilePath,
		SourceFilePath: tempFilePath,
	})
}

func CompressAllFilesInDir(tempDirPath string, targetImageQuality int, viewJpegScriptPath string, config *wlptypes.DirToDirOperationConfig) error {
	// For each page
	filesToCompress, err := os.ReadDir(config.SourceDirPath)
	if err != nil {
		slog.Error("CompressAllFilesInDir : Error reading files in dir : %s", slog.String("sourcedir", config.SourceDirPath))
		return err
	}

	for _, file := range filesToCompress {
		compressionErr := CompressSinglePageFile(tempDirPath, targetImageQuality, viewJpegScriptPath, &wlptypes.FileToFileOperationConfig{
			SourceFilePath: path.Join(config.SourceDirPath, file.Name()),
			TargetFilePath: path.Join(config.TargetDirPath, file.Name()),
		})
		if compressionErr != nil {
			return compressionErr
		}
	}
	return nil
}
