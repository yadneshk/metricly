package pollster

import (
	"fmt"
	"log/slog"
	"strings"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

type metricData struct {
	value  float64
	labels []string
}

type MetriclyCollector struct {
	metrics map[string]*prometheus.Desc
	data    map[string]metricData
	mutex   sync.Mutex
}

func CreateMetricCollector() *MetriclyCollector {
	return &MetriclyCollector{
		metrics: make(map[string]*prometheus.Desc),
		data:    make(map[string]metricData),
	}
}

func (mc *MetriclyCollector) Describe(ch chan<- *prometheus.Desc) {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()

	for _, desc := range mc.metrics {
		ch <- desc
	}

}

func (mc *MetriclyCollector) Collect(ch chan<- prometheus.Metric) {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()

	// Split the key to extract the metric name and labels
	for name, data := range mc.data {
		parts := strings.Split(name, "|")
		metricName := parts[0]
		labels := parts[1:]
		// if _, exists := mc.metrics[name]; exists {

		ch <- prometheus.MustNewConstMetric(
			mc.metrics[metricName],
			prometheus.GaugeValue,
			data.value,
			labels...,
		)
		// ch <- metric
		// }
	}
}

func (mc *MetriclyCollector) addMetric(name string, description string, labels []string) {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()

	// prepend exporter name to every metric name
	// name = fmt.Sprintf("metricly_%s", name)

	// if _, exists := mc.metrics[name]; !exists {
	mc.metrics[name] = prometheus.NewDesc(
		name,
		description,
		labels,
		nil,
	)
	slog.Debug(fmt.Sprintf("Adding metric %s to registry %T\n", name, prometheus.DefaultRegisterer))
	// }
}

func (mc *MetriclyCollector) updateMetric(name string, value float64, labels []string) {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()

	if len(labels) > 0 {
		// Only Network and Disk collectors use labels while updating metrics
		// because multiple metrics get reported for same resource

		// Create a key from the metric name and labels
		name = fmt.Sprintf("%s|%s", name, strings.Join(labels, "|"))
	}

	// prepend exporter name to every metric name
	// name = fmt.Sprintf("metricly_%s", name)

	// if _, exists := mc.metrics[name]; exists {
	mc.data[name] = metricData{
		value:  value,
		labels: labels,
	}
	// }

}
