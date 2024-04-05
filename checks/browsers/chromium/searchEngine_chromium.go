package chromium

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"

	"github.com/InfoSec-Agent/InfoSec-Agent/filemock"

	"github.com/InfoSec-Agent/InfoSec-Agent/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/utils"
)

// SearchEngineChromium inspects the default search engine setting in Chromium-based browsers.
//
// Parameters:
//   - browser: A string representing the name of the Chromium-based browser to inspect. This could be "Chrome", "Edge", etc.
//
// Returns:
//   - A checks.Check object representing the result of the check. The result contains the name of the default search engine used in the specified browser. If an error occurs during the check, the result will contain a description of the error.
//
// This function works by locating the preferences file in the user's home directory, which contains the browser's settings. It opens and reads this file, then parses it as JSON to access the settings. It specifically looks for the "default_search_provider_data" key in the JSON data, which holds the name of the default search engine. If this key is found, its value is returned as the result of the check. If any error occurs during this process, such as an error reading the file or parsing the JSON, this error is returned as the result of the check.
func SearchEngineChromium(browser string) checks.Check {
	var browserPath string
	var returnBrowserName string

	// Set the browser path and the return browser name based on the browser to check
	// Currently, supports checking of Google Chrome and Microsoft Edge
	if browser == chrome {
		returnBrowserName = "SearchEngineChrome"
		browserPath = chromePath
	}
	if browser == edge {
		returnBrowserName = "SearchEngineEdge"
		browserPath = edgePath
	}
	// Holds the return value and sets the default value to chrome in case you never changed your search engine
	defaultSE := "google.com"
	var user string
	var err error
	user, err = os.UserHomeDir()
	if err != nil {
		return checks.NewCheckErrorf(returnBrowserName, "Error: ", err)
	}

	// Get the current user's home directory, where the preferences can be found
	preferencesDir := filepath.Join(user, "AppData", "Local", browserPath, "User Data", "Default", "Preferences")
	// TODO: var needs to be filemock.File (this requires implementation of all other functions used for files in project)
	var file *os.File
	file, err = os.Open(filepath.Clean(preferencesDir))
	if err != nil {
		return checks.NewCheckErrorf(returnBrowserName, "Error: ", err)
	}
	defer func(file filemock.File) {
		err = utils.CloseFile(file)
		if err != nil {
			log.Println("Error closing file")
		}
	}(file)

	// Byte array holding the preferences json data used to unmarshal the data later
	var byteValue []byte
	byteValue, err = io.ReadAll(file)
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
