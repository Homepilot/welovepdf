package main

import (
	"context"
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

	runtime.LogPrintf(a.ctx, "Got files !!")
	result.files = files
	return files
}

func (a *App) MergePdfFiles(filePathes []string) bool {
	targetFilePath, err := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{})
	if err != nil {
		runtime.LogErrorf(a.ctx, "Error retrieving targetPath: %s", err.Error())
		return false
	}

	mergeError := pdfApi.MergeCreateFile(filePathes, targetFilePath+".pdf", pdfApi.LoadConfiguration())
	if err != nil {
		runtime.LogErrorf(a.ctx, "Error retrieving targetPath: %s", mergeError.Error())
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
