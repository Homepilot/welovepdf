package main

import (
	"embed"
	"os"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed binary
var gs []byte

//go:embed all:frontend/dist
var assets embed.FS

var baseDirectory string

func main() {
	// Set globals
	homeDirPath, err := os.UserHomeDir()
	baseDirectory = homeDirPath + "/Documents/welovepdf"
	EnsureDirectory(baseDirectory)

	// Create an instance of the app structure
	app := NewApp()
	pdfUtils := NewPdfUtils()

	// Create application with options
	startErr := wails.Run(&options.App{
		Title:  "We   ❤   PDF",
		Width:  777,
		Height: 555,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 42, G: 47, B: 38, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
			pdfUtils,
		},
	})

	if startErr != nil {
		println("Error:", err.Error())
	}
}
