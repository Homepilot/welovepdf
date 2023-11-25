package models

import (
	"context"
	"strings"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx          context.Context
	outputDir    string
	tempDir      string
	logoIcon     []byte
	compressIcon []byte
	resizeA4Icon []byte
}

// NewApp creates a new App application struct
func NewApp(
	outputDir string,
	tempDir string,
	logoIcon []byte,
	compressIcon []byte,
	resizeA4Icon []byte,
) *App {
	return &App{
		outputDir:    outputDir,
		tempDir:      tempDir,
		logoIcon:     logoIcon,
		compressIcon: compressIcon,
		resizeA4Icon: resizeA4Icon,
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

	if strings.HasSuffix(targetFilePath, ".pdf") {
		return targetFilePath
	}

	return targetFilePath + ".pdf"
}

type PromptSelectConfig struct {
	Title   string
	Message string
	Buttons []string
	Icon    string
}

func (a *App) PromptUserSelect(config *PromptSelectConfig) string {
	var cancelBtnLabel = "Annuler"
	config.Buttons = append(config.Buttons, cancelBtnLabel)

	dialogOptions := runtime.MessageDialogOptions{
		Title:        config.Title,
		Message:      config.Message,
		Buttons:      config.Buttons,
		CancelButton: "Annuler",
		Icon:         a.logoIcon,
	}

	if config.Icon == "compress" {
		dialogOptions.Icon = a.compressIcon
	}

	if config.Icon == "resizeA4" {
		dialogOptions.Icon = a.resizeA4Icon
	}

	selection, err := runtime.MessageDialog(a.ctx, dialogOptions)

	if err != nil {
		runtime.LogErrorf(a.ctx, "Error retrieving target compression mode : %s", err.Error())
		return ""
	}

	if selection == cancelBtnLabel {
		return ""
	}

	return selection
}