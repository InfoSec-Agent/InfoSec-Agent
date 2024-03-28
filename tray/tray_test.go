package tray

import (
	"io"
	"os"
	"testing"
	"time"

	"github.com/InfoSec-Agent/InfoSec-Agent/localization"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/stretchr/testify/require"

	"github.com/getlantern/systray"
)

// TestMain initializes the system tray application and runs the tests
//
// Parameters: m *testing.M - The testing framework
//
// Returns: _
func TestMain(m *testing.M) {
	// Initialize systray

	go localization.Init("../")
	time.Sleep(100 * time.Millisecond)

	go systray.Run(OnReady, OnQuit)

	// Wait for the system tray application to initialize
	time.Sleep(100 * time.Millisecond)

	// Run tests
	exitCode := m.Run()

	// Clean up systray
	systray.Quit()

	os.Exit(exitCode)
}

// TestChangeScanInterval tests the ChangeScanInterval function with valid and invalid inputs
//
// Parameters: t *testing.T - The testing framework
//
// Returns: _
func TestChangeScanInterval(t *testing.T) {
	// Define test cases with input values and expected results
	testCases := []struct {
		input           string
		expectedMessage string
	}{
		// Valid input
		{"24", "Scan interval changed to 24 hours\n"},
		// Invalid input (non-numeric)
		{"abc", "Invalid input. Using default interval of 24 hours.Scan interval changed to 24 hours\n"},
		// Invalid input (negative)
		{"-1", "Invalid input. Using default interval of 24 hours.Scan interval changed to 24 hours\n"},
		// Invalid input (zero)
		{"0", "Invalid input. Using default interval of 24 hours.Scan interval changed to 24 hours\n"},
		// Valid large input
		{"1000", "Scan interval changed to 1000 hours\n"},
	}

	// Iterate over test cases
	for _, tc := range testCases {
		// Capture standard output to check the printed message
		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		// Run the function with mocked user input
		go ChangeScanInterval(tc.input)

		// Wait for the function to complete
		time.Sleep(100 * time.Millisecond)

		// Restore standard output
		err := w.Close()
		if err != nil {
			t.Errorf("Error closing pipe: %v", err)
		}
		os.Stdout = oldStdout
		capturedOutput, _ := io.ReadAll(r)

		// Assert that the printed message matches the expected message
		require.Equal(t, string(capturedOutput), tc.expectedMessage)
	}
}

// TestScanNow tests the ScanNow function
//
// Parameters: t *testing.T - The testing framework
//
// Returns: _
func TestScanNow(t *testing.T) {
	// Set up initial scanCounter value
	initialScanCounter := 0

	tickerAdvanced := make(chan struct{})

	// Listen for ticker advancement
	go func() {
		<-scanTicker.C
		tickerAdvanced <- struct{}{}
	}()

	// Run the function
	ScanNow()

	// Assert that scanCounter was incremented
	finalScanCounter := scanCounter
	expectedScanCounter := initialScanCounter + 1
	require.Equal(t, finalScanCounter, expectedScanCounter)
}

// TestOnQuit tests the OnQuit function
//
// Parameters: t *testing.T - The testing framework
//
// Returns: _
func TestOnQuit(t *testing.T) {
	// Mock OS signals channel
	sigc := make(chan os.Signal, 1)

	// Create a channel to signal completion
	quitCompleted := make(chan struct{})

	// Run OnQuit in a separate goroutine
	go func() {
		OnQuit()
		quitCompleted <- struct{}{}
	}()

	// Simulate quitting the application
	sigc <- os.Interrupt

	// Wait for OnQuit to complete
	select {
	case <-quitCompleted:
		// OnQuit completed
	case <-time.After(1 * time.Second):
		t.Error("OnQuit did not complete within the timeout")
	}
}

// TestTranslation tests the localization package, ensuring that strings are translated correctly
//
// Parameters: t *testing.T - The testing framework
//
// Returns: _
func TestTranslation(t *testing.T) {
	var localizer = localization.Localizers()[0]
	s1 := localizer.MustLocalize(&i18n.LocalizeConfig{
		MessageID: "ScanIntervalTitle",
	})
	// Change the language, then check if the translation is different
	localizer = localization.Localizers()[1]
	s2 := localizer.MustLocalize(&i18n.LocalizeConfig{
		MessageID: "ScanIntervalTitle",
	})
	require.NotEqual(t, s1, s2)
}

// TestChangeLang tests the tray.ChangeLang function on valid and invalid inputs
//
// Parameters: t *testing.T - The testing framework
//
// Returns: _
func TestChangeLang(t *testing.T) {
	testCases := []struct {
		input         string
		expectedIndex int
	}{
		// Valid input
		{"German", 0},
		{"British English", 1},
		{"American English", 2},
		{"Spanish", 3},
		{"French", 4},
		{"Dutch", 5},
		{"Portuguese", 6},
		// Invalid input, should return the default index (British English)
		{"Italian", 1},
	}

	for _, tc := range testCases {
		ChangeLanguage(tc.input)
		require.Equal(t, tc.expectedIndex, language)
	}
}

// TestRefreshMenu tests the RefreshMenu function
//
// Parameters: t *testing.T - The testing framework
//
// Returns: _
func TestRefreshMenu(t *testing.T) {
	value1 := menuItems[0].MenuTitle
	translation1 := localization.Localize(language, value1)
	ChangeLanguage("Spanish")
	// Refresh the menu, then check if the translation is different
	// RefreshMenu(menuItems)
	value2 := menuItems[0].MenuTitle
	translation2 := localization.Localize(language, value2)

	require.NotEqual(t, translation1, translation2)
}

// TestOpenReportingPageWhenAlreadyOpen tests the OpenReportingPage function when the reporting page is already open.
// It should not be able to open another reporting page when one is already running.
//
// Parameters: t *testing.T - The testing framework
//
// Returns: _
func TestOpenReportingPageWhenAlreadyOpen(t *testing.T) {
	rpPage = true
	err := openReportingPage("../")
	require.Error(t, err)
	require.Equal(t, "reporting-page is already running", err.Error())
}
