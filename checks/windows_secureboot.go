package checks

import (
	"golang.org/x/sys/windows/registry"
	"log"
)

// SecureBoot checks if Windows secure boot is enabled
//
// Parameters: _
//
// Returns: If Windows secure boot is enabled or not
func SecureBoot() Check {
	// Get secure boot information from the registry
	windowsSecureBoot, err := registry.OpenKey(registry.LOCAL_MACHINE,
		`SYSTEM\CurrentControlSet\Control\SecureBoot\State`, registry.READ)

	if err != nil {
		return NewCheckError("SecureBoot", err)
	}

	defer func(key registry.Key) {
		err := key.Close()
		if err != nil {
			log.Printf("error closing registry key: %v", err)
		}
	}(windowsSecureBoot)

	// Read the status of secure boot
	secureBootStatus, _, err := windowsSecureBoot.GetIntegerValue("UEFISecureBootEnabled")
	if err != nil {
		return NewCheckError("SecureBoot", err)
	}

	// Using the status, determine if secure boot is enabled or not
	if secureBootStatus == 1 {
		return NewCheckResult("SecureBoot", "Secure boot is enabled")
	}

	return NewCheckResult("SecureBoot", "Secure boot is disabled")
}
