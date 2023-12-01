package models

import (
	"context"
	"embed"
	"log/slog"
	"os"
	"path"
	"strings"
	"welovepdf/pkg/utils"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx          context.Context
	logger       *utils.CustomLogger
	config       *utils.AppConfig
	LogoIcon     []byte
	compressIcon []byte
	resizeA4Icon []byte
}

// NewApp creates a new App application struct
func NewApp(assetsDir embed.FS, logger *utils.CustomLogger, config *utils.AppConfig) *App {
	newApp := &App{
		logger: logger,
		config: config,
	}

	return newApp.ensureRequiredDirectories().loadIconAssets(assetsDir)
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx

	a.logger.Debug("Application starting")
}

func (a *App) Ready(ctx context.Context) {
	a.logger.Info("Application successfully started")
}

func (a *App) Shutdown(ctx context.Context) {
	defer a.logger.Close()
	a.logger.Info("Application closed with no errors")
}

func (a *App) BeforeClose(ctx context.Context) bool {
	a.logger.Debug("BeforeClose fired")
	defer a.logger.Debug("BeforeClose done")

	tempDirRemovalErr := os.RemoveAll(a.config.TempDirPath)
	if tempDirRemovalErr != nil {
		a.logger.Warn("BeforeClose : temp dir removed", slog.String("reason", tempDirRemovalErr.Error()))
	} else {
		a.logger.Debug("BeforeClose : temp dir removed")
	}

	outputDirRemovalErr := os.Remove(a.config.OutputDirPath)
	if outputDirRemovalErr != nil && !os.IsExist(outputDirRemovalErr) {
		a.logger.Warn("BeforeClose : error removing output dir", slog.String("reason", outputDirRemovalErr.Error()))
	}
	return false
}

func (a *App) ensureRequiredDirectories() *App {
	localAssetsDirPath := a.config.LocalAssetsDirPath

	err1 := utils.EnsureDirectory(localAssetsDirPath)
	err2 := utils.EnsureDirectory(path.Join(localAssetsDirPath, "bin"))
	err3 := utils.EnsureDirectory(path.Join(localAssetsDirPath, "code"))
	err4 := utils.EnsureDirectory(a.config.OutputDirPath)
	err5 := utils.EnsureDirectory(a.config.TempDirPath)

	if err1 != nil ||
		err2 != nil ||
		err3 != nil ||
		err4 != nil ||
		err5 != nil {
		errMsg := "Error ensuring required directories for app"
		a.logger.Error(errMsg)
		panic(errMsg)
	}

	return a
}

func (a *App) loadIconAssets(assetsDir embed.FS) *App {
	logoIcon, err := assetsDir.ReadFile("assets/images/logo_light.svg")

	if err != nil {
		a.logger.Error("Error loading Application assets", slog.String("reason", err.Error()))
		panic("Error loading App assets")
	}

	compressIcon, err1 := assetsDir.ReadFile("./images/compress.svg")
	if err1 != nil {
		compressIcon = logoIcon
	}
	resizeA4Icon, err2 := assetsDir.ReadFile("./images/resize_A4.svg")
	if err2 != nil {
		resizeA4Icon = logoIcon
	}

	a.LogoIcon = logoIcon
	a.compressIcon = compressIcon
	a.resizeA4Icon = resizeA4Icon

	return a
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
		a.logger.Error("Error in OpenMultipleFilesDialog", slog.String("reason", err.Error()))
		result.error = err.Error()
		return []string{}
	}

	result.files = files
	return files
}

func (a *App) OpenSaveFileDialog() string {
	targetFilePath, err := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
		DefaultDirectory: a.config.OutputDirPath,
	})

	if err != nil {
		a.logger.Error("Save dialog :error retrieving targetPath", slog.String("reason", err.Error()))
		return ""
	}

	if strings.HasSuffix(targetFilePath, ".pdf") {
		return utils.SanitizeFilePath(targetFilePath)
	}

	return utils.SanitizeFilePath(targetFilePath) + ".pdf"
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
		Icon:         a.LogoIcon,
	}

	if config.Icon == "compress" {
		dialogOptions.Icon = a.compressIcon
	}

	if config.Icon == "resizeA4" {
		dialogOptions.Icon = a.resizeA4Icon
	}

	selection, err := runtime.MessageDialog(a.ctx, dialogOptions)
	if err != nil {
		a.logger.Error("Error retrieving user select value", slog.String("reason", err.Error()))
		return ""
	}

	if selection == cancelBtnLabel {
		return ""
	}

	return selection
}

func (a *App) SearchFileInUserDir(filename string, size int, lastModifiedAt int) string {
	baseSearchConfig := &utils.SearchFileConfig{
		Filename:           filename,
		FileSize:           size,
		FileLastModifiedAt: lastModifiedAt,
	}

	a.logger.Debug("SearchFilepathByName started", slog.String("filename", filename), slog.Int("size", size), slog.Int("lastModif", lastModifiedAt))
	dirsToCheck, err := getDirectoriesToCheck(a.config.UserHomeDir)
	if err != nil {
		a.logger.Error("error reading user home dir", slog.String("reason", err.Error()))
		return ""
	}

	for _, dirName := range dirsToCheck {
		searchConfig := baseSearchConfig
		searchConfig.RootDirPath = path.Join(a.config.UserHomeDir, dirName)
		matchingFilePath := utils.SearchFileInDirectoryTree(searchConfig)
		if matchingFilePath != "" {
			return matchingFilePath
		}
	}

	return ""
}

func getDirectoriesToCheck(userHomeDir string) ([]string, error) {
	dirsToExclude := map[string]bool{"Library": true}
	dirsToCheck := []string{"Desktop", "Downloads", "Documents", "Pictures"}
	includedDirs := map[string]bool{}
	for _, dirName := range dirsToCheck {
		includedDirs[dirName] = true
	}

	homeDirContent, err := os.ReadDir(userHomeDir)
	if err != nil {
		return []string{}, err
	}

	for _, dir := range homeDirContent {
		_dirName := dir.Name()
		isDirAndNotHidden := !strings.HasPrefix(_dirName, ".") && dir.IsDir()
		isToExclude := dirsToExclude[_dirName] || includedDirs[_dirName]
		if isDirAndNotHidden && !isToExclude {
			dirsToCheck = append(dirsToCheck, _dirName)
		}
	}
	return dirsToCheck, nil
}
