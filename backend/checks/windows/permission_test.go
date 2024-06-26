package windows_test

import (
	"errors"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks/windows"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/mocking"

	"github.com/stretchr/testify/assert"
)

// TestPermission is a function that tests the Permission function's ability to correctly return permissions.
//
// Parameters:
//   - t *testing.T: The testing framework provided by the Go testing package. It provides methods for reporting test failures and logging additional information.
//
// Returns: None
//
// This function tests the Permission function with different scenarios. It uses a mock implementation of the RegistryKey interface to simulate different sets of permissions. Each test case checks if the Permission function correctly identifies the presence or absence of specific permissions based on the simulated registry keys. The function asserts that the returned permissions match the expected results.
func TestPermission(t *testing.T) {
	tests := []struct {
		name       string
		permission string
		key        mocking.RegistryKey
		want       checks.Check
		error      bool
	}{
		{
			name:       "NonPackagedWebcamPermissionExists",
			permission: "webcam",
			key: &mocking.MockRegistryKey{
				SubKeys: []mocking.MockRegistryKey{
					{KeyName: "Software\\Microsoft\\Windows\\CurrentVersion\\CapabilityAccessManager\\ConsentStore\\webcam",
						StringValues: map[string]string{"Value": "Allow"},
						SubKeys: []mocking.MockRegistryKey{
							{KeyName: "NonPackaged",
								SubKeys: []mocking.MockRegistryKey{
									{KeyName: "microsoft.webcam", StringValues: map[string]string{"Value": "Allow"}},
								},
							},
						},
					},
				},
			},
			want: checks.NewCheckResult(checks.WebcamID, 0, "microsoft webcam"),
		},
		{
			name:       "WebcamPermissionExists",
			permission: "webcam",
			key: &mocking.MockRegistryKey{
				SubKeys: []mocking.MockRegistryKey{
					{KeyName: "Software\\Microsoft\\Windows\\CurrentVersion\\CapabilityAccessManager\\ConsentStore\\webcam",
						StringValues: map[string]string{"Value": "Allow"},
						SubKeys: []mocking.MockRegistryKey{
							{KeyName: "microsoft.webcam", StringValues: map[string]string{"Value": "Allow"}},
						},
					},
				},
			},
			want: checks.NewCheckResult(checks.WebcamID, 0, "microsoft webcam"),
		},
		{
			name:       "No apps with webcam permission",
			permission: "webcam",
			key: &mocking.MockRegistryKey{
				SubKeys: []mocking.MockRegistryKey{
					{KeyName: "Software\\Microsoft\\Windows\\CurrentVersion\\CapabilityAccessManager\\ConsentStore\\webcam",
						StringValues: map[string]string{"Value": "Deny"},
						SubKeys:      []mocking.MockRegistryKey{},
					},
				},
			},
			want: checks.NewCheckResult(checks.WebcamID, 1),
		},
		{
			name:       "Error reading string value",
			permission: "webcam",
			key: &mocking.MockRegistryKey{
				SubKeys: []mocking.MockRegistryKey{
					{KeyName: "Software\\Microsoft\\Windows\\CurrentVersion\\CapabilityAccessManager\\ConsentStore\\webcam",
						StringValues: map[string]string{"Value": ""},
						SubKeys: []mocking.MockRegistryKey{
							{KeyName: "microsoft.webcam", StringValues: map[string]string{"Value": "Allow"}},
						},
					},
				},
			},
			want:  checks.NewCheckError(checks.WebcamID, errors.New("error")),
			error: true,
		},
		{
			name:       "Error reading sub key",
			permission: "webcam",
			key: &mocking.MockRegistryKey{
				SubKeys: []mocking.MockRegistryKey{
					{KeyName: "Software\\Microsoft\\Windows\\CurrentVersion\\CapabilityAccessManager\\ConsentStore\\webcam",
						StringValues: map[string]string{"Value": "Allow"},
						SubKeys: []mocking.MockRegistryKey{
							{KeyName: "microsoft.webcam", StringValues: map[string]string{"test": "test"}},
						},
					},
				},
			},
			want:  checks.NewCheckError(checks.WebcamID, errors.New("error")),
			error: true,
		},
		{
			name:       "Error opening sub key",
			permission: "webcam",
			key: &mocking.MockRegistryKey{
				SubKeys: []mocking.MockRegistryKey{
					{KeyName: "Software\\Microsoft\\Windows\\CurrentVersion\\CapabilityAccessManager\\ConsentStore\\webcam",
						StringValues: map[string]string{"Value": "Allow"},
						SubKeys: []mocking.MockRegistryKey{
							{KeyName: "microsoft.webcam", StringValues: map[string]string{"Value": "Allow"},
								IntegerValues: map[string]uint64{"test": 1}},
						},
					},
				},
			},
			want:  checks.NewCheckError(checks.WebcamID, errors.New("error")),
			error: true,
		},
		{
			name:       "Error reading sub key sub key names",
			permission: "webcam",
			key: &mocking.MockRegistryKey{
				SubKeys: []mocking.MockRegistryKey{
					{KeyName: "Software\\Microsoft\\Windows\\CurrentVersion\\CapabilityAccessManager\\ConsentStore\\webcam",
						StringValues: map[string]string{"Value": "Allow"},
						SubKeys: []mocking.MockRegistryKey{
							{KeyName: "NonPackaged", StringValues: map[string]string{"Value": "Allow"},
								SubKeys: []mocking.MockRegistryKey{
									{KeyName: "microsoft.webcam", StringValues: map[string]string{"Value": "", "test": "test"}},
								},
							},
						},
					},
				},
			},
			want:  checks.NewCheckError(checks.WebcamID, errors.New("error")),
			error: true,
		},
		{
			name:       "Error reading sub key string value",
			permission: "webcam",
			key: &mocking.MockRegistryKey{
				SubKeys: []mocking.MockRegistryKey{
					{KeyName: "Software\\Microsoft\\Windows\\CurrentVersion\\CapabilityAccessManager\\ConsentStore\\webcam",
						StringValues: map[string]string{"Value": "Allow"},
						SubKeys: []mocking.MockRegistryKey{
							{KeyName: "Packaged", StringValues: map[string]string{"Value": ""},
								SubKeys: []mocking.MockRegistryKey{
									{KeyName: "microsoft.webcam", StringValues: map[string]string{"Value": "", "test": "test"}},
								},
							},
						},
					},
				},
			},
			want: checks.NewCheckResult(checks.WebcamID, 1),
		},
		{
			name:       "Error reading sub key string value",
			permission: "webcam",
			key: &mocking.MockRegistryKey{
				SubKeys: []mocking.MockRegistryKey{
					{KeyName: "Software\\Microsoft\\Windows\\CurrentVersion\\CapabilityAccessManager\\ConsentStore\\webcam",
						StringValues: map[string]string{"Value": "Allow"},
						SubKeys: []mocking.MockRegistryKey{
							{KeyName: "Packaged", StringValues: map[string]string{"Value": "test"},
								SubKeys: []mocking.MockRegistryKey{
									{KeyName: "microsoft.webcam", StringValues: map[string]string{"Value": "", "test": "test"}},
								},
							},
						},
					},
				},
			},
			want: checks.NewCheckResult(checks.WebcamID, 1),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := windows.Permission(checks.WebcamID, tc.permission, tc.key)
			if tc.error {
				require.Equal(t, -1, result.ResultID)
			} else {
				require.Equal(t, tc.want, result)
			}
		})
	}
}

// TestFormatPermission is a function that tests the Permission function's ability to correctly format the returned permissions.
//
// Parameters:
//   - t *testing.T: The testing framework provided by the Go testing package. It provides methods for reporting test failures and logging additional information.
//
// Returns: None
//
// This function tests the Permission function with a specific scenario where the permission string is in a certain format. It uses a mock implementation of the RegistryKey interface to simulate a specific permission format. The test case checks if the Permission function correctly formats the returned permission string by removing any '#' characters. The function asserts that the returned permission string matches the expected format.
func TestFormatPermission(t *testing.T) {
	key := &mocking.MockRegistryKey{
		SubKeys: []mocking.MockRegistryKey{
			{KeyName: "Software\\Microsoft\\Windows\\CurrentVersion\\CapabilityAccessManager\\ConsentStore\\location",
				StringValues: map[string]string{"Value": "Allow"},
				SubKeys: []mocking.MockRegistryKey{
					{KeyName: "NonPackaged",
						StringValues: map[string]string{"Value": "Allow"},
						SubKeys: []mocking.MockRegistryKey{
							{KeyName: "test#test#test.exe"},
						},
					},
				},
			},
		},
	}
	c := windows.Permission(checks.LocationID, "location", key)
	assert.NotContains(t, c.Result, "#")
	assert.Contains(t, c.Result[0], "test")
}

// TestNonExistingPermission is a function that tests the Permission function's behavior when the requested permission does not exist.
//
// Parameters:
//   - t *testing.T: The testing framework provided by the Go testing package. It provides methods for reporting test failures and logging additional information.
//
// Returns: None
//
// This function tests the Permission function with a scenario where the requested permission does not exist in the simulated registry keys. It uses a mock implementation of the RegistryKey interface to simulate this scenario. The test case checks if the Permission function correctly returns an error when the requested permission does not exist. The function asserts that the returned Check instance contains the expected error message.
func TestNonExistingPermission(t *testing.T) {
	key := &mocking.MockRegistryKey{
		SubKeys: []mocking.MockRegistryKey{
			{KeyName: "Software\\Microsoft\\Windows\\CurrentVersion\\CapabilityAccessManager\\ConsentStore\\location",
				StringValues: map[string]string{"Value": "Allow"},
				SubKeys: []mocking.MockRegistryKey{
					{KeyName: "NonPackaged",
						StringValues: map[string]string{"Value": "Allow"},
						SubKeys: []mocking.MockRegistryKey{
							{KeyName: "test test"},
						},
					},
				},
			},
		},
	}
	c := windows.Permission(99, "hello", key)
	assert.Equal(t, c.Result, []string(nil))
	assert.EqualError(t, c.Error, "error opening registry key: error opening registry key: key not found")
}

// TestRemoveDuplicateStrRemovesDuplicates validates the functionality of the RemoveDuplicateStr function by ensuring it correctly removes duplicate string values from a given slice.
//
// This test function creates a slice with duplicate string values and passes it to the RemoveDuplicateStr function.
// It asserts that the returned slice contains only the unique string values from the input slice, in the order of their first occurrence.
//
// Parameters:
//   - t *testing.T: The testing framework used for assertions.
//
// No return values.
func TestRemoveDuplicateStrRemovesDuplicates(t *testing.T) {
	input := []string{"a", "b", "a", "c", "b"}
	expected := []string{"a", "b", "c"}
	result := windows.RemoveDuplicateStr(input)
	require.Equal(t, expected, result)
}

// TestRemoveDuplicateStrEmptyInput validates the behavior of the RemoveDuplicateStr function when provided with an empty input.
//
// This test function creates an empty string slice and passes it to the RemoveDuplicateStr function.
// It asserts that the returned slice is also empty, confirming that the function handles empty input correctly.
//
// Parameters:
//   - t *testing.T: The testing framework used for assertions.
//
// No return values.
func TestRemoveDuplicateStrEmptyInput(t *testing.T) {
	var input []string
	var expected []string
	result := windows.RemoveDuplicateStr(input)
	require.Equal(t, expected, result)
}
