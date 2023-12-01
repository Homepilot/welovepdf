package main

import (
	"embed"
	"log/slog"

	"welovepdf/pkg/models"
	"welovepdf/pkg/utils"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
)

//go:embed all:assets
var assets embed.FS

//go:embed all:frontend/dist
var frontendAssets embed.FS

func main() {
	appConfig := utils.GetAppConfigFromAssetsDir(assets)

	// Create an instance of the app structure
	logger := utils.NewLogger(appConfig.Logger.LogsDirPath, appConfig.Logger.LogtailToken)
	app := models.NewApp(assets, logger, appConfig)

	frontendLogger := models.NewFrontendLogger(logger)
	pdfService := models.NewPdfService(assets, logger, appConfig)

	// Create application with options
	startErr := wails.Run(&options.App{
		Title:            "We Love PDF",
		Width:            700,
		Height:           777,
		MinWidth:         700,
		MinHeight:        777,
		BackgroundColour: &options.RGBA{R: 42, G: 47, B: 38, A: 1},
		OnStartup:        app.Startup,
		OnDomReady:       app.Ready,
		OnBeforeClose:    app.BeforeClose,
		OnShutdown:       app.Shutdown,
		AssetServer: &assetserver.Options{
			Assets: frontendAssets,
		},
		Bind: []interface{}{
			app,
			pdfService,
			frontendLogger,
		},
		Mac: &mac.Options{
			About: &mac.AboutInfo{
				Title:   "We   ❤   PDF",
				Message: "Homepilot by iad @ 2023",
				Icon:    app.LogoIcon,
			},
		},
	})

	if startErr != nil {
		logger.Error("Error starting app", slog.String("reason", startErr.Error()))
		panic(startErr)
	}
}
