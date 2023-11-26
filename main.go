package main

import (
	"context"
	"embed"
	"log"
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
var GS_BINARY_PATH string

func main() {
	initGlobals()
	ensureRequiredDirectories()
	utils.EnsureGhostScriptSetup(GS_BINARY_PATH, gsBinary)

	// Create an instance of the app structure
	app := models.NewApp(OUTPUT_DIR, TEMP_DIR, logoLightIcon, compressIcon, resizeA4Icon)
	pdfHandler := models.NewPdfHandler(OUTPUT_DIR, TEMP_DIR, GS_BINARY_PATH)

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
		log.Fatalf("Error retrieving the user's home directory : %s", err.Error())
	}

	var localAssetsDir = path.Join(userHomeDir, ".welovepdf")
	localBinDir = path.Join(localAssetsDir, "bin")

	OUTPUT_DIR = utils.GetTodaysOutputDir(userHomeDir)
	TEMP_DIR = path.Join(localAssetsDir, "temp")
	GS_BINARY_PATH = path.Join(localBinDir, "ghostscript_welovepdf")

}

func ensureRequiredDirectories() {
	err := utils.EnsureDirectory(localBinDir)
	if err != nil {
		log.Fatalf("Error creating local bin directory : %s", err.Error())
	}

	err = utils.EnsureDirectory(OUTPUT_DIR)
	if err != nil {
		log.Fatalf("Error creating target directory : %s", err.Error())
	}

	err = utils.EnsureDirectory(TEMP_DIR)
	if err != nil {
		log.Fatalf("Error creating temp directory : %s", err.Error())
	}
}

func onAppClose(_ context.Context) {
	log.Println("OnAppClose fired")

	_ = os.RemoveAll(TEMP_DIR)

	outputDirContent, _ := os.ReadDir(OUTPUT_DIR)
	for i := 0; i < len(outputDirContent); i += 1 {
		if outputDirContent[i].IsDir() {
			log.Println("OnAppClose done")
			return
		}

		if !strings.HasPrefix(outputDirContent[i].Name(), ".") {
			log.Println("OnAppClose done")
			return
		}

	}

	log.Println("found no files or directory in output dir, remove output dir")
	_ = os.RemoveAll(OUTPUT_DIR)
	log.Println("OnAppClose done")
}
