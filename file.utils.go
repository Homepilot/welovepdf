package main

import (
	"log"
	"os"
	"path"
	"strings"

	"github.com/google/uuid"
)

func ensureDirectory(dirPath string) error {
	stats, err := os.Stat(dirPath)
	if err == nil && stats.IsDir() {
		log.Println("Target directory successfully found")
		return err
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

func getFileNameFromPath(inputFilePath string) string {
	pathParts := strings.Split(inputFilePath, "/")
	return pathParts[len(pathParts)-1]
}

func getFileExtensionFromPath(inputFilePath string) string {
	pathParts := strings.Split(inputFilePath, ".")
	return pathParts[len(pathParts)-1]
}

func getFileNameWoExtensionFromPath(inputFilePath string) string {
	pathParts := strings.Split(getFileNameFromPath(inputFilePath), ".")
	return strings.Join(pathParts[:len(pathParts)-1], ".")
}

func getNewTempFilePath(extension string) string {
	newId := uuid.New().String()
	return path.Join(TEMP_DIR, newId+"."+extension)
}
