package utils

import (
	"fmt"
	"log/slog"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

func getFileExtensionFromPath(inputFilePath string) string {
	pathParts := strings.Split(inputFilePath, ".")
	ext := pathParts[len(pathParts)-1]
	if len(ext) < 2 || len(ext) > 3 {
		return ""
	}
	return ext
}

func EnsureDirectory(dirPath string) error {
	stats, err := os.Stat(dirPath)
	if err == nil && stats.IsDir() {
		return nil
	}

	if !os.IsNotExist((err)) {
		return err
	}

	creationErr := os.MkdirAll(dirPath, 0755)

	if creationErr != nil {
		return err
	}

	return nil
}

func WriteContentToFileIfNotExists(filePath string, content []byte) error {
	_, err := os.Stat(filePath)
	if err == nil {
		return nil
	}
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	err = file.Chmod(0755)
	if err != nil {
		return err
	}

	_, err = file.Write(content)
	return err
}

func GetFileNameFromPath(inputFilePath string) string {
	pathParts := strings.Split(inputFilePath, "/")
	return pathParts[len(pathParts)-1]
}

func GetFileNameWoExtensionFromPath(inputFilePath string) string {
	pathParts := strings.Split(GetFileNameFromPath(inputFilePath), ".")
	return strings.Join(pathParts[:len(pathParts)-1], ".")
}

func GetNewTempFilePath(tempDirPath string, extension string) string {
	newId := uuid.New().String()
	return path.Join(tempDirPath, newId+"."+extension)
}

func AddSuffixToFileName(fileName string, suffix string) string {
	nameParts := strings.Split(GetFileNameFromPath(fileName), ".")
	if len(nameParts) < 2 {
		return fileName + suffix
	}
	nameParts[len(nameParts)-2] = nameParts[len(nameParts)-2] + suffix
	return strings.Join(nameParts, ".")
}

func getCurrentDateStr() string {
	currentTime := time.Now()
	dateStr := strings.Split(currentTime.String(), " ")[0]
	return strings.Join(strings.Split(dateStr, "-"), "")
}

func GetTodaysOutputDir(userHomeDir string) string {
	return path.Join(userHomeDir, "Documents", "welovepdf", getCurrentDateStr())
}

type FindFileConfig struct {
	RootDirPath        string
	Filename           string
	FileSize           int
	FileLastModifiedAt int
}

func FindFileInDirectoryTree(config *FindFileConfig) string {
	matchingFilePath := ""
	lastModifiedStr := fmt.Sprintf("%d", config.FileLastModifiedAt)

	slog.Debug(fmt.Sprintf("Looking for file %s in dir %s", config.Filename, config.RootDirPath), slog.String("searched", config.Filename), slog.Int("size", int(config.FileSize)), slog.String("modified", lastModifiedStr))

	err := filepath.WalkDir(config.RootDirPath, func(currentPath string, dirEntry os.DirEntry, err error) error {
		if err != nil {
			slog.Warn("error walking through directory", slog.String("reason", err.Error()))
			return nil
		}

		if strings.HasPrefix(dirEntry.Name(), ".") {
			return nil
		}

		if !dirEntry.IsDir() &&
			dirEntry.Name() == config.Filename {
			fileInfo, fileInfoErr := dirEntry.Info()
			if fileInfoErr != nil {
				return nil
			}
			modifiedStr := fmt.Sprintf("%d", fileInfo.ModTime().Unix())
			slog.Debug("Found file w/ matching name", slog.String("searched", config.Filename), slog.String("actual", fileInfo.Name()), slog.Int("size", int(fileInfo.Size())), slog.String("modified", modifiedStr))

			if fileInfo.Size() == int64(config.FileSize) &&
				strings.Compare(modifiedStr, lastModifiedStr[:len(modifiedStr)]) == 0 {
				slog.Debug("Path added")
				matchingFilePath = currentPath
			}
			return nil
		}
		return nil
	})

	if err != nil {
		slog.Debug(fmt.Sprintf("Error looking for file %s in dir %s", config.Filename, config.RootDirPath), slog.String("reason", err.Error()))
		return ""
	}

	if matchingFilePath == "" {
		slog.Debug(fmt.Sprintf("No file matching %s in dir %s", config.Filename, config.RootDirPath))
		return ""
	}

	slog.Info("found matching file", slog.String("filepath", matchingFilePath))
	return matchingFilePath
}
