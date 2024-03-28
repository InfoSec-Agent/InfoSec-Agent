package firefox

import (
	"encoding/binary"
	"io"
	"os"
	"path/filepath"
	"regexp"

	"github.com/InfoSec-Agent/InfoSec-Agent/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/utils"
	"github.com/pierrec/lz4"
)

// SearchEngineFireFox checks the standard search engine in firefox.
//
// Parameters: _
//
// Returns: The standard search engine for firefox
func SearchEngineFirefox() checks.Check {
	// Determine the directory in which the Firefox profile is stored
	ffdirectory, err := utils.FirefoxFolder()
	if err != nil {
		return checks.NewCheckErrorf("SearchEngineFirefox", "No firefox directory found", err)
	}
	filePath := ffdirectory[0] + "/search.json.mozlz4"

	// Create a temporary file to copy the compressed json to
	tempSearch := filepath.Join(os.TempDir(), "tempSearch.json.mozlz4")
	defer os.Remove(tempSearch)

	// Copy the compressed json to a temporary location
	copyError := utils.CopyFile(filePath, tempSearch)
	if copyError != nil {
		return checks.NewCheckErrorf("SearchEngineFirefox", "Unable to make a copy of the file", copyError)
	}

	fileInfo, err := os.Stat(tempSearch)
	if err != nil {
		return checks.NewCheckErrorf("SearchEngineFirefox", "Unable to retrieve information about the file", err)
	}
	fileSize := fileInfo.Size()

	// Holds the information from the copied file
	file, err := os.Open(filepath.Clean(tempSearch))
	if err != nil {
		return checks.NewCheckErrorf("SearchEngineFirefox", "Unable to open the file", err)
	}
	defer utils.CloseFile(file)

	// Holds the custom magig number for the mozzila lz4 compression
	magicNumber := make([]byte, 8)

	// Retrieves the magicNumber from the file
	_, err = file.Read(magicNumber)
	if err != nil {
		return checks.NewCheckErrorf("SearchEngineFirefox", "Unable to read the file", err)
	}

	// Holds the size of the file after decompressing it
	uncompressSize := make([]byte, 4)

	// Skip the first 8 bytes to take the bytes 8-11 that hold the size after decompression
	_, err = file.Seek(8, io.SeekStart)
	if err != nil {
		return checks.NewCheckErrorf("SearchEngineFirefox", "Unable to skip the first 8 bytes", err)
	}

	// Retrieves bytes 8-11 to find the size of the file
	_, err = file.Read(uncompressSize)
	if err != nil {
		return checks.NewCheckErrorf("SearchEngineFirefox", "Unable to read the file", err)
	}

	// Transforms the size of the file after decompression from Little Endian to a normal 32-bit integer
	unCompressedSize := binary.LittleEndian.Uint32(uncompressSize)

	// Skip the first 12 bytes because that is the start of the data
	_, err = file.Seek(12, io.SeekStart)
	if err != nil {
		return checks.NewCheckErrorf("SearchEngineFirefox", "Unable to skip the first 12 bytes", err)
	}

	// Byte slice to hold the compressed data without the header (first 12 bytes)
	compressedData := make([]byte, fileSize-12)

	_, err = file.Read(compressedData)
	if err != nil {
		return checks.NewCheckErrorf("SearchEngineFirefox", "Unable to read the file", err)
	}

	data := make([]byte, unCompressedSize)
	_, err = lz4.UncompressBlock(compressedData, data)
	if err != nil {
		return checks.NewCheckErrorf("SearchEngineFirefox", "Unable to uncompress", err)
	}
	output := string(data)
	var result string
	// Regex to check if the defaultEngineId is empty which means that the engine is Google
	re := regexp.MustCompile(`"defaultEngineId":""`)
	// Performs the regex on the string and returns either google or goes into the next check
	matches := re.FindString(output)
	if matches != "" {
		result = "Google"
	} else {
		// This pattern looks for the values of the other known search engines and returns them
		pattern := `defaultEngineId":"(?:ddg|bing|ebay|wikipedia|amazon)@search\.mozilla\.org`
		re := regexp.MustCompile(pattern)
		matches := re.FindString(output)
		if matches == "" {
			return checks.NewCheckResult("SearchEngineFirefox", "Other Search Engine")
		}
		result = matches[18:]
	}
	return checks.NewCheckResult("SearchEngineFirefox", result)
}
