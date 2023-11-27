package models

import (
	"log"
	"log/slog"
	"os"
	"path"
	"welovepdf/pkg/utils"

	"github.com/google/uuid"
	pdfcpu "github.com/pdfcpu/pdfcpu/pkg/api"
)

type PdfHandler struct {
	logger     *slog.Logger
	outputDir  string
	tempDir    string
	binaryPath string
}

func NewPdfHandler(
	logger *slog.Logger,
	outputDir string,
	tempDir string,
	binaryPath string,
) *PdfHandler {
	logger.Info("PdfHandler w/ binaryPath : %s", binaryPath)
	return &PdfHandler{
		logger:     logger,
		outputDir:  outputDir,
		tempDir:    tempDir,
		binaryPath: binaryPath,
	}
}

func (p *PdfHandler) MergePdfFiles(targetFilePath string, filePathes []string, canResize bool) bool {
	if !canResize {
		return utils.MergePdfFiles(targetFilePath, filePathes)
	}

	tempFilePath := utils.GetNewTempFilePath(p.tempDir, "pdf")
	defer os.Remove(tempFilePath)

	isSuccess := utils.MergePdfFiles(tempFilePath, filePathes)
	if !isSuccess {
		p.logger.Error("Error merging PDF files", "files", filePathes)
		return false
	}

	result := utils.ResizePdfToA4(&utils.FileToFileOperationConfig{
		TargetFilePath: targetFilePath,
		SourceFilePath: tempFilePath,
		BinaryPath:     p.binaryPath,
	})

	if !result {
		p.logger.Error("Error merging files", "files", filePathes)
		return false
	}
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
	p.logger.Info("CompressPdfFile: operation starting", "targetQuality", targetImageQuality)
	resultFilePath := path.Join(p.outputDir, utils.GetFileNameWoExtensionFromPath(filePath)+"_compressed.pdf")

	pageCount, err := pdfcpu.PageCountFile(filePath)
	if err == nil && pageCount == 1 {
		result := utils.CompressSinglePageFile(p.tempDir, targetImageQuality, &utils.FileToFileOperationConfig{
			SourceFilePath: filePath,
			TargetFilePath: resultFilePath,
			BinaryPath:     p.binaryPath,
		})
		if !result {
			p.logger.Error("Error compressing single page file", "file", filePath, "targetQuality", targetImageQuality)
			return false
		}
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
		log.Printf("Error splitting file, compression aborted, error: %s\n", err.Error())
		return false
	}

	// 2. Compress the splitted files
	isCompressionSuccess := utils.CompressAllFilesInDir(p.tempDir, targetImageQuality, &utils.DirToDirOperationConfig{
		SourceDirPath: tempDirPath1,
		TargetDirPath: tempDirPath2,
		BinaryPath:    p.binaryPath,
	})
	if !isCompressionSuccess {
		log.Printf("Error compressing files in dir : %s to dir %s", tempDirPath1, tempDirPath2)
		return false
	}

	// 3. Merge all compressed files back into 1
	isMergeSuccess := utils.MergeAllFilesInDir(resultFilePath, tempDirPath2)
	if isMergeSuccess != true {
		log.Println("Error during final merge !")
		return false
	}

	log.Printf("File compression successful: %s", resultFilePath)

	return true
}

func (p *PdfHandler) ConvertImageToPdf(filePath string, canResize bool) bool {
	log.Printf("can resize : %t", canResize)
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
		return false
	}

	if !canResize {
		return true
	}

	isSuccess = utils.ResizePdfToA4(&utils.FileToFileOperationConfig{
		TargetFilePath: targetFilePath,
		SourceFilePath: tempFilePath,
		BinaryPath:     p.binaryPath,
	})

	if !isSuccess {
		return true // file is converted, even though not resized
	}

	return true
}

func (p *PdfHandler) ResizePdfFileToA4(filePath string) bool {
	return utils.ResizePdfToA4(&utils.FileToFileOperationConfig{
		BinaryPath:     p.binaryPath,
		SourceFilePath: filePath,
		TargetFilePath: path.Join(p.outputDir, utils.AddSuffixToFileName(utils.GetFileNameFromPath(filePath), "_resized")),
	})
}

func (p *PdfHandler) CreateTempFilesFromUpload(fileAsBase64 []byte) string {
	newFilePath := utils.GetNewTempFilePath(p.tempDir, "pdf")
	err := os.WriteFile(newFilePath, []byte(fileAsBase64), 0755)
	if err != nil {
		log.Printf("Error saving data to file : %s", err.Error())
		return ""
	}
	return newFilePath
}
