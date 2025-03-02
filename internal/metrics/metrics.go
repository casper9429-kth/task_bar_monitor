package metrics

import "log"

// SystemMetrics holds all system metrics in one structure
type SystemMetrics struct {
	CPUUsage     float64
	MemoryUsage  float64
	DiskUsage    float64
	NetworkUsage NetworkUsage
}

// GetAllMetrics returns all system metrics in one call
func GetAllMetrics() (SystemMetrics, error) {
	metrics := SystemMetrics{}
	var err error

	// Get CPU usage
	metrics.CPUUsage, err = GetCPUUsage()
	if err != nil {
		log.Printf("Failed to get CPU usage: %v", err)
	}

	// Get memory usage
	metrics.MemoryUsage, err = GetMemoryUsage()
	if err != nil {
		log.Printf("Failed to get memory usage: %v", err)
	}

	// Get disk usage
	metrics.DiskUsage, err = GetDiskUsage()
	if err != nil {
		log.Printf("Failed to get disk usage: %v", err)
	}

	// Get network usage
	metrics.NetworkUsage, err = GetNetworkUsage()
	if err != nil {
		log.Printf("Failed to get network usage: %v", err)
	}

	return metrics, nil
}
