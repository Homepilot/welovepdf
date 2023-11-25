package main

import (
	"log"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/google/uuid"
	pdfcpu "github.com/pdfcpu/pdfcpu/pkg/api"
)

type PdfUtils struct{}

func NewPdfUtils() *PdfUtils {
	return &PdfUtils{}
}

func (p *PdfUtils) MergePdfFiles(targetFilePath string, filePathes []string) bool {
	return mergePdfFiles(targetFilePath, filePathes)
}

func (p *PdfUtils) OptimizePdfFile(filePath string) bool {
	log.Println("OptimizePdfFile: operation starting")

	nameParts := strings.Split(getFileNameFromPath(filePath), ".")
	nameParts[len(nameParts)-2] = nameParts[len(nameParts)-2] + "_compressed"
	targetFilePath := path.Join(OUTPUT_DIR, strings.Join(nameParts, "."))

	err := pdfcpu.OptimizeFile(filePath, targetFilePath, pdfcpu.LoadConfiguration())
	if err != nil {
		log.Printf("Error retrieving targetPath: %s", err.Error())
		return false
	}

	log.Println("Optimization succeeded")

	return true
}

func (p *PdfUtils) CompressFile(filePath string, targetImageQuality int) bool {
	resultFilePath := path.Join(OUTPUT_DIR, getFileNameWoExtensionFromPath(filePath)+"_compressed.pdf")

	pageCount, err := pdfcpu.PageCountFile(filePath)
	if pageCount == 1 {
		return compressSinglePageFile(filePath, resultFilePath, targetImageQuality)
	}

	fileId := uuid.New().String()
	tempDirPath1 := path.Join(TEMP_DIR, fileId, "compress_jpg")
	tempDirPath2 := path.Join(TEMP_DIR, fileId, "compress_pdf")
	ensureDirectory(tempDirPath1)
	ensureDirectory(tempDirPath2)
	defer os.RemoveAll(tempDirPath1)
	defer os.RemoveAll(tempDirPath2)

	// 1. Split file into 1 file per page
	err = pdfcpu.SplitFile(filePath, tempDirPath1, 1, nil)
	if err != nil {
		log.Printf("Error splitting file, compression aborted, error: %s\n", err.Error())
		return false
	}

	// 2. Compress the splitted files
	isCompressionSuccess := compressAllFilesInDir(tempDirPath1, tempDirPath2, targetImageQuality)
	if !isCompressionSuccess {
		log.Printf("Error compressing files in dir : %s to dir %s", tempDirPath1, tempDirPath2)
		return false
	}

	// 3. Merge all compressed files back into 1
	isMergeSuccess := mergeAllFilesInDir(resultFilePath, tempDirPath2)
	if isMergeSuccess != true {
		log.Println("Error during final merge !")
		return false
	}

	log.Printf("File compression successful: %s", resultFilePath)

	return true
}

func (p *PdfUtils) ConvertImageToPdf(filePath string, canResize bool) bool {
	targetFilePath := path.Join(OUTPUT_DIR, getFileNameWoExtensionFromPath(filePath)+"_resized.pdf")
	tempFilePath := targetFilePath
	if canResize {
		tempFilePath = getNewTempFilePath("pdf")
		defer os.Remove(tempFilePath)
	}

	isSuccess := convertImageToPdf(filePath, tempFilePath)
	if !isSuccess {
		return false
	}

	if !canResize {
		return true
	}

	isSuccess = resizePdfToA4(tempFilePath, targetFilePath)
	if !isSuccess {
		return true // file is converted, even though not resized
	}

	return true
}

func compressAllFilesInDir(sourceDirPath string, targetDirPath string, targetImageQuality int) bool {
	// For each page
	filesToCompress, err := os.ReadDir(sourceDirPath)
	if err != nil {
		log.Printf("Error reading files in dir : %s", sourceDirPath)
		return false
	}

	log.Printf("found %d compressed files to compress", len(filesToCompress))
	for _, file := range filesToCompress {
		isCompressionSuccess := compressSinglePageFile(path.Join(sourceDirPath, file.Name()), targetDirPath, targetImageQuality)
		if isCompressionSuccess != true {
			return false
		}
	}
	return true
}

func mergeAllFilesInDir(sourceDirPath string, targetFilePath string) bool {
	filesToMerge, err := os.ReadDir(sourceDirPath)
	if err != nil {
		log.Printf("Error reading temp dir to merge: %s", err.Error())
		return false
	}

	log.Printf("found %d compressed files to merge", len(filesToMerge))
	filesPathesToMerge := []string{}
	for _, v := range filesToMerge {
		filesPathesToMerge = append(filesPathesToMerge, path.Join(sourceDirPath, v.Name()))
	}

	return mergePdfFiles(targetFilePath, filesPathesToMerge)

}

func convertImageToPdf(filePath string, targetFilePath string) bool {
	log.Println("convertImageToPdf: operation starting")

	conversionError := pdfcpu.ImportImagesFile([]string{filePath}, targetFilePath, nil, nil)

	if conversionError != nil {
		log.Printf("Error importing image: %s", conversionError.Error())
		return false
	}

	log.Println("Conversion to PDF succeeded")

	return true
}

func compressSinglePageFile(filePath string, targetDirPath string, targetImageQuality int) bool {
	tempFilePath := getNewTempFilePath("jpg")
	defer os.Remove(tempFilePath)

	isSuccess := convertToLowQualityJpeg(filePath, tempFilePath, targetImageQuality)
	if !isSuccess {
		log.Printf("Error converting file to JPG: %s", tempFilePath)
		return false
	}

	isSuccess = convertImageToPdf(tempFilePath, targetDirPath)
	if !isSuccess {
		log.Printf("Error converting file back to PDF: %s", tempFilePath)
		return false
	}
	return true
}

func resizePdfToA4(sourceFilePath string, targetFilePath string) bool {
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

func mergePdfFiles(targetFilePath string, filePathes []string) bool {
	err := pdfcpu.MergeCreateFile(filePathes, targetFilePath, pdfcpu.LoadConfiguration())
	if err != nil {
		log.Printf("Error merging files: %s", err.Error())
		return false
	}

	log.Println("Merge succeeded")

	return true
}
