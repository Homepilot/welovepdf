package commands

import (
	"log/slog"
	wlptypes "welovepdf/pkg/types"
	"welovepdf/pkg/utils"
)

func BuildResizePdfFileToA4(logger *utils.CustomLogger, resizeFileTo4 wlptypes.FileToFileOperation) func(targetFilePath string, srcFilePath string) bool {
	return func(targetFilePath string, srcFilePath string) bool {
		logger.Debug("MergePdfFiles: operation starting")
		err := resizeFileTo4(&wlptypes.FileToFileOperationConfig{
			TargetFilePath: utils.SanitizeFilePath(targetFilePath),
			SourceFilePath: srcFilePath,
		})
		if err != nil {
			logger.Error("MergePdfFiles operation failed", slog.String("reason", err.Error()))
			return false
		}
		logger.Debug("MergePdfFiles: operation succeeded")
		return true
	}
}
