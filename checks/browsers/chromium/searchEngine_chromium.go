package chromium

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"

	"github.com/InfoSec-Agent/InfoSec-Agent/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/utils"
)

// SearchEngineChromium checks the standard search engine in chromium based browsers.
//
// Parameters: _
//
// Returns: The standard search engine for chromium based browsers
func SearchEngineChromium(browser string) checks.Check {
	var browserPath string
	var returnBrowserName string
	// Set the browser path and the return browser name based on the browser to check
	// Currently, supports checking of Google Chrome and Microsoft Edge
	if browser == "Chrome" {
		returnBrowserName = "SearchEngineChrome"
		browserPath = "Google/Chrome"
	}
	if browser == "Edge" {
		returnBrowserName = "SearchEngineEdge"
		browserPath = "Microsoft/Edge"
	}
	// Holds the return value and sets the default value to chrome in case you never changed your search engine
	defaultSE := "google.com"
	user, err := os.UserHomeDir()
	if err != nil {
		return checks.NewCheckErrorf(returnBrowserName, "Error: ", err)
	}

	// Get the current user's home directory, where the preferences can be found
	preferencesDir := filepath.Join(user, "AppData", "Local", browserPath, "User Data", "Default", "Preferences")
	file, err := os.Open(filepath.Clean(preferencesDir))
	if err != nil {
		return checks.NewCheckErrorf(returnBrowserName, "Error: ", err)
	}
	defer utils.CloseFile(file)

	// Byte array holding the preferences json data used to unmarshal the data later
	byteValue, err := io.ReadAll(file)
	if err != nil {
		return checks.NewCheckErrorf(returnBrowserName, " Can't read data,Error: ", err)
	}
	// Holds the unmarshaled data of the json for acces to the key value pairs
	var dev map[string]interface{}
	err = json.Unmarshal(byteValue, &dev)
	if err != nil {
		return checks.NewCheckErrorf(returnBrowserName, "Error: ", err)
	}

	// Iterate through the json dev map to look for our search engine key
	for key, value := range dev {
		if key == "default_search_provider_data" {
			text := fmt.Sprintf("%v", value)
			// Regex pattern to find the string keyword: and everything after that until we hit a space
			pattern := `keyword:\s*(\S+)`
			regex := regexp.MustCompile(pattern)
			matches := regex.FindString(text)
			if matches == "" {
				return checks.NewCheckErrorf(returnBrowserName, "Error: ", err)
			}
			// Removes the word keyword: from the result
			defaultSE = matches[8:]
		}
	}
	// Returns the default search engine used
	return checks.NewCheckResult(returnBrowserName, defaultSE)
}
