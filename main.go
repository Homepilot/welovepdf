package main

import (
	"context"
	"embed"
	"os"
	"path"

	"welovepdf/pkg/models"
	"welovepdf/pkg/utils"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
)

//go:embed assets/bin/gs
var gsBinary []byte

//go:embed assets/images/compress.svg
var compressIcon []byte

//go:embed assets/images/logo_light.svg
var logoLightIcon []byte

//go:embed all:frontend/dist
var assets embed.FS

var localBinDir string
var OUTPUT_DIR string
var TEMP_DIR string
var GS_BINARY_PATH string

func main() {
	initGlobals()
	ensureRequiredDirectories()
	utils.EnsureGhostScriptSetup(GS_BINARY_PATH, gsBinary)

	// Create an instance of the app structure
	app := models.NewApp(OUTPUT_DIR, TEMP_DIR, compressIcon)
	pdfHandler := models.NewPdfHandler(OUTPUT_DIR, TEMP_DIR)

	// Create application with options
	startErr := wails.Run(&options.App{
		Title:      "We   ❤   PDF",
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
		println("Error:", startErr.Error())
	}
}

func initGlobals() {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		utils.LogFatalAndPanic("Error retrieving the user's home directory", err)
	}

	var localAssetsDir = path.Join(userHomeDir, "Documents", ".welovepdf")
	localBinDir = path.Join(localAssetsDir, "bin")

	OUTPUT_DIR = utils.GetTodaysOutputDir(userHomeDir)
	TEMP_DIR = path.Join(localBinDir, "temp")
	GS_BINARY_PATH = path.Join(localBinDir, "ghostscript_welovepdf")

}

func ensureRequiredDirectories() {
	err := utils.EnsureDirectory(localBinDir)
	if err != nil {
		utils.LogFatalAndPanic("Error creating local bin directory", err)
	}

	err = utils.EnsureDirectory(OUTPUT_DIR)
	if err != nil {
		utils.LogFatalAndPanic("Error creating target directory", err)
	}

	err = utils.EnsureDirectory(TEMP_DIR)
	if err != nil {
		utils.LogFatalAndPanic("Error creating temp directory", err)
	}
}

func onAppClose(_ context.Context) {
	tempDirContent, _ := os.ReadDir(TEMP_DIR)
	if len(tempDirContent) > 0 {
		_ = os.Remove(TEMP_DIR)
	}

	outputDirContent, _ := os.ReadDir(OUTPUT_DIR)
	if len(outputDirContent) > 0 {
		return
	}
	_ = os.Remove(TEMP_DIR)
}
