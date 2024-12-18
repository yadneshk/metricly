package cpu

import (
	collector "metricly/internal/collector"
	helper "metricly/internal/pollster/tests"
	"os"
	"testing"
)

func TestReadCpuStats(t *testing.T) {
	collectorSource := "cpustat.txt"

	mntContent := `cpu  2255 34 2290 22625563 6290 127 456 0 0 0
cpu0 1132 17 1145 11312780 3154 63 228 0 0 0
cpu1 1123 17 1145 11312783 3154 63 228 0 0 0`

	err := helper.SetupCollectorSources(collectorSource, mntContent)
	if err != nil {
		t.Fatalf("failed to setup collector file: %v", err)
	}
	defer os.Remove(collectorSource)
	procStat = collectorSource

	cpuStats, err := readCPUStats()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cpuStats.User != 2255 {
		t.Errorf("expected User=2255, got %d", cpuStats.User)
	}
	if cpuStats.Nice != 34 {
		t.Errorf("expected Nice=34, got %d", cpuStats.Nice)
	}
	if cpuStats.System != 2290 {
		t.Errorf("expected System=2290, got %d", cpuStats.System)
	}
	if cpuStats.Idle != 22625563 {
		t.Errorf("expected Idle=22625563, got %d", cpuStats.Idle)
	}
	if cpuStats.Iowait != 6290 {
		t.Errorf("expected Iowait=6290, got %d", cpuStats.Iowait)
	}
	if cpuStats.Irq != 127 {
		t.Errorf("expected Irq=127, got %d", cpuStats.Irq)
	}
	if cpuStats.Softirq != 456 {
		t.Errorf("expected Softirq=456, got %d", cpuStats.Softirq)
	}
	if cpuStats.Steal != 0 {
		t.Errorf("expected Steal=0, got %d", cpuStats.Steal)
	}
	if cpuStats.Total != 22637015 {
		t.Errorf("expected Total=22637015, got %d", cpuStats.Total)
	}

}

func TestCalculateCPUUsage(t *testing.T) {
	prev := cpuUsage{
		User:   100,
		System: 200,
		Idle:   500,
		Total:  800,
	}
	curr := cpuUsage{
		User:   200,
		System: 300,
		Idle:   700,
		Total:  1400,
	}

	// Calculate and validate percentages
	totalUsage := calculateTotalUsage(prev, curr)

	if totalUsage != 66.66 {
		t.Errorf("expected totalUsage=66.666667, got %f", totalUsage)
	}

	userUsage := calculateUserUsage(prev, curr)
	if userUsage != 16.66 {
		t.Errorf("expected userUsage=16.67, got %.2f", userUsage)
	}

	systemUsage := calculateSystemUsage(prev, curr)
	if systemUsage != 16.66 {
		t.Errorf("expected systemUsage=16.67, got %.2f", systemUsage)
	}

	stealUsage := calculateStealUsage(prev, curr)
	if stealUsage != 0.0 {
		t.Errorf("expected stealUsage=0.0, got %.2f", stealUsage)
	}
}

func TestReportCpuUsage(t *testing.T) {
	collectorSource := "cpustats.txt"

	mntContent := `cpu  100 200 300 400 50 60 70 80 90
cpu0 50 100 150 200 25 30 35 40 45`

	err := helper.SetupCollectorSources(collectorSource, mntContent)
	if err != nil {
		t.Fatalf("failed to setup collector file: %v", err)
	}
	defer os.Remove(collectorSource)
	procStat = collectorSource

	prevCPU, _ = readCPUStats()

	mntContent = `cpu  200 300 400 500 60 70 80 90 100
cpu0 100 150 200 250 30 35 40 45 50`

	err = helper.SetupCollectorSources(collectorSource, mntContent)
	if err != nil {
		t.Fatalf("failed to setup collector file: %v", err)
	}

	mc := collector.CreateMetricCollector()
	RegisterCPUMetrics(mc)

	ReportCpuUsage(mc)

	helper.VerifyMetric(t, mc, "metricly_cpu_total", 77.27)
	helper.VerifyMetric(t, mc, "metricly_cpu_system", 22.72)
	helper.VerifyMetric(t, mc, "metricly_cpu_user", 22.72)
	helper.VerifyMetric(t, mc, "metricly_cpu_steal", 2.27)
}
