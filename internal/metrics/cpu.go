package metrics

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/v3/cpu" // Updated import path
)

// GetCPUUsage returns the current CPU usage as a percentage
func GetCPUUsage() (float64, error) {
	percent, err := cpu.Percent(time.Second, false)
	if err != nil {
		return 0, err
	}

	if len(percent) == 0 {
		return 0, fmt.Errorf("no CPU usage data available")
	}

	return percent[0], nil
}
