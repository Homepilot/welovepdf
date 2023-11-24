package main

import (
	"embed"
	"log"
	"os"
	"path/filepath"

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

var baseDirectory string
var GS_BINARY_PATH string

func main() {
	// Set globals
	homeDirPath, _ := os.UserHomeDir()
	GS_BINARY_PATH = filepath.Join(homeDirPath, ".welovepdf/gs_binary")
	ensureGhostScriptSetup()
	baseDirectory = homeDirPath + "/Documents/welovepdf"
	ensureDirectory(baseDirectory)

	// Create an instance of the app structure
	app := NewApp()
	pdfUtils := NewPdfUtils()

	// Create application with options
	startErr := wails.Run(&options.App{
		Title:  "We   ❤   PDF",
		Width:  700,
		Height: 777,
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

func ensureGhostScriptSetup() {
	_, err := os.Stat(GS_BINARY_PATH)

	if err == nil {
		log.Println("GhostScript already setup")
		return
	}

	if !os.IsNotExist((err)) {
		log.Fatalf("Error setting up GhostScript: %s", err.Error())
		panic(err)
	}

	ensureDirectory(filepath.Dir(GS_BINARY_PATH))
	file, err := os.Create(GS_BINARY_PATH)
	if err != nil {
		log.Fatalf("Error creating GhostScript binary file: %s", err.Error())
		panic(err)
	}

	defer file.Close()

	err = file.Chmod(755)
	if err != nil {
		log.Fatalf("Error make GhostScript binary file executable: %s", err.Error())
		panic(err)
	}

	_, err = file.Write(gsBinary)
	if err != nil {
		log.Fatalf("Error writing GhostScript binary to target file: %s", err.Error())
		panic(err)
	}

	// Remove gsBinary from memory
	gsBinary = []byte{}
}
