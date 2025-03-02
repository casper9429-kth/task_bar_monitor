package metrics

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	lastRxBytes    uint64
	lastTxBytes    uint64
	lastSampleTime time.Time
	networkMutex   sync.Mutex
)

// NetworkUsage contains network usage statistics
type NetworkUsage struct {
	DownloadSpeed float64 // Download speed in KB/s
	UploadSpeed   float64 // Upload speed in KB/s
	TotalDownload uint64  // Total downloaded in KB
	TotalUpload   uint64  // Total uploaded in KB
}

// GetNetworkUsage returns the current network usage statistics
func GetNetworkUsage() (NetworkUsage, error) {
	networkMutex.Lock()
	defer networkMutex.Unlock()

	// Read /proc/net/dev for network statistics
	file, err := os.Open("/proc/net/dev")
	if err != nil {
		return NetworkUsage{}, err
	}
	defer file.Close()

	var rxBytes, txBytes uint64 = 0, 0
	scanner := bufio.NewScanner(file)

	// Skip the header lines
	scanner.Scan() // skip first line
	scanner.Scan() // skip second line

	// Process each network interface line
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(strings.TrimSpace(line))

		// Skip loopback interface
		if strings.HasPrefix(fields[0], "lo:") {
			continue
		}

		// Parse received and transmitted bytes
		rx, err := strconv.ParseUint(fields[1], 10, 64)
		if err != nil {
			continue
		}

		tx, err := strconv.ParseUint(fields[9], 10, 64)
		if err != nil {
			continue
		}

		rxBytes += rx
		txBytes += tx
	}

	now := time.Now()

	// Calculate speeds
	var downloadSpeed, uploadSpeed float64
	timeDiff := now.Sub(lastSampleTime).Seconds()

	if lastRxBytes > 0 && lastTxBytes > 0 && timeDiff > 0 {
		// Convert to KB/s
		downloadSpeed = float64(rxBytes-lastRxBytes) / 1024 / timeDiff
		uploadSpeed = float64(txBytes-lastTxBytes) / 1024 / timeDiff
	}

	// Update last values
	lastRxBytes = rxBytes
	lastTxBytes = txBytes
	lastSampleTime = now

	// Return network usage statistics
	return NetworkUsage{
		DownloadSpeed: downloadSpeed,
		UploadSpeed:   uploadSpeed,
		TotalDownload: rxBytes / 1024,
		TotalUpload:   txBytes / 1024,
	}, nil
}

// GetNetworkUsageString returns a formatted string with network usage information
func GetNetworkUsageString() (string, error) {
	usage, err := GetNetworkUsage()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("↓ %.1f KB/s  ↑ %.1f KB/s", usage.DownloadSpeed, usage.UploadSpeed), nil
}
