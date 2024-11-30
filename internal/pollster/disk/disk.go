package disk

import (
	"bufio"
	"fmt"
	"log/slog"
	collector "metricly/internal/collector"
	"metricly/pkg/common"
	"os"
	"strings"
	"syscall"
	"time"
)

var (
	procDiskStats = "/proc/diskstats"
	procMounts    = "/proc/mounts"
)

// diskStats holds metrics for a single disk.
type diskStats struct {
	ReadsCompleted        uint64
	SectorsRead           uint64
	WriteCompleted        uint64
	SectorsWritten        uint64
	IOInProgress          uint64
	IOTimeSpentMillis     uint64
	WeightedIOTimeSpentMs uint64
	ReadThroughputBytes   uint64 // Calculated in bytes
	WriteThroughputBytes  uint64 // Calculated in bytes
}

type diskSpaceStat struct {
	Total     uint64  // Total disk space in bytes
	Used      uint64  // Used disk space in bytes
	Available uint64  // Available disk space in bytes
	Usage     float64 // Usage percentage
}

// parseDiskStats parses /proc/diskstats for metrics.
func parseDiskStats() (map[string]diskStats, error) {

	if procDiskStatsEnv := os.Getenv("PROC_DISK_STATS"); procDiskStatsEnv != "" {
		procDiskStats = procDiskStatsEnv
	}

	file, err := os.Open(procDiskStats)
	if err != nil {
		return nil, fmt.Errorf("failed to open %s: %v", procDiskStats, err)
	}
	defer file.Close()

	diskStatsMap := make(map[string]diskStats)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())

		if len(fields) < 14 {
			// Diskstats file must have at least 14 fields.
			continue
		}

		// Parse disk name and stats
		deviceName := fields[2]
		readCompleted := common.ParseUint(fields[3])
		sectorsRead := common.ParseUint(fields[5])
		writeCompleted := common.ParseUint(fields[7])
		sectorsWritten := common.ParseUint(fields[9])
		ioInProgress := common.ParseUint(fields[11])
		ioTimeSpentMillis := common.ParseUint(fields[12])
		weightedIOTimeSpentMs := common.ParseUint(fields[13])

		// Convert sectors to bytes (1 sector = 512 bytes)
		readThroughputBytes := sectorsRead * 512
		writeThroughputBytes := sectorsWritten * 512

		// Add metrics to diskStatsMap
		diskStatsMap[deviceName] = diskStats{
			ReadsCompleted:        readCompleted,
			SectorsRead:           sectorsRead,
			WriteCompleted:        writeCompleted,
			SectorsWritten:        sectorsWritten,
			IOInProgress:          ioInProgress,
			IOTimeSpentMillis:     ioTimeSpentMillis,
			WeightedIOTimeSpentMs: weightedIOTimeSpentMs,
			ReadThroughputBytes:   readThroughputBytes,
			WriteThroughputBytes:  writeThroughputBytes,
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to parse %s: %v", procDiskStats, err)
	}

	return diskStatsMap, nil
}

// ReadDiskSpaceStats retrieves disk space statistics for the specified mount point
func readDiskSpaceStats(mountPoint []string) (map[string]diskSpaceStat, error) {
	var stat syscall.Statfs_t
	diskSpaceMap := make(map[string]diskSpaceStat)

	for _, mount := range mountPoint {
		if err := syscall.Statfs(mount, &stat); err != nil {
			slog.Warn(fmt.Sprintf("failed to retrieve disk space stats for %s: %v", mountPoint, err))
			continue
		}

		total := stat.Blocks * uint64(stat.Bsize)
		available := stat.Bavail * uint64(stat.Bsize)
		used := total - (stat.Bfree * uint64(stat.Bsize))
		usage := float64(used) / float64(total) * 100

		diskSpaceMap[mount] = diskSpaceStat{
			Total:     total,
			Available: available,
			Used:      used,
			Usage:     usage,
		}

	}

	return diskSpaceMap, nil

}

// GetMountPoints retrieves a list of mount points from /proc/mounts
func getMountPoints() ([]string, error) {

	if procMountsEnv := os.Getenv("PROC_DISK_MOUNTS"); procMountsEnv != "" {
		procMounts = procMountsEnv
	}

	file, err := os.Open(procMounts)
	if err != nil {
		return nil, fmt.Errorf("failed to open %s: %v", procMounts, err)
	}
	defer file.Close()

	var mountPoints []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}

		// The second field in each line represents the mount point
		mountPoint := fields[1]

		// Filter out pseudo-filesystems (optional)
		if strings.HasPrefix(fields[2], "tmpfs") || strings.HasPrefix(fields[2], "proc") {
			continue
		}

		mountPoints = append(mountPoints, mountPoint)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading /proc/mounts: %v", err)
	}

	return mountPoints, nil
}

// RegisterDiskMetrics registers disk metrics.
func RegisterDiskMetrics(mc *collector.MetriclyCollector) {
	mc.AddMetric("disk_reads_completed_total", "Total disk reads completed", []string{"device"})
	mc.AddMetric("disk_writes_completed_total", "Total disk writes completed", []string{"device"})
	mc.AddMetric("disk_read_throughput_bytes", "Disk read throughput in bytes", []string{"device"})
	mc.AddMetric("disk_write_throughput_bytes", "Disk write throughput in bytes", []string{"device"})
	mc.AddMetric("disk_io_in_progress", "Current disk IO operations in progress", []string{"device"})
	mc.AddMetric("disk_io_time_spent_seconds", "Time spent on IO operations in seconds", []string{"device"})
	mc.AddMetric("disk_weighted_io_time_seconds", "Weighted time spent on IO in seconds", []string{"device"})
	mc.AddMetric("disk_total_bytes", "Total disk space in bytes", []string{"mount_point"})
	mc.AddMetric("disk_used_bytes", "Used disk space in bytes", []string{"mount_point"})
	mc.AddMetric("disk_available_bytes", "Available disk space in bytes", []string{"mount_point"})
	mc.AddMetric("disk_usage_percentage", "Disk usage percentage", []string{"mount_point"})
}

// ReportDiskMetrics reports disk metrics periodically.
func ReportDiskUsage(mc *collector.MetriclyCollector) {
	start := time.Now()
	// get disk I/O usage
	diskStatsMap, err := parseDiskStats()
	if err != nil {
		fmt.Printf("Error reading disk stats: %v\n", err)
		return
	}

	for device, stats := range diskStatsMap {
		mc.UpdateMetric(
			"disk_reads_completed_total",
			float64(stats.ReadsCompleted),
			[]string{device},
		)

		mc.UpdateMetric(
			"disk_writes_completed_total",
			float64(stats.WriteCompleted),
			[]string{device},
		)

		mc.UpdateMetric(
			"disk_read_throughput_bytes",
			float64(stats.ReadThroughputBytes),
			[]string{device},
		)

		mc.UpdateMetric(
			"disk_write_throughput_bytes",
			float64(stats.WriteThroughputBytes),
			[]string{device},
		)

		mc.UpdateMetric(
			"disk_io_in_progress",
			float64(stats.IOInProgress),
			[]string{device},
		)

		mc.UpdateMetric(
			"disk_io_time_spent_seconds",
			float64(stats.IOTimeSpentMillis)/1000.0,
			[]string{device},
		)

		mc.UpdateMetric(
			"disk_weighted_io_time_seconds",
			float64(stats.WeightedIOTimeSpentMs)/1000.0,
			[]string{device},
		)
	}

	// get disk space usage
	mounts, err := getMountPoints()
	if err != nil {
		slog.Warn(fmt.Sprintf("failed to retrieve disk mounts: %s", err))
		return
	}

	diskSpaceStats, err := readDiskSpaceStats(mounts)
	if err != nil {
		slog.Warn(fmt.Sprintf("failed to retrieve disk stats: %s", err))
		return
	}
	for mount, stats := range diskSpaceStats {
		mc.UpdateMetric(
			"disk_total_bytes",
			float64(stats.Total),
			[]string{mount},
		)
		mc.UpdateMetric(
			"disk_used_bytes",
			float64(stats.Used),
			[]string{mount},
		)
		mc.UpdateMetric(
			"disk_available_bytes",
			float64(stats.Available),
			[]string{mount},
		)
		mc.UpdateMetric(
			"disk_usage_percentage",
			float64(stats.Usage),
			[]string{mount},
		)
	}
	slog.Info(fmt.Sprintf("Collected Disk metrics in %s", time.Since(start)))
}
