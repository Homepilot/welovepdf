package models

import (
	"context"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx          context.Context
	outputDir    string
	tempDir      string
	compressIcon []byte
}

// NewApp creates a new App application struct
func NewApp(
	outputDir string,
	tempDir string,
	compressIcon []byte,
) *App {
	return &App{
		outputDir:    outputDir,
		tempDir:      tempDir,
		compressIcon: compressIcon,
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) Startup(ctx context.Context) {
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
	targetFilePath, err := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
		DefaultDirectory: a.outputDir,
	})

	if err != nil {
		runtime.LogErrorf(a.ctx, "Error retrieving targetPath: %s", err.Error())
		return ""
	}

	return targetFilePath
}

func (a *App) ChooseCompressionMode() string {
	selection, err := runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
		Title:        "Mode de compression",
		Message:      "Choississez un mode de compression",
		Buttons:      []string{"Optimisation", "Compression", "Compression extrême", "Annuler"},
		CancelButton: "Annuler",
		Icon:         a.compressIcon,
	})

	if err != nil {
		runtime.LogErrorf(a.ctx, "Error retrieving target compression mode : %s", err.Error())
		return ""
	}

	if selection == "Annuler" {
		return ""
	}

	return selection
}
