package pollster

import (
	"os"
	"testing"
)

func TestReadMemoryStats(t *testing.T) {
	// Step 1: Create a temporary file to simulate /proc/meminfo
	tmpFile, err := os.Create("meminfo.txt")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name()) // Clean up after the test

	// Step 2: Write mock data to the temporary file
	mockMemInfo := `MemTotal:       16384000 kB
MemFree:        8192000 kB
MemAvailable:   12288000 kB
HugePages_Total:       64
HugePages_Free:        32
HugePages_Rsvd:        16
HugePages_Surp:         8`

	if _, err := tmpFile.WriteString(mockMemInfo); err != nil {
		t.Fatalf("failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	// Step 3: Override the global `procMemInfo` variable
	originalProcMemInfo := procMemInfo
	procMemInfo = tmpFile.Name()
	defer func() { procMemInfo = originalProcMemInfo }() // Restore after test

	// Step 4: Call `readMemoryStats` and validate the results
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
	// Step 1: Create a temporary file to simulate /proc/meminfo
	tmpFile, err := os.Create("meminfo.txt")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name()) // Clean up after the test

	// Step 2: Write mock data to the temporary file
	mockMemInfo := `MemTotal:       16384000 kB
MemFree:        8192000 kB
MemAvailable:   12288000 kB
HugePages_Total:       64
HugePages_Free:        32
HugePages_Rsvd:        16
HugePages_Surp:         8`
	if _, err := tmpFile.WriteString(mockMemInfo); err != nil {
		t.Fatalf("failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	// Step 3: Override the global `procMemInfo` variable
	originalProcMemInfo := procMemInfo
	procMemInfo = tmpFile.Name()
	defer func() { procMemInfo = originalProcMemInfo }() // Restore after test

	// Step 4: Create a mock MetriclyCollector
	mc := CreateMetricCollector()
	RegisterMemoryMetrics(mc)

	// Step 5: Call `ReportMemoryUsage` and validate the updated metrics
	ReportMemoryUsage(mc)
	// fmt.Println(mc)
	// Validate the metrics
	if metric, exists := mc.data["metricly_memory_total_bytes"]; exists {
		if metric.value != 16384000*1024 {
			t.Errorf("expected memory_total_bytes=16384000, got %f", metric.value)
		}
	} else {
		t.Error("metricly_memory_total_bytes not found in metrics data")
	}

	if metric, exists := mc.data["metricly_memory_free_bytes"]; exists {
		if metric.value != 8192000*1024 {
			t.Errorf("expected memory_free_bytes=8192000, got %f", metric.value)
		}
	} else {
		t.Error("metricly_memory_free_bytes not found in metrics data")
	}

	if metric, exists := mc.data["metricly_memory_available_bytes"]; exists {
		if metric.value != 12288000*1024 {
			t.Errorf("expected memory_available_bytes=12288000, got %f", metric.value)
		}
	} else {
		t.Error("metricly_memory_available_bytes not found in metrics data")
	}

	if metric, exists := mc.data["metricly_memory_hugepages_total"]; exists {
		if metric.value != 64 {
			t.Errorf("expected memory_hugepages_total=64, got %f", metric.value)
		}
	} else {
		t.Error("metricly_memory_hugepages_total not found in metrics data")
	}
}
