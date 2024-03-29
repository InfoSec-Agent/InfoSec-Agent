package firefox

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/InfoSec-Agent/InfoSec-Agent/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/utils"

	_ "modernc.org/sqlite"
)

// CookieFirefox checks the cookies in the Firefox browser.
//
// Parameters: _
//
// Returns: A list of cookies in the Firefox browser
func CookieFirefox() checks.Check {
	var output []string
	// Determine the directory in which the Firefox profile is stored
	ffdirectory, err := utils.FirefoxFolder()
	if err != nil {
		return checks.NewCheckErrorf("CookieFirefox", "No firefox directory found", err)
	}

	// Copy the database, so problems don't arise when the file gets locked
	tempCookieDbff := filepath.Join(os.TempDir(), "tempCookieDbff.sqlite")

	// Clean up the temporary file when the function returns
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			log.Println("error removing file: ", err)
		}
	}(tempCookieDbff)

	// Copy the database to a temporary location
	copyError := utils.CopyFile(ffdirectory[0]+"\\cookies.sqlite", tempCookieDbff)
	if copyError != nil {
		return checks.NewCheckErrorf("CookieFirefox", "Unable to make a copy of the file", copyError)
	}

	db, err := sql.Open("sqlite", tempCookieDbff)
	if err != nil {
		return checks.NewCheckError("CookieFirefox", err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Println("error closing database: ", err)
		}
	}(db)

	// Query the name, origin and when the cookie was created from the database
	rows, err := db.Query("SELECT name, host, creationTime FROM moz_cookies")
	if err != nil {
		return checks.NewCheckError("CookieFirefox", err)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Println("error closing rows: ", err)
		}
	}(rows)

	// Iterate over each found cookie
	for rows.Next() {
		var name, host string
		var creationTime int64
		// Scan the row into variables
		if err := rows.Scan(&name, &host, &creationTime); err != nil {
			return checks.NewCheckError("CookieFirefox", err)
		}
		// Append the cookie to the result list
		timeofCreation := time.UnixMicro(creationTime)
		timeString := timeofCreation.String()
		output = append(output, name, host, timeString)
	}
	return checks.NewCheckResult("CookieFirefox", output...)
}
