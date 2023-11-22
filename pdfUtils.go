package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"

	pdfApi "github.com/pdfcpu/pdfcpu/pkg/api"
)

type PdfUtils struct{}

func NewPdfUtils() *PdfUtils {
	return &PdfUtils{}
}

func (p *PdfUtils) MergePdfFiles(targetFilePath string, filePathes []string) bool {
	log.Println("MergePdfFiles: operation starting")
	// EnsureTargetDirPath()

	err := pdfApi.MergeCreateFile(filePathes, targetFilePath, pdfApi.LoadConfiguration())
	if err != nil {
		log.Printf("Error merging files: %s", err.Error())
		return false
	}

	log.Println("Merge succeeded")

	return true
}

func (p *PdfUtils) OptimizePdfFile(filePath string) bool {
	log.Println("OptimizePdfFile: operation starting")

	nameParts := strings.Split(GetFileNameFromPath(filePath), ".")
	nameParts[len(nameParts)-2] = nameParts[len(nameParts)-2] + "_compressed"
	targetFilePath := GetTargetDirectoryPath() + "/" + strings.Join(nameParts, ".")

	err := pdfApi.OptimizeFile(filePath, targetFilePath, pdfApi.LoadConfiguration())
	if err != nil {
		log.Printf("Error retrieving targetPath: %s", err.Error())
		return false
	}

	log.Println("Optimization succeeded")

	return true
}

func (p *PdfUtils) ConvertImageToPdf(filePath string, targetDir ...string) bool {
	log.Println("ConvertImageToPdf: operation starting")

	targetDirPath := GetTargetDirectoryPath()
	if len(targetDir) > 0 {
		targetDirPath = targetDir[0]
	}

	targetFilePath := targetDirPath + "/" + GetFileNameWoExtensionFromPath(filePath) + ".pdf"
	conversionError := pdfApi.ImportImagesFile([]string{filePath}, targetFilePath, nil, nil)

	if conversionError != nil {
		log.Printf("Error importing image: %s", conversionError.Error())
		return false
	}

	log.Println("Conversion to PDF succeeded")

	return true
}

func (p *PdfUtils) CompressSinglePageFile(filePath string, targetDirPath string, targetImageQuality int64) bool {
	tempFilePath := path.Join(targetDirPath, GetFileNameWoExtensionFromPath(filePath)+"_compressed.jpeg")
	// log.Println=4 -dPDFSETTINGS=/screen -dNOPAUSE -dQUIET -dBATCH -sOutputFile=output.pdf input.pdf
	convertHQCmd := exec.Command("./binary/gs", "-sDEVICE=jpeg", "-o", tempFilePath, "-dJPEGQ="+fmt.Sprintf("%d", targetImageQuality), "-dNOPAUSE", "-dBATCH", "-dUseCropBox", "-dTextAlphaBits=4", "-dGraphicsAlphaBits=4", "-r140", filePath)
	err := convertHQCmd.Run()
	if err != nil {
		log.Printf("Error converting file to JPEG: %s", err.Error())
		return false
	}

	log.Printf("Success converting file to JPEG: %s \n", tempFilePath)

	isSuccess := p.ConvertImageToPdf(tempFilePath, targetDirPath)

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

func (p *PdfUtils) CompressFile(filePath string, targetImageQuality int64) bool {
	tempDirPath1 := baseDirectory + "/temp/compress"
	tempDirPath2 := baseDirectory + "/temp/compress2"
	os.RemoveAll(tempDirPath1)
	os.RemoveAll(tempDirPath2)
	EnsureDirectory(tempDirPath1)
	EnsureDirectory(tempDirPath2)

	err := pdfApi.SplitFile(filePath, tempDirPath1, 1, nil)
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
		isCompressionSuccess := p.CompressSinglePageFile(path.Join(tempDirPath1, file.Name()), tempDirPath2, targetImageQuality)
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

	outFilePath := path.Join(GetTargetDirectoryPath(), GetFileNameWoExtensionFromPath(filePath)+"_compressed.pdf")
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

// func (p *PdfUtils) CompressFile(filePath string) bool {
// 	tempDirPath1 := baseDirectory + "/temp/compress"
// 	tempDirPath2 := baseDirectory + "/temp/compress2"
// 	os.RemoveAll(tempDirPath1)
// 	os.RemoveAll(tempDirPath2)
// 	EnsureDirectory(tempDirPath1)
// 	EnsureDirectory(tempDirPath2)

// 	doc, err := fitz.New(filePath)
// 	if err != nil {
// 		log.Printf("Error opening PDf file: %s", err.Error())
// 		return false
// 	}

// 	defer doc.Close()

// 	// Extract pages as images
// 	for n := 0; n < doc.NumPage(); n++ {
// 		img, err := doc.Image(n)
// 		if err != nil {
// 			log.Printf("Error converting page %d to image: %s", n, err.Error())
// 			return false
// 		}

// 		f, err := os.Create(path.Join(tempDirPath1, GetFileNameWoExtensionFromPath(filePath)+fmt.Sprintf("%d.jpg", n)))
// 		if err != nil {
// 			log.Printf("Error creating JPEG file for page %d: %s", n, err.Error())
// 			return false
// 		}

// 		err = jpeg.Encode(f, img, &jpeg.Options{jpeg.DefaultQuality})
// 		if err != nil {
// 			log.Printf("Error encoding page %d to JPEG file: %s", n, err.Error())
// 			return false
// 		}

// 		f.Close()
// 	}
// 	return true
// }
