package commands

import (
	"log/slog"
	wlptypes "welovepdf/pkg/types"
	"welovepdf/pkg/utils"
)

func BuildMergePdfFiles(logger *utils.CustomLogger, mergePdfFiles wlptypes.FilesToFileOperation) func(targetFilePath string, filePathes []string) bool {
	return func(targetFilePath string, filePathes []string) bool {
		logger.Debug("MergePdfFiles: operation starting")
		err := mergePdfFiles(&wlptypes.FilesToFileOperationConfig{
			TargetFilePath:    utils.SanitizeFilePath(targetFilePath),
			SourceFilesPathes: filePathes,
		})
		if err != nil {
			logger.Error("MergePdfFiles operation failed", slog.String("reason", err.Error()))
			return false
		}
		logger.Debug("MergePdfFiles: operation succeeded")
		return true
	}
}
