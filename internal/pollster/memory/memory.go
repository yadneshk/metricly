package memory

import (
	"bufio"
	"fmt"
	"log/slog"
	collector "metricly/internal/collector"
	"metricly/pkg/common"
	"os"
	"strings"
	"time"
)

var (
	procMemInfo = "/proc/meminfo"
)

type memoryStats struct {
	MemTotal       uint64
	MemFree        uint64
	MemAvailabe    uint64
	HugePagesTotal uint64
	HugePagesFree  uint64
	HugePagesRsvd  uint64
	HugePagesSurp  uint64
}

func readMemoryStats() (memoryStats, error) {

	if procMemInfoEnv := os.Getenv("PROC_MEMORY_INFO"); procMemInfoEnv != "" {
		procMemInfo = procMemInfoEnv
	}

	memInfo, err := os.Open(procMemInfo)
	if err != nil {
		return memoryStats{}, err
	}
	defer memInfo.Close()
	var memStats memoryStats
	scanner := bufio.NewScanner(memInfo)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)

		// the smallest slice is of 2 elements ["HugePages_Total:", "0"]
		if len(fields) < 2 {
			continue
		}

		key := strings.TrimSuffix(fields[0], ":")
		value := common.ParseUint(fields[1])

		if len(fields) > 2 && fields[2] == "kB" {
			// convert value to bytes
			value *= 1024
		}

		switch key {
		case "MemTotal":
			memStats.MemTotal = value
		case "MemFree":
			memStats.MemFree = value
		case "MemAvailable":
			memStats.MemAvailabe = value
		case "HugePages_Total":
			memStats.HugePagesTotal = value
		case "HugePages_Free":
			memStats.HugePagesFree = value
		case "HugePages_Rsvd":
			memStats.HugePagesRsvd = value
		case "HugePages_Surp":
			memStats.HugePagesSurp = value
		}
	}

	if err := scanner.Err(); err != nil {
		return memoryStats{}, fmt.Errorf("failed to parse /proc/meminfo %v", err)
	}

	return memStats, nil
}

func RegisterMemoryMetrics(mc *collector.MetriclyCollector) {
	// constLabelMap := make(map[string]string)
	mc.AddMetric("memory_total_bytes", "Total memory usage", []string{})
	mc.AddMetric("memory_free_bytes", "Free memory", []string{})
	mc.AddMetric("memory_available_bytes", "available memory", []string{})
	mc.AddMetric("memory_hugepages_total", "Total number of hugepages", []string{})
	mc.AddMetric("memory_hugepages_free", "Free hugepages", []string{})
	mc.AddMetric("memory_hugepages_rsvd", "Reserved hugepages", []string{})
	mc.AddMetric("memory_hugepages_surp", "Surplus hugepages", []string{})
}

func ReportMemoryUsage(mc *collector.MetriclyCollector) {
	start := time.Now()
	memStats, err := readMemoryStats()
	if err != nil {
		slog.Warn(fmt.Sprint(err))
		return
	}

	mc.UpdateMetric(
		"memory_total_bytes",
		float64(memStats.MemTotal),
		[]string{},
	)

	mc.UpdateMetric(
		"memory_free_bytes",
		float64(memStats.MemFree),
		[]string{},
	)

	mc.UpdateMetric(
		"memory_available_bytes",
		float64(memStats.MemAvailabe),
		[]string{},
	)

	mc.UpdateMetric(
		"memory_hugepages_total",
		float64(memStats.HugePagesTotal),
		[]string{},
	)

	mc.UpdateMetric(
		"memory_hugepages_free",
		float64(memStats.HugePagesFree),
		[]string{},
	)

	mc.UpdateMetric(
		"memory_hugepages_rsvd",
		float64(memStats.HugePagesRsvd),
		[]string{},
	)

	mc.UpdateMetric(
		"memory_hugepages_surp",
		float64(memStats.HugePagesSurp),
		[]string{},
	)
	slog.Info(fmt.Sprintf("Collected Memory metrics in %s", time.Since(start)))
}
