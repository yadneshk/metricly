package pollster

import (
	"bufio"
	"log/slog"
	"metricly/pkg/common"
	"os"
	"strings"
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

	if procNetDevEnv := os.Getenv("PROC_NET_DEV"); procNetDev != "" {
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
			bytesTx:       common.ParseUint(fields[8]),
			packetsRx:     common.ParseUint(fields[2]),
			packetsTx:     common.ParseUint(fields[9]),
			errorsRx:      common.ParseUint(fields[3]),
			errorsTx:      common.ParseUint(fields[10]),
			dropsRx:       common.ParseUint(fields[4]),
			dropsTx:       common.ParseUint(fields[11]),
		})
		// RegisterNetworkMetrics(mc, stat.interfaceName)
	}
	return stats, nil
}

func RegisterNetworkMetrics(mc *MetriclyCollector) {
	mc.addMetric("network_rx_bytes", "total bytes received", []string{"hostname", "interface"})
	mc.addMetric("network_tx_bytes", "total bytes transmitted", []string{"hostname", "interface"})
	mc.addMetric("network_rx_packets", "total packets received", []string{"hostname", "interface"})
	mc.addMetric("network_tx_packets", "total packets transmitted", []string{"hostname", "interface"})
	mc.addMetric("network_rx_bytes", "total bytes received", []string{"hostname", "interface"})
	mc.addMetric("network_tx_bytes", "total bytes transmitted", []string{"hostname", "interface"})
	mc.addMetric("network_rx_packets", "total packets received", []string{"hostname", "interface"})
	mc.addMetric("network_tx_packets", "total packets transmitted", []string{"hostname", "interface"})
	mc.addMetric("network_rx_errors", "total errors received", []string{"hostname", "interface"})
	mc.addMetric("network_tx_errors", "total errors transmitted", []string{"hostname", "interface"})
	mc.addMetric("network_rx_drops", "total drops received", []string{"hostname", "interface"})
	mc.addMetric("network_tx_drops", "total drops transmitted", []string{"hostname", "interface"})
}

func ReportNetworkUsage(mc *MetriclyCollector) {

	slog.Info("Polling Network metrics...")
	currNWStat, _ := readNetworkStats()

	for _, stat := range currNWStat {

		// constLabelsMap := map[string]string{
		// "interface": stat.interfaceName,
		// }

		mc.updateMetric(
			"network_rx_bytes",
			float64(stat.bytesRx),
			[]string{common.GetHostname(), stat.interfaceName},
		)

		mc.updateMetric(
			"network_tx_bytes",
			float64(stat.bytesTx),
			[]string{common.GetHostname(), stat.interfaceName},
		)

		mc.updateMetric(
			"network_rx_packets",
			float64(stat.packetsRx),
			[]string{common.GetHostname(), stat.interfaceName},
		)

		mc.updateMetric(
			"network_tx_packets",
			float64(stat.packetsTx),
			[]string{common.GetHostname(), stat.interfaceName},
		)

		mc.updateMetric(
			"network_rx_errors",
			float64(stat.errorsRx),
			[]string{common.GetHostname(), stat.interfaceName},
		)

		mc.updateMetric(
			"network_tx_errors",
			float64(stat.errorsTx),
			[]string{common.GetHostname(), stat.interfaceName},
		)

		mc.updateMetric(
			"network_rx_drops",
			float64(stat.dropsRx),
			[]string{common.GetHostname(), stat.interfaceName},
		)

		mc.updateMetric(
			"network_tx_drops",
			float64(stat.dropsTx),
			[]string{common.GetHostname(), stat.interfaceName},
		)
	}
	slog.Info("Polling Network metrics complete")
}
