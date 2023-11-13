package main

import (
	"os/exec"
	"strings"

	pdfApi "github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type PdfTools struct {
	app App
}

func NewPdfTools(a App) *PdfTools {
	return &PdfTools{app: a}
}

func (t *PdfTools) MergePdfFiles(targetFilePath string, filePathes []string) bool {

	err := pdfApi.MergeCreateFile(filePathes, targetFilePath+".pdf", pdfApi.LoadConfiguration())
	if err != nil {
		runtime.LogErrorf(t.app.ctx, "Error retrieving targetPath: %s", err.Error())
		return false
	}

	return true
}

func (t *PdfTools) CompressPdfFile(filePath string) bool {
	pathParts := strings.Split(filePath, ".")
	pathParts[len(pathParts)-2] = pathParts[len(pathParts)-2] + "_compressed"
	targetFilePath := strings.Join(pathParts, ".")

	err := pdfApi.OptimizeFile(filePath, targetFilePath, pdfApi.LoadConfiguration())
	if err != nil {
		runtime.LogErrorf(t.app.ctx, "Error retrieving targetPath: %s", err.Error())
		return false
	}

	return true
}

func (t *PdfTools) ConvertImageToPdf(filePath string) bool {
	pathParts := strings.Split(filePath, ".")
	pathParts[len(pathParts)-1] = "pdf"
	targetFilePath := strings.Join(pathParts, ".")

	conversionError := pdfApi.ImportImagesFile([]string{filePath}, targetFilePath, nil, nil)

	if conversionError != nil {
		runtime.LogErrorf(t.app.ctx, "Error importing image: %s", conversionError.Error())
		return false
	}

	return true
}

func (t *PdfTools) CompressFile(filePath string) bool {
	// gs -sDEVICE=pdfwrite -dCompatibilityLevel=1.4 -dPDFSETTINGS=/screen -dNOPAUSE -dQUIET -dBATCH -sOutputFile=output.pdf input.pdf
	cmd := exec.Command("gs",
		"-sDEVICE=pdfwrite",
		"-dCompatibilityLevel=1.4",
		"-dSubsetFonts=true",
		"-dUseFlateCompression=true",
		"-dOptimize=true",
		"-dProcessColorModel=/DeviceRGB",
		"-dDownsampleGrayImages=true ",
		"-dGrayImageDownsampleType=/Bicubic",
		"-dGrayImageResolution=75",
		"-dAutoFilterGrayImages=false",
		"-dDownsampleMonoImages=true",
		"-dMonoImageDownsampleType=/Bicubic",
		"-dCompressPages=true",
		"-dMonoImageResolution=75",
		"-dDownsampleColorImages=true ",
		"-dCompressStreams=true ",
		"-dColorImageDownsampleType=/Bicubic",
		"-dColorImageResolution=75",
		"-dImageQuality=80",
		"-dImageUpperPPI=100",
		"-dAutoFilterColorImages=false",
		"-dPDFSETTINGS=/default",
		"-dNOPAUSE",
		"-dQUIET",
		"-dBATCH",
		"-dSAFER",
		"-sOutputFile="+filePath+"_compressed.pdf",
		"-dCompressFonts=true",
		"-r150",
		filePath,
	)
	// cmd := exec.Command("ps2pdf",
	// 	"-dPDFSETTINGS=/screen",
	// 	filePath,
	// 	filePath+"_compressed.pdf",
	// )

	out, err := cmd.Output()
	if err != nil {
		runtime.LogErrorf(t.app.ctx, "Error compressing file: %s", err.Error())
		return false
	}

	runtime.LogInfof(t.app.ctx, "Success compressing file: %s", out)

	return true
}
