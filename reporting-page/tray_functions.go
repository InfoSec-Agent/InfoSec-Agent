package main

import (
	"github.com/InfoSec-Agent/InfoSec-Agent/tray"
)

// Tray represents the system tray, used for calling system tray functions from the front-end
type Tray struct {
}

// NewTray creates a new Tray struct
//
// Parameters: _
//
// Returns: a pointer to a new Tray struct
func NewTray() *Tray {
	return &Tray{}
}

// ScanNow calls the ScanNow function from the tray package
//
// Parameters: t (*Tray) - a pointer to the Tray struct
//
// Returns: _
func (t *Tray) ScanNow() {
	tray.ScanNow()
}

// ChangeLanguage calls the ChangeLanguage function from the tray package
//
// Parameters: t (*Tray) - a pointer to the Tray struct
//
// Returns: _
func (t *Tray) ChangeLanguage() {
	tray.ChangeLanguage()
	tray.RefreshMenu()
}

// ChangeScanInterval calls the ChangeScanInterval function from the tray package
//
// Parameters: t (*Tray) - a pointer to the Tray struct
//
// Returns: _
func (t *Tray) ChangeScanInterval() {
	tray.ChangeScanInterval()
}
