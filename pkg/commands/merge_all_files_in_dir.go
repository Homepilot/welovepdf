package commands

import (
	"fmt"
	"log/slog"
	"os"
	"path"
	"strings"
	wlptypes "welovepdf/pkg/types"
)

func buildMergeAllFilesInDir(mergePdfFiles wlptypes.FilesToFileOperation) func(config *wlptypes.DirToFileOperationConfig) error {
	return func(config *wlptypes.DirToFileOperationConfig) error {
		filesToMerge, err := os.ReadDir(config.SourceDirPath)
		if err != nil {
			slog.Error("Error reading temp dir to merge", slog.String("reason", err.Error()))
			return err
		}
		if len(filesToMerge) < 1 {
			return fmt.Errorf("No files to merge, aborting")
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
			slog.Error("Error merging files", slog.String("reason", err.Error()))
			return err

		}

		return nil
	}
}
