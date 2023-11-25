package utils

import (
	"log"
	"os"
	"path"

	pdfcpu "github.com/pdfcpu/pdfcpu/pkg/api"
)

func MergePdfFiles(targetFilePath string, filePathes []string) bool {
	err := pdfcpu.MergeCreateFile(filePathes, targetFilePath, pdfcpu.LoadConfiguration())
	if err != nil {
		log.Printf("Error merging files: %s", err.Error())
		return false
	}

	log.Println("Merge succeeded")

	return true
}

func MergeAllFilesInDir(sourceDirPath string, targetFilePath string) bool {
	filesToMerge, err := os.ReadDir(sourceDirPath)
	if err != nil {
		log.Printf("Error reading temp dir to merge: %s", err.Error())
		return false
	}

	if len(filesToMerge) < 1 {
		log.Println("No files to merge, aborting")
		return false
	}

	log.Printf("found %d compressed files to merge", len(filesToMerge))
	filesPathesToMerge := []string{}
	for v := 0; v < len(filesToMerge); v += 1 {
		filesPathesToMerge = append(filesPathesToMerge, path.Join(sourceDirPath, filesToMerge[v].Name()))
	}

	return MergePdfFiles(targetFilePath, filesPathesToMerge)

}

func ConvertImageToPdf(filePath string, targetFilePath string) bool {
	log.Println("convertImageToPdf: operation starting")

	conversionError := pdfcpu.ImportImagesFile([]string{filePath}, targetFilePath, nil, nil)

	if conversionError != nil {
		log.Printf("Error importing image: %s", conversionError.Error())
		return false
	}

	log.Println("Conversion to PDF succeeded")

	return true
}

func CompressSinglePageFile(tempDirPath string, targetImageQuality int, compressionConfig *FileToFileOperationConfig) bool {
	tempFilePath := GetNewTempFilePath(tempDirPath, "jpg")
	defer os.Remove(tempFilePath)

	isSuccess := convertToLowQualityJpeg(targetImageQuality, &FileToFileOperationConfig{
		SourceFilePath: compressionConfig.SourceFilePath,
		TargetFilePath: tempFilePath,
		BinaryPath:     compressionConfig.BinaryPath,
	})

	if !isSuccess {
		log.Printf("Error converting file to JPG: %s", tempFilePath)
		return false
	}

	isSuccess = ConvertImageToPdf(tempFilePath, compressionConfig.TargetFilePath)
	if !isSuccess {
		log.Printf("Error converting file back to PDF: %s", tempFilePath)
		return false
	}
	return true
}

func CompressAllFilesInDir(tempDirPath string, targetImageQuality int, config *DirToDirOperationConfig) bool {
	// For each page
	filesToCompress, err := os.ReadDir(config.SourceDirPath)
	if err != nil {
		log.Printf("Error reading files in dir : %s", config.SourceDirPath)
		return false
	}

	log.Printf("found %d compressed files to compress", len(filesToCompress))
	for _, file := range filesToCompress {
		isCompressionSuccess := CompressSinglePageFile(tempDirPath, targetImageQuality, &FileToFileOperationConfig{
			SourceFilePath: path.Join(config.SourceDirPath, file.Name()),
			TargetFilePath: path.Join(config.TargetDirPath, file.Name()),
			BinaryPath:     config.BinaryPath,
		})
		if isCompressionSuccess != true {
			return false
		}
	}
	return true
}
