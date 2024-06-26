package windows_test

import (
	"errors"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks/windows"
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/mocking"
)

// TestSecureBoot is a function that tests the behavior of the SecureBoot function with various inputs.
//
// Parameters:
//   - t *testing.T: The testing framework provided by the Go testing package. It provides methods for reporting test failures and logging additional information.
//
// Returns: None
//
// This function tests the SecureBoot function with different scenarios. It uses a mock implementation of the RegistryKey interface to simulate the behavior of the Secure Boot registry key. Each test case checks if the SecureBoot function correctly identifies the status of Secure Boot (enabled, disabled, or unknown) based on the simulated registry key value. The function asserts that the returned Check instance contains the expected results.
func TestSecureBoot(t *testing.T) {
	tests := []struct {
		name  string
		key   mocking.RegistryKey
		want  checks.Check
		error bool
	}{
		{
			name: "SecureBootEnabled",
			key: &mocking.MockRegistryKey{SubKeys: []mocking.MockRegistryKey{{
				KeyName:       "SYSTEM\\CurrentControlSet\\Control\\SecureBoot\\State",
				IntegerValues: map[string]uint64{"UEFISecureBootEnabled": 1}, Err: nil}}},
			want: checks.NewCheckResult(checks.SecureBootID, 1),
		},
		{
			name: "SecureBootDisabled",
			key: &mocking.MockRegistryKey{SubKeys: []mocking.MockRegistryKey{{
				KeyName:       "SYSTEM\\CurrentControlSet\\Control\\SecureBoot\\State",
				IntegerValues: map[string]uint64{"UEFISecureBootEnabled": 0}, Err: nil}}},
			want: checks.NewCheckResult(checks.SecureBootID, 0),
		},
		{
			name: "SecureBootUnknown",
			key: &mocking.MockRegistryKey{SubKeys: []mocking.MockRegistryKey{{
				KeyName:       "SYSTEM\\CurrentControlSet\\Control\\SecureBoot\\State",
				IntegerValues: map[string]uint64{"UEFISecureBootEnabled": 2}, Err: nil}}},
			want: checks.NewCheckResult(checks.SecureBootID, 2),
		},
		{
			name: "SecureBootUnknown",
			key: &mocking.MockRegistryKey{SubKeys: []mocking.MockRegistryKey{{
				KeyName:       "SYSTEM\\CurrentControlSet\\Control\\SecureBoot\\State",
				IntegerValues: map[string]uint64{"UEFISecureBootEnabled": 2}, Err: nil}}},
			want: checks.NewCheckResult(checks.SecureBootID, 2),
		},
		{
			name:  "Error opening registry key",
			key:   &mocking.MockRegistryKey{SubKeys: []mocking.MockRegistryKey{}},
			want:  checks.NewCheckError(checks.SecureBootID, errors.New("error")),
			error: true,
		},
		{
			name: "Error reading integer value",
			key: &mocking.MockRegistryKey{SubKeys: []mocking.MockRegistryKey{{
				KeyName: "SYSTEM\\CurrentControlSet\\Control\\SecureBoot\\State"}}},
			want:  checks.NewCheckError(checks.SecureBootID, errors.New("error")),
			error: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := windows.SecureBoot(tt.key)
			if tt.error {
				require.Equal(t, -1, got.ResultID)
			} else {
				require.Equal(t, tt.want, got)
			}
		})
	}
}
