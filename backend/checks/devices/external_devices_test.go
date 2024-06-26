package devices_test

import (
	"errors"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks/devices"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/mocking"
)

// TestExternalDevices is a unit test for the ExternalDevices function in the "checks" package.
// It runs a series of test scenarios with varying inputs to validate the function's behavior.
//
// Parameters:
//   - t (*testing.T): A pointer to an instance of the testing framework, used for reporting test results.
//
// Returns: None. Failures are reported through the *testing.T parameter.
//
// The function is part of the test suite for the "checks" package. It ensures that the ExternalDevices function accurately detects external devices connected to the system and handles error scenarios appropriately.
func TestExternalDevices(t *testing.T) {
	tests := []struct {
		name          string
		executorClass *mocking.MockCommandExecutor
		want          checks.Check
	}{
		{
			name:          "No external devices connected",
			executorClass: &mocking.MockCommandExecutor{Output: "", Err: nil},
			want:          checks.NewCheckResult(checks.ExternalDevicesID, 0),
		},
		{
			name:          "External devices connected",
			executorClass: &mocking.MockCommandExecutor{Output: "\r\nFriendlyName\r\n-\r\nHD WebCam\r\n\r\n\r\n\r\n", Err: nil},
			want:          checks.NewCheckResult(checks.ExternalDevicesID, 1, "HD WebCam", ""),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := devices.ExternalDevices(tt.executorClass)
			require.Equal(t, tt.want, got)
		})
	}
}

// TestCheckDeviceClasses is a testing function that ensures the correct functionality of the CheckDeviceClass function.
// It runs multiple test cases with different inputs to validate the expected behavior of the function in various situations.
//
// Parameters:
//   - t (*testing.T): A pointer to the testing framework instance, used for logging and reporting test results.
//
// Returns: None. Any test failures are reported via the *testing.T parameter.
//
// This function is a component of the "checks" package test suite. It confirms that the CheckDeviceClass function accurately identifies devices of a given class connected to the system and properly handles any errors.
func TestCheckDeviceClasses(t *testing.T) {
	tests := []struct {
		name          string
		deviceClass   []string
		executorClass *mocking.MockCommandExecutor
		want          []string
		error         bool
	}{
		{
			name:          "No devices of the specified class",
			deviceClass:   []string{"Mouse"},
			executorClass: &mocking.MockCommandExecutor{Output: "\r\nFriendlyName\r\n-\r\n\r\n\r\n\r\n", Err: nil},
			want:          []string{""},
		},
		{
			name:        "Devices of the specified class",
			deviceClass: []string{"Camera"},
			executorClass: &mocking.MockCommandExecutor{
				Output: "\r\nFriendlyName\r\n-\r\nHD WebCam\r\n\r\n\r\n\r\n", Err: nil},
			want: []string{"HD WebCam", ""},
		},
		{
			name:        "Error executing command",
			deviceClass: []string{"Camera"},
			executorClass: &mocking.MockCommandExecutor{
				Output: "", Err: errors.New("error")},
			want:  nil,
			error: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := devices.CheckDeviceClasses(tt.deviceClass, tt.executorClass)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ExternalDevices() = %v, want %v", got, tt.want)
			}
		})
	}
}
