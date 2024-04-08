// Package (reporting page) main contains the entry point of the reporting page application
//
// Exported function(s): NewApp, NewTray
//
// Exported struct(s): App
package main

import (
	"embed"
	"github.com/InfoSec-Agent/InfoSec-Agent/localization"
	"github.com/InfoSec-Agent/InfoSec-Agent/logger"
	"github.com/InfoSec-Agent/InfoSec-Agent/tray"
	"github.com/InfoSec-Agent/InfoSec-Agent/usersettings"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

//go:embed all:frontend/dist
var assets embed.FS

// main is the entry point of the reporting page program, starts the Wails application
//
// Parameters: _
//
// Returns: _
func main() {
	logger.Setup()
	logger.Log.Info("Reporting page starting")

	// Create a new instance of the app and tray struct
	app := NewApp()
	systemTray := NewTray(logger.Log)
	database := NewDataBase()
	customLogger := logger.Log
	localization.Init("../")
	lang := usersettings.LoadUserSettings("../usersettings").Language
	tray.Language = lang

	// Create a Wails application with the specified options
	err := wails.Run(&options.App{
		Title:       "reporting-page",
		Width:       1024,
		Height:      768,
		StartHidden: true,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
			systemTray,
			database,
		},
		Logger: customLogger,
		Windows: &windows.Options{
			Theme: windows.SystemDefault,
			CustomTheme: &windows.ThemeSettings{
				DarkModeTitleBar:   windows.RGB(20, 20, 20),
				DarkModeTitleText:  windows.RGB(200, 200, 200),
				DarkModeBorder:     windows.RGB(20, 0, 20),
				LightModeTitleBar:  windows.RGB(200, 200, 200),
				LightModeTitleText: windows.RGB(20, 20, 20),
				LightModeBorder:    windows.RGB(200, 200, 200),
			},
		},
	})
	if err != nil {
		logger.Log.ErrorWithErr("Error creating Wails application:", err)
	}
}
