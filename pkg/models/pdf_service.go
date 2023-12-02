package models

import (
	"embed"
	"log/slog"
	"os"
	"path"
	"welovepdf/pkg/commands"
	"welovepdf/pkg/ghostscript"
	wlptypes "welovepdf/pkg/types"
	"welovepdf/pkg/utils"
)

type PdfServiceAssets struct {
	BinaryPath string
	ScriptPath string
}

type PdfService struct {
	logger      *utils.CustomLogger
	gsCommander *ghostscript.GhostScriptCommander
	outputDir   string
	tempDir     string
	binaryPath  string
	scriptPath  string
}

func NewPdfService(
	assetsDir embed.FS,
	logger *utils.CustomLogger,
	config *utils.AppConfig,
) *PdfService {
	binaryPath := path.Join(config.LocalAssetsDirPath, "bin/gs")
	scriptPath := path.Join(config.LocalAssetsDirPath, "code/viewjpeg.ps")

	pdfService := &PdfService{
		logger:      logger,
		outputDir:   config.OutputDirPath,
		tempDir:     config.TempDirPath,
		binaryPath:  binaryPath,
		scriptPath:  scriptPath,
		gsCommander: ghostscript.NewGhostscriptClient(binaryPath, scriptPath),
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

func (p *PdfService) CompressFile(filePath string, targetImageQuality int) bool {
	pageCount, err := p.gsCommander.GetPdfPageCount(filePath)
	if err != nil {
		p.logger.Error("CompressFile : Operation failed")
		p.logger.Debug("Error getting page count, operation failed", slog.String("reason", err.Error()))
		return false
	}
	if pageCount < 1 {
		p.logger.Error("Page count below 1, operation aborted")
		return false
	}

	targetFilePath := path.Join(p.outputDir, utils.SanitizeFilePath(utils.GetFileNameFromPath(filePath)))
	compressSinglePageFile := commands.BuildCompressSinglePageFile(p.gsCommander.ConvertPdfToJpeg, p.gsCommander.ConvertJpegToPdf, p.tempDir)

	if pageCount == 1 {
		err := compressSinglePageFile(targetImageQuality, &wlptypes.FileToFileOperationConfig{
			SourceFilePath: filePath,
			TargetFilePath: targetFilePath,
		})
		if err != nil {
			p.logger.Error("File: operation failed", slog.String("reason", err.Error()))
			return false
		}
		return true
	}

	compressPdfFile := commands.BuildCompressMultiPagePdfFile(
		p.logger,
		compressSinglePageFile,
		p.gsCommander.GetPdfPageCount,
		p.gsCommander.SplitPdfFile,
		p.gsCommander.MergePdfFiles,
		p.tempDir,
	)

	err = compressPdfFile(targetImageQuality, &wlptypes.FileToFileOperationConfig{
		SourceFilePath: filePath,
		TargetFilePath: targetFilePath,
	})

	if err != nil {
		p.logger.Error("CompressFile: operation failed", slog.String("reason", err.Error()))
		return false
	}

	return true
}

func (p *PdfService) ConvertImageToPdf(sourceFilePath string) bool {
	p.logger.Debug("ResizePdfFileToA4 : operation started")
	convertImageToPdf := commands.BuildConvertImageToPdf(p.logger, p.tempDir, p.gsCommander.ConvertJpegToPdf)
	targetFilePath := utils.ComputeTargetFilePath(p.outputDir, sourceFilePath, "pdf", "_resized")
	p.logger.Debug("MergePdfFiles: operation starting")
	return convertImageToPdf(&wlptypes.FileToFileOperationConfig{
		SourceFilePath: sourceFilePath,
		TargetFilePath: utils.ComputeTargetFilePath(p.outputDir, targetFilePath, "pdf", ""),
	})
}

func (p *PdfService) MergePdfFiles(targetFilePath string, filePathes []string) bool {
	mergePfFiles := commands.BuildMergePdfFiles(p.logger, p.gsCommander.MergePdfFiles)
	p.logger.Debug("MergePdfFiles: operation starting")
	return mergePfFiles(utils.ComputeTargetFilePath(p.outputDir, targetFilePath, "pdf", ""), filePathes)
}

func (p *PdfService) ResizePdfFileToA4(sourceFilePath string) bool {
	resizePdfFileToA4 := commands.BuildResizePdfFileToA4(p.logger, p.gsCommander.ResizePdfToA4)
	targetFilePath := utils.ComputeTargetFilePath(p.outputDir, sourceFilePath, "pdf", "_resized")
	p.logger.Debug("MergePdfFiles: operation starting")
	return resizePdfFileToA4(utils.ComputeTargetFilePath(p.outputDir, targetFilePath, "pdf", ""), sourceFilePath)
}

func (p *PdfService) RemoveFile(filePath string) bool {
	err := os.Remove(filePath)
	if err != nil {
		p.logger.Error("error removing file", slog.String("reason", err.Error()))
		return false
	}
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
