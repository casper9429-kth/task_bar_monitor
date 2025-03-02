package metrics

import (
	"fmt"

	"github.com/shirou/gopsutil/v3/mem"
)

// GetMemoryUsage returns the current memory usage as a percentage
func GetMemoryUsage() (float64, error) {
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		return 0, err
	}

	return memInfo.UsedPercent, nil
}

// GetMemoryUsageDetails returns detailed memory usage information
func GetMemoryUsageDetails() (string, error) {
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		return "", err
	}

	used := memInfo.Used / 1024 / 1024   // Convert to MB
	total := memInfo.Total / 1024 / 1024 // Convert to MB

	return fmt.Sprintf("%d MB / %d MB (%.1f%%)", used, total, memInfo.UsedPercent), nil
}
