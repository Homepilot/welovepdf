package utils

import (
	"log/slog"
	"os"
	"path"
)

func CompressSinglePageFile(tempDirPath string, targetImageQuality int, compressionConfig *FileToFileOperationConfig) error {
	tempFilePath := GetNewTempFilePath(tempDirPath, "jpg")
	defer os.Remove(tempFilePath)

	err := convertToLowQualityJpeg(targetImageQuality, &FileToFileOperationConfig{
		SourceFilePath: compressionConfig.SourceFilePath,
		TargetFilePath: tempFilePath,
		BinaryPath:     compressionConfig.BinaryPath,
	})

	if err != nil {
		slog.Error("CompressSinglePageFile : Error converting file to JPG: %s", slog.String("filepath", tempFilePath))
		return err
	}

	return ConvertImageToPdf(&FileToFileOperationConfig{
		BinaryPath:     compressionConfig.BinaryPath,
		TargetFilePath: compressionConfig.TargetFilePath,
		SourceFilePath: tempFilePath,
	})
}

func CompressAllFilesInDir(tempDirPath string, targetImageQuality int, config *DirToDirOperationConfig) error {
	// For each page
	filesToCompress, err := os.ReadDir(config.SourceDirPath)
	if err != nil {
		slog.Error("CompressAllFilesInDir : Error reading files in dir : %s", slog.String("sourcedir", config.SourceDirPath))
		return err
	}

	for _, file := range filesToCompress {
		compressionErr := CompressSinglePageFile(tempDirPath, targetImageQuality, &FileToFileOperationConfig{
			SourceFilePath: path.Join(config.SourceDirPath, file.Name()),
			TargetFilePath: path.Join(config.TargetDirPath, file.Name()),
			BinaryPath:     config.BinaryPath,
		})
		if compressionErr != nil {
			return compressionErr
		}
	}
	return nil
}
