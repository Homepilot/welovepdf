package models

import (
	"context"
	"log/slog"
	"os"
	"path"
	"welovepdf/pkg/utils"
)

// App struct
type App struct {
	ctx          context.Context
	config       *utils.AppConfig
	customLogger *utils.Logger
	LogoIcon     []byte
}

// NewApp creates a new App application struct
func NewApp(logger *utils.Logger, config *utils.AppConfig) *App {
	newApp := &App{
		config:       config,
		customLogger: logger,
	}

	return newApp.ensureRequiredDirectories()
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx

	slog.Debug("Application starting")
}

func (a *App) Ready(ctx context.Context) {
	slog.Info("Application successfully started")
}

func (a *App) Shutdown(ctx context.Context) {
	slog.Info("Application closed with no errors")
	a.customLogger.Close()
}

func (a *App) BeforeClose(ctx context.Context) bool {
	slog.Debug("BeforeClose fired")
	defer slog.Debug("BeforeClose done")

	tempDirRemovalErr := os.RemoveAll(a.config.TempDirPath)
	if tempDirRemovalErr != nil {
		slog.Warn("BeforeClose : temp dir removed", slog.String("reason", tempDirRemovalErr.Error()))
	} else {
		slog.Debug("BeforeClose : temp dir removed")
	}

	outputDirRemovalErr := os.Remove(a.config.OutputDirPath)
	if outputDirRemovalErr != nil && !os.IsExist(outputDirRemovalErr) {
		slog.Warn("BeforeClose : error removing output dir", slog.String("reason", outputDirRemovalErr.Error()))
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
		slog.Error(errMsg)
		panic(errMsg)
	}

	return a
}
