package commands

import (
	"log/slog"
	"os"
	"path"
	"strings"
	wlptypes "welovepdf/pkg/types"
	"welovepdf/pkg/utils"
)

func buildMergeAllFilesInDir(logger *utils.CustomLogger, mergePdfFiles wlptypes.FilesToFileOperation) func(config *wlptypes.DirToFileOperationConfig) bool {
	return func(config *wlptypes.DirToFileOperationConfig) bool {
		filesToMerge, err := os.ReadDir(config.SourceDirPath)
		if err != nil {
			logger.Error("Error reading temp dir to merge", slog.String("reason", err.Error()))
			return false
		}
		if len(filesToMerge) < 1 {
			logger.Warn("No files to merge, aborting")
			return false
		}

		filesPathesToMerge := []string{}
		for _, file := range filesToMerge {
			if !file.IsDir() && strings.HasSuffix(file.Name(), ".pdf") {
				filesPathesToMerge = append(filesPathesToMerge, path.Join(config.SourceDirPath, file.Name()))
			}
		}

		err = mergePdfFiles(&wlptypes.FilesToFileOperationConfig{
			SourceFilesPathes: filesPathesToMerge,
			TargetFilePath:    config.TargetFilePath,
		})

		if err != nil {
			logger.Error("Error merging files", slog.String("reason", err.Error()))
			return false

		}

		return true
	}
}
