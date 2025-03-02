package ui

import (
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/getlantern/systray"
	"github.com/yourusername/ubuntu-system-monitor/internal/metrics"
	"github.com/yourusername/ubuntu-system-monitor/internal/settings"
)

// Indicator represents the system tray indicator
type Indicator struct {
	cpuItem           *systray.MenuItem
	memoryItem        *systray.MenuItem
	networkItem       *systray.MenuItem
	diskItem          *systray.MenuItem
	settingsItem      *systray.MenuItem
	quitItem          *systray.MenuItem
	settings          *settings.Config
	settingsWin       *SettingsWindow
	ready             bool
	mutex             sync.Mutex
	stopChan          chan struct{}
	onSettingsChanged func() // Callback for settings changes
}

// NewIndicator creates a new system tray indicator
func NewIndicator(s *settings.Config) *Indicator {
	log.Println("Creating new indicator")
	return &Indicator{
		settings: s,
		ready:    false,
		stopChan: make(chan struct{}),
	}
}

// NewIndicatorWithCallback creates a new system tray indicator with a settings change callback
func NewIndicatorWithCallback(s *settings.Config, callback func()) *Indicator {
	log.Println("Creating new indicator with settings callback")
	return &Indicator{
		settings:          s,
		ready:             false,
		stopChan:          make(chan struct{}),
		onSettingsChanged: callback,
	}
}

// Start initializes the system tray indicator
func (i *Indicator) Start() {
	log.Println("Starting indicator")
	// Run systray synchronously to ensure proper initialization
	systray.Run(i.onReady, i.onExit)
}

// onReady is called when the systray is ready
func (i *Indicator) onReady() {
	log.Println("Indicator ready, setting up menu items")
	// Set a simple title initially
	systray.SetTitle("System Monitor")
	systray.SetTooltip("Ubuntu System Monitor")

	// Create menu items
	i.cpuItem = systray.AddMenuItem("CPU: Loading...", "CPU Usage")
	i.memoryItem = systray.AddMenuItem("Memory: Loading...", "Memory Usage")
	i.networkItem = systray.AddMenuItem("Network: Loading...", "Network Usage")
	i.diskItem = systray.AddMenuItem("Disk: Loading...", "Disk Usage")

	systray.AddSeparator()
	i.settingsItem = systray.AddMenuItem("Settings", "Configure the application")
	systray.AddSeparator()
	i.quitItem = systray.AddMenuItem("Quit", "Quit the application")

	// Apply initial visibility settings
	i.updateItemVisibility()

	// Mark as ready
	i.mutex.Lock()
	i.ready = true
	i.mutex.Unlock()

	log.Println("Menu setup complete, starting handlers")

	// Start handling events
	go i.handleEvents()

	// Initialize with default values to show something immediately
	i.updateMetricsDisplay(0, 0, 0, metrics.NetworkUsage{})
}

// handleEvents processes menu item events
func (i *Indicator) handleEvents() {
	// Initialize the settings window (must be done on the same goroutine that will show it)
	i.settingsWin = NewSettingsWindow(i.settings, i.onSettingsSaved)

	for {
		select {
		case <-i.settingsItem.ClickedCh:
			log.Println("Settings clicked")
			i.settingsWin.Show()
		case <-i.quitItem.ClickedCh:
			log.Println("Quit clicked")
			systray.Quit()
			return
		case <-i.stopChan:
			log.Println("Stop signal received")
			systray.Quit()
			return
		}
	}
}

// Stop stops the indicator
func (i *Indicator) Stop() {
	log.Println("Stopping indicator")
	close(i.stopChan)
}

// onSettingsSaved is called when settings are saved
func (i *Indicator) onSettingsSaved() {
	log.Println("Settings saved, updating visibility")
	i.updateItemVisibility()

	// Force an immediate metrics update to reflect any changes
	log.Println("Refreshing metrics display after settings change")
	metrics, err := metrics.GetAllMetrics()
	if err != nil {
		log.Printf("Error getting metrics: %v", err)
	} else {
		i.UpdateMetrics(metrics.CPUUsage, metrics.MemoryUsage, metrics.DiskUsage, metrics.NetworkUsage)
	}

	// Call the settings changed callback if provided
	if i.onSettingsChanged != nil {
		log.Println("Calling settings changed callback")
		i.onSettingsChanged()
	}
}

// updateItemVisibility updates the visibility of menu items based on settings
func (i *Indicator) updateItemVisibility() {
	log.Println("Updating item visibility")
	if i.settings.ShowCPU {
		i.cpuItem.Show()
	} else {
		i.cpuItem.Hide()
	}

	if i.settings.ShowMemory {
		i.memoryItem.Show()
	} else {
		i.memoryItem.Hide()
	}

	if i.settings.ShowNetwork {
		i.networkItem.Show()
	} else {
		i.networkItem.Hide()
	}

	if i.settings.ShowDisk {
		i.diskItem.Show()
	} else {
		i.diskItem.Hide()
	}
}

// UpdateMetrics updates the menu items with the latest metrics
func (i *Indicator) UpdateMetrics(cpuUsage, memUsage, diskUsage float64, netUsage metrics.NetworkUsage) {
	// Only update if ready
	i.mutex.Lock()
	ready := i.ready
	i.mutex.Unlock()

	if !ready {
		log.Println("Indicator not ready yet, skipping metrics update")
		return
	}

	log.Printf("Updating metrics - CPU: %.1f%%, Memory: %.1f%%, Disk: %.1f%%, Net: ↓%.1f KB/s ↑%.1f KB/s",
		cpuUsage, memUsage, diskUsage, netUsage.DownloadSpeed, netUsage.UploadSpeed)

	// Update in the main UI thread for safety
	go func() {
		i.updateMetricsDisplay(cpuUsage, memUsage, diskUsage, netUsage)
	}()
}

// updateMetricsDisplay updates the UI components with metrics
func (i *Indicator) updateMetricsDisplay(cpuUsage, memUsage, diskUsage float64, netUsage metrics.NetworkUsage) {
	// Create title with all requested metrics
	var titleParts []string

	if i.settings.ShowCPUInTitle && i.settings.ShowCPU {
		cpuText := fmt.Sprintf("C:%.1f%%", cpuUsage)
		titleParts = append(titleParts, cpuText)
	}

	if i.settings.ShowMemoryInTitle && i.settings.ShowMemory {
		memText := fmt.Sprintf("M:%.1f%%", memUsage)
		titleParts = append(titleParts, memText)
	}

	if i.settings.ShowDiskInTitle && i.settings.ShowDisk {
		diskText := fmt.Sprintf("D:%.1f%%", diskUsage)
		titleParts = append(titleParts, diskText)
	}

	if i.settings.ShowNetworkInTitle && i.settings.ShowNetwork {
		var netText string
		if i.settings.ShowBothNetworkSpeeds {
			// Show both upload and download speeds
			netText = fmt.Sprintf("N:↓%.1f ↑%.1f", netUsage.DownloadSpeed, netUsage.UploadSpeed)
		} else {
			// Show only download speed to save space
			netText = fmt.Sprintf("N:↓%.1f", netUsage.DownloadSpeed)
		}
		titleParts = append(titleParts, netText)
	}

	// If no metrics selected for title, show a default
	if len(titleParts) == 0 {
		systray.SetTitle("Sys Monitor")
	} else {
		// Join all parts with a separator
		titleText := strings.Join(titleParts, " | ")
		log.Printf("Setting title to: %s", titleText)
		systray.SetTitle(titleText)
	}

	// Set tooltips with more detailed information
	tooltipText := ""
	sep := ""

	if i.settings.ShowCPU {
		cpuTooltip := fmt.Sprintf("CPU: %.1f%%", cpuUsage)
		tooltipText += cpuTooltip
		sep = " | "
	}

	if i.settings.ShowMemory {
		memTooltip := fmt.Sprintf("MEM: %.1f%%", memUsage)
		tooltipText += sep + memTooltip
		sep = " | "
	}

	if i.settings.ShowDisk {
		diskTooltip := fmt.Sprintf("DISK: %.1f%%", diskUsage)
		tooltipText += sep + diskTooltip
		sep = " | "
	}

	if i.settings.ShowNetwork {
		netTooltip := fmt.Sprintf("NET: ↓%.1f ↑%.1f KB/s", netUsage.DownloadSpeed, netUsage.UploadSpeed)
		tooltipText += sep + netTooltip
	}

	if tooltipText == "" {
		tooltipText = "System Monitor"
	}

	log.Printf("Setting tooltip to: %s", tooltipText)
	systray.SetTooltip(tooltipText)

	// Update individual menu items
	if i.settings.ShowCPU {
		i.cpuItem.SetTitle(fmt.Sprintf("CPU: %.1f%%", cpuUsage))
	}

	if i.settings.ShowMemory {
		i.memoryItem.SetTitle(fmt.Sprintf("Memory: %.1f%%", memUsage))
	}

	if i.settings.ShowDisk {
		i.diskItem.SetTitle(fmt.Sprintf("Disk: %.1f%%", diskUsage))
	}

	if i.settings.ShowNetwork {
		i.networkItem.SetTitle(fmt.Sprintf("Network: ↓%.1f KB/s ↑%.1f KB/s",
			netUsage.DownloadSpeed, netUsage.UploadSpeed))
	}
}

// onExit is called when the systray is exiting
func (i *Indicator) onExit() {
	log.Println("Exiting system monitor")
}

// getDefaultIcon returns a default icon for the systray
func getDefaultIcon() []byte {
	// A simple 16x16 icon with transparency
	return []byte{
		0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A, 0x00, 0x00, 0x00, 0x0D, 0x49, 0x48, 0x44, 0x52,
		0x00, 0x00, 0x00, 0x10, 0x00, 0x00, 0x00, 0x10, 0x08, 0x06, 0x00, 0x00, 0x00, 0x1F, 0xF3, 0xFF,
		0x61, 0x00, 0x00, 0x00, 0x01, 0x73, 0x52, 0x47, 0x42, 0x00, 0xAE, 0xCE, 0x1C, 0xE9, 0x00, 0x00,
		0x00, 0xB4, 0x49, 0x44, 0x41, 0x54, 0x38, 0x11, 0xED, 0xD2, 0xB1, 0x6A, 0xC2, 0x50, 0x18, 0xC5,
		0xF1, 0xF3, 0xDD, 0x24, 0x28, 0x4E, 0x59, 0x9D, 0x9C, 0x7C, 0x04, 0x71, 0x72, 0x11, 0x1C, 0x8B,
		0xAE, 0xEE, 0x0E, 0x3E, 0x42, 0x97, 0xD0, 0xB5, 0x83, 0x93, 0x83, 0x8F, 0xE0, 0x23, 0x38, 0x75,
		0xF0, 0x15, 0x3A, 0x39, 0x09, 0x0E, 0xA5, 0xB9, 0xB9, 0x1F, 0x69, 0x0C, 0xB1, 0xE9, 0x50, 0xF0,
		0x3F, 0x1E, 0xEE, 0xBD, 0xF7, 0xE3, 0x72, 0xA0, 0x69, 0x1A, 0x9C, 0x73, 0xB8, 0xDA, 0xF3, 0x7D,
		0x17, 0x6B, 0xCD, 0xC6, 0x18, 0x3C, 0x77, 0xAE, 0x58, 0x42, 0xA1, 0xF3, 0x76, 0x58, 0xF0, 0xB7,
		0x81, 0x3E, 0x5C, 0x96, 0x25, 0x4C, 0x3F, 0x12, 0x0E, 0xF1, 0x34, 0x8E, 0x41, 0xA1, 0x18, 0x89,
		0xC5, 0xD5, 0x17, 0x81, 0x49, 0xC7, 0xF6, 0xA8, 0x39, 0x80, 0x20, 0x51, 0xEF, 0x42, 0xC1, 0x6C,
		0xB6, 0xC0, 0x5B, 0x96, 0xE1, 0xBD, 0xEF, 0x23, 0x0C, 0x43, 0x0C, 0xFB, 0xDD, 0xFC, 0xDB, 0xE7,
		0x91, 0xA9, 0x5F, 0xC7, 0xF0, 0xBE, 0x0A, 0x27, 0xA2, 0xB8, 0xC9, 0xA6, 0xBF, 0x8B, 0xCF, 0x72,
		0x74, 0x30, 0x1F, 0x5E, 0xE4, 0x87, 0x98, 0x4C, 0xE7, 0x07, 0x83, 0xA5, 0xE5, 0xEF, 0xFE, 0x1D,
		0xCF, 0x87, 0x88, 0x23, 0x07, 0xC2, 0x86, 0x33, 0xF4, 0x89, 0xA2, 0x08, 0xD7, 0x74, 0xA6, 0x44,
		0x5F, 0xE8, 0x4E, 0x44, 0xD6, 0x54,
	}
}
