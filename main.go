package main

import (
	"context"
	"embed"
	"log/slog"
	"os"
	"path"
	"strings"

	"welovepdf/pkg/models"
	"welovepdf/pkg/utils"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
)

//go:embed assets/bin/gs
var gsBinary []byte

//go:embed assets/images/resize_A4.svg
var resizeA4Icon []byte

//go:embed assets/images/compress.svg
var compressIcon []byte

//go:embed assets/images/logo_light.svg
var logoLightIcon []byte

//go:embed all:frontend/dist
var assets embed.FS

var localBinDir string
var OUTPUT_DIR string
var TEMP_DIR string
var LOGS_DIR string
var GS_BINARY_PATH string

var logtailSourceToken string = "eEUGryMaDBULCSbtc7358Z6s"
var logger *utils.CustomLogger

func main() {
	initGlobals()
	logger := utils.NewLogger(LOGS_DIR, logtailSourceToken)
	ensureRequiredDirectories()
	utils.EnsureGhostScriptSetup(GS_BINARY_PATH, gsBinary)

	// Create an instance of the app structure
	app := models.NewApp(logger, OUTPUT_DIR, TEMP_DIR, logoLightIcon, compressIcon, resizeA4Icon)
	pdfHandler := models.NewPdfHandler(logger, OUTPUT_DIR, TEMP_DIR, GS_BINARY_PATH)

	// Create application with options
	startErr := wails.Run(&options.App{
		Title:      "We Love PDF",
		Width:      700,
		Height:     777,
		OnShutdown: onAppClose,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 42, G: 47, B: 38, A: 1},
		OnStartup:        app.Startup,
		Bind: []interface{}{
			app,
			pdfHandler,
			logger,
		},
		Mac: &mac.Options{
			About: &mac.AboutInfo{
				Title:   "We   ❤   PDF",
				Message: "by Homepilot @ 2023",
				Icon:    logoLightIcon,
			},
		},
	})

	if startErr != nil {
		logger.Error("Error starting app", slog.String("reason", startErr.Error()))
		panic(startErr)
	}

	logger.Info("Application successfully started")
}

func initGlobals() {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		logger.Error("Error retrieving the user's home directory", slog.String("reason", err.Error()))
	}

	var localAssetsDir = path.Join(userHomeDir, ".welovepdf")
	localBinDir = path.Join(localAssetsDir, "bin")

	OUTPUT_DIR = utils.GetTodaysOutputDir(userHomeDir)
	TEMP_DIR = path.Join(localAssetsDir, "temp")
	LOGS_DIR = path.Join(localAssetsDir, "logs")
	GS_BINARY_PATH = path.Join(localBinDir, "ghostscript_welovepdf")

}

func ensureRequiredDirectories() {
	err := utils.EnsureDirectory(localBinDir)
	if err != nil {
		logger.Error("Error creating local bin directory", slog.String("reason", err.Error()))
	}

	err = utils.EnsureDirectory(OUTPUT_DIR)
	if err != nil {
		logger.Error("Error creating target directory", slog.String("reason", err.Error()))
	}

	err = utils.EnsureDirectory(TEMP_DIR)
	if err != nil {
		logger.Error("Error creating temp directory", slog.String("reason", err.Error()))
	}

	err = utils.EnsureDirectory(LOGS_DIR)
	if err != nil {
		logger.Error("Error creating logs directory", slog.String("reason", err.Error()))
	}
}

func onAppClose(_ context.Context) {
	logger.Info("OnAppClose fired")

	_ = os.RemoveAll(TEMP_DIR)

	outputDirContent, _ := os.ReadDir(OUTPUT_DIR)
	for i := 0; i < len(outputDirContent); i += 1 {
		if outputDirContent[i].IsDir() {
			logger.Info("OnAppClose done")
			return
		}

		if !strings.HasPrefix(outputDirContent[i].Name(), ".") {
			logger.Info("OnAppClose done")
			return
		}

	}

	logger.Info("OnAppClose done")
	logger.Close()

	_ = os.RemoveAll(OUTPUT_DIR)
}
