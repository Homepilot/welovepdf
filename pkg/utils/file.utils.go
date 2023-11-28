package utils

import (
	"log"
	"os"
	"path"
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
		log.Println("Target directory successfully found")
		return nil
	}

	if !os.IsNotExist((err)) {
		log.Printf("Error ensuring target directory: %s", err.Error())
		return err
	}

	creationErr := os.MkdirAll(dirPath, 0755)

	if creationErr != nil {
		log.Printf("Error creating target folder: %s", creationErr.Error())
		return err
	}

	return nil
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
