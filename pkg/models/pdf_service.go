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
	gsCommander *ghostscript.GhostScriptCommander
	outputDir   string
	tempDir     string
	binaryPath  string
	scriptPath  string
}

func NewPdfService(
	assetsDir embed.FS,
	config *utils.AppConfig,
) *PdfService {
	binaryPath := path.Join(config.LocalAssetsDirPath, "bin/gs")
	scriptPath := path.Join(config.LocalAssetsDirPath, "code/viewjpeg.ps")

	pdfService := &PdfService{
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
		slog.Error("Error loading PDF Service assets")
		panic("Error loading PDF Service assets")
	}

	err1 = utils.WriteContentToFileIfNotExists(p.binaryPath, gsBinaryContent)
	err2 = utils.WriteContentToFileIfNotExists(p.scriptPath, viewJpegScriptContent)

	if err1 != nil {
		slog.Error("Error writing GS binary to file", slog.String("reason", err1.Error()))
		panic(err1)
	}

	if err2 != nil {
		slog.Error("Error writing viewJPEG script to file", slog.String("reason", err2.Error()))
		panic(err2)
	}

	return p
}

func (p *PdfService) CompressFile(filePath string, targetImageQuality int) bool {
	targetFilePath := utils.ComputeTargetFilePath(p.outputDir, filePath, "pdf", "_compressed")
	compressSinglePageFile := commands.BuildCompressSinglePageFile(p.gsCommander.ConvertPdfToJpeg, p.gsCommander.ConvertJpegToPdf, p.tempDir)

	compressPdfFile := commands.BuildCompressMultiPagePdfFile(
		compressSinglePageFile,
		p.gsCommander.SplitPdfFile,
		p.gsCommander.MergePdfFiles,
		p.tempDir,
	)

	err := compressPdfFile(targetImageQuality, &wlptypes.FileToFileOperationConfig{
		SourceFilePath: filePath,
		TargetFilePath: targetFilePath,
	})

	if err != nil {
		slog.Error("CompressFile: operation failed", slog.String("reason", err.Error()))
		return false
	}

	return true
}

func (p *PdfService) ConvertImageToPdf(sourceFilePath string) bool {
	slog.Debug("ResizePdfFileToA4 : operation started")
	convertImageToPdf := commands.BuildConvertImageToPdf(p.tempDir, p.gsCommander.ConvertJpegToPdf, utils.ConvertImageToJpeg)
	targetFilePath := utils.ComputeTargetFilePath(p.outputDir, sourceFilePath, "pdf", "_converted")
	slog.Debug("MergePdfFiles: operation starting")
	return convertImageToPdf(&wlptypes.FileToFileOperationConfig{
		SourceFilePath: sourceFilePath,
		TargetFilePath: utils.ComputeTargetFilePath(p.outputDir, targetFilePath, "pdf", ""),
	})
}

func (p *PdfService) MergePdfFiles(targetFilePath string, filePathes []string) bool {
	mergePfFiles := commands.BuildMergePdfFiles(p.gsCommander.MergePdfFiles)
	slog.Debug("MergePdfFiles: operation starting")
	return mergePfFiles(utils.ComputeTargetFilePath(p.outputDir, targetFilePath, "pdf", ""), filePathes)
}

func (p *PdfService) ResizePdfFileToA4(sourceFilePath string) bool {
	resizePdfFileToA4 := commands.BuildResizePdfFileToA4(p.gsCommander.ResizePdfToA4)
	targetFilePath := utils.ComputeTargetFilePath(p.outputDir, sourceFilePath, "pdf", "_resized")
	slog.Debug("MergePdfFiles: operation starting")
	return resizePdfFileToA4(utils.ComputeTargetFilePath(p.outputDir, targetFilePath, "pdf", ""), sourceFilePath)
}

func (p *PdfService) RemoveFile(filePath string) bool {
	err := os.Remove(filePath)
	if err != nil {
		slog.Error("error removing file", slog.String("reason", err.Error()))
		return false
	}
	return true
}
