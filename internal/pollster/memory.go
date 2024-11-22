package pollster

import (
	"bufio"
	"fmt"
	"log"
	"metricly/pkg/common"
	"os"
	"strings"
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
	memInfo, err := os.Open("/proc/meminfo")
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
		case "HugePagesTotal":
			memStats.HugePagesTotal = value
		case "HugePagesFree":
			memStats.HugePagesFree = value
		case "HugePagesRsvd":
			memStats.HugePagesRsvd = value
		case "HugePagesSurp":
			memStats.HugePagesSurp = value
		}
	}

	if err := scanner.Err(); err != nil {
		return memoryStats{}, fmt.Errorf("failed to parse /proc/meminfo %v", err)
	}

	return memStats, nil
}

func ReportMemoryUsage(mc *MetricCollector) {

	memStats, err := readMemoryStats()
	if err != nil {
		log.Println(err)
		return
	}
	labels := make(map[string]string)

	mc.UpdateMetric(
		"memory_total_bytes",
		float64(memStats.MemTotal),
		"Total memory bytes",
		labels,
	)

	mc.UpdateMetric(
		"memory_free_bytes",
		float64(memStats.MemFree),
		"Free memory bytes",
		labels,
	)

	mc.UpdateMetric(
		"memory_available_bytes",
		float64(memStats.MemAvailabe),
		"Total memory bytes",
		labels,
	)

	mc.UpdateMetric(
		"memory_hugepages_total",
		float64(memStats.HugePagesTotal),
		"Total memory bytes",
		labels,
	)

	mc.UpdateMetric(
		"memory_hugepages_free",
		float64(memStats.HugePagesFree),
		"Total memory bytes",
		labels,
	)

	mc.UpdateMetric(
		"memory_hugepages_rsvd",
		float64(memStats.HugePagesRsvd),
		"Total memory bytes",
		labels,
	)

	mc.UpdateMetric(
		"memory_hugepages_surp",
		float64(memStats.HugePagesSurp),
		"Total memory bytes",
		labels,
	)

}
