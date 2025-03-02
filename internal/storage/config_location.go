package storage

import (
	"os"
	"path/filepath"
)

const (
	// AppName is the name of the application, used for config directories
	AppName = "ubuntu-system-monitor"

	// ConfigFileName is the name of the configuration file
	ConfigFileName = "config.json"
)

// GetConfigDir returns the configuration directory for the application
func GetConfigDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	configDir := filepath.Join(homeDir, ".config", AppName)

	// Ensure the directory exists
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return "", err
	}

	return configDir, nil
}

// GetConfigPath returns the full path to the configuration file
func GetConfigPath() (string, error) {
	configDir, err := GetConfigDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(configDir, ConfigFileName), nil
}

// EnsureDirectoryExists ensures that a directory exists
func EnsureDirectoryExists(dirPath string) error {
	return os.MkdirAll(dirPath, 0755)
}
