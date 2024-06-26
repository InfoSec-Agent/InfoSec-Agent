// Package database provides functionality for interacting with the JSON database.
//
// Exported function(s): GetData
//
// Exported type(s): JSONData, Issue, Result, Data
package database

import (
	"encoding/json"
	"os"
	"strconv"

	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/logger"
)

// JSONData represents the structure of the JSON data
type JSONData map[string]Issue

// Issue represents the structure of an issue in the JSON data.
//
// Fields:
//   - Type (string): The type of the issue.
//   - Information (string): Additional information about the issue.
//   - Results (map[string]Result): A map of results for the issue, where the key is the result ID and the value is a Result object.
type Issue struct {
	Type        string            `json:"Type"`
	Information string            `json:"Information"`
	Results     map[string]Result `json:"Results"`
}

// Result represents the structure of a result in the JSON data.
//
// Fields:
//   - Severity (int): The severity of the result.
//   - Name (string): The name of the result.
//   - Solution ([]string): A list of solutions for the result.
//   - Screenshots ([]string): A list of screenshots for the result.
type Result struct {
	Severity    int      `json:"Severity"`
	Name        string   `json:"Name"`
	Solution    []string `json:"Solution"`
	Screenshots []string `json:"Screenshots"`
}

// Data represents a simplified structure with issueID and severity
//
// Fields:
//   - IssueID (int): The unique identifier for the issue.
//   - Severity (int): The severity of the issue.
type Data struct {
	IssueID  int
	Severity int
}

// GetData is a function that computes the severity for all found checks and puts them into a list of Data.
//
// Parameters:
//   - jsonFilePath (string): The path to the JSON file containing the severity data.
//   - checkResults ([]checks.Check): A list of check results.
//
// Returns:
//   - []Data: A list of Data objects containing the issueID and severity for each check result.
//   - error: An error object if an error occurred during the process.
//
// The function first opens the JSON file and decodes the JSON data into a JSONData object.
// If there's an error opening the file or decoding the JSON data, it returns nil and the error.
// The function then initializes a slice of Data and iterates through all check results to compute the severities.
// For each check result, it retrieves the issueID and resultID, and calls the GetSeverity function to compute the severity.
// If there's an error retrieving the severity, it logs the error and skips to the next check result.
// If the severity is successfully computed, it creates a new Data object with the issueID and severity, and appends it to the dataList.
// After iterating through all check results, the function returns the dataList and nil for the error.
func GetData(jsonFilePath string, checkResults []checks.Check) ([]Data, error) {
	logger.Log.Trace("Getting data from JSON file " + jsonFilePath)
	byteValue, err := os.ReadFile(jsonFilePath)
	if err != nil {
		logger.Log.ErrorWithErr("Error reading JSON file "+jsonFilePath, err)
		return nil, err
	}

	var data map[string]map[string]interface{}
	err = json.Unmarshal(byteValue, &data)
	if err != nil {
		logger.Log.ErrorWithErr("Error parsing JSON", err)
		return nil, err
	}

	// Initialize the list of Data
	var dataList []Data

	// Iterate through all check results to compute the severities
	for _, checkResult := range checkResults {
		if checkResult.Error != nil {
			logger.Log.Debug("Error reading scan result for IssueID " + strconv.Itoa(checkResult.IssueID) + ": " + checkResult.ErrorMSG)
			continue
		}

		// Convert IssueID and ResultID to string to access JSON data
		issueKey := strconv.Itoa(checkResult.IssueID)
		resultKey := strconv.Itoa(checkResult.ResultID)

		// Get the severity from JSON
		issueData, ok := data[issueKey]
		if !ok {
			logger.Log.Debug("IssueID " + strconv.Itoa(checkResult.IssueID) + " not found in JSON")
			continue
		}
		resultData, ok := issueData[resultKey].(map[string]interface{})
		if !ok {
			logger.Log.Debug("ResultID not found in JSON: " + strconv.Itoa(checkResult.ResultID))
			continue
		}
		sev, ok := resultData["Severity"].(float64)
		if !ok {
			logger.Log.Debug("Severity not found or invalid for IssueID " + strconv.Itoa(checkResult.IssueID) + "ResultID:" + strconv.Itoa(checkResult.ResultID))
			continue
		}
		dataList = append(dataList, Data{IssueID: checkResult.IssueID, Severity: int(sev)})
	}

	return dataList, nil
}
