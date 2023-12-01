package models

import (
	"embed"
	"log/slog"
	"os"
	"path"
	"welovepdf/pkg/utils"

	"github.com/google/uuid"
)

type PdfServiceAssets struct {
	BinaryPath string
	ScriptPath string
}

type PdfService struct {
	logger     *utils.CustomLogger
	outputDir  string
	tempDir    string
	binaryPath string
	scriptPath string
}

func NewPdfService(
	assetsDir embed.FS,
	logger *utils.CustomLogger,
	config *utils.AppConfig,
) *PdfService {
	pdfService := &PdfService{
		logger:     logger,
		outputDir:  config.OutputDirPath,
		tempDir:    config.TempDirPath,
		binaryPath: path.Join(config.LocalAssetsDirPath, "bin/gs"),
		scriptPath: path.Join(config.LocalAssetsDirPath, "code/viewjpeg.ps"),
	}

	return pdfService.init(assetsDir)
}

func (p *PdfService) init(assetsDir embed.FS) *PdfService {
	gsBinaryContent, err1 := assetsDir.ReadFile("assets/bin/gs")
	viewJpegScriptContent, err2 := assetsDir.ReadFile("assets/code/viewjpeg.ps")

	if err1 != nil || err2 != nil {
		p.logger.Error("Error loading PDF Service assets")
		panic("Error loading PDF Service assets")
	}

	err1 = utils.WriteContentToFileIfNotExists(p.binaryPath, gsBinaryContent)
	err2 = utils.WriteContentToFileIfNotExists(p.scriptPath, viewJpegScriptContent)

	if err1 != nil {
		p.logger.Error("Error writing GS binary to file", slog.String("reason", err1.Error()))
		panic(err1)
	}

	if err2 != nil {
		p.logger.Error("Error writing viewJPEG script to file", slog.String("reason", err2.Error()))
		panic(err2)
	}

	return p
}

func (p *PdfService) MergePdfFiles(targetFilePath string, filePathes []string) bool {
	p.logger.Debug("MergePdfFiles: operation starting")
	err := utils.MergePdfFiles(&utils.FilesToFileOperationConfig{
		BinaryPath:        p.binaryPath,
		TargetFilePath:    utils.ComputeTargetFilePath(p.outputDir, targetFilePath, "pdf", ""),
		SourceFilesPathes: filePathes,
	})
	if err != nil {
		p.logger.Error("MergePdfFiles operation failed", slog.String("reason", err.Error()))
		return false
	}
	p.logger.Debug("MergePdfFiles: operation succeeded")
	return true
}

// TODO Split/refacto
func (p *PdfService) CompressFile(filePath string, targetImageQuality int) bool {
	p.logger.Debug("CompressFile: operation starting", slog.Int("targetQuality", targetImageQuality))
	resultFilePath := utils.ComputeTargetFilePath(p.outputDir, filePath, "pdf", "_compressed")

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
	err = utils.CompressAllFilesInDir(p.tempDir, targetImageQuality, p.scriptPath, &utils.DirToDirOperationConfig{
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

func (p *PdfService) RotateImageFile(filePath string, canResize bool) bool {
	slog.Debug("IN FUCKING FUNCTION")
	p.logger.Debug("ConvertImageToPdf : operation started")
	// tempFilePath := utils.GetNewTempFilePath(p.tempDir, "pdf")
	// defer os.Remove(tempFilePath)

	// err := utils.ConvertImageToPdf(p.tempDir, p.scriptPath, &utils.FileToFileOperationConfig{
	// 	BinaryPath:     p.binaryPath,
	// 	TargetFilePath: tempFilePath,
	// 	SourceFilePath: filePath,
	// })
	// if err != nil {
	// 	p.logger.Error("ConvertImageToPdf : operation failed at ConvertImageToPdf", slog.String("reason", err.Error()))
	// 	return false
	// }

	// if !canResize {
	// 	err := os.Rename(tempFilePath, targetFilePath)
	// 	if err != nil {
	// 		p.logger.Error("Error renaming file after conversion", slog.String("reasonMsg", err.Error()))
	// 		return false
	// 	}
	// 	p.logger.Debug("ConvertImageToPdf : operation succeeded")
	// 	return true
	// }

	// err = utils.ResizePdfToA4(&utils.FileToFileOperationConfig{
	// 	BinaryPath:     p.binaryPath,
	// 	TargetFilePath: targetFilePath,
	// 	SourceFilePath: tempFilePath,
	// })
	targetFilePath := utils.ComputeTargetFilePath(p.outputDir, filePath, ".jpg", "_rotated")
	err := utils.RotateImageClockwise90(filePath, targetFilePath)
	if err != nil {
		p.logger.Error("Rotate image failed !!", slog.String("reason", err.Error()))
		return false // file is converted, even though not resized
	}

	p.logger.Debug("ConvertImageToPdf : operation succeeded")
	return true
}

func (p *PdfService) ResizePdfFileToA4(filePath string) bool {
	p.logger.Debug("ResizePdfFileToA4 : operation started")
	targetFilePath := utils.ComputeTargetFilePath(p.outputDir, filePath, "pdf", "_resized")
	err := utils.ResizePdfToA4(&utils.FileToFileOperationConfig{
		BinaryPath:     p.binaryPath,
		SourceFilePath: filePath,
		TargetFilePath: targetFilePath,
	})
	if err != nil {
		p.logger.Error("ResizePdfFileToA4 : operation failed", slog.String("reason", err.Error()))
		return false
	}
	p.logger.Debug("ResizePdfFileToA4 : operation succeeded")
	return true
}

func (p *PdfService) ConvertImageToPdf(filePath string) bool {
	p.logger.Debug("ResizePdfFileToA4 : operation started")
	targetFilePath := utils.ComputeTargetFilePath(p.outputDir, filePath, "pdf", "")
	err := utils.ConvertImageToPdf(p.tempDir, p.scriptPath, &utils.FileToFileOperationConfig{
		BinaryPath:     p.binaryPath,
		SourceFilePath: filePath,
		TargetFilePath: targetFilePath,
	})
	if err != nil {
		p.logger.Error("ResizePdfFileToA4 : operation failed", slog.String("reason", err.Error()))
		return false
	}
	p.logger.Debug("ResizePdfFileToA4 : operation succeeded")
	return true
}

func (p *PdfService) RemoveFile(filePath string) bool {
	err := os.Remove(filePath)
	if err != nil {
		p.logger.Error("error removing file", slog.String("reason", err.Error()))
		return false
	}
	return true
}
