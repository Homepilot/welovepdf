package utils

import (
	"log"
	"os"
	"path"

	pdfcpu "github.com/pdfcpu/pdfcpu/pkg/api"
)

func MergePdfFilesLegacy(targetFilePath string, filePathes []string) error {
	return pdfcpu.MergeCreateFile(filePathes, targetFilePath, pdfcpu.LoadConfiguration())
}

func SplitFile(filePath string, targetDirPath string) error {
	return pdfcpu.SplitFile(filePath, targetDirPath, 1, nil)
}

func ConvertImageToPdf(filePath string, targetFilePath string) error {
	return pdfcpu.ImportImagesFile([]string{filePath}, targetFilePath, nil, nil)
}

func CompressSinglePageFile(tempDirPath string, targetImageQuality int, compressionConfig *FileToFileOperationConfig) error {
	tempFilePath := GetNewTempFilePath(tempDirPath, "jpg")
	defer os.Remove(tempFilePath)

	err := convertToLowQualityJpeg(targetImageQuality, &FileToFileOperationConfig{
		SourceFilePath: compressionConfig.SourceFilePath,
		TargetFilePath: tempFilePath,
		BinaryPath:     compressionConfig.BinaryPath,
	})

	if err != nil {
		log.Printf("CompressSinglePageFile : Error converting file to JPG: %s", tempFilePath)
		return err
	}

	return ConvertImageToPdf(tempFilePath, compressionConfig.TargetFilePath)
}

func CompressAllFilesInDir(tempDirPath string, targetImageQuality int, config *DirToDirOperationConfig) error {
	// For each page
	filesToCompress, err := os.ReadDir(config.SourceDirPath)
	if err != nil {
		log.Printf("CompressAllFilesInDir : Error reading files in dir : %s", config.SourceDirPath)
		return err
	}

	log.Printf("found %d compressed files to compress", len(filesToCompress))
	for _, file := range filesToCompress {
		compressionErr := CompressSinglePageFile(tempDirPath, targetImageQuality, &FileToFileOperationConfig{
			SourceFilePath: path.Join(config.SourceDirPath, file.Name()),
			TargetFilePath: path.Join(config.TargetDirPath, file.Name()),
			BinaryPath:     config.BinaryPath,
		})
		if compressionErr != nil {
			return compressionErr
		}
	}
	return nil
}
