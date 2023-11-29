package utils

import (
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

func findFileInDirectoryTree(rootDirPath string, filename string) (string, error) {

	var files []string

	err := filepath.Walk(rootDirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if !info.IsDir() && strings.HasSuffix(path, filename) {
			files = append(files, path)
		}

		return nil
	})

	if err != nil {
		return "", err
	}

	if len(files) == 0 {
		slog.Warn("no files found during search")
		return "", nil
	}

	for _, file := range files {
		slog.Info(file)
	}
	return files[0], nil
}
