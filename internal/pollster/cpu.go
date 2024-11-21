package pollster

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type cpuUsage struct {
	User    uint64
	Nice    uint64
	System  uint64
	Idle    uint64
	Iowait  uint64
	Irq     uint64
	Softirq uint64
	Steal   uint64
	Total   uint64
}

// ReadCPUStats reads CPU statistics from /proc/stat
func readCPUStats() (cpuUsage, error) {
	file, err := os.Open("/proc/stat")
	if err != nil {
		return cpuUsage{}, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "cpu ") {
			// the first line aggregates the numbers in all of the other “cpuN” lines
			fields := strings.Fields(line)
			if len(fields) < 8 {
				// the first 8 fields are important to calculate the total time spent by all CPUs
				return cpuUsage{}, fmt.Errorf("unexpected format in /proc/stat")
			}

			usage := cpuUsage{}
			usage.User = parseUint(fields[1])
			usage.Total += usage.User
			usage.Nice = parseUint(fields[2])
			usage.Total += usage.Nice
			usage.System = parseUint(fields[3])
			usage.Total += usage.System
			usage.Idle = parseUint(fields[4])
			usage.Total += usage.Idle
			usage.Iowait = parseUint(fields[5])
			usage.Total += usage.Iowait
			usage.Irq = parseUint(fields[6])
			usage.Total += usage.Irq
			usage.Softirq = parseUint(fields[7])
			usage.Total += usage.Softirq
			usage.Steal = parseUint(fields[8])
			usage.Total += usage.Steal

			return usage, nil
		}
	}
	return cpuUsage{}, fmt.Errorf("cpu stats not found in /proc/stat")
}

// parseUint safely parses a string to uint64
func parseUint(s string) uint64 {
	value, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0
	}
	return value
}

// CalculateCPUUsage calculates the CPU usage percentage
func calculateTotalUsage(prev, curr cpuUsage) float64 {

	totalDelta := curr.Total - prev.Total
	idleDelta := curr.Idle - prev.Idle

	if totalDelta == 0 {
		return 0.0
	}

	return 100.0 * float64(totalDelta-idleDelta) / float64(totalDelta)
}

func calculateUserUsage(prev, curr cpuUsage) float64 {

	totalDelta := curr.Total - prev.Total
	userDelta := curr.User - prev.User

	if totalDelta == 0 {
		return 0.0
	}

	return 100.0 * float64(userDelta) / float64(totalDelta)
}

func calculateSystemUsage(prev, curr cpuUsage) float64 {

	totalDelta := curr.Total - prev.Total
	systemDelta := curr.System - prev.System

	if totalDelta == 0 {
		return 0.0
	}

	return 100.0 * float64(systemDelta) / float64(totalDelta)
}

func calculateStealUsage(prev, curr cpuUsage) float64 {

	totalDelta := curr.Total - prev.Total
	stealDelta := curr.Steal - prev.Steal

	if totalDelta == 0 {
		return 0.0
	}

	return 100.0 * float64(stealDelta) / float64(totalDelta)
}

// collectCPUUsage collects the CPU usage as a percentage over a defined time interval.
func ReportCpuUsage(mc *MetricCollector) {
	// Capture initial CPU stats
	prevCPU, _ := readCPUStats()

	// Wait for a small interval to calculate the usage delta
	time.Sleep(1 * time.Second)

	// Capture current CPU stats
	currCPU, _ := readCPUStats()

	// Calculate CPU usage percentage
	mc.UpdateMetric(
		"cpu_total",
		calculateTotalUsage(prevCPU, currCPU),
		"CPU usage percentage",
		map[string]string{"hostname": "mynode"},
	)

	// Calculate user CPU usage percentage
	mc.UpdateMetric(
		"cpu_user",
		calculateUserUsage(prevCPU, currCPU),
		"User process CPU usage percentage",
		map[string]string{"hostname": "mynode"},
	)

	// Calculate system (kernel level) CPU usage percentage
	mc.UpdateMetric(
		"cpu_system",
		calculateSystemUsage(prevCPU, currCPU),
		"System process CPU usage percentage",
		map[string]string{"hostname": "mynode"},
	)

	// Calculate steal percentage
	mc.UpdateMetric(
		"cpu_steal",
		calculateStealUsage(prevCPU, currCPU),
		"CPU steal percentage",
		map[string]string{"hostname": "mynode"},
	)

}
