package integration_test

import (
	i "github.com/InfoSec-Agent/InfoSec-Agent/backend/integration_testing"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/logger"

	"testing"
)

var testsInvalid = []func(t *testing.T){
	i.TestIntegrationExtensionsChromiumWithoutAdBlocker,
	i.TestIntegrationHistoryChromiumWithPhishing,
	i.TestIntegrationCISRegistrySettingsIncorrect,
	i.TestIntegrationExtensionsFirefoxWithoutAdBlocker,
	i.TestIntegrationHistoryFirefoxWithPhishing,
	i.TestIntegrationOpenPortsPorts,
	i.TestIntegrationSmbCheckBadSetup,
	i.TestIntegrationPasswordManagerNotPresent,
	i.TestIntegrationAdvertisementActive,
	i.TestIntegrationAutomatedLoginActive,
	i.TestIntegrationGuestAccountActive,
	i.TestIntegrationLoginMethodPINOnly,
	i.TestIntegrationPermissionWithApps,
	i.TestIntegrationRemoteDesktopEnabled,
	i.TestIntegrationSecureBootDisabled,
	i.TestIntegrationStartupWithApps,
	i.TestIntegrationUACDisabled,
	i.TestIntegrationCookiesFirefoxWithCookies,
	i.TestIntegrationCookiesChromiumWithCookies,
	i.TestIntegrationRemoteRPCEnabled,
	i.TestIntegrationNetBIOSEnabled,
	i.TestIntegrationWPADEnabled,
	i.TestIntegrationCredentialGuardDisabled,
	i.TestIntegrationFirewallDisabled,
	i.TestIntegrationPasswordComplexityInvalid,
	i.TestIntegrationScreenLockDisabled,
}

func TestAllInvalid(t *testing.T) {
	logger.SetupTests()
	for _, test := range testsInvalid {
		test(t)
	}
}
