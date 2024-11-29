package memory

import (
	pollster "metricly/internal/collector"
	helper "metricly/internal/pollster/tests"
	"os"
	"testing"
)

func TestReadMemoryStats(t *testing.T) {

	collectorSource := "meminfo.txt"
	mntContent := `MemTotal:       16384000 kB
MemFree:        8192000 kB
MemAvailable:   12288000 kB
HugePages_Total:       64
HugePages_Free:        32
HugePages_Rsvd:        16
HugePages_Surp:         8`
	err := helper.SetupCollectorSources(collectorSource, mntContent)
	if err != nil {
		t.Fatalf("failed to setup collector file: %v", err)
	}
	defer os.Remove(collectorSource)
	procMemInfo = collectorSource

	memStats, err := readMemoryStats()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Validate the parsed values
	if memStats.MemTotal != 16384000*1024 {
		t.Errorf("expected MemTotal=16384000, got %d", memStats.MemTotal)
	}
	if memStats.MemFree != 8192000*1024 {
		t.Errorf("expected MemFree=8192000, got %d", memStats.MemFree)
	}
	if memStats.MemAvailabe != 12288000*1024 {
		t.Errorf("expected MemAvailable=12288000, got %d", memStats.MemAvailabe)
	}
	if memStats.HugePagesTotal != 64 {
		t.Errorf("expected HugePagesTotal=64, got %d", memStats.HugePagesTotal)
	}
	if memStats.HugePagesFree != 32 {
		t.Errorf("expected HugePagesFree=32, got %d", memStats.HugePagesFree)
	}
	if memStats.HugePagesRsvd != 16 {
		t.Errorf("expected HugePagesRsvd=16, got %d", memStats.HugePagesRsvd)
	}
	if memStats.HugePagesSurp != 8 {
		t.Errorf("expected HugePagesSurp=8, got %d", memStats.HugePagesSurp)
	}
}

func TestReportMemoryUsage(t *testing.T) {

	collectorSource := "meminfo.txt"
	mntContent := `MemTotal:       16384000 kB
MemFree:        8192000 kB
MemAvailable:   12288000 kB
HugePages_Total:       64
HugePages_Free:        32
HugePages_Rsvd:        16
HugePages_Surp:         8`
	err := helper.SetupCollectorSources(collectorSource, mntContent)
	if err != nil {
		t.Fatalf("failed to setup collector file: %v", err)
	}
	defer os.Remove(collectorSource)
	procMemInfo = collectorSource

	// start testing target function
	mc := pollster.CreateMetricCollector()
	RegisterMemoryMetrics(mc)

	ReportMemoryUsage(mc)

	helper.VerifyMetric(t, mc, "metricly_memory_total_bytes", 16384000*1024)
	helper.VerifyMetric(t, mc, "metricly_memory_free_bytes", 8192000*1024)
	helper.VerifyMetric(t, mc, "metricly_memory_available_bytes", 12288000*1024)
	helper.VerifyMetric(t, mc, "metricly_memory_hugepages_total", 64)

}
