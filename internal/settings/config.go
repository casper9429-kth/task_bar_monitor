package settings

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Config represents the application settings
type Config struct {
	ShowCPU               bool     `json:"showCPU"`
	ShowMemory            bool     `json:"showMemory"`
	ShowNetwork           bool     `json:"showNetwork"`
	ShowDisk              bool     `json:"showDisk"`
	ShowCPUInTitle        bool     `json:"showCPUInTitle"`
	ShowMemoryInTitle     bool     `json:"showMemoryInTitle"`
	ShowNetworkInTitle    bool     `json:"showNetworkInTitle"`
	ShowBothNetworkSpeeds bool     `json:"showBothNetworkSpeeds"` // Option for showing both upload and download
	ShowDiskInTitle       bool     `json:"showDiskInTitle"`
	RefreshInterval       int      `json:"refreshInterval"`
	ShowMetrics           []string `json:"showMetrics"` // For compatibility with UI
	configPath            string
}

// DefaultSettings returns the default application settings
func DefaultSettings() *Config {
	homeDir, _ := os.UserHomeDir()
	configDir := filepath.Join(homeDir, ".config", "ubuntu-system-monitor")

	return &Config{
		ShowCPU:               true,
		ShowMemory:            true,
		ShowNetwork:           true,
		ShowDisk:              true,
		ShowCPUInTitle:        true,  // By default, show CPU in title
		ShowMemoryInTitle:     false, // Others are optional
		ShowNetworkInTitle:    false,
		ShowBothNetworkSpeeds: false, // Off by default to save space
		ShowDiskInTitle:       false,
		RefreshInterval:       2,
		ShowMetrics:           []string{"cpu", "memory", "network", "disk"},
		configPath:            filepath.Join(configDir, "config.json"),
	}
}

// LoadSettings loads the settings from the config file
func LoadSettings() (*Config, error) {
	settings := DefaultSettings()

	// Create config directory if it doesn't exist
	configDir := filepath.Dir(settings.configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return settings, err
	}

	// Try to open the config file
	file, err := os.Open(settings.configPath)
	if err != nil {
		if os.IsNotExist(err) {
			// If file doesn't exist, save default settings
			settings.Save()
			return settings, nil
		}
		return settings, err
	}
	defer file.Close()

	// Decode the config file
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(settings); err != nil {
		return settings, err
	}

	return settings, nil
}

// Save saves the settings to the config file
func (s *Config) Save() error {
	// Create config directory if it doesn't exist
	configDir := filepath.Dir(s.configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	// Open the config file for writing
	file, err := os.Create(s.configPath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Encode the settings to JSON
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(s)
}
