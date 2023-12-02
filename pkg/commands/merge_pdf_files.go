package commands

import (
	"log/slog"
	wlptypes "welovepdf/pkg/types"
	"welovepdf/pkg/utils"
)

func BuildMergePdfFiles(mergePdfFiles wlptypes.FilesToFileOperation) func(targetFilePath string, filePathes []string) bool {
	return func(targetFilePath string, filePathes []string) bool {
		slog.Debug("MergePdfFiles: operation starting")
		err := mergePdfFiles(&wlptypes.FilesToFileOperationConfig{
			TargetFilePath:    utils.SanitizeFilePath(targetFilePath),
			SourceFilesPathes: filePathes,
		})
		if err != nil {
			slog.Error("MergePdfFiles operation failed", slog.String("reason", err.Error()))
			return false
		}
		slog.Debug("MergePdfFiles: operation succeeded")
		return true
	}
}
