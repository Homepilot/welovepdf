package models

import (
	"log/slog"
	"os"
	"path"
	"welovepdf/pkg/utils"

	"github.com/google/uuid"
)

type PdfService struct {
	logger     *utils.CustomLogger
	outputDir  string
	tempDir    string
	binaryPath string
}

func NewPdfService(
	logger *utils.CustomLogger,
	outputDir string,
	tempDir string,
	binaryPath string,
) *PdfService {
	logger.Debug("PdfService w/ binaryPath", slog.String("binarypath", binaryPath))
	return &PdfService{
		logger:     logger,
		outputDir:  outputDir,
		tempDir:    tempDir,
		binaryPath: binaryPath,
	}
}

func (p *PdfService) MergePdfFiles(targetFilePath string, filePathes []string, canResize bool) bool {
	p.logger.Debug("MergePdfFiles: operation starting")

	tempFilePath := utils.GetNewTempFilePath(p.tempDir, "pdf")
	defer os.Remove(tempFilePath)

	err := utils.MergePdfFiles(&utils.FilesToFileOperationConfig{
		BinaryPath:        p.binaryPath,
		TargetFilePath:    tempFilePath,
		SourceFilesPathes: filePathes,
	})
	if err != nil {
		p.logger.Error("MergePdfFiles operation failed", slog.String("reason", err.Error()))
		return false
	}

	if !canResize {
		err := os.Rename(tempFilePath, targetFilePath)
		if err != nil {
			p.logger.Error("Error renamig temp file", slog.String("reason", err.Error()))
			return false
		}
		p.logger.Debug("MergePdfFiles: operation succeeded")
		return true
	}

	err = utils.ResizePdfToA4(&utils.FileToFileOperationConfig{
		TargetFilePath: targetFilePath,
		SourceFilePath: tempFilePath,
		BinaryPath:     p.binaryPath,
	})

	if err != nil {
		p.logger.Error("MergePdfFiles : Error resizing merged file", slog.Int("nbOfFiles", len(filePathes)), slog.String("reason", err.Error()))
		return false
	}
	p.logger.Debug("MergePdfFiles: operation succeeded")
	return true
}

// TODO Split/refacto
func (p *PdfService) CompressFile(filePath string, targetImageQuality int) bool {
	p.logger.Debug("CompressFile: operation starting", slog.Int("targetQuality", targetImageQuality))
	resultFilePath := path.Join(p.outputDir, utils.GetFileNameWoExtensionFromPath(filePath)+"_compressed.pdf")

	// pageCount, err := pdfcpu.PageCountFile(filePath)
	// if err == nil && pageCount == 1 {
	// 	err := utils.CompressSinglePageFile(p.tempDir, targetImageQuality, &utils.FileToFileOperationConfig{
	// 		SourceFilePath: filePath,
	// 		TargetFilePath: resultFilePath,
	// 		BinaryPath:     p.binaryPath,
	// 	})
	// 	if err != nil {
	// 		p.logger.Error("CompressFile : error compressing single page file", slog.String("file", filePath), slog.Int("targetQuality", targetImageQuality), slog.String("reason", err.Error()))
	// 		return false
	// 	}
	// 	p.logger.Debug("CompressFile operation succeeded")
	// 	return true
	// }

	fileId := uuid.New().String()
	tempDirPath1 := path.Join(p.tempDir, fileId, "compress_jpg")
	tempDirPath2 := path.Join(p.tempDir, fileId, "compress_pdf")
	utils.EnsureDirectory(tempDirPath1)
	utils.EnsureDirectory(tempDirPath2)
	defer os.RemoveAll(tempDirPath1)
	defer os.RemoveAll(tempDirPath2)

	// 1. Split file into 1 file per page
	err := utils.SplitPdfFile(&utils.FileToDirOperationConfig{
		BinaryPath:     p.binaryPath,
		SourceFilePath: filePath,
		TargetDirPath:  tempDirPath1,
	})
	if err != nil {
		p.logger.Error("CompressFile : error splitting file, compression aborted", slog.String("reason", err.Error()))
		return false
	}

	// 2. Compress the splitted files
	err = utils.CompressAllFilesInDir(p.tempDir, targetImageQuality, &utils.DirToDirOperationConfig{
		SourceDirPath: tempDirPath1,
		TargetDirPath: tempDirPath2,
		BinaryPath:    p.binaryPath,
	})
	if err != nil {
		p.logger.Error("CompressFile : error compressing files in dir", slog.String("tempDirPath1", tempDirPath1), slog.String("tempDirPath2", tempDirPath2), slog.String("reason", err.Error()))
		return false
	}

	// 3. Merge all compressed files back into 1
	err = utils.MergeAllFilesInDir(&utils.DirToFileOperationConfig{
		BinaryPath:     p.binaryPath,
		SourceDirPath:  tempDirPath2,
		TargetFilePath: resultFilePath,
	})
	if err != nil {
		p.logger.Error("CompressFile : error during final merge !", slog.String("reason", err.Error()))
		return false
	}

	p.logger.Debug("CompressFile operation succeeded")

	return true
}

func (p *PdfService) ConvertImageToPdf(filePath string, canResize bool) bool {
	p.logger.Debug("ConvertImageToPdf : operation started")
	tempFilePath := utils.GetNewTempFilePath(p.tempDir, "pdf")
	defer os.Remove(tempFilePath)

	err := utils.ConvertImageToPdf(&utils.FileToFileOperationConfig{
		BinaryPath:     p.binaryPath,
		TargetFilePath: tempFilePath,
		SourceFilePath: filePath,
	})
	if err != nil {
		p.logger.Error("ConvertImageToPdf : operation failed at ConvertImageToPdf", slog.String("reason", err.Error()))
		return false
	}

	targetFilePath := path.Join(p.outputDir, utils.GetFileNameWoExtensionFromPath(filePath)+".pdf")
	if !canResize {
		err := os.Rename(tempFilePath, targetFilePath)
		if err != nil {
			p.logger.Error("Error renaming file after conversion", slog.String("reasonMsg", err.Error()))
			return false
		}
		p.logger.Debug("ConvertImageToPdf : operation succeeded")
		return true
	}

	err = utils.ResizePdfToA4(&utils.FileToFileOperationConfig{
		BinaryPath:     p.binaryPath,
		TargetFilePath: targetFilePath,
		SourceFilePath: tempFilePath,
	})
	if err != nil {
		p.logger.Error("ConvertImageToPdf : operation failed at ResizePdfFileToA4", slog.String("reason", err.Error()))
		return true // file is converted, even though not resized
	}

	p.logger.Debug("ConvertImageToPdf : operation succeeded")
	return true
}

func (p *PdfService) ResizePdfFileToA4(filePath string) bool {
	p.logger.Debug("CreateTempFilesFromUpload : operation started")
	err := utils.ResizePdfToA4(&utils.FileToFileOperationConfig{
		BinaryPath:     p.binaryPath,
		SourceFilePath: filePath,
		TargetFilePath: path.Join(p.outputDir, utils.AddSuffixToFileName(utils.GetFileNameFromPath(filePath), "_resized")),
	})
	if err != nil {
		p.logger.Error("ResizePdfFileToA4 : operation failed", slog.String("reason", err.Error()))
		return false
	}
	p.logger.Debug("CreateTempFilesFromUpload : operation succeeded")
	return true
}

func (p *PdfService) CreateTempFilesFromUpload(fileAsBase64 []byte) string {
	p.logger.Debug("CreateTempFilesFromUpload : operation started")
	newFilePath := utils.GetNewTempFilePath(p.tempDir, "pdf")
	err := os.WriteFile(newFilePath, []byte(fileAsBase64), 0755)
	if err != nil {
		p.logger.Error("Error saving data to file", slog.String("reason", err.Error()))
		return ""
	}
	p.logger.Debug("CreateTempFilesFromUpload : operation succeeded")
	return newFilePath
}
