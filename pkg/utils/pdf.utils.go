package utils

import (
	"log/slog"
	"os"
	"path"
	"path/filepath"
	"strings"
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

	return convertJpegToPdf(&FileToFileOperationConfig{
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

func ConvertImageToPdf(tempDir string, config *FileToFileOperationConfig) error {
	fileExt := strings.ToLower(filepath.Ext(config.SourceFilePath))
	isJpeg := fileExt == "jpg" || fileExt == "jpeg"
	if isJpeg {
		return convertJpegToPdf(config)
	}

	tempFilePath := GetNewTempFilePath(tempDir, "jpg")
	defer os.Remove(tempFilePath)
	err := convertImageToJpeg(config.SourceFilePath, tempFilePath)
	if err != nil {
		return err
	}
	config.SourceFilePath = tempFilePath
	return convertJpegToPdf(config)
}
