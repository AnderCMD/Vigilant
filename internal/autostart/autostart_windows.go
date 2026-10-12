// Package autostart manages the "Start with Windows" registry entry.
package autostart

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"golang.org/x/sys/windows/registry"
)

const (
	runKey    = `Software\Microsoft\Windows\CurrentVersion\Run`
	valueName = "Vigilant"
)

func quotedExe() (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", err
	}
	return `"` + exe + `"`, nil
}

// Enabled reports whether Vigilant is registered to start at login and the
// registered path points at this executable.
func Enabled() bool {
	k, err := registry.OpenKey(registry.CURRENT_USER, runKey, registry.QUERY_VALUE)
	if err != nil {
		return false
	}
	defer k.Close()
	v, _, err := k.GetStringValue(valueName)
	if err != nil {
		return false
	}
	exe, err := quotedExe()
	return err == nil && strings.EqualFold(v, exe)
}

// Set registers or removes the login entry for the current user.
func Set(enabled bool) error {
	k, _, err := registry.CreateKey(registry.CURRENT_USER, runKey, registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("open Run key: %w", err)
	}
	defer k.Close()
	if !enabled {
		if err := k.DeleteValue(valueName); err != nil && !errors.Is(err, registry.ErrNotExist) {
			return err
		}
		return nil
	}
	exe, err := quotedExe()
	if err != nil {
		return err
	}
	return k.SetStringValue(valueName, exe)
}
