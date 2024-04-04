// Package firefox is responsible for running checks on Firefox.
//
// Exported function(s): CookieFirefox, ExtensionFirefox, HistoryFirefox, PasswordFirefox
package firefox

import (
	"encoding/json"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/InfoSec-Agent/InfoSec-Agent/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/utils"

	"github.com/andrewarchi/browser/firefox"
)

// ExtensionFirefox checks the extensions in the Firefox browser, and specifically for an adblocker.
//
// Parameters: _
//
// Returns: A list of found extensions, and if an adblocker is installed
func ExtensionFirefox() (checks.Check, checks.Check) {
	var resultID int
	// Determine the directory in which the Firefox profile is stored
	ffdirectory, err := utils.FirefoxFolder()
	if err != nil {
		return checks.NewCheckErrorf(27, "No firefox directory found", err),
			checks.NewCheckErrorf(28, "No firefox directory found", err)
	}

	addBlocker := false // Variable used for checking if an adblocker is used
	var output []string
	// Open the extensions.json file, which contains a list of all installed Firefox extensions
	content, err := os.Open(ffdirectory[0] + "\\extensions.json")
	if err != nil {
		return checks.NewCheckError(27, err), checks.NewCheckError(28, err)
	}
	defer func(content *os.File) {
		err = content.Close()
		if err != nil {
			log.Println("error closing file: ", err)
		}
	}(content)

	// Create a struct for the JSON file
	var extensions firefox.Extensions
	decoder := json.NewDecoder(content)
	err = decoder.Decode(&extensions)
	if err != nil {
		return checks.NewCheckError(27, err), checks.NewCheckError(28, err)
	}

	// In the result list, add: the name of the addon, type of the addon, the creator, and whether it is active or not
	for _, addon := range extensions.Addons {
		output = append(output, addon.DefaultLocale.Name+","+addon.Type+","+addon.DefaultLocale.Creator+","+
			""+strconv.FormatBool(addon.Active))
		// Determine if the addon is an adblocker
		if adblockerFirefox(addon.DefaultLocale.Name) {
			addBlocker = true
			resultID += 1
		}
	}
	adBlockused := strconv.FormatBool(addBlocker)
	return checks.NewCheckResult(27, 0, output...),
		checks.NewCheckResult(28, resultID, adBlockused)
}

// adblockerFirefox checks if the given extension is an adblocker
//
// Parameters: extensionName (string) - The name of the extension to check
//
// Returns: If the extension is an adblocker (bool)
func adblockerFirefox(extensionName string) bool {
	// List of known/popular adblockers to match against
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
