package app

import (
	"log"
	"sync"
	"time"

	"github.com/casper9429-kth/task_bar_monitor/internal/metrics"
	"github.com/casper9429-kth/task_bar_monitor/internal/settings"
	"github.com/casper9429-kth/task_bar_monitor/internal/ui"
)

// Application represents the main application
type Application struct {
	settings    *settings.Config
	tray        ui.TrayInterface
	monitoring  bool
	wg          sync.WaitGroup
	refreshChan chan struct{} // Channel to signal settings updates
}

// NewApplication creates a new application instance
func NewApplication() *Application {
	return &Application{
		refreshChan: make(chan struct{}, 1),
	}
}

// Run starts the application
func (a *Application) Run() error {
	// Load settings
	s, err := settings.LoadSettings()
	if err != nil {
		log.Printf("Failed to load settings: %v", err)
		// Continue with default settings
		s = settings.DefaultSettings()
	}
	a.settings = s

	// Initialize the UI with SysTray implementation
	a.tray = ui.NewTrayWithCallback(a.settings, a.onSettingsChanged)

	// Start monitoring in the background
	a.Start()

	// Start the system tray synchronously (this will block until tray is quit)
	a.tray.Start()

	return nil
}

// onSettingsChanged is called when settings are updated
func (a *Application) onSettingsChanged() {
	log.Println("Settings changed, updating application...")

	// Send signal to update the refresh interval
	select {
	case a.refreshChan <- struct{}{}:
		log.Println("Refresh signal sent")
	default:
		// Channel already has a pending update
	}
}

// Start begins monitoring system metrics
func (a *Application) Start() {
	a.monitoring = true
	a.wg.Add(1)
	go a.monitor()
}

// Stop stops monitoring system metrics
func (a *Application) Stop() {
	a.monitoring = false
	a.wg.Wait()

	// Stop the tray
	if a.tray != nil {
		a.tray.Stop()
	}
}

// monitor continuously monitors system metrics
func (a *Application) monitor() {
	defer a.wg.Done()

	// Set up a ticker for regular updates
	ticker := time.NewTicker(time.Duration(a.settings.RefreshInterval) * time.Second)
	defer ticker.Stop()

	log.Printf("Starting metrics monitoring with refresh interval: %d seconds", a.settings.RefreshInterval)

	// Immediately collect and update metrics once
	updateMetrics(a)

	// Then continue on the ticker schedule
	for a.monitoring {
		select {
		case <-ticker.C:
			updateMetrics(a)
		case <-a.refreshChan:
			// Settings have changed, update ticker with new interval
			ticker.Stop()
			ticker = time.NewTicker(time.Duration(a.settings.RefreshInterval) * time.Second)
			log.Printf("Refresh interval updated to: %d seconds", a.settings.RefreshInterval)

			// Immediately update metrics with new settings
			updateMetrics(a)
		}
	}
}

// updateMetrics collects and updates metrics
func updateMetrics(a *Application) {
	// Get CPU usage
	cpuUsage := 0.0
	if a.settings.ShowCPU {
		usage, err := metrics.GetCPUUsage()
		if err != nil {
			log.Printf("Failed to get CPU usage: %v", err)
		} else {
			cpuUsage = usage
		}
	}

	// Get memory usage
	memUsage := 0.0
	if a.settings.ShowMemory {
		usage, err := metrics.GetMemoryUsage()
		if err != nil {
			log.Printf("Failed to get memory usage: %v", err)
		} else {
			memUsage = usage
		}
	}

	// Get disk usage
	diskUsage := 0.0
	if a.settings.ShowDisk {
		usage, err := metrics.GetDiskUsage()
		if err != nil {
			log.Printf("Failed to get disk usage: %v", err)
		} else {
			diskUsage = usage
		}
	}

	// Get network usage
	netUsage, err := metrics.GetNetworkUsage()
	if err != nil {
		log.Printf("Failed to get network usage: %v", err)
	}

	// Update the tray with the latest metrics
	a.tray.UpdateMetrics(cpuUsage, memUsage, diskUsage, netUsage)
}
