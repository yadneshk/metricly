package server

import (
	"context"
	"metricly/internal/pollster"
	"time"
)

func StartMetricsCollection(ctx context.Context, interval time.Duration, cc *pollster.MetriclyCollector) {

	
	pollster.RegisterCPUMetrics(cc)
	pollster.RegisterNetworkMetrics(cc)
	pollster.RegisterMemoryMetrics(cc)

	// Helper function to periodically execute metric reporting
	startPolling := func(reportFunc func(*pollster.MetriclyCollector)) {
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
	startPolling(pollster.ReportCpuUsage)
	startPolling(pollster.ReportMemoryUsage)
	startPolling(pollster.ReportNetworkUsage)
}
