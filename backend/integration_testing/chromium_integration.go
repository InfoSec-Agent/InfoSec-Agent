package integration

import (
	"testing"

	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks/browsers"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks/browsers/chromium"
	"github.com/stretchr/testify/require"
)

func TestIntegrationExtensionsChromiumWithAdBlocker(t *testing.T) {
	tests := []struct {
		name    string
		browser string
	}{
		{
			name:    "Chrome with ad blocker",
			browser: "Chrome",
		},
		{
			name:    "Edge with ad blocker",
			browser: "Edge",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := chromium.ExtensionsChromium(tt.browser, browsers.RealDefaultDirGetter{}, chromium.RealExtensionIDGetter{}, chromium.ChromeExtensionNameGetter{})
			require.NotEqual(t, -1, result.ResultID)
			require.NotEmpty(t, result)
			require.Equal(t, 0, result.ResultID)
		})
	}
}

func TestIntegrationExtensionsChromiumWithoutAdBlocker(t *testing.T) {
	tests := []struct {
		name    string
		browser string
	}{
		{
			name:    "Chrome without ad blocker",
			browser: "Chrome",
		},
		{
			name:    "Edge without ad blocker",
			browser: "Edge",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := chromium.ExtensionsChromium(tt.browser, browsers.RealDefaultDirGetter{}, chromium.RealExtensionIDGetter{}, chromium.ChromeExtensionNameGetter{})
			require.NotEqual(t, -1, result.ResultID)
			require.NotEmpty(t, result)
			require.Equal(t, 1, result.ResultID)
		})
	}
}

func TestIntegrationHistoryChromiumWithoutPhishing(t *testing.T) {
	tests := []struct {
		name    string
		browser string
	}{
		{
			name:    "Chrome without phishing",
			browser: "Chrome",
		},
		{
			name:    "Edge without phishing",
			browser: "Edge",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := chromium.HistoryChromium(tt.browser, browsers.RealDefaultDirGetter{}, chromium.RealCopyDBGetter{}, chromium.RealQueryDatabaseGetter{}, chromium.RealProcessQueryResultsGetter{}, browsers.RealPhishingDomainGetter{})
			require.NotEqual(t, -1, result.ResultID)
			require.NotEmpty(t, result)
			require.Equal(t, 1, result.ResultID)
		})
	}
}

func TestIntegrationHistoryChromiumWithPhishing(t *testing.T) {
	tests := []struct {
		name    string
		browser string
	}{
		{
			name:    "Chrome with phishing",
			browser: "Chrome",
		},
		{
			name:    "Edge with phishing",
			browser: "Edge",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := chromium.HistoryChromium(tt.browser, browsers.RealDefaultDirGetter{}, chromium.RealCopyDBGetter{}, chromium.RealQueryDatabaseGetter{}, chromium.RealProcessQueryResultsGetter{}, browsers.RealPhishingDomainGetter{})
			require.NotEqual(t, -1, result.ResultID)
			require.NotEmpty(t, result)
			require.Equal(t, 0, result.ResultID)
		})
	}
}

func TestIntegrationSearchEngineChromiumWithSearchEngine(t *testing.T) {
	tests := []struct {
		name    string
		browser string
	}{
		{
			name:    "Chrome with some search engine",
			browser: "Chrome",
		},
		{
			name:    "Edge with some search engine",
			browser: "Edge",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := chromium.SearchEngineChromium(tt.browser, false, nil, browsers.RealDefaultDirGetter{})
			require.NotEmpty(t, result)
			require.NotEqual(t, -1, result.ResultID)
			require.Equal(t, 0, result.ResultID)
		})
	}
}

func TestIntegrationCookiesChromiumWithCookies(t *testing.T) {
	tests := []struct {
		name    string
		browser string
	}{
		{
			name:    "Chrome with cookies",
			browser: "Chrome",
		},
		{
			name:    "Edge with cookies",
			browser: "Edge",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := chromium.CookiesChromium(tt.browser, browsers.RealDefaultDirGetter{}, browsers.RealCopyFileGetter{}, browsers.RealQueryCookieDatabaseGetter{})
			require.NotEmpty(t, result)
			require.NotEqual(t, -1, result.ResultID)
			require.Equal(t, 1, result.ResultID)
		})
	}
}

func TestIntegrationCookiesChromiumWithoutCookies(t *testing.T) {
	tests := []struct {
		name    string
		browser string
	}{
		{
			name:    "Chrome without cookies",
			browser: "Chrome",
		},
		{
			name:    "Edge without cookies",
			browser: "Edge",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := chromium.CookiesChromium(tt.browser, browsers.RealDefaultDirGetter{}, browsers.RealCopyFileGetter{}, browsers.RealQueryCookieDatabaseGetter{})
			require.Empty(t, result.Result)
			require.Equal(t, 0, result.ResultID)
		})
	}
}
