package ui

import (
	"log"

	"github.com/yourusername/ubuntu-system-monitor/internal/metrics"
	"github.com/yourusername/ubuntu-system-monitor/internal/settings"
)

// TrayInterface defines the interface for system tray implementations
type TrayInterface interface {
	Start()
	UpdateMetrics(cpuUsage, memUsage, diskUsage float64, netUsage metrics.NetworkUsage)
	Stop()
}

// NewTray creates a new system tray
func NewTray(s *settings.Config) TrayInterface {
	log.Println("Creating system tray indicator")
	return NewIndicator(s)
}

// NewTrayWithCallback creates a new system tray with a settings change callback
func NewTrayWithCallback(s *settings.Config, onSettingsChanged func()) TrayInterface {
	log.Println("Creating system tray indicator with settings callback")
	return NewIndicatorWithCallback(s, onSettingsChanged)
}
