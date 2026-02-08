package config

import (
	"os"

	"golang.org/x/sys/windows/registry"
)

const AutoStartName = "Aviator"

// IsAutoStartEnabled checks if the registry key exists
func IsAutoStartEnabled() bool {
	k, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Run`, registry.QUERY_VALUE)
	if err != nil {
		return false
	}
	defer k.Close()

	if _, _, err := k.GetStringValue(AutoStartName); err != nil {
		return false // Not found
	}
	return true
}

// SetAutoStart adds or removes the registry key
func SetAutoStart(enabled bool) error {
	k, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Run`, registry.QUERY_VALUE|registry.SET_VALUE)
	if err != nil {
		return err
	}
	defer k.Close()

	if enabled {
		exePath, err := os.Executable()
		if err != nil {
			return err
		}
		// Write path to registry
		return k.SetStringValue(AutoStartName, exePath)
	}

	// Remove from registry (ignore error if not exists)
	if err := k.DeleteValue(AutoStartName); err != nil && err != registry.ErrNotExist {
		return err
	}
	return nil
}
