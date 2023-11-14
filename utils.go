package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)

func EnsureDirectory(dirPath string) {
	stats, err := os.Stat(dirPath)
	if err == nil && stats.IsDir() {
		fmt.Println("Target directory successfully found")
		return
	}

	if !os.IsNotExist((err)) {
		fmt.Printf("Error ensuring target directory: %s", err.Error())
		return
	}

	creationErr := os.MkdirAll(dirPath, 0755)

	if creationErr != nil {
		fmt.Printf("Error creating target folder: %s", creationErr.Error())
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
