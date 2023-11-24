package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	pdfcpu "github.com/pdfcpu/pdfcpu/pkg/api"
)

type PdfUtils struct{}

func NewPdfUtils() *PdfUtils {
	return &PdfUtils{}
}

func (p *PdfUtils) MergePdfFiles(targetFilePath string, filePathes []string) bool {
	log.Println("MergePdfFiles: operation starting")

	err := pdfcpu.MergeCreateFile(filePathes, targetFilePath, pdfcpu.LoadConfiguration())
	if err != nil {
		log.Printf("Error merging files: %s", err.Error())
		return false
	}

	log.Println("Merge succeeded")

	return true
}

func (p *PdfUtils) OptimizePdfFile(filePath string) bool {
	log.Println("OptimizePdfFile: operation starting")

	nameParts := strings.Split(getFileNameFromPath(filePath), ".")
	nameParts[len(nameParts)-2] = nameParts[len(nameParts)-2] + "_compressed"
	targetFilePath := getTargetDirectoryPath() + "/" + strings.Join(nameParts, ".")

	err := pdfcpu.OptimizeFile(filePath, targetFilePath, pdfcpu.LoadConfiguration())
	if err != nil {
		log.Printf("Error retrieving targetPath: %s", err.Error())
		return false
	}

	log.Println("Optimization succeeded")

	return true
}

func (p *PdfUtils) CompressFile(filePath string, targetImageQuality int64) bool {
	tempDirPath1 := baseDirectory + "/temp/compress"
	tempDirPath2 := baseDirectory + "/temp/compress2"
	os.RemoveAll(tempDirPath1)
	os.RemoveAll(tempDirPath2)
	ensureDirectory(tempDirPath1)
	ensureDirectory(tempDirPath2)

	err := pdfcpu.SplitFile(filePath, tempDirPath1, 1, nil)
	if err != nil {
		log.Printf("Error splitting file, compression aborted, error: %s\n", err.Error())
		return false
	}
	log.Println("Split succeeded")
	// For each page
	filesToCompress, err := os.ReadDir(tempDirPath1)
	if err != nil {
		log.Printf("Error reading directory to compress: %s", err.Error())
		return false
	}

	log.Printf("found %d compressed files to compress", len(filesToCompress))
	for _, file := range filesToCompress {
		isCompressionSuccess := p.compressSinglePageFile(path.Join(tempDirPath1, file.Name()), tempDirPath2, targetImageQuality)
		if isCompressionSuccess != true {
			return false
		}
	}

	err = os.RemoveAll(tempDirPath1)
	if err != nil {
		log.Printf("Error removing uncompressed temp dir")
	}

	filesToMerge, err := os.ReadDir(tempDirPath2)

	if err != nil {
		log.Printf("Error reading temp dir to merge: %s", err.Error())
		return false
	}

	log.Printf("found %d compressed files to merge", len(filesToMerge))
	filesPathesToMerge := []string{}
	for _, v := range filesToMerge {
		filesPathesToMerge = append(filesPathesToMerge, path.Join(tempDirPath2, v.Name()))
	}

	outFilePath := path.Join(getTargetDirectoryPath(), getFileNameWoExtensionFromPath(filePath)+"_compressed.pdf")
	isMergeSuccess := p.MergePdfFiles(outFilePath, filesPathesToMerge)

	// err = os.RemoveAll(tempDirPath2)
	if err != nil {
		log.Printf("Error removing compressed temp dir")
	}

	if isMergeSuccess != true {
		log.Println("Error during final merge !")
		return false
	}

	log.Printf("File compression successful: %s", outFilePath)

	// Remove all temp files
	os.RemoveAll(tempDirPath1)
	os.RemoveAll(tempDirPath2)

	return true
}

func (p *PdfUtils) ConvertImageToPdf(filePath string, canResize bool) bool {

	convertedFilePath := filepath.Join(baseDirectory, getFileNameWoExtensionFromPath(filePath)+"_converted.pdf")
	isSuccess := p.convertImageToPdf(filePath, convertedFilePath)
	if isSuccess == false {
		return false
	}

	if !canResize {
		return true
	}

	targetFilePath := filepath.Join(baseDirectory, getFileNameWoExtensionFromPath(filePath)+"_resized.pdf")
	isSuccess = resizePDFToA4(convertedFilePath, targetFilePath)
	if isSuccess == false {
		return true // file is converted, even though not resized
	}

	os.Remove(convertedFilePath)
	return true
}

func (p *PdfUtils) convertImageToPdf(filePath string, targetFilePath string) bool {
	log.Println("convertImageToPdf: operation starting")

	conversionError := pdfcpu.ImportImagesFile([]string{filePath}, targetFilePath, nil, nil)

	if conversionError != nil {
		log.Printf("Error importing image: %s", conversionError.Error())
		return false
	}

	log.Println("Conversion to PDF succeeded")

	return true
}

func (p *PdfUtils) compressSinglePageFile(filePath string, targetDirPath string, targetImageQuality int64) bool {
	tempFilePath := path.Join(targetDirPath, getFileNameWoExtensionFromPath(filePath)+"_compressed.jpeg")

	convertToLowQualityJpegCmd := exec.Command(GS_BINARY_PATH, "-sDEVICE=jpeg", "-o", tempFilePath, "-dJPEGQ="+fmt.Sprintf("%d", targetImageQuality), "-dNOPAUSE", "-dBATCH", "-dUseCropBox", "-dTextAlphaBits=4", "-dGraphicsAlphaBits=4", "-r140", filePath)
	err := convertToLowQualityJpegCmd.Run()
	if err != nil {
		log.Printf("Error converting file to JPEG: %s", err.Error())
		return false
	}

	log.Printf("Success converting file to JPEG: %s \n", tempFilePath)

	isSuccess := p.convertImageToPdf(tempFilePath, targetDirPath)

	if !isSuccess {
		log.Printf("Error converting file back to PDF: %s", tempFilePath)
	}

	removeErr := os.Remove(tempFilePath)
	if removeErr != nil {
		log.Printf("Error removing tempFile: %s \n", tempFilePath)
	}

	log.Printf("One page compression succeeded")
	return true
}

func resizePDFToA4(sourceFilePath string, targetFilePath string) bool {
	log.Printf("Starting resize w/ source : %s, target : %s", sourceFilePath, targetFilePath)
	resizePdfToA4Cmd := exec.Command(
		GS_BINARY_PATH,
		"-o",
		targetFilePath,
		"-sDEVICE=pdfwrite",
		"-sPAPERSIZE=a4",
		"-dFIXEDMEDIA",
		"-dPDFFitPage",
		"-dCompatibilityLevel=1.4",
		sourceFilePath)

	err := resizePdfToA4Cmd.Run()
	if err != nil {
		log.Printf("Error resizing file : %s", err.Error())
		return false
	}

	log.Println("Resize succeeded")
	return true
}
