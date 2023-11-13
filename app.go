package main

import (
	"context"
	"os/exec"
	"strings"

	pdfApi "github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

type SelectFilesResult struct {
	files []string
	error string
}

func (a *App) SelectMultipleFiles(fileType string, selectFilesPrompt string) []string {
	pdfFilters := []runtime.FileFilter{
		{
			DisplayName: "PDF (*.pdf)",
			Pattern:     "*.pdf;*.PDF",
		},
	}
	imageFilters := []runtime.FileFilter{
		{
			DisplayName: "Images (*.png;*.jpg)",
			Pattern:     "*.png;*.jpg;*.jpeg;*.PNG;*.JPG;*.JPEG",
		},
	}

	filters := pdfFilters
	if fileType == "IMAGE" {
		filters = imageFilters
	}

	result := SelectFilesResult{}

	files, err := runtime.OpenMultipleFilesDialog(a.ctx, runtime.OpenDialogOptions{
		Title:   selectFilesPrompt,
		Filters: filters,
	})
	if err != nil {
		runtime.LogPrintf(a.ctx, "Got an error !!")
		result.error = err.Error()
		return []string{}
	}

	result.files = files
	return files
}

func (a *App) OpenSaveFileDialog() string {
	targetFilePath, err := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{})
	if err != nil {
		runtime.LogErrorf(a.ctx, "Error retrieving targetPath: %s", err.Error())
		return ""
	}

	return targetFilePath
}

func (a *App) MergePdfFiles(targetFilePath string, filePathes []string) bool {

	err := pdfApi.MergeCreateFile(filePathes, targetFilePath+".pdf", pdfApi.LoadConfiguration())
	if err != nil {
		runtime.LogErrorf(a.ctx, "Error retrieving targetPath: %s", err.Error())
		return false
	}

	return true
}

func (a *App) CompressPdfFile(filePath string) bool {
	pathParts := strings.Split(filePath, ".")
	pathParts[len(pathParts)-2] = pathParts[len(pathParts)-2] + "_compressed"
	targetFilePath := strings.Join(pathParts, ".")

	err := pdfApi.OptimizeFile(filePath, targetFilePath, pdfApi.LoadConfiguration())
	if err != nil {
		runtime.LogErrorf(a.ctx, "Error retrieving targetPath: %s", err.Error())
		return false
	}

	return true
}

func (a *App) ConvertImageToPdf(filePath string) bool {
	pathParts := strings.Split(filePath, ".")
	pathParts[len(pathParts)-1] = "pdf"
	targetFilePath := strings.Join(pathParts, ".")

	conversionError := pdfApi.ImportImagesFile([]string{filePath}, targetFilePath, nil, nil)

	if conversionError != nil {
		runtime.LogErrorf(a.ctx, "Error importing image: %s", conversionError.Error())
		return false
	}

	return true
}

func (a *App) CompressFile(filePath string) bool {
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
		runtime.LogErrorf(a.ctx, "Error compressing file: %s", err.Error())
		return false
	}

	runtime.LogInfof(a.ctx, "Success compressing file: %s", out)

	return true
}
