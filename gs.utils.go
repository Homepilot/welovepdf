package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

func ensureGhostScriptSetup(binaryContent []byte) {
	if isGhostScriptSetup() {
		log.Println("GhostScript already setup")
		return
	}

	file, err := os.Create(GS_BINARY_PATH)
	if err != nil {
		logFatalAndPanic("Error creating GhostScript binary file: %s", err)
	}
	defer file.Close()

	err = file.Chmod(755)
	if err != nil {
		logFatalAndPanic("Error make GhostScript binary file executable: %s", err)
	}

	_, err = file.Write(binaryContent)
	if err != nil {
		logFatalAndPanic("Error writing GhostScript binary to target file: %s", err)
	}
}

func isGhostScriptSetup() bool {
	_, err := os.Stat(GS_BINARY_PATH)

	if err == nil {
		return true
	}
	if !os.IsNotExist(err) {
		logFatalAndPanic("Error setting up GhostScript: %s", err)
	}

	return false
}

func convertToLowQualityJpeg(sourceFilePath string, targetFilePath string, targetImageQuality int) bool {
	convertToLowQualityJpegCmd := exec.Command(GS_BINARY_PATH, "-sDEVICE=jpeg", "-o", sourceFilePath, "-dJPEGQ="+fmt.Sprintf("%d", targetImageQuality), "-dNOPAUSE", "-dBATCH", "-dUseCropBox", "-dTextAlphaBits=4", "-dGraphicsAlphaBits=4", "-r140", targetFilePath)
	err := convertToLowQualityJpegCmd.Run()
	if err != nil {
		log.Printf("Error converting file to JPEG: %s", err.Error())
		return false
	}

	return true
}
