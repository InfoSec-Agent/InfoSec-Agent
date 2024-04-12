package firefox

import (
	"encoding/json"
	"os"
	"strconv"

	"github.com/InfoSec-Agent/InfoSec-Agent/logger"

	"github.com/InfoSec-Agent/InfoSec-Agent/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/utils"
	"github.com/andrewarchi/browser/firefox"
)

// PasswordFirefox checks the passwords in the Firefox browser.
//
// Parameters: _
//
// Returns:
func PasswordFirefox() checks.Check {
	// Determine the directory in which the Firefox profile is stored
	ffdirectory, err := utils.RealProfileFinder{}.FirefoxFolder()

	var output []string
	// Open the logins.json file, which contains a list of all saved Firefox passwords
	content, err := os.Open(ffdirectory[0] + "\\logins.json")
	if err != nil {
		return checks.NewCheckError(99, err)
	}
	defer func(content *os.File) {
		err = content.Close()
		if err != nil {
			logger.Log.ErrorWithErr("Error closing file: ", err)
		}
	}(content)

	// Creates a struct for the JSON file
	var extensions firefox.Extensions
	decoder := json.NewDecoder(content)
	err = decoder.Decode(&extensions)
	if err != nil {
		return checks.NewCheckError(99, err)
	}

	// TODO: Final functionality currently not implemented yet, should return an analysis on the used passwords
	// Below code is a placeholder
	for _, addon := range extensions.Addons {
		output = append(output,
			addon.DefaultLocale.Name+addon.Type+addon.DefaultLocale.Creator+strconv.FormatBool(addon.Active))
	}
	return checks.NewCheckResult(99, 0, output...)
}
