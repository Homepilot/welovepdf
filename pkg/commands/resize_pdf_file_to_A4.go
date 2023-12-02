package commands

import (
	"log/slog"
	wlptypes "welovepdf/pkg/types"
	"welovepdf/pkg/utils"
)

func BuildResizePdfFileToA4(resizeFileTo4 wlptypes.FileToFileOperation) func(targetFilePath string, srcFilePath string) bool {
	return func(targetFilePath string, srcFilePath string) bool {
		slog.Debug("MergePdfFiles: operation starting")
		err := resizeFileTo4(&wlptypes.FileToFileOperationConfig{
			TargetFilePath: utils.SanitizeFilePath(targetFilePath),
			SourceFilePath: srcFilePath,
		})
		if err != nil {
			slog.Error("MergePdfFiles operation failed", slog.String("reason", err.Error()))
			return false
		}
		slog.Debug("MergePdfFiles: operation succeeded")
		return true
	}
}
