package utils

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

func EnsureGhostScriptSetup(gsBinaryPath string, binaryContent []byte) {
	if IsGhostScriptSetup(gsBinaryPath) {
		log.Println("GhostScript already setup")
		return
	}

	log.Println("setting up GhostScript")
	file, err := os.Create(gsBinaryPath)
	if err != nil {
		log.Fatalf("Error creating GhostScript binary file: %s", err.Error())
	}
	defer file.Close()

	err = file.Chmod(0755)
	if err != nil {
		log.Fatalf("Error make GhostScript binary file executable: %s", err.Error())
	}
	log.Println("GhostScript binary file permissions set")

	_, err = file.Write(binaryContent)
	if err != nil {
		log.Fatalf("Error writing GhostScript binary to target file: %s", err.Error())
	}

	log.Println("Ghostscript binary successfully setup")
}

func IsGhostScriptSetup(gsBinaryPath string) bool {
	_, err := os.Stat(gsBinaryPath)

	if err == nil {
		return true
	}
	if !os.IsNotExist(err) {
		log.Fatalf("Error setting up GhostScript: %s", err.Error())
	}

	return false
}

func convertToLowQualityJpeg(targetImageQuality int, config *FileToFileOperationConfig) bool {
	log.Printf("converting w/ GS using quality %d, binaryPath '%s', source '%s', target '%s'", targetImageQuality, config.BinaryPath, config.SourceFilePath, config.TargetFilePath)
	convertToLowQualityJpegCmd := exec.Command(config.BinaryPath, "-sDEVICE=jpeg", "-o", config.SourceFilePath, "-dJPEGQ="+fmt.Sprintf("%d", targetImageQuality), "-dNOPAUSE", "-dBATCH", "-dUseCropBox", "-dTextAlphaBits=4", "-dGraphicsAlphaBits=4", "-r140", config.TargetFilePath)
	err := convertToLowQualityJpegCmd.Run()
	if err != nil {
		log.Printf("Error converting file to JPEG: %s", err.Error())
		return false
	}

	return true
}

func ResizePdfToA4(config *FileToFileOperationConfig) bool {
	log.Printf("binaryPath : %s", config.BinaryPath)
	log.Printf("Starting resize w/ source : %s, target : %s", config.SourceFilePath, config.TargetFilePath)
	resizePdfToA4Cmd := exec.Command(
		config.BinaryPath,
		"-o",
		config.TargetFilePath,
		"-sDEVICE=pdfwrite",
		"-sPAPERSIZE=a4",
		"-dFIXEDMEDIA",
		"-dPDFFitPage",
		"-dCompatibilityLevel=1.4",
		config.SourceFilePath)

	err := resizePdfToA4Cmd.Run()
	if err != nil {
		log.Printf("Error resizing file : %s", err.Error())
		return false
	}

	log.Println("Resize succeeded")
	return true
}
