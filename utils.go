package main

import (
	"log"
	"os"
	"strings"
	"time"
)

func EnsureDirectory(dirPath string) {
	stats, err := os.Stat(dirPath)
	if err == nil && stats.IsDir() {
		log.Println("Target directory successfully found")
		return
	}

	if !os.IsNotExist((err)) {
		log.Printf("Error ensuring target directory: %s", err.Error())
		return
	}

	creationErr := os.MkdirAll(dirPath, 0755)

	if creationErr != nil {
		log.Printf("Error creating target folder: %s", creationErr.Error())
		return
	}
}

func GetTargetDirectoryPath() string {
	targetDirPath := baseDirectory + "/" + GetCurrentDateString()
	EnsureDirectory(targetDirPath)
	return targetDirPath
}

func GetCurrentDateString() string {
	currentTime := time.Now()
	dateStr := strings.Split(currentTime.String(), " ")[0]
	formattedDateStr := strings.Join(strings.Split(dateStr, "-"), "")
	return formattedDateStr
}

func GetFileNameFromPath(filePath string) string {
	pathParts := strings.Split(filePath, "/")
	return pathParts[len(pathParts)-1]
}

func GetFileNameWoExtensionFromPath(filePath string) string {
	pathParts := strings.Split(GetFileNameFromPath(filePath), ".")
	return strings.Join(pathParts[:len(pathParts)-1], ".")
}
