package server

import (
	"context"
	collector "metricly/internal/collector"
	cpu "metricly/internal/pollster/cpu"
	disk "metricly/internal/pollster/disk"
	memory "metricly/internal/pollster/memory"
	network "metricly/internal/pollster/network"
	"time"
)

func StartMetricsCollection(ctx context.Context, interval time.Duration, cc *collector.MetriclyCollector) {

	cpu.RegisterCPUMetrics(cc)
	network.RegisterNetworkMetrics(cc)
	memory.RegisterMemoryMetrics(cc)
	disk.RegisterDiskMetrics(cc)

	// Helper function to periodically execute metric reporting
	startPolling := func(reportFunc func(*collector.MetriclyCollector)) {
		go func() {
			ticker := time.NewTicker(interval)
			defer ticker.Stop()
			for {
				select {
				case <-ctx.Done():
					return
				case <-ticker.C:
					reportFunc(cc)
				}
			}
		}()
	}

	// Start collectors for CPU, memory, and network metrics
	startPolling(cpu.ReportCpuUsage)
	startPolling(memory.ReportMemoryUsage)
	startPolling(network.ReportNetworkUsage)
	startPolling(disk.ReportDiskUsage)
}
