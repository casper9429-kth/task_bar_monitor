package settings

import (
	"encoding/json"
	"log"
	"os"
)

// SaveToCustomPath saves settings to a custom location
func SaveToCustomPath(s *Config, customPath string) error {
	// Open the config file for writing
	file, err := os.Create(customPath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Encode the settings to JSON
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(s); err != nil {
		return err
	}

	log.Printf("Settings saved to %s", customPath)
	return nil
}

// LoadFromCustomPath loads settings from a custom location
func LoadFromCustomPath(customPath string) (*Config, error) {
	settings := DefaultSettings()

	// Try to open the config file
	file, err := os.Open(customPath)
	if err != nil {
		if os.IsNotExist(err) {
			// If file doesn't exist, use default settings
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

// IsMetricEnabled checks if a specific metric is enabled in the settings
func IsMetricEnabled(s *Config, metricName string) bool {
	switch metricName {
	case "cpu":
		return s.ShowCPU
	case "memory":
		return s.ShowMemory
	case "network":
		return s.ShowNetwork
	case "disk":
		return s.ShowDisk
	default:
		return false
	}
}

// ToggleMetric enables or disables a specific metric
func ToggleMetric(s *Config, metricName string, enabled bool) {
	switch metricName {
	case "cpu":
		s.ShowCPU = enabled
	case "memory":
		s.ShowMemory = enabled
	case "network":
		s.ShowNetwork = enabled
	case "disk":
		s.ShowDisk = enabled
	}

	// Update the ShowMetrics slice for compatibility
	s.ShowMetrics = []string{}
	if s.ShowCPU {
		s.ShowMetrics = append(s.ShowMetrics, "cpu")
	}
	if s.ShowMemory {
		s.ShowMetrics = append(s.ShowMetrics, "memory")
	}
	if s.ShowNetwork {
		s.ShowMetrics = append(s.ShowMetrics, "network")
	}
	if s.ShowDisk {
		s.ShowMetrics = append(s.ShowMetrics, "disk")
	}
}

// UpdateMetricsFromSlice updates the individual metric booleans based on the ShowMetrics slice
func UpdateMetricsFromSlice(s *Config) {
	// Reset all metrics to false
	s.ShowCPU = false
	s.ShowMemory = false
	s.ShowNetwork = false
	s.ShowDisk = false

	// Enable metrics based on the slice
	for _, metric := range s.ShowMetrics {
		switch metric {
		case "cpu":
			s.ShowCPU = true
		case "memory":
			s.ShowMemory = true
		case "network":
			s.ShowNetwork = true
		case "disk":
			s.ShowDisk = true
		}
	}
}

// GetEnabledMetrics returns a slice of enabled metrics
func GetEnabledMetrics(s *Config) []string {
	var enabled []string
	if s.ShowCPU {
		enabled = append(enabled, "cpu")
	}
	if s.ShowMemory {
		enabled = append(enabled, "memory")
	}
	if s.ShowNetwork {
		enabled = append(enabled, "network")
	}
	if s.ShowDisk {
		enabled = append(enabled, "disk")
	}
	return enabled
}
