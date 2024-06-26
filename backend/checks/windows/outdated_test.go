package windows_test

import (
	"errors"
	"testing"

	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks/windows"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/html"

	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/mocking"
)

// TestWindowsOutdated is a function that tests the behavior of the Outdated function with various inputs.
//
// Parameters:
//   - t *testing.T: The testing framework provided by the Go testing package. It provides methods for reporting test failures and logging additional information.
//
// Returns: None
//
// This function tests the Outdated function with different scenarios. It uses a mock implementation of the WindowsVersion interface to simulate the behavior of retrieving the Windows version information. Each test case checks if the Outdated function correctly identifies whether the Windows version is up-to-date, outdated, or unsupported based on the simulated Windows version information. The function asserts that the returned Check instance contains the expected results.
func TestWindowsOutdated(t *testing.T) {
	win10HTML := windows.GetURLBody("https://learn.microsoft.com/en-us/windows/release-health/release-information")
	latestWin10Build := windows.FindWindowsBuild(win10HTML)

	win11HTML := windows.GetURLBody("https://learn.microsoft.com/en-us/windows/release-health/windows11-release-information")
	latestWin11Build := windows.FindWindowsBuild(win11HTML)

	expected1 := []string{"11", "WinVersion: 22000.000"}
	expected2 := []string{"10", "WinVersion: 0.0"}

	tests := []struct {
		name         string
		mockExecutor *mocking.MockCommandExecutor
		want         checks.Check
		error        bool
	}{
		{
			name:         "Windows 11 up-to-date",
			mockExecutor: &mocking.MockCommandExecutor{Output: "Microsoft Windows [Version 10.0." + latestWin11Build + "]", Err: nil},
			want:         checks.NewCheckResult(checks.WindowsOutdatedID, 0, "11"),
		},
		{
			name:         "Windows 11 outdated",
			mockExecutor: &mocking.MockCommandExecutor{Output: "Microsoft Windows [Version 10.0.22000.000]", Err: nil},
			want:         checks.NewCheckResult(checks.WindowsOutdatedID, 1, expected1...),
		},
		{
			name:         "Windows 10 up-to-date",
			mockExecutor: &mocking.MockCommandExecutor{Output: "Microsoft Windows [Version 10.0." + latestWin10Build + "]", Err: nil},
			want:         checks.NewCheckResult(checks.WindowsOutdatedID, 0, "10"),
		},
		{
			name:         "Windows 10 outdated",
			mockExecutor: &mocking.MockCommandExecutor{Output: "Microsoft Windows [Version 10.0.0.0]", Err: nil},
			want:         checks.NewCheckResult(checks.WindowsOutdatedID, 1, expected2...),
		},
		{
			name:         "Unsupported Windows version",
			mockExecutor: &mocking.MockCommandExecutor{Output: "Microsoft Windows [Version 9.0.0.0]", Err: nil},
			want:         checks.NewCheckResult(checks.WindowsOutdatedID, 2),
		},
		{
			name:         "Error executing command",
			mockExecutor: &mocking.MockCommandExecutor{Output: "", Err: errors.New("command error")},
			want:         checks.NewCheckError(checks.WindowsOutdatedID, errors.New("command error")),
		},
		{
			name:         "Match not long enough",
			mockExecutor: &mocking.MockCommandExecutor{Output: "Microsoft Windows [Version 10.0.0]", Err: nil},
			want:         checks.NewCheckError(checks.WindowsOutdatedID, errors.New("error parsing Windows version string")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := windows.Outdated(tt.mockExecutor)
			if tt.error {
				require.Equal(t, -1, got.ResultID)
			} else {
				require.Equal(t, tt.want, got)
			}
		})
	}
}

func TestGetUrlBody(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want *html.Node
	}{
		{
			name: "Invalid URL",
			url:  "invalid",
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := windows.GetURLBody(tt.url)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestFindWindowsBuild(t *testing.T) {
	tests := []struct {
		name string
		html *html.Node
		want string
	}{
		{
			name: "Valid HTML with build number",
			html: &html.Node{
				Type: html.ElementNode, Data: "tbody", FirstChild: &html.Node{Type: html.ElementNode, Data: "tr", FirstChild: &html.Node{Type: html.ElementNode, Data: "td",
					NextSibling: &html.Node{Type: html.ElementNode, Data: "td", NextSibling: &html.Node{Type: html.ElementNode, Data: "td", NextSibling: &html.Node{Type: html.ElementNode, Data: "td",
						NextSibling: &html.Node{Type: html.ElementNode, Data: "td", FirstChild: &html.Node{
							Type: html.TextNode, Data: "10.0.19042"}}}}}}}},
			want: "10.0.19042",
		},
		{
			name: "Valid HTML without build number",
			html: &html.Node{
				Type: html.ElementNode, Data: "tbody", FirstChild: &html.Node{Type: html.ElementNode, Data: "tr", FirstChild: &html.Node{Type: html.ElementNode, Data: "td",
					NextSibling: &html.Node{Type: html.ElementNode, Data: "td", NextSibling: &html.Node{Type: html.ElementNode, Data: "td", NextSibling: &html.Node{Type: html.ElementNode, Data: "td",
						NextSibling: &html.Node{Type: html.ElementNode, Data: "td", FirstChild: &html.Node{
							Type: html.TextNode, Data: ""}}}}}}}},
			want: "",
		},
		{
			name: "Valid HTML with table in different location",
			html: &html.Node{
				Type: html.ElementNode, Data: "table", FirstChild: &html.Node{
					Type: html.TextNode, Data: "", NextSibling: &html.Node{
						Type: html.ElementNode, Data: "tbody", FirstChild: &html.Node{Type: html.ElementNode, Data: "tr", FirstChild: &html.Node{Type: html.ElementNode, Data: "td",
							NextSibling: &html.Node{Type: html.ElementNode, Data: "td", NextSibling: &html.Node{Type: html.ElementNode, Data: "td", NextSibling: &html.Node{Type: html.ElementNode, Data: "td",
								NextSibling: &html.Node{Type: html.ElementNode, Data: "td", FirstChild: &html.Node{
									Type: html.TextNode, Data: "10.0.19042"}}}}}}}}}},
			want: "10.0.19042",
		},
		{
			name: "Invalid HTML",
			html: &html.Node{Type: html.ElementNode, Data: "div"},
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := windows.FindWindowsBuild(tt.html)
			require.Equal(t, tt.want, got)
		})
	}
}
