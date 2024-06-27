// Package gamification handles the gamification within the application, to reward users for performing security checks and staying secure.
//
// Exported function(s): UpdateGameState, PointCalculationGetter.PointCalculation, LighthouseStateTransition, SufficientActivity
//
// Exported types(s): GameState, PointCalculationGetter
package gamification

import (
	"strconv"
	"time"

	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/database"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/logger"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/usersettings"
)

// GameState is a struct that represents the state of the gamification.
// This consists of the user's points, a history of all previous points, and a lighthouse state.
//
// Fields:
//   - Points (int): The number of points the user has.
//   - PointsHistory ([]int): A list of all previous points the user has had.
//   - TimeStamps ([]time.Time): A list of timestamps for when the user has received points.
//   - LighthouseState (int): The current state of the lighthouse.
type GameState struct {
	Points           int
	PointsHistory    []int
	TimeStamps       []time.Time
	LighthouseState  int
	ProgressBarState int
}

// UpdateGameState updates the game state based on the scan results and the current game state.
//
// Parameters:
//   - scanResults ([]checks.Check): The results of the scans.
//   - databasePath (string): The path to the database file.
//   - getter (PointCalculationGetter): An object that implements the PointCalculationGetter interface.
//   - userGetter (usersettings.SaveUserSettingsGetter): An object that implements the SaveUserSettingsGetter interface.
//
// Returns:
//   - The updated game state with the new points amount and new lighthouse state.
func UpdateGameState(scanResults []checks.Check, databasePath string, getter PointCalculationGetter, userGetter usersettings.SaveUserSettingsGetter) (GameState, error) {
	logger.Log.Trace("Updating game state")
	gs := GameState{Points: 0, PointsHistory: nil, TimeStamps: nil, LighthouseState: 0}

	// Loading the game state from the user settings and putting it in the game state struct
	userSettings := usersettings.LoadUserSettings()
	gs.Points = userSettings.Points
	gs.PointsHistory = userSettings.PointsHistory
	gs.TimeStamps = userSettings.TimeStamps
	gs.LighthouseState = userSettings.LighthouseState
	gs.ProgressBarState = userSettings.ProgressBarState

	gs, err := getter.PointCalculation(gs, scanResults, databasePath)
	if err != nil {
		logger.Log.ErrorWithErr("Error calculating points", err)
		return gs, err
	}
	gs = LighthouseStateTransition(gs)

	// Saving the game state in the user settings
	current := usersettings.LoadUserSettings()
	current.Points = gs.Points
	current.PointsHistory = gs.PointsHistory
	current.TimeStamps = gs.TimeStamps
	current.LighthouseState = gs.LighthouseState
	current.ProgressBarState = gs.ProgressBarState
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
	logger.Log.Trace("Calculating gamification points ")
	newPoints := 0

	dataList, err := database.GetData(jsonFilePath, scanResults)
	if err != nil {
		logger.Log.ErrorWithErr("Error getting data from database", err)
		return gs, err
	}

	for _, data := range dataList {
		sev := data.Severity
		if sev >= 0 && sev < 4 {
			newPoints += sev
		}
	}
	logger.Log.Trace("Calculated gamification points: " + strconv.Itoa(newPoints))
	gs.Points = newPoints
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
	modResult := gs.Points % 10
	gs.ProgressBarState = 100 - (modResult * 10)
	switch {
	case gs.Points <= 10 && SufficientActivity(gs):
		gs.LighthouseState = 4 // The best state
		gs.ProgressBarState = 100
	case (gs.Points <= 20 && SufficientActivity(gs)) || gs.Points <= 10:
		gs.LighthouseState = 3
		if !SufficientActivity(gs) {
			gs.ProgressBarState = 99
		}
	case (gs.Points <= 30 && SufficientActivity(gs)) || gs.Points <= 20:
		gs.LighthouseState = 2
	case (gs.Points <= 40 && SufficientActivity(gs)) || gs.Points <= 30:
		gs.LighthouseState = 1
	default:
		gs.LighthouseState = 0
	}
	logger.Log.Trace("Calculated lighthouse state: " + strconv.Itoa(gs.LighthouseState))
	return gs
}

// SufficientActivity checks if the user has been active enough to transition to another lighthouse state
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

	// The oldest record is the first timestamp made
	oldestRecord := gs.TimeStamps[0]
	return time.Since(oldestRecord) > requiredDuration
}
