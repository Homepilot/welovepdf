package main

import (
	"log"
	"os"
	"strings"
	"time"
)

func ensureDirectory(dirPath string) {
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

func getTargetDirectoryPath() string {
	targetDirPath := baseDirectory + "/" + getCurrentDateString()
	ensureDirectory(targetDirPath)
	return targetDirPath
}

func getCurrentDateString() string {
	currentTime := time.Now()
	dateStr := strings.Split(currentTime.String(), " ")[0]
	formattedDateStr := strings.Join(strings.Split(dateStr, "-"), "")
	return formattedDateStr
}

func getFileNameFromPath(filePath string) string {
	pathParts := strings.Split(filePath, "/")
	return pathParts[len(pathParts)-1]
}

func getFileNameWoExtensionFromPath(filePath string) string {
	pathParts := strings.Split(getFileNameFromPath(filePath), ".")
	return strings.Join(pathParts[:len(pathParts)-1], ".")
}
