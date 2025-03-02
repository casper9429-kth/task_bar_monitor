package metrics

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/shirou/gopsutil/v3/disk"
)

// GetDiskUsage returns the current disk usage as a percentage
func GetDiskUsage() (float64, error) {
	usage, err := disk.Usage("/")
	if err != nil {
		return 0, err
	}

	return usage.UsedPercent, nil
}

// GetDiskUsageDetails returns detailed disk usage information
func GetDiskUsageDetails() (string, error) {
	cmd := exec.Command("df", "-h", "/")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	lines := strings.Split(string(output), "\n")
	if len(lines) < 2 {
		return "", fmt.Errorf("unexpected df output format")
	}

	fields := strings.Fields(lines[1])
	if len(fields) < 5 {
		return "", fmt.Errorf("unexpected df output fields")
	}

	return fmt.Sprintf("%s Used, %s Total", fields[2], fields[1]), nil
}

// GetAvailableDiskSpace returns the available disk space in human readable format
func GetAvailableDiskSpace() (string, error) {
	usage, err := disk.Usage("/")
	if err != nil {
		return "", err
	}

	// Convert to GB
	available := float64(usage.Free) / 1024 / 1024 / 1024
	total := float64(usage.Total) / 1024 / 1024 / 1024

	return fmt.Sprintf("%.1f GB free of %.1f GB", available, total), nil
}
