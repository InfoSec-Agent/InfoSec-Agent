package checks_test

import (
	"github.com/InfoSec-Agent/InfoSec-Agent/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/registrymock"
	"reflect"
	"testing"
)

func TestSecureBoot(t *testing.T) {
	tests := []struct {
		name string
		key  registrymock.RegistryKey
		want checks.Check
	}{
		{
			name: "SecureBootEnabled",
			key:  &registrymock.MockRegistryKey{SubKeys: []registrymock.MockRegistryKey{{KeyName: "SYSTEM\\CurrentControlSet\\Control\\SecureBoot\\State", IntegerValues: map[string]uint64{"UEFISecureBootEnabled": 1}, Err: nil}}},
			want: checks.NewCheckResult("SecureBoot", "Secure boot is enabled"),
		},
		{
			name: "SecureBootDisabled",
			key:  &registrymock.MockRegistryKey{SubKeys: []registrymock.MockRegistryKey{{KeyName: "SYSTEM\\CurrentControlSet\\Control\\SecureBoot\\State", IntegerValues: map[string]uint64{"UEFISecureBootEnabled": 0}, Err: nil}}},
			want: checks.NewCheckResult("SecureBoot", "Secure boot is disabled"),
		},
		{
			name: "SecureBootUnknown",
			key:  &registrymock.MockRegistryKey{SubKeys: []registrymock.MockRegistryKey{{KeyName: "SYSTEM\\CurrentControlSet\\Control\\SecureBoot\\State", IntegerValues: map[string]uint64{"UEFISecureBootEnabled": 2}, Err: nil}}},
			want: checks.NewCheckResult("SecureBoot", "Secure boot status is unknown"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := checks.SecureBoot(tt.key); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SecureBoot() = %v, want %v", got, tt.want)
			}
		})
	}
}
