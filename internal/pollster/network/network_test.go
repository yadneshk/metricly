package network

import (
	"fmt"
	"log"
	pollster "metricly/internal/collector"
	helper "metricly/internal/pollster/tests"
	"metricly/pkg/common"
	"os"
	"testing"
	"time"
)

func TestReadNetworkStats(t *testing.T) {

	mntContent := `Inter-|   Receive                                                |  Transmit
 face |bytes    packets errs drop fifo frame compressed multicast|bytes    packets errs drop fifo colls carrier compressed
    lo: 266527100  184168    0    0    0     0          0         0 266527100  184168    0    0    0     0       0          0
wlp0s20f3: 6540158835 5786650    0    2    0     0          0         0 1421100604 2521626    0  278    0     0       0          0`
	collectorSource := "network.txt"

	err := helper.SetupCollectorSources(collectorSource, mntContent)
	if err != nil {
		t.Fatalf("failed to setup collector file: %v", err)
	}
	defer os.Remove(collectorSource)
	procNetDev = collectorSource

	// start testing target function
	stats, err := readNetworkStats()
	if err != nil {
		t.Fatalf("Failed to read network stats: %v", err)
	}

	// Validate the parsed results
	if stats[0].bytesRx != 266527100 {
		t.Errorf("Expected eth0 BytesReceived = 266527100, got %d", stats[0].bytesRx)
	}
	if stats[0].bytesTx != 266527100 {
		t.Errorf("Expected eth0 BytesTransmitted = 266527100, got %d", stats[0].bytesTx)
	}
	if stats[1].packetsRx != 5786650 {
		t.Errorf("Expected wlan0 PacketsReceived = 5786650, got %d", stats[1].packetsRx)
	}
	if stats[1].errorsRx != 0 {
		t.Errorf("Expected wlan0 ErrorsReceived = 1, got %d", stats[1].errorsRx)
	}
}

func TestReportNetworkUsage(t *testing.T) {
	mntContent := `Inter-|   Receive                                                |  Transmit
 face |bytes    packets errs drop fifo frame compressed multicast|bytes    packets errs drop fifo colls carrier compressed
    lo: 266527100  184168    0    0    0     0          0         0 266527100  184168    0    0    0     0       0          0
wlp0s20f3: 6540158835 5786650    0    2    0     0          0         0 1421100604 2521626    0  278    0     0       0          0`
	collectorSource := "network.txt"

	err := helper.SetupCollectorSources(collectorSource, mntContent)
	if err != nil {
		t.Fatalf("failed to setup collector file: %v", err)
	}
	// defer os.Remove(collectorSource)
	procNetDev = collectorSource

	mc := pollster.CreateMetricCollector()
	RegisterNetworkMetrics(mc)

	go func() {
		time.Sleep(time.Millisecond * 5)
		fi, _ := os.OpenFile(collectorSource, os.O_WRONLY|os.O_TRUNC, 0644)
		mntContent := `Inter-|   Receive                                                |  Transmit
		face |bytes    packets errs drop fifo frame compressed multicast|bytes    packets errs drop fifo colls carrier compressed
		   lo: 266527100  184168    0    0    0     0          0         0 266527100  184168    0    0    0     0       0          0
	   wlp0s20f3: 6540158935 5786650    2    2    0     0          0         0 1421100604 2521626    0  278    0     0       0          0`
		if _, err := fi.Write([]byte(mntContent)); err != nil {
			log.Fatalf("failed to update mock stat file: %v", err)
		}
		fi.Close()
	}()
	ReportNetworkUsage(mc)

	if metric, exists := mc.Data[fmt.Sprintf("network_rx_bytes_total|%s|wlp0s20f3", common.GetHostname())]; exists {
		if metric.Value != 100 {
			t.Errorf("expected metric network_rx_bytes_total=100, got %f", metric.Value)
		}
	}

	if metric, exists := mc.Data[fmt.Sprintf("network_rx_errors_total|%s|wlp0s20f3", common.GetHostname())]; exists {
		if metric.Value != 2 {
			t.Errorf("expected metric network_rx_errors_total=2, got %f", metric.Value)
		}

	}

}
