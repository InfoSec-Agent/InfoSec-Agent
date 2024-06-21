package apiconnection_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	apiconnection "github.com/InfoSec-Agent/InfoSec-Agent/backend/api_connection"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/usersettings"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseCheckResult(t *testing.T) {
	tests := []struct {
		name  string
		check checks.Check
		want  apiconnection.IssueData
	}{
		{
			name: "Check error",
			check: checks.Check{
				IssueID: 1,
				Error:   errors.New("error"),
			},
			want: apiconnection.IssueData{IssueID: 1, Detected: false},
		},
		{
			name: "Check no error",
			check: checks.Check{
				IssueID:  3,
				ResultID: 1,
				Error:    nil,
			},
			want: apiconnection.IssueData{IssueID: 3, Detected: true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := apiconnection.ParseCheckResult(tt.check)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestParseScanResults(t *testing.T) {
	tests := []struct {
		name     string
		metaData apiconnection.Metadata
		checks   []checks.Check
		want     apiconnection.ParseResult
	}{
		{
			name: "Empty checks",
			metaData: apiconnection.Metadata{
				WorkStationID: 1,
				User:          "test",
				Date:          "2021-01-01",
			},
			checks: []checks.Check{},
			want: apiconnection.ParseResult{
				Metadata: apiconnection.Metadata{
					WorkStationID: 1,
					User:          "test",
					Date:          "2021-01-01",
				},
				Results: nil,
			},
		},
		{
			name: "Non-empty checks",
			metaData: apiconnection.Metadata{
				WorkStationID: 1,
				User:          "test",
				Date:          "2021-01-01",
			},
			checks: []checks.Check{
				{
					IssueID:  3,
					ResultID: 1,
					Error:    nil,
				}},
			want: apiconnection.ParseResult{
				Metadata: apiconnection.Metadata{
					WorkStationID: 1,
					User:          "test",
					Date:          "2021-01-01",
				},
				Results: []apiconnection.IssueData{
					{
						IssueID:  3,
						Detected: true,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := apiconnection.ParseScanResults(tt.metaData, tt.checks)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestParseString(t *testing.T) {
	tests := []struct {
		name string
		p    apiconnection.ParseResult
		want string
	}{
		{
			name: "Empty results",
			p: apiconnection.ParseResult{
				Metadata: apiconnection.Metadata{
					WorkStationID: 1,
					User:          "test",
					Date:          "2021-01-01",
				},
				Results: nil,
			},
			want: "Metadata: {1 test 2021-01-01}, Results: []",
		},
		{
			name: "Non-empty results",
			p: apiconnection.ParseResult{
				Metadata: apiconnection.Metadata{
					WorkStationID: 1,
					User:          "test",
					Date:          "2021-01-01",
				},
				Results: []apiconnection.IssueData{
					{
						IssueID:  3,
						Detected: true,
					},
				},
			},
			want: "Metadata: {1 test 2021-01-01}, Results: [{3 true []}]",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.p.String()
			require.Equal(t, tt.want, got)
		})
	}
}

// Mock implementation of UserSettings
type MockUserSettings struct{}

func (m *MockUserSettings) LoadUserSettings() usersettings.UserSettings {
	return usersettings.UserSettings{
		IntegrationKey: "mock-integration-key",
	}
}

var oldLoadUserSettings func() usersettings.UserSettings

// Create a ParseResult instance
type ParseResult struct {
	Status string `json:"status"`
}

func TestSendResultsToAPI(t *testing.T) {
	// Create a test server
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "Bearer mock-integration-key", r.Header.Get("Authorization"))

		var result ParseResult
		err := json.NewDecoder(r.Body).Decode(&result)
		assert.NoError(t, err)

		// Send a response
		w.WriteHeader(http.StatusOK)
	}))
	defer testServer.Close()

	// Override the URL for the test
	url := testServer.URL

	result := ParseResult{
		Status: "success",
	}

	// Convert result to JSON
	jsonData, err := json.Marshal(result)
	assert.NoError(t, err)

	// Act
	buffer := bytes.NewBuffer(jsonData)
	req, err := http.NewRequest("POST", url, buffer)
	assert.NoError(t, err)

	settings := usersettings.LoadUserSettings()
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+settings.IntegrationKey)
	req.Header.Set("Content-Length", fmt.Sprintf("%d", buffer.Len()))

	client := &http.Client{
		Timeout: 60 * time.Second,
	}
	resp, err := client.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	// Assert
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
