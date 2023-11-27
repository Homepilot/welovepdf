package models

import (
	"os"
	"path"
	"welovepdf/pkg/utils"

	"github.com/google/uuid"
	pdfcpu "github.com/pdfcpu/pdfcpu/pkg/api"
)

type PdfHandler struct {
	logger     *utils.CustomLogger
	outputDir  string
	tempDir    string
	binaryPath string
}

func NewPdfHandler(
	logger *utils.CustomLogger,
	outputDir string,
	tempDir string,
	binaryPath string,
) *PdfHandler {
	logger.Info("PdfHandler w/ binaryPath", "binarypath", binaryPath)
	return &PdfHandler{
		logger:     logger,
		outputDir:  outputDir,
		tempDir:    tempDir,
		binaryPath: binaryPath,
	}
}

func (p *PdfHandler) MergePdfFiles(targetFilePath string, filePathes []string, canResize bool) bool {
	p.logger.Info("MergePdfFiles: operation starting")
	if !canResize {
		result := utils.MergePdfFiles(targetFilePath, filePathes)
		if !result {
			p.logger.Error("MergePdfFiles operation failed")
			return false
		}
		p.logger.Info("MergePdfFiles: operation succeeded")
		return true
	}

	tempFilePath := utils.GetNewTempFilePath(p.tempDir, "pdf")
	defer os.Remove(tempFilePath)

	isSuccess := utils.MergePdfFiles(tempFilePath, filePathes)
	if !isSuccess {
		p.logger.Error("Error merging PDF files", "files", filePathes)
		p.logger.Info("MergePdfFiles: operation failed")
		return false
	}

	result := utils.ResizePdfToA4(&utils.FileToFileOperationConfig{
		TargetFilePath: targetFilePath,
		SourceFilePath: tempFilePath,
		BinaryPath:     p.binaryPath,
	})

	if !result {
		p.logger.Error("Error merging files", "files", filePathes)
		p.logger.Info("MergePdfFiles: operation failed")
		return false
	}
	p.logger.Info("MergePdfFiles: operation succeeded")
	return true
}

func (p *PdfHandler) OptimizePdfFile(filePath string) bool {
	p.logger.Info("OptimizePdfFile: operation starting")

	newFileName := utils.AddSuffixToFileName(utils.GetFileNameFromPath(filePath), "_compressed")
	targetFilePath := path.Join(p.outputDir, newFileName)

	err := pdfcpu.OptimizeFile(filePath, targetFilePath, pdfcpu.LoadConfiguration())
	if err != nil {
		p.logger.Error("Error optimizing file", "file", filePath, "reason", err.Error())
		return false
	}

	p.logger.Info("OptimizePdfFile: operation succeeded")
	return true
}

func (p *PdfHandler) CompressFile(filePath string, targetImageQuality int) bool {
	p.logger.Info("CompressFile: operation starting", "targetQuality", targetImageQuality)
	resultFilePath := path.Join(p.outputDir, utils.GetFileNameWoExtensionFromPath(filePath)+"_compressed.pdf")

	pageCount, err := pdfcpu.PageCountFile(filePath)
	if err == nil && pageCount == 1 {
		result := utils.CompressSinglePageFile(p.tempDir, targetImageQuality, &utils.FileToFileOperationConfig{
			SourceFilePath: filePath,
			TargetFilePath: resultFilePath,
			BinaryPath:     p.binaryPath,
		})
		if !result {
			p.logger.Error("CompressFile : error compressing single page file", "file", filePath, "targetQuality", targetImageQuality)
			p.logger.Error("CompressFile operation failed")
			return false
		}
		p.logger.Info("CompressFile operation succeeded")
		return true
	}

	fileId := uuid.New().String()
	tempDirPath1 := path.Join(p.tempDir, fileId, "compress_jpg")
	tempDirPath2 := path.Join(p.tempDir, fileId, "compress_pdf")
	utils.EnsureDirectory(tempDirPath1)
	utils.EnsureDirectory(tempDirPath2)
	defer os.RemoveAll(tempDirPath1)
	defer os.RemoveAll(tempDirPath2)

	// 1. Split file into 1 file per page
	err = pdfcpu.SplitFile(filePath, tempDirPath1, 1, nil)
	if err != nil {
		p.logger.Error("CompressFile : error splitting file, compression aborted", "reason", err.Error())
		p.logger.Error("CompressFile operation failed")
		return false
	}

	// 2. Compress the splitted files
	isCompressionSuccess := utils.CompressAllFilesInDir(p.tempDir, targetImageQuality, &utils.DirToDirOperationConfig{
		SourceDirPath: tempDirPath1,
		TargetDirPath: tempDirPath2,
		BinaryPath:    p.binaryPath,
	})
	if !isCompressionSuccess {
		p.logger.Error("Error compressing files in dir", "tempDirPath1", tempDirPath1, "tempDirPath2", tempDirPath2)
		p.logger.Error("CompressFile operation failed")
		return false
	}

	// 3. Merge all compressed files back into 1
	isMergeSuccess := utils.MergeAllFilesInDir(resultFilePath, tempDirPath2)
	if isMergeSuccess != true {
		p.logger.Error("CompressFile : error during final merge !")
		p.logger.Error("CompressFile operation failed")
		return false
	}

	p.logger.Info("CompressFile operation succeeded")

	return true
}

func (p *PdfHandler) ConvertImageToPdf(filePath string, canResize bool) bool {
	p.logger.Info("ConvertImageToPdf : operation started")
	targetFilePath := path.Join(p.outputDir, utils.GetFileNameWoExtensionFromPath(filePath)+".pdf")
	if canResize {
		targetFilePath = path.Join(p.outputDir, utils.AddSuffixToFileName(utils.GetFileNameWoExtensionFromPath(filePath), "_resized.pdf"))
	}
	tempFilePath := targetFilePath
	if canResize {
		tempFilePath = utils.GetNewTempFilePath(p.tempDir, "pdf")
		defer os.Remove(tempFilePath)
	}

	isSuccess := utils.ConvertImageToPdf(filePath, tempFilePath)
	if !isSuccess {
		p.logger.Error("ConvertImageToPdf : operation failed")
		return false
	}

	if !canResize {
		p.logger.Info("ConvertImageToPdf : operation succeeded")
		return true
	}

	isSuccess = utils.ResizePdfToA4(&utils.FileToFileOperationConfig{
		TargetFilePath: targetFilePath,
		SourceFilePath: tempFilePath,
		BinaryPath:     p.binaryPath,
	})

	if !isSuccess {
		p.logger.Error("ConvertImageToPdf : operation failed at ResizePdfFileToA4")
		return true // file is converted, even though not resized
	}

	p.logger.Info("ConvertImageToPdf : operation succeeded")
	return true
}

func (p *PdfHandler) ResizePdfFileToA4(filePath string) bool {
	p.logger.Info("CreateTempFilesFromUpload : operation started")
	result := utils.ResizePdfToA4(&utils.FileToFileOperationConfig{
		BinaryPath:     p.binaryPath,
		SourceFilePath: filePath,
		TargetFilePath: path.Join(p.outputDir, utils.AddSuffixToFileName(utils.GetFileNameFromPath(filePath), "_resized")),
	})
	if !result {
		p.logger.Error("ResizePdfFileToA4 : operation failed")
		return false
	}
	p.logger.Info("CreateTempFilesFromUpload : operation succeeded")
	return true
}

func (p *PdfHandler) CreateTempFilesFromUpload(fileAsBase64 []byte) string {
	p.logger.Info("CreateTempFilesFromUpload : operation started")
	newFilePath := utils.GetNewTempFilePath(p.tempDir, "pdf")
	err := os.WriteFile(newFilePath, []byte(fileAsBase64), 0755)
	if err != nil {
		p.logger.Error("Error saving data to file", "reason", err.Error())
		return ""
	}
	p.logger.Info("CreateTempFilesFromUpload : operation succeeded")
	return newFilePath
}
