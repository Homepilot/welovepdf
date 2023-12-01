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

type UserPrompter struct {
	ctx          context.Context
	logger       *utils.CustomLogger
	config       *utils.AppConfig
	LogoIcon     []byte
	compressIcon []byte
	resizeA4Icon []byte
}

type SelectFilesResult struct {
	files []string
	error string
}

func NewUserPrompter(assetsDir embed.FS, logger *utils.CustomLogger, config *utils.AppConfig) *UserPrompter {
	newUserPrompter := &UserPrompter{
		logger: logger,
		config: config,
	}

	return newUserPrompter.loadIconAssets(assetsDir)
}

// Init is called when the app starts. The context is saved
// so we can call the runtime methods
func (up *UserPrompter) Init(ctx context.Context) {
	up.ctx = ctx

	up.logger.Debug("UserPrompter setup OK")
}

func (up *UserPrompter) SelectMultipleFiles(fileType string, selectFilesPrompt string) []string {
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

	files, err := runtime.OpenMultipleFilesDialog(up.ctx, runtime.OpenDialogOptions{
		Title:   selectFilesPrompt,
		Filters: filters,
	})
	if err != nil {
		up.logger.Error("Error in OpenMultipleFilesDialog", slog.String("reason", err.Error()))
		result.error = err.Error()
		return []string{}
	}

	result.files = files
	return files
}

func (up *UserPrompter) OpenSaveFileDialog() string {
	targetFilePath, err := runtime.SaveFileDialog(up.ctx, runtime.SaveDialogOptions{
		DefaultDirectory: up.config.OutputDirPath,
	})

	if err != nil {
		up.logger.Error("Save dialog :error retrieving targetPath", slog.String("reason", err.Error()))
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

func (up *UserPrompter) PromptUserSelect(config *PromptSelectConfig) string {
	var cancelBtnLabel = "Annuler"
	config.Buttons = append(config.Buttons, cancelBtnLabel)

	dialogOptions := runtime.MessageDialogOptions{
		Title:        config.Title,
		Message:      config.Message,
		Buttons:      config.Buttons,
		CancelButton: "Annuler",
		Icon:         up.LogoIcon,
	}

	if config.Icon == "compress" {
		dialogOptions.Icon = up.compressIcon
	}

	if config.Icon == "resizeA4" {
		dialogOptions.Icon = up.resizeA4Icon
	}

	selection, err := runtime.MessageDialog(up.ctx, dialogOptions)
	if err != nil {
		up.logger.Error("Error retrieving user select value", slog.String("reason", err.Error()))
		return ""
	}

	if selection == cancelBtnLabel {
		return ""
	}

	return selection
}

func (up *UserPrompter) SearchFileInUserDir(filename string, size int, lastModifiedAt int) string {
	baseSearchConfig := &utils.SearchFileConfig{
		Filename:           filename,
		FileSize:           size,
		FileLastModifiedAt: lastModifiedAt,
	}

	up.logger.Debug("SearchFilepathByName started", slog.String("filename", filename), slog.Int("size", size), slog.Int("lastModif", lastModifiedAt))
	dirsToCheck, err := getDirectoriesToCheck(up.config.UserHomeDir)
	if err != nil {
		up.logger.Error("error reading user home dir", slog.String("reason", err.Error()))
		return ""
	}

	for _, dirName := range dirsToCheck {
		searchConfig := baseSearchConfig
		searchConfig.RootDirPath = path.Join(up.config.UserHomeDir, dirName)
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

func (up *UserPrompter) loadIconAssets(assetsDir embed.FS) *UserPrompter {
	logoIcon, err := assetsDir.ReadFile("assets/images/logo_light.svg")

	if err != nil {
		up.logger.Error("Error loading Application assets", slog.String("reason", err.Error()))
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

	up.LogoIcon = logoIcon
	up.compressIcon = compressIcon
	up.resizeA4Icon = resizeA4Icon

	return up
}
