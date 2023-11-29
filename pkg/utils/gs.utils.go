package utils

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/exec"
	"path"
	"strings"
)

func convertToLowQualityJpeg(targetImageQuality int, config *FileToFileOperationConfig) error {
	log.Printf("converting w/ GS using quality %d, binaryPath '%s', source '%s', target '%s'", targetImageQuality, config.BinaryPath, config.SourceFilePath, config.TargetFilePath)
	convertToLowQualityJpegCmd := exec.Command(config.BinaryPath, "-sDEVICE=jpeg", "-o", config.TargetFilePath, "-dJPEGQ="+fmt.Sprintf("%d", targetImageQuality), "-dNOPAUSE", "-dBATCH", "-dUseCropBox", "-dTextAlphaBits=4", "-dGraphicsAlphaBits=4", "-r140", config.SourceFilePath)
	err := convertToLowQualityJpegCmd.Run()
	return err
}

func convertJpegToPdf(viewJpegFilePath string, config *FileToFileOperationConfig) error {
	convertCmd := exec.Command(
		config.BinaryPath,
		"-dNOSAFER",
		"-sDEVICE=pdfwrite",
		"-o",
		config.TargetFilePath,
		viewJpegFilePath,
		"-c",
		"("+config.SourceFilePath+")",
		"viewJPEG",
	)
	slog.Info("the printed string", slog.String("the string", convertCmd.String()))

	err := convertCmd.Run()
	return err
}

func ResizePdfToA4(config *FileToFileOperationConfig) error {
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
	return err
}

func MergePdfFiles(config *FilesToFileOperationConfig) error {
	mergePdfFilesCmd := exec.Command(
		config.BinaryPath,
		"-dNOPAUSE",
		"-sDEVICE=pdfwrite",
		"-sOUTPUTFILE="+config.TargetFilePath,
		"-dBATCH",
	)
	mergePdfFilesCmd.Args = append(mergePdfFilesCmd.Args, config.SourceFilesPathes...)

	err := mergePdfFilesCmd.Run()
	return err
}

func MergeAllFilesInDir(config *DirToFileOperationConfig) error {
	filesToMerge, err := os.ReadDir(config.SourceDirPath)
	if err != nil {
		log.Printf("Error reading temp dir to merge: %s", err.Error())
		return err
	}
	if len(filesToMerge) < 1 {
		log.Println("No files to merge, aborting")
		return nil
	}

	log.Printf("found %d compressed files to merge", len(filesToMerge))
	filesPathesToMerge := []string{}
	for _, file := range filesToMerge {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".pdf") {
			filesPathesToMerge = append(filesPathesToMerge, path.Join(config.SourceDirPath, file.Name()))
		}
	}

	return MergePdfFiles(&FilesToFileOperationConfig{
		BinaryPath:        config.BinaryPath,
		SourceFilesPathes: filesPathesToMerge,
		TargetFilePath:    config.TargetFilePath,
	})
}

func SplitPdfFile(config *FileToDirOperationConfig) error {
	splitPdfFileCmd := exec.Command(
		config.BinaryPath,
		"-sDEVICE=pdfwrite",
		"-dSAFER",
		"-o",
		path.Join(config.TargetDirPath, "outfile.%d.pdf"),
		config.SourceFilePath,
	)

	err := splitPdfFileCmd.Run()
	return err
}
