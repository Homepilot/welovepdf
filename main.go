package main

import (
	"context"
	"embed"
	"log"
	"os"
	"path"
	"strings"
	"time"

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
	ensureGhostScriptSetup(gsBinary)

	// Create an instance of the app structure
	app := NewApp()
	pdfUtils := NewPdfUtils()

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
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
			pdfUtils,
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
		logFatalAndPanic("Error retrieving the user's home directory", err)
	}

	var localAssetsDir = path.Join(userHomeDir, "Documents", ".welovepdf")
	localBinDir = path.Join(localAssetsDir, "bin")

	OUTPUT_DIR = path.Join(userHomeDir, "Documents", "welovepdf", getCurrentDateString())
	TEMP_DIR = path.Join(localBinDir, "temp")
	GS_BINARY_PATH = path.Join(localBinDir, "ghostscript_welovepdf")

}

func ensureRequiredDirectories() {
	err := ensureDirectory(localBinDir)
	if err != nil {
		logFatalAndPanic("Error creating local bin directory", err)
	}

	err = ensureDirectory(OUTPUT_DIR)
	if err != nil {
		logFatalAndPanic("Error creating target directory", err)
	}

	err = ensureDirectory(TEMP_DIR)
	if err != nil {
		logFatalAndPanic("Error creating temp directory", err)
	}
}

func onAppClose(ctx context.Context) {

}

func logFatalAndPanic(msg string, err error) {
	log.Fatalf("%s : %s", msg, err.Error())
	panic(err)
}

func getCurrentDateString() string {
	currentTime := time.Now()
	dateStr := strings.Split(currentTime.String(), " ")[0]
	formattedDateStr := strings.Join(strings.Split(dateStr, "-"), "")
	return formattedDateStr
}
