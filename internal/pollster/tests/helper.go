package tests

import (
	"fmt"
	collector "metricly/internal/collector"
	"os"
	"testing"
)

func SetupCollectorSources(fileName, fileContent string) error {

	colltrFile, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("failed to create mounts file: %s", err)
	}

	if _, err = colltrFile.Write([]byte(fileContent)); err != nil {
		return fmt.Errorf("failed to write mock content: %s", err)
	}
	colltrFile.Close()

	return nil
}

func VerifyMetric(t *testing.T, mc *collector.MetriclyCollector, metricName string, metricValue float64) {

	// Validate metrics
	if metric, exists := mc.Data[metricName]; exists {
		if metric.Value != metricValue {
			t.Fatalf("expected %s=%f, got %f", metricName, metricValue, metric.Value)
		}
	} else {
		t.Fatalf("%s not found in metrics data", metricName)
	}
}
