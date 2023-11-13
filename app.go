package main

import (
	"context"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx           context.Context
	baseDirPath   string
	targetDirPath string
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	a.baseDirPath = homeDir + "/Documents/welovepdf"
	a.EnsureBaseDirPath()
	a.targetDirPath = a.baseDirPath + "/" + GetDateString()
	cmd := exec.Command("open /Users/gregoire/Documents")
	cmd.Run()
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
	a.EnsureTargetDirPath()

	targetFilePath, err := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
		DefaultDirectory: a.targetDirPath,
	})

	if err != nil {
		runtime.LogErrorf(a.ctx, "Error retrieving targetPath: %s", err.Error())
		return ""
	}

	return targetFilePath
}

func (a *App) EnsureTargetDirPath() {
	stats, err := os.Stat(a.targetDirPath)
	if err == nil && stats.IsDir() {
		runtime.LogInfo(a.ctx, "Target directory successfully found")
		return
	}
	if !os.IsNotExist((err)) {
		runtime.LogErrorf(a.ctx, "Error ensuring target directory: %s", err.Error())
		return
	}

	creationErr := os.MkdirAll(a.targetDirPath, os.ModePerm)

	if creationErr != nil {
		runtime.LogErrorf(a.ctx, "Error creating target folder: %s", creationErr.Error())
		return
	}

	runtime.LogInfo(a.ctx, "Target folder successfully created")
}

func (a *App) EnsureBaseDirPath() {
	stats, err := os.Stat(a.baseDirPath)
	if err == nil && stats.IsDir() {
		runtime.LogInfo(a.ctx, "Base folder successfully found")
		return
	}
	if !os.IsNotExist((err)) {
		runtime.LogErrorf(a.ctx, "Error ensuring base folder: %s", err.Error())
		return
	}

	creationErr := os.MkdirAll(a.baseDirPath, os.ModePerm)

	if creationErr != nil {
		runtime.LogErrorf(a.ctx, "Error creating base folder: %s", creationErr.Error())
		return
	}

	runtime.LogInfo(a.ctx, "Base folder successfully created")
}

func GetDateString() string {
	currentTime := time.Now()
	dateStr := strings.Split(currentTime.String(), " ")[0]
	formattedDateStr := strings.Join(strings.Split(dateStr, "-"), "")
	return formattedDateStr
}
