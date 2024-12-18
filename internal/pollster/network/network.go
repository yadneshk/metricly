package network

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
	procNetDev = "/proc/net/dev"
)

type networkStats struct {
	interfaceName string
	bytesRx       uint64
	bytesTx       uint64
	packetsRx     uint64
	packetsTx     uint64
	errorsRx      uint64
	errorsTx      uint64
	dropsRx       uint64
	dropsTx       uint64
}

func readNetworkStats() ([]networkStats, error) {

	if procNetDevEnv := os.Getenv("PROC_NETWORK_DEV"); procNetDevEnv != "" {
		procNetDev = procNetDevEnv
	}

	nwStats, err := os.Open(procNetDev)
	if err != nil {
		return []networkStats{}, err
	}

	var stats []networkStats
	scanner := bufio.NewScanner(nwStats)

	// skip first two info lines
	for i := 0; i < 2; i++ {
		scanner.Scan()
	}

	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)

		stats = append(stats, networkStats{
			interfaceName: strings.TrimSuffix(fields[0], ":"),
			bytesRx:       common.ParseUint(fields[1]),
			bytesTx:       common.ParseUint(fields[9]),
			packetsRx:     common.ParseUint(fields[2]),
			packetsTx:     common.ParseUint(fields[10]),
			errorsRx:      common.ParseUint(fields[3]),
			errorsTx:      common.ParseUint(fields[11]),
			dropsRx:       common.ParseUint(fields[4]),
			dropsTx:       common.ParseUint(fields[12]),
		})
	}
	return stats, nil
}

func RegisterNetworkMetrics(mc *collector.MetriclyCollector) {
	mc.AddMetric("network_rx_bytes", "total bytes received", []string{"interface"})
	mc.AddMetric("network_tx_bytes", "total bytes transmitted", []string{"interface"})
	mc.AddMetric("network_rx_packets", "total packets received", []string{"interface"})
	mc.AddMetric("network_tx_packets", "total packets transmitted", []string{"interface"})
	mc.AddMetric("network_rx_errors", "total errors received", []string{"interface"})
	mc.AddMetric("network_tx_errors", "total errors transmitted", []string{"interface"})
	mc.AddMetric("network_rx_drops", "total drops received", []string{"interface"})
	mc.AddMetric("network_tx_drops", "total drops transmitted", []string{"interface"})
}

func subtractCurrPrev(prev, curr networkStats) networkStats {
	return networkStats{
		interfaceName: prev.interfaceName,
		bytesRx:       curr.bytesRx - prev.bytesRx,
		bytesTx:       curr.bytesTx - prev.bytesTx,
		packetsRx:     curr.packetsRx - prev.packetsRx,
		packetsTx:     curr.packetsTx - prev.packetsTx,
		errorsRx:      curr.errorsRx - prev.errorsRx,
		errorsTx:      curr.errorsTx - prev.errorsTx,
		dropsRx:       curr.dropsRx - prev.dropsRx,
		dropsTx:       curr.dropsTx - prev.dropsTx,
	}
}

func calculatePerSecondMetrics(prev, curr []networkStats) ([]networkStats, error) {

	if len(prev) != len(curr) {
		return nil, fmt.Errorf("invalid previous and current network metrics: failed to find rate")
	}

	rateResult := make([]networkStats, len(prev))

	for i := range len(prev) {
		rateResult[i] = subtractCurrPrev(prev[i], curr[i])
	}
	return rateResult, nil

}

func ReportNetworkUsage(mc *collector.MetriclyCollector) {

	start := time.Now()
	prevNWStat, _ := readNetworkStats()
	time.Sleep(1 * time.Second)
	currNWStat, _ := readNetworkStats()

	increaseNWStats, err := calculatePerSecondMetrics(prevNWStat, currNWStat)
	if err != nil {
		slog.Warn(fmt.Sprint(err))
		return
	}

	for _, stat := range increaseNWStats {

		mc.UpdateMetric(
			"network_rx_bytes",
			float64(stat.bytesRx),
			[]string{stat.interfaceName},
		)

		mc.UpdateMetric(
			"network_tx_bytes",
			float64(stat.bytesTx),
			[]string{stat.interfaceName},
		)

		mc.UpdateMetric(
			"network_rx_packets",
			float64(stat.packetsRx),
			[]string{stat.interfaceName},
		)

		mc.UpdateMetric(
			"network_tx_packets",
			float64(stat.packetsTx),
			[]string{stat.interfaceName},
		)

		mc.UpdateMetric(
			"network_rx_errors",
			float64(stat.errorsRx),
			[]string{stat.interfaceName},
		)

		mc.UpdateMetric(
			"network_tx_errors",
			float64(stat.errorsTx),
			[]string{stat.interfaceName},
		)

		mc.UpdateMetric(
			"network_rx_drops",
			float64(stat.dropsRx),
			[]string{stat.interfaceName},
		)

		mc.UpdateMetric(
			"network_tx_drops",
			float64(stat.dropsTx),
			[]string{stat.interfaceName},
		)
	}
	slog.Info(fmt.Sprintf("Collected Network metrics in %s", time.Since(start)))
}
