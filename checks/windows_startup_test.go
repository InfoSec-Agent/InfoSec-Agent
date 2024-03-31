package checks_test

import (
	"reflect"
	"testing"

	"github.com/InfoSec-Agent/InfoSec-Agent/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/registrymock"
)

// TestStartup tests the Startup function on (in)valid input
//
// Parameters: t (testing.T) - the testing framework
//
// Returns: _
func TestStartup(t *testing.T) {
	// [not done yet]
	tests := []struct {
		name string
		key1 registrymock.RegistryKey
		key2 registrymock.RegistryKey
		key3 registrymock.RegistryKey
		want checks.Check
	}{{
		name: "No startup programs found",
		key1: &registrymock.MockRegistryKey{SubKeys: []registrymock.MockRegistryKey{{KeyName: "SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Explorer\\StartupApproved\\Run"}}},
		key2: &registrymock.MockRegistryKey{SubKeys: []registrymock.MockRegistryKey{{KeyName: "SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Explorer\\StartupApproved\\Run"}}},
		key3: &registrymock.MockRegistryKey{SubKeys: []registrymock.MockRegistryKey{{KeyName: "SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Explorer\\StartupApproved\\Run32"}}},
		want: checks.NewCheckResult("Startup", "No startup programs found"),
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := checks.Startup()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Startup() = %v, want %v", got, tt.want)
			}
		})
	}
}
