package ui

import (
	"log"

	"github.com/andlabs/ui"
	"github.com/casper9429-kth/task_bar_monitor/internal/settings"
)

// SettingsWindow represents the settings window UI
type SettingsWindow struct {
	window                *ui.Window
	cpuCheck              *ui.Checkbox
	memoryCheck           *ui.Checkbox
	networkCheck          *ui.Checkbox
	diskCheck             *ui.Checkbox
	cpuTitleCheck         *ui.Checkbox
	memoryTitleCheck      *ui.Checkbox
	networkTitleCheck     *ui.Checkbox
	networkFullSpeedCheck *ui.Checkbox
	diskTitleCheck        *ui.Checkbox
	refreshIntervalEntry  *ui.Spinbox
	saveButton            *ui.Button
	cancelButton          *ui.Button
	appSettings           *settings.Config
	onSaved               func()
}

// NewSettingsWindow creates a new settings window
func NewSettingsWindow(appSettings *settings.Config, onSaved func()) *SettingsWindow {
	sw := &SettingsWindow{
		appSettings: appSettings,
		onSaved:     onSaved,
	}
	sw.initUI()
	return sw
}

// initUI initializes the settings window UI
func (sw *SettingsWindow) initUI() {
	// Create a smaller window now that we don't have tabs
	sw.window = ui.NewWindow("System Monitor Settings", 450, 430, false)
	sw.window.SetMargined(true)
	sw.window.OnClosing(func(*ui.Window) bool {
		sw.window.Hide()
		return false // Don't destroy the window, just hide it
	})

	// Main vertical box for all content
	mainBox := ui.NewVerticalBox()
	mainBox.SetPadded(true)

	// METRIC VISIBILITY GROUP
	visibilityGroup := ui.NewGroup("Show Metrics in Menu")
	visibilityGroup.SetMargined(true)

	visibilityVBox := ui.NewVerticalBox()
	visibilityVBox.SetPadded(true)

	// CPU checkbox
	sw.cpuCheck = ui.NewCheckbox("Show CPU Usage")
	sw.cpuCheck.SetChecked(sw.appSettings.ShowCPU)
	visibilityVBox.Append(sw.cpuCheck, false)

	// Memory checkbox
	sw.memoryCheck = ui.NewCheckbox("Show Memory Usage")
	sw.memoryCheck.SetChecked(sw.appSettings.ShowMemory)
	visibilityVBox.Append(sw.memoryCheck, false)

	// Network checkbox
	sw.networkCheck = ui.NewCheckbox("Show Network Usage")
	sw.networkCheck.SetChecked(sw.appSettings.ShowNetwork)
	visibilityVBox.Append(sw.networkCheck, false)

	// Disk checkbox
	sw.diskCheck = ui.NewCheckbox("Show Disk Usage")
	sw.diskCheck.SetChecked(sw.appSettings.ShowDisk)
	visibilityVBox.Append(sw.diskCheck, false)

	visibilityGroup.SetChild(visibilityVBox)
	mainBox.Append(visibilityGroup, false)

	// TASKBAR DISPLAY GROUP
	taskbarGroup := ui.NewGroup("Taskbar Display")
	taskbarGroup.SetMargined(true)

	taskbarVBox := ui.NewVerticalBox()
	taskbarVBox.SetPadded(true)

	// CPU in taskbar
	hboxCPU := ui.NewHorizontalBox()
	hboxCPU.SetPadded(true)
	sw.cpuTitleCheck = ui.NewCheckbox("Show CPU in Taskbar")
	sw.cpuTitleCheck.SetChecked(sw.appSettings.ShowCPUInTitle)
	hboxCPU.Append(sw.cpuTitleCheck, false)
	taskbarVBox.Append(hboxCPU, false)

	// Memory in taskbar
	hboxMem := ui.NewHorizontalBox()
	hboxMem.SetPadded(true)
	sw.memoryTitleCheck = ui.NewCheckbox("Show Memory in Taskbar")
	sw.memoryTitleCheck.SetChecked(sw.appSettings.ShowMemoryInTitle)
	hboxMem.Append(sw.memoryTitleCheck, false)
	taskbarVBox.Append(hboxMem, false)

	// Network in taskbar
	hboxNet := ui.NewHorizontalBox()
	hboxNet.SetPadded(true)
	sw.networkTitleCheck = ui.NewCheckbox("Show Network in Taskbar")
	sw.networkTitleCheck.SetChecked(sw.appSettings.ShowNetworkInTitle)
	hboxNet.Append(sw.networkTitleCheck, false)
	taskbarVBox.Append(hboxNet, false)

	// Network details option
	hboxNetDetails := ui.NewHorizontalBox()
	hboxNetDetails.SetPadded(true)
	sw.networkFullSpeedCheck = ui.NewCheckbox("Show Both Upload & Download")
	sw.networkFullSpeedCheck.SetChecked(sw.appSettings.ShowBothNetworkSpeeds)
	hboxNetDetails.Append(ui.NewLabel("    "), false) // Indent
	hboxNetDetails.Append(sw.networkFullSpeedCheck, false)
	taskbarVBox.Append(hboxNetDetails, false)

	// Disk in taskbar
	hboxDisk := ui.NewHorizontalBox()
	hboxDisk.SetPadded(true)
	sw.diskTitleCheck = ui.NewCheckbox("Show Disk in Taskbar")
	sw.diskTitleCheck.SetChecked(sw.appSettings.ShowDiskInTitle)
	hboxDisk.Append(sw.diskTitleCheck, false)
	taskbarVBox.Append(hboxDisk, false)

	taskbarGroup.SetChild(taskbarVBox)
	mainBox.Append(taskbarGroup, false)

	// UPDATE SETTINGS GROUP - moved from advanced tab
	updateGroup := ui.NewGroup("Update Settings")
	updateGroup.SetMargined(true)

	updateVBox := ui.NewVerticalBox()
	updateVBox.SetPadded(true)

	// Refresh interval control
	refreshHBox := ui.NewHorizontalBox()
	refreshHBox.SetPadded(true)
	refreshHBox.Append(ui.NewLabel("Refresh Interval (seconds):"), false)

	// Create a spinbox with valid range (1-10 seconds)
	sw.refreshIntervalEntry = ui.NewSpinbox(1, 10)
	sw.refreshIntervalEntry.SetValue(sw.appSettings.RefreshInterval)
	refreshHBox.Append(sw.refreshIntervalEntry, true)
	updateVBox.Append(refreshHBox, false)

	// Add info text
	infoText := ui.NewLabel("Lower values update more frequently but may use more CPU.")
	updateVBox.Append(infoText, false)

	updateGroup.SetChild(updateVBox)
	mainBox.Append(updateGroup, false)

	// BUTTONS
	buttonsBox := ui.NewHorizontalBox()
	buttonsBox.SetPadded(true)

	sw.saveButton = ui.NewButton("Save")
	sw.saveButton.OnClicked(func(*ui.Button) {
		sw.saveSettings()
	})

	sw.cancelButton = ui.NewButton("Cancel")
	sw.cancelButton.OnClicked(func(*ui.Button) {
		sw.window.Hide()
	})

	// Add some space at the left to push buttons to the right
	buttonsBox.Append(ui.NewLabel(""), true)
	buttonsBox.Append(sw.cancelButton, false)
	buttonsBox.Append(sw.saveButton, false)

	mainBox.Append(buttonsBox, false)

	sw.window.SetChild(mainBox)
}

// saveSettings saves the settings from the UI to the config
func (sw *SettingsWindow) saveSettings() {
	// Update metric visibilities
	sw.appSettings.ShowCPU = sw.cpuCheck.Checked()
	sw.appSettings.ShowMemory = sw.memoryCheck.Checked()
	sw.appSettings.ShowNetwork = sw.networkCheck.Checked()
	sw.appSettings.ShowDisk = sw.diskCheck.Checked()

	// Update taskbar title visibilities
	sw.appSettings.ShowCPUInTitle = sw.cpuTitleCheck.Checked()
	sw.appSettings.ShowMemoryInTitle = sw.memoryTitleCheck.Checked()
	sw.appSettings.ShowNetworkInTitle = sw.networkTitleCheck.Checked()
	sw.appSettings.ShowBothNetworkSpeeds = sw.networkFullSpeedCheck.Checked()
	sw.appSettings.ShowDiskInTitle = sw.diskTitleCheck.Checked()

	// Update refresh interval
	sw.appSettings.RefreshInterval = sw.refreshIntervalEntry.Value()

	// Update the ShowMetrics slice for compatibility
	sw.appSettings.ShowMetrics = []string{}
	if sw.appSettings.ShowCPU {
		sw.appSettings.ShowMetrics = append(sw.appSettings.ShowMetrics, "cpu")
	}
	if sw.appSettings.ShowMemory {
		sw.appSettings.ShowMetrics = append(sw.appSettings.ShowMetrics, "memory")
	}
	if sw.appSettings.ShowNetwork {
		sw.appSettings.ShowMetrics = append(sw.appSettings.ShowMetrics, "network")
	}
	if sw.appSettings.ShowDisk {
		sw.appSettings.ShowMetrics = append(sw.appSettings.ShowMetrics, "disk")
	}

	// Save to file
	err := sw.appSettings.Save()
	if err != nil {
		log.Printf("Failed to save settings: %v", err)
	} else {
		log.Println("Settings saved successfully")
	}

	// Hide window and notify
	sw.window.Hide()
	if sw.onSaved != nil {
		sw.onSaved()
	}
}

// Show displays the settings window
func (sw *SettingsWindow) Show() {
	sw.window.Show()
}
