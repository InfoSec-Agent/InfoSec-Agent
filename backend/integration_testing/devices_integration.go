package integration

import (
	"testing"

	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks/devices"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/mocking"
	"github.com/stretchr/testify/require"
	"golang.org/x/sys/windows/registry"
)

func TestIntegrationBluetoothNoDevices(t *testing.T) {
	result := devices.Bluetooth(mocking.NewRegistryKeyWrapper(registry.LOCAL_MACHINE))
	// Check if function did not return an error
	require.NotEmpty(t, result)
	require.Empty(t, result.Result)
	require.Equal(t, 0, result.ResultID)
}

func TestIntegrationExternalDevicesNoDevices(t *testing.T) {
	result := devices.ExternalDevices(&mocking.RealCommandExecutor{})
	require.NotEmpty(t, result)
	require.Empty(t, result.Result)
	require.Equal(t, 0, result.ResultID)
}
