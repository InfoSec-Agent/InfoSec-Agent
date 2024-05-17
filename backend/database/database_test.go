package database_test

import (
	"database/sql"
	"os"
	"testing"

	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/database"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/logger"
	"github.com/stretchr/testify/require"
)

// TestMain sets up the necessary environment for the database package tests and executes them.
//
// This function initializes the logger for the tests and runs the tests.
//
// Parameters:
//   - m *testing.M: The testing framework that manages and runs the tests.
//
// Returns: None. The function calls os.Exit with the exit code returned by m.Run().
func TestMain(m *testing.M) {
	logger.SetupTests()

	// Run tests
	exitCode := m.Run()
	os.Exit(exitCode)
}

// TestGetSeverity tests the GetSeverity function to ensure it returns the correct severity level for a given issue ID and result ID pair.
//
// This test function creates a new SQLite database connection and calls the GetSeverity function with known issue IDs and result IDs.
// It then asserts that the returned severity level matches the expected value for each issue ID and result ID pair.
//
// Parameters:
//   - t *testing.T: The testing framework used for assertions.
//
// No return values.
func TestGetSeverity(t *testing.T) {
	// logger.SetupTests()

	// Arrange database connection
	db, err := sql.Open("sqlite", "../../reporting-page/database.db")
	if err != nil {
		t.Errorf("Error occurred: %v", err)
	}

	// Test for valid issue ID and result ID
	severity, err := database.GetSeverity(db, 1, 1)
	require.NoError(t, err)
	require.Equal(t, 4, severity)

	// Test for invalid issue ID and result ID
	_, err = database.GetSeverity(db, 0, 0)
	require.Error(t, err)
	require.Equal(t, sql.ErrNoRows.Error(), err.Error())
}

// TestGetJSONKey tests the GetJSONKey function to ensure it returns the correct JSON key for a given issue ID and result ID pair.
//
// This test function creates a new SQLite database connection and calls the GetJSONKey function with known issue IDs and result IDs.
// It then asserts that the returned JSON key matches the expected value for each issue ID and result ID pair.
//
// Parameters:
//   - t *testing.T: The testing framework used for assertions.
//
// No return values.
func TestGetJSONKey(t *testing.T) {
	// logger.SetupTests()

	// Arrange database connection
	db, err := sql.Open("sqlite", "../../reporting-page/database.db")
	if err != nil {
		t.Errorf("Error occurred: %v", err)
	}

	// Test for valid issue ID and result ID
	jsonKey, err := database.GetJSONKey(db, 1, 1)
	require.NoError(t, err)
	require.Equal(t, 11, jsonKey)

	// Test for invalid issue ID and result ID
	_, err = database.GetSeverity(db, 0, 0)
	require.Error(t, err)
	require.Equal(t, sql.ErrNoRows.Error(), err.Error())
}

// TestGetData tests the GetData function to ensure it returns the correct database data for a given list of checks.
//
// This test function creates a list of check results and calls the GetDataBaseData function.
// It then asserts that the returned database data matches the expected data for the given checks.
//
// Parameters:
//   - t *testing.T: The testing framework used for assertions.
//
// No return values.
func TestGetData(t *testing.T) {
	// logger.SetupTests()

	scanResult := []checks.Check{
		{
			IssueID:  1,
			ResultID: 1,
			Result:   []string{"Issue 1"},
			Error:    nil,
			ErrorMSG: "",
		},
	}
	expectedData := []database.Data{
		{
			CheckID:  1,
			Severity: 4,
			JSONKey:  11,
		},
	}
	emptyScanResult := []checks.Check{}
	emptyExpectedData := []database.Data{}
	invalidScanResult := []checks.Check{
		{
			IssueID:  0,
			ResultID: 0,
			Result:   []string{"Issue 0"},
			Error:    nil,
			ErrorMSG: "",
		},
	}
	invalidExpectedData := []database.Data{
		{
			CheckID:  0,
			Severity: 0,
			JSONKey:  0,
		},
	}
	wrongPathExpectedData := []database.Data{
		{
			CheckID:  1,
			Severity: 0,
			JSONKey:  0,
		},
	}
	testCases := []struct {
		scanResult   []checks.Check
		expectedData []database.Data
	}{
		{scanResult, expectedData},
		{emptyScanResult, emptyExpectedData},
		{invalidScanResult, invalidExpectedData},
	}

	for _, tc := range testCases {
		data, err := database.GetData(tc.scanResult, "../../reporting-page/database.db")
		if err != nil {
			t.Errorf("Error occurred: %v", err)
		}
		require.Equal(t, tc.expectedData, data)
		require.Equal(t, tc.expectedData, data)
	}

	// Test for invalid database path
	result, _ := database.GetData(scanResult, "")
	require.Equal(t, wrongPathExpectedData, result)
}