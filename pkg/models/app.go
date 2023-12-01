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
	ctx      context.Context
	logger   *utils.CustomLogger
	config   *utils.AppConfig
	LogoIcon []byte
}

// NewApp creates a new App application struct
func NewApp(logger *utils.CustomLogger, config *utils.AppConfig) *App {
	newApp := &App{
		logger: logger,
		config: config,
	}

	return newApp.ensureRequiredDirectories()
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
