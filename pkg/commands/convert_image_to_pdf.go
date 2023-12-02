package commands

import (
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	wlptypes "welovepdf/pkg/types"
	"welovepdf/pkg/utils"
)

func BuildConvertImageToPdf(
	logger *utils.CustomLogger,
	tempDirPath string,
	convertJpegToPdf wlptypes.FileToFileOperation,
) func(config *wlptypes.FileToFileOperationConfig) bool {
	return func(config *wlptypes.FileToFileOperationConfig) bool {
		logger.Debug("ConvertImageToPdf: operation starting")
		fileExt := strings.Replace(strings.ToLower(filepath.Ext(config.SourceFilePath)), ".", "", 1)
		isJpeg := fileExt == "jpg" || fileExt == "jpeg"
		if isJpeg {
			err := convertJpegToPdf(&wlptypes.FileToFileOperationConfig{
				TargetFilePath: utils.SanitizeFilePath(config.TargetFilePath),
				SourceFilePath: config.SourceFilePath,
			})
			if err != nil {
				logger.Error("ConvertImageToPdf operation failed", slog.String("reason", err.Error()))
				return false
			}
			return true
		}

		tempFilePath := utils.GetNewTempFilePath(tempDirPath, "jpg")
		defer os.Remove(tempFilePath)
		err := utils.ConvertImageToJpeg(config.SourceFilePath, tempFilePath)
		if err != nil {
			logger.Error("ConvertImageToPdf operation failed", slog.String("reason", err.Error()))
			return false
		}
		config.SourceFilePath = tempFilePath
		err = convertJpegToPdf(&wlptypes.FileToFileOperationConfig{
			TargetFilePath: utils.SanitizeFilePath(config.TargetFilePath),
			SourceFilePath: config.SourceFilePath,
		})
		if err != nil {
			logger.Error("ConvertImageToPdf operation failed", slog.String("reason", err.Error()))
			return false
		}
		logger.Debug("MergePdfFiles: operation succeeded")
		return true
	}
}
