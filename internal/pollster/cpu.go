package pollster

import (
	"bufio"
	"fmt"
	"log"
	"metricly/pkg/common"
	"os"
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

var (
	procStat = "/proc/stat"
)

// ReadCPUStats reads CPU statistics from /proc/stat
func readCPUStats() (cpuUsage, error) {
	file, err := os.Open(procStat)
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
			usage.User = common.ParseUint(fields[1])
			usage.Total += usage.User
			usage.Nice = common.ParseUint(fields[2])
			usage.Total += usage.Nice
			usage.System = common.ParseUint(fields[3])
			usage.Total += usage.System
			usage.Idle = common.ParseUint(fields[4])
			usage.Total += usage.Idle
			usage.Iowait = common.ParseUint(fields[5])
			usage.Total += usage.Iowait
			usage.Irq = common.ParseUint(fields[6])
			usage.Total += usage.Irq
			usage.Softirq = common.ParseUint(fields[7])
			usage.Total += usage.Softirq
			usage.Steal = common.ParseUint(fields[8])
			usage.Total += usage.Steal

			return usage, nil
		}
	}
	return cpuUsage{}, fmt.Errorf("cpu stats not found in /proc/stat")
}

func truncate(value float64) float64 {
	return float64(int(value*100)) / 100
}

// CalculateCPUUsage calculates the CPU usage percentage
func calculateTotalUsage(prev, curr cpuUsage) float64 {

	totalDelta := curr.Total - prev.Total
	idleDelta := curr.Idle - prev.Idle

	if totalDelta == 0 {
		return 0.0
	}

	return truncate(100.0 * float64(totalDelta-idleDelta) / float64(totalDelta))
}

func calculateUserUsage(prev, curr cpuUsage) float64 {

	totalDelta := curr.Total - prev.Total
	userDelta := curr.User - prev.User

	if totalDelta == 0 {
		return 0.0
	}

	return truncate(100.0 * float64(userDelta) / float64(totalDelta))
}

func calculateSystemUsage(prev, curr cpuUsage) float64 {

	totalDelta := curr.Total - prev.Total
	systemDelta := curr.System - prev.System

	if totalDelta == 0 {
		return 0.0
	}

	return truncate(100.0 * float64(systemDelta) / float64(totalDelta))
}

func calculateStealUsage(prev, curr cpuUsage) float64 {

	totalDelta := curr.Total - prev.Total
	stealDelta := curr.Steal - prev.Steal

	if totalDelta == 0 {
		return 0.0
	}

	return truncate(100.0 * float64(stealDelta) / float64(totalDelta))
}

func RegisterCPUMetrics(mc *MetriclyCollector) {
	mc.addMetric("cpu_total", "CPU usage percentage", []string{"hostname"})
	mc.addMetric("cpu_user", "User process CPU usage percentage", []string{"hostname"})
	mc.addMetric("cpu_system", "System process CPU usage percentage", []string{"hostname"})
	mc.addMetric("cpu_steal", "CPU steal percentage", []string{"hostname"})
}

// collectCPUUsage collects the CPU usage as a percentage over a defined time interval.
func ReportCpuUsage(mc *MetriclyCollector) {
	log.Println("Polling CPU metrics...")

	// Capture initial CPU stats
	prevCPU, _ := readCPUStats()
	// log.Println(conf.CollectionInterval)
	time.Sleep(10 * time.Second)

	// Capture current CPU stats
	currCPU, _ := readCPUStats()

	mc.updateMetric("cpu_total", calculateTotalUsage(prevCPU, currCPU), []string{common.GetHostname()})
	mc.updateMetric("cpu_user", calculateUserUsage(prevCPU, currCPU), []string{common.GetHostname()})
	mc.updateMetric("cpu_system", calculateSystemUsage(prevCPU, currCPU), []string{common.GetHostname()})
	mc.updateMetric("cpu_steal", calculateStealUsage(prevCPU, currCPU), []string{common.GetHostname()})

	log.Println("Polling CPU metrics complete")
}
