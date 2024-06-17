// Package gamification handles the gamification within the application, to reward users for performing security checks and staying secure.
//
// Exported function(s): UpdateGameState, PointCalculation, LighthouseStateTransition
package gamification

import (
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/database"
	"time"

	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/logger"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/usersettings"
)

// GameState is a struct that represents the state of the gamification.
// This consists of the user's points, a history of all previous points, and a lighthouse state.
type GameState struct {
	Points          int
	PointsHistory   []int
	TimeStamps      []time.Time
	LighthouseState int
}

// UpdateGameState updates the game state based on the scan results and the current game state.
//
// Parameters:
//   - scanResults ([]checks.Check): The results of the scans.
//   - databasePath (string): The path to the database file.
//   - getter (PointCalculationGetter): An object that implements the PointCalculationGetter interface.
//   - userGetter (usersettings.SaveUserSettingsGetter): An object that implements the SaveUserSettingsGetter interface.
//
// Returns: The updated game state with the new points amount and new lighthouse state.
func UpdateGameState(scanResults []checks.Check, databasePath string, getter PointCalculationGetter, userGetter usersettings.SaveUserSettingsGetter) (GameState, error) {
	gs := GameState{Points: 0, PointsHistory: nil, TimeStamps: nil, LighthouseState: 0}

	// Loading the game state from the user settings and putting it in the game state struct
	userSettings := usersettings.LoadUserSettings()
	gs.Points = userSettings.Points
	gs.PointsHistory = userSettings.PointsHistory
	gs.TimeStamps = userSettings.TimeStamps
	gs.LighthouseState = userSettings.LighthouseState

	gs, err := getter.PointCalculation(gs, scanResults, databasePath)
	if err != nil {
		logger.Log.ErrorWithErr("Error calculating points:", err)
		return gs, err
	}
	gs = LighthouseStateTransition(gs)

	// Saving the game state in the user settings
	current := usersettings.LoadUserSettings()
	current.Points = gs.Points
	current.PointsHistory = gs.PointsHistory
	current.TimeStamps = gs.TimeStamps
	current.LighthouseState = gs.LighthouseState
	err = userGetter.SaveUserSettings(current)
	if err != nil {
		logger.Log.Warning("Gamification settings not saved to file")
	}
	return gs, nil
}

// PointCalculationGetter is an interface that defines a method for calculating points
// based on the game state, scan results, and a file path.
//
// The PointCalculation method takes a GameState struct, a slice of checks.Check, and a string
// representing a file path. It returns an updated GameState and an error.
//
// This interface is implemented by any type that needs to calculate points for the gamification system.
type PointCalculationGetter interface {
	PointCalculation(gs GameState, scanResults []checks.Check, filePath string) (GameState, error)
}

// RealPointCalculationGetter is a struct that implements the PointCalculationGetter interface.
//
// It provides a real-world implementation of the PointCalculation method, which calculates the number of points
// for the user based on the check results.
type RealPointCalculationGetter struct{}

// PointCalculation calculates the number of points for the user based on the check results.
//
// Parameters:
//   - gs (GameState): The current game state, which includes the user's points and lighthouse state.
//   - scanResults ([]checks.Check): The results of the scans.
//   - databasePath (string): The path to the database file.
//
// Returns:
//   - GameState: The updated game state with the new points amount.
func (r RealPointCalculationGetter) PointCalculation(gs GameState, scanResults []checks.Check, jsonFilePath string) (GameState, error) {
	gs.Points = 0

	dataList, err := database.GetData(jsonFilePath, scanResults)
	if err != nil {
		logger.Log.ErrorWithErr("Error getting data from database:", err)
		return gs, err
	}
	for _, data := range dataList {
		sev := data.Severity
		if sev >= 0 && sev < 4 {
			gs.Points += sev
		}
	}
	gs.PointsHistory = append(gs.PointsHistory, gs.Points)
	gs.TimeStamps = append(gs.TimeStamps, time.Now())

	return gs, nil
}

// LighthouseStateTransition determines the lighthouse state based on the user's points (the less points, the better)
//
// Parameters:
//   - gs (GameState): The current game state, which includes the user's points and lighthouse state.
//
// Returns:
//   - GameState: The updated game state with the new lighthouse state.
func LighthouseStateTransition(gs GameState) GameState {
	switch {
	case gs.Points < 10 && SufficientActivity(gs):
		gs.LighthouseState = 5 // The best state
	case gs.Points < 20 && SufficientActivity(gs):
		gs.LighthouseState = 4
	case gs.Points < 30 && SufficientActivity(gs):
		gs.LighthouseState = 3
	case gs.Points < 40 && SufficientActivity(gs):
		gs.LighthouseState = 2
	case gs.Points < 50 && SufficientActivity(gs):
		gs.LighthouseState = 1
	default:
		gs.LighthouseState = 1
	}
	return gs
}

// sufficientActivity checks if the user has been active enough to transition to another lighthouse state
//
// Parameters:
//   - gs (GameState): The game state of the user.
//   - duration (time.Duration): The duration of time that the user needs to be active.
//
// Returns:
//   - bool: whether the user has been active enough.
func SufficientActivity(gs GameState) bool {
	// The duration threshold of which the user has been active
	// Note that we define active as having performed a security check more than [1 week ago]
	requiredDuration := 7 * 24 * time.Hour

	if len(gs.TimeStamps) == 0 {
		return false
	}
	oldestRecord := gs.TimeStamps[0] // The oldest record is the first timestamp made

	return time.Since(oldestRecord) > requiredDuration
}
