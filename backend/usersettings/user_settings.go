// Package usersettings contains functions for loading and saving user settings in the users AppData directory.
//
// Exported function(s): LoadUserSettings, SaveUserSettingsGetter.SaveUserSettings
//
// Exported type(s): UserSettings, SaveUserSettingsGetter
package usersettings

import (
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/logger"

	"encoding/json"
	"errors"
	"os"
	"time"
)

// UserSettings represents the settings for a user in the system.
//
// Fields:
//   - Language (int): An integer representing the user's preferred language. The specific language each integer represents can vary based on the system's language settings.
//   - ScanInterval (int): An integer representing the interval (in hours) at which the system should perform scans.
//   - Integration (bool): A boolean indicating whether the user has enabled integration with other systems or services.
//   - NextScan (time.Time): A time.Time value indicating when the next scan should occur.
//   - Points (int): An integer representing the user's current points amount.
//   - PointsHistory ([]int): A slice of integers representing the user's points history for each scan.
//   - TimeStamps ([]time.Time): A slice of time.Time values representing the time stamps for each scan.
//   - LighthouseState: An integer representing the gamification lighthouse state.
type UserSettings struct {
	Language         int         `json:"Language"`         // User's preferred language
	ScanInterval     int         `json:"ScanInterval"`     // Interval for system scans (in hours)
	Integration      bool        `json:"Integration"`      // Integration status with other systems
	NextScan         time.Time   `json:"NextScan"`         // Time for the next system scan
	Points           int         `json:"Points"`           // Current points amount
	PointsHistory    []int       `json:"PointsHistory"`    // Points history for each scan
	TimeStamps       []time.Time `json:"TimeStamps"`       // Time stamps for each scan
	LighthouseState  int         `json:"LighthouseState"`  // User's game state
	ProgressBarState int         `json:"ProgressBarState"` // User's progress bar state
	IntegrationKey   string      `json:"IntegrationKey"`   // Integration key for external systems
}

var DefaultUserSettings = UserSettings{Language: 1, ScanInterval: 7, Integration: false, NextScan: time.Now().Add((time.Hour * 24) * 7), Points: 0, PointsHistory: nil, TimeStamps: nil, LighthouseState: 0, IntegrationKey: ""}

// LoadUserSettings loads the user settings from a JSON file in the Windows AppData folder.
//
// The function uses the os package to get the path to the AppData folder, and reads the user settings from a file named "user_settings.json" in this folder.
// If there is an error while getting the AppData folder path, reading the JSON data from the file, or unmarshalling the JSON data to a UserSettings struct,
//
// Parameters: None
//
// Returns:
//   - UserSettings: The loaded user settings. This is a UserSettings struct.
func LoadUserSettings() UserSettings {
	logger.Log.Debug("Loading user settings")

	logger.Log.Trace("Getting user config directory")
	appDataPath, err := os.UserConfigDir()
	if err != nil {
		logger.Log.Warning("Error getting user config directory, using default settings")
		return DefaultUserSettings
	}
	dirPath := appDataPath + `\InfoSec-Agent`
	logger.Log.Trace("Creating/reading directory at:" + dirPath)
	err = os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		logger.Log.Warning("Error creating directory, using default settings")
		return DefaultUserSettings
	}

	filePath := dirPath + `\user_settings.json`

	logger.Log.Trace("Reading user settings from file:" + filePath)
	data, err := os.ReadFile(filePath)
	if err != nil {
		logger.Log.Warning("Error reading user settings file, using default settings")
		return DefaultUserSettings
	}

	var settings UserSettings
	logger.Log.Trace("Unmarshalling user settings JSON")
	err = json.Unmarshal(data, &settings)
	if err != nil {
		logger.Log.Warning("Error unmarshalling user settings JSON, using default settings")
		return DefaultUserSettings
	}
	logger.Log.Debug("Loaded user settings")
	return settings
}

// SaveUserSettingsGetter is an interface that defines a method for saving user settings.
//
// The SaveUserSettings method takes a UserSettings struct as input, which contains the user settings to be saved.
// It returns an error if any occurred while saving the user settings. If no error occurred, the method returns nil.
//
// This interface is implemented by any type that needs to save user settings for the system.
type SaveUserSettingsGetter interface {
	SaveUserSettings(settings UserSettings) error
}

// RealSaveUserSettingsGetter is a struct that implements the SaveUserSettingsGetter interface.
//
// It provides a real-world implementation of the SaveUserSettings method, which saves the user settings to a JSON file in the Windows AppData\Roaming folder.
type RealSaveUserSettingsGetter struct{}

// SaveUserSettings saves the user settings to a JSON file in the Windows AppData\Roaming folder.
//
// The function takes a UserSettings struct as input, which contains the user settings to be saved.
// It uses the os package to get the path to the AppData folder, and saves the user settings to a file named "user_settings.json" in this folder.
//
// Parameters:
//   - settings (UserSettings): The user settings to be saved.
//
// Returns:
//   - An error if any occurred while saving the user settings. If no error occurred, the function returns nil.
func (r RealSaveUserSettingsGetter) SaveUserSettings(settings UserSettings) error {
	logger.Log.Debug("Saving user settings")

	logger.Log.Trace("Getting user config directory")
	appDataPath, err := os.UserConfigDir()
	if err != nil {
		err = errors.New("Error getting user config directory: " + err.Error())
		logger.Log.Error(err.Error())
		return err
	}

	dirPath := appDataPath + `\InfoSec-Agent`
	logger.Log.Trace("Creating/reading directory at:" + dirPath)
	err = os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		err = errors.New("Error creating directory: " + err.Error())
		logger.Log.Error(err.Error())
		return err
	}

	logger.Log.Trace("Marshalling data to JSON")
	file, err := json.MarshalIndent(settings, "", " ")
	if err != nil {
		err = errors.New("Error marshalling user settings JSON: " + err.Error())
		logger.Log.Error(err.Error())
		return err
	}

	filePath := dirPath + `\user_settings.json`
	logger.Log.Trace("Writing user settings to file:" + filePath)
	err = os.WriteFile(filePath, file, 0600)
	if err != nil {
		err = errors.New("Error writing user settings to file: " + err.Error())
		logger.Log.Error(err.Error())
		return err
	}

	logger.Log.Debug("User settings saved successfully")
	return nil
}
