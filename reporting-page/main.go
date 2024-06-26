// The InfoSec Agent Reporting Page is a Wails application that presents the results of the InfoSec Agent scans in a user-friendly format.
// It serves as a frontend for the [InfoSec-Agent] package.
//
// This package is responsible for initializing the application and creating the link between the Go backend and JavaScript frontend.
//
// [InfoSec-Agent]: https://github.com/InfoSec-Agent/InfoSec-Agent/
package main

import (
	"embed"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/InfoSec-Agent/InfoSec-Agent/backend/config"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/localization"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/logger"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/tray"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/usersettings"
	"github.com/rodolfoag/gow32"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
	"golang.org/x/sys/windows/registry"
)

//go:embed all:frontend/dist
var assets embed.FS

// FileLoader is a struct that implements the http.Handler interface to serve files from the frontend/src/assets/images directory.
//
// Fields
//   - Handler (http.Handler): The handler that serves the requested file.
type FileLoader struct {
	http.Handler
}

// NewFileLoader creates a new instance of the FileLoader struct.
func NewFileLoader() *FileLoader {
	return &FileLoader{}
}

// ServeHTTP serves the requested file from the images directory.
// It first constructs the full path to the requested file and checks if the file is within the allowed directory.
//
// This function reads the file data and writes it to the response writer.
// If an error occurs during the file reading or writing process, it logs the error and returns an appropriate HTTP error response.
//
// Parameters:
//   - res (http.ResponseWriter): The response writer to write the file data to.
//   - req (*http.Request): The request object containing information about the requested file.
//
// Returns: None. This function does not return a value as it serves the requested file.
func (h *FileLoader) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	requestedPath := req.URL.Path
	cleanPath := filepath.Clean(requestedPath) // Clean the path to avoid directory traversal
	newPath := strings.Replace(cleanPath, `\`, ``, 1)

	// Ensure the requested path is relative and does not try to traverse directories
	if cleanPath == "." || strings.Contains(cleanPath, "..") {
		logger.Log.Error("Invalid file path: " + requestedPath)
		http.Error(res, "Invalid file path", http.StatusBadRequest)
		return
	}

	// Construct the full path to the file
	fullPath := filepath.Join(config.ReportingPageImageDir, cleanPath)

	// Check if the file is within the allowed directory
	if !strings.HasPrefix(fullPath, filepath.Clean(config.ReportingPageImageDir)+string(os.PathSeparator)) {
		logger.Log.Error("Access to the file path denied:" + newPath)
		http.Error(res, "Access denied", http.StatusForbidden)
		return
	}

	fileData, err := os.ReadFile(newPath)
	if err != nil {
		logger.Log.ErrorWithErr("Could not load file: "+newPath, err)
		http.Error(res, "File not found", http.StatusNotFound)
		return
	}

	if _, err = res.Write(fileData); err != nil {
		logger.Log.ErrorWithErr("Could not write file: "+newPath, err)
		http.Error(res, "Failed to serve file", http.StatusInternalServerError)
	}
}

// main is the entry point of the reporting page application. It initializes the localization settings, creates a new instance of the application, tray, and database, and starts the Wails application.
//
// This function first calls the Init function from the localization package to set up the localization settings for the application. It then creates new instances of the App, Tray, and DataBase structs.
// After setting up these instances, it creates a Wails application with specific options including the title, dimensions, startup behavior, asset server, background color, startup function, and bound interfaces.
// It also sets up the Windows-specific options for the Wails application, including the theme and custom theme settings.
// If an error occurs during the creation or startup of the Wails application, it is logged and the program terminates.
//
// Parameters: None.
//
// Returns: None. This function does not return a value as it is the entry point of the application.
func main() {
	// Create a mutex to ensure only one instance of the application is running
	// If the mutex already exists, it means another instance of the application is running, so we exit
	// This also ensures program is not running when uninstalling the application
	_, mutexErr := gow32.CreateMutex("InfoSec-Agent-Reporting-Page")
	if mutexErr != nil {
		return
	}

	// Setup log file
	logger.Setup("reporting-page-log.txt", config.LogLevel, config.LogLevelSpecific)
	logger.Log.Info("Reporting page starting")

	// Change directory to the reporting page directory.
	// When opening the reporting page from the Windows notification we start in C:\Windows\System32, and we cannot run the reporting page from there.
	var path string
	k, err := registry.OpenKey(
		registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\App Paths\InfoSec-Agent-Reporting-Page.exe`,
		registry.QUERY_VALUE,
	)
	if err != nil {
		logger.Log.ErrorWithErr("Error opening registry key", err)
	} else {
		logger.Log.Debug("Found registry key")
		defer func(k registry.Key) {
			err = k.Close()
			if err != nil {
				logger.Log.ErrorWithErr("Error closing registry key", err)
			}
		}(k)

		// Get reporting page executable path
		path, _, err = k.GetStringValue("Path")
		if err != nil {
			logger.Log.ErrorWithErr("Error getting path string", err)
		}
	}
	changeDirectory(path)

	// Create a new instance of the app and tray struct
	app := NewApp()
	systemTray := NewTray(logger.Log)
	customLogger := logger.Log
	localization.Init("../")
	lang := usersettings.LoadUserSettings().Language
	tray.Language = lang

	logger.Log.Debug("Starting wails application")
	// Create a Wails application with the specified options
	err = wails.Run(&options.App{
		Title:       "InfoSec Agent Reporting Page",
		Width:       1024,
		Height:      768,
		StartHidden: false,
		AlwaysOnTop: false,
		AssetServer: &assetserver.Options{
			Assets:  assets,
			Handler: NewFileLoader(),
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
			systemTray,
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
		logger.Log.ErrorWithErr("Error creating Wails application", err)
	}
}

// changeDirectory attempts to change the current working directory to the specified path.
// If the config ReportingPagePath contains "build/bin", it sets the path to "reporting-page" to run the reporting page from the development environment.
// If changing the directory is successful, a debug message is logged indicating the new directory.
// If an error occurs during the directory change, the error is logged with an error level.
//
// Parameters:
//   - path (string): The path to change the current working directory to.
//
// Returns: None.
func changeDirectory(path string) {
	// If the config ReportingPagePath contains "build/bin", we are running the reporting page from the development environment
	if strings.Contains(config.ReportingPagePath, "build/bin") {
		path = "reporting-page"
	}

	err := os.Chdir(path)
	if err != nil {
		logger.Log.ErrorWithErr("Error changing directory", err)
	} else {
		logger.Log.Debug("Changed directory to " + path)
	}
}
