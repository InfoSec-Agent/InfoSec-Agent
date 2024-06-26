package firefox

import (
	"encoding/json"
	"os"
	"strconv"
	"strings"

	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks/browsers"

	"github.com/InfoSec-Agent/InfoSec-Agent/backend/logger"

	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks"
	"github.com/andrewarchi/browser/firefox"
)

// ExtensionFirefox inspects the extensions installed in the Firefox browser and checks for the presence of an ad blocker.
//
// Parameters:
//   - profileFinder: An object that implements the FirefoxProfileFinder interface. It is used to find the Firefox profile directory.
//
// Returns:
//   - A checks.Check object representing the result of the check. The result contains a list of installed extensions in the Firefox browser. Each extension is represented as a string that includes the name, type, creator, and its active status.
//   - A checks.Check object representing the result of the ad blocker check. The result is a boolean indicating whether an adblocker is installed or not.
//
// This function works by locating the Firefox profile directory and opening the extensions.json file, which contains a list of all installed Firefox extensions. It decodes the JSON file into a struct and iterates over the addons in the struct. For each addon, it appends a string to the result list that includes the name, type, creator, and active status of the addon. It also checks if the addon is an ad blocker by calling the adblockerFirefox function with the addon's name. If an adblocker is found, it sets a boolean variable to true. The function returns two checks.Check objects: one with the list of extensions and one with the boolean indicating the presence of an adblocker.
func ExtensionFirefox(profileFinder browsers.FirefoxProfileFinder) (checks.Check, checks.Check) {
	var resultID int
	// Determine the directory in which the Firefox profile is stored
	ffdirectory, err := profileFinder.FirefoxFolder()
	if err != nil {
		return checks.NewCheckErrorf(checks.ExtensionFirefoxID, "No firefox directory found", err),
			checks.NewCheckErrorf(checks.AdblockFirefoxID, "No firefox directory found", err)
	}

	resultID = 1
	var output []string
	// Open the extensions.json file, which contains a list of all installed Firefox extensions
	content, err := os.Open(ffdirectory[0] + "\\extensions.json")
	if err != nil {
		return checks.NewCheckError(checks.ExtensionFirefoxID, err), checks.NewCheckError(checks.AdblockFirefoxID, err)
	}
	defer func(content *os.File) {
		err = content.Close()
		if err != nil {
			logger.Log.ErrorWithErr("Error closing file", err)
		}
	}(content)

	// Create a struct for the JSON file
	var extensions firefox.Extensions
	decoder := json.NewDecoder(content)
	err = decoder.Decode(&extensions)
	if err != nil {
		return checks.NewCheckError(checks.ExtensionFirefoxID, err), checks.NewCheckError(checks.AdblockFirefoxID, err)
	}

	// In the result list, add: the name of the addon, type of the addon, the creator, and whether it is active or not
	for _, addon := range extensions.Addons {
		output = append(output, addon.DefaultLocale.Name+","+addon.Type+","+addon.DefaultLocale.Creator+","+
			""+strconv.FormatBool(addon.Active))
		// Determine if the addon is an ad blocker
		if AdblockerFirefox(addon.DefaultLocale.Name) {
			resultID = 0
		}
	}
	return checks.NewCheckResult(checks.ExtensionFirefoxID, 0, output...),
		checks.NewCheckResult(checks.AdblockFirefoxID, resultID)
}

// AdblockerFirefox determines whether the provided Firefox extension functions as an adblocker.
//
// Parameters:
//   - extensionName: A string representing the name of the Firefox extension to be evaluated.
//
// Returns:
//   - A boolean value indicating whether the provided extension is recognized as an ad blocker. The function returns true if the extension name matches any known adblocker names, and false otherwise.
//
// This function works by comparing the provided extension name to a list of known ad blocker names. The comparison is case-insensitive. If a match is found, the function returns true. If no match is found, the function returns false. This function is used in the context of the ExtensionFirefox function to identify whether any of the installed Firefox extensions function as ad blockers.
func AdblockerFirefox(extensionName string) bool {
	// List of known/popular ad blockers to match against
	adblockerNames := []string{
		"adblocker ultimate",
		"adguard adblocker",
		"adblocker for youtube",
		"ublock origin",
		"adblock plus",
		"adblock for firefox",
	}
	for _, adblockerName := range adblockerNames {
		if strings.Contains(strings.ToLower(extensionName), adblockerName) {
			return true
		}
	}
	return false
}
