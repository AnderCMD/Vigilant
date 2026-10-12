// Package assets embeds the application icons.
package assets

import _ "embed"

// IconActive is the tray icon shown while Vigilant is running.
//
//go:embed icon.ico
var IconActive []byte

// IconPaused is the tray icon shown while Vigilant is paused.
//
//go:embed icon-paused.ico
var IconPaused []byte
