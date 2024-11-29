package pollster

import (
	"fmt"
	"log/slog"
	"metricly/pkg/common"
	"strings"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

type metricData struct {
	Value  float64
	Labels []string
}

type MetriclyCollector struct {
	Metrics map[string]*prometheus.Desc
	Data    map[string]metricData
	Mutex   sync.Mutex
}

func CreateMetricCollector() *MetriclyCollector {
	return &MetriclyCollector{
		Metrics: make(map[string]*prometheus.Desc),
		Data:    make(map[string]metricData),
	}
}

func (mc *MetriclyCollector) Describe(ch chan<- *prometheus.Desc) {
	mc.Mutex.Lock()
	defer mc.Mutex.Unlock()

	for _, desc := range mc.Metrics {
		ch <- desc
	}

}

func (mc *MetriclyCollector) Collect(ch chan<- prometheus.Metric) {
	mc.Mutex.Lock()
	defer mc.Mutex.Unlock()

	// Split the key to extract the metric name and labels
	for name, data := range mc.Data {
		parts := strings.Split(name, "|")
		metricName := parts[0]
		labels := parts[1:]
		// if _, exists := mc.metrics[name]; exists {

		ch <- prometheus.MustNewConstMetric(
			mc.Metrics[metricName],
			prometheus.GaugeValue,
			data.Value,
			labels...,
		)
		// ch <- metric
		// }
	}
}

func (mc *MetriclyCollector) AddMetric(name string, description string, labels []string) {
	mc.Mutex.Lock()
	defer mc.Mutex.Unlock()

	// prepend exporter name to every metric name
	name = fmt.Sprintf("metricly_%s", name)

	// if _, exists := mc.metrics[name]; !exists {
	mc.Metrics[name] = prometheus.NewDesc(
		name,
		description,
		labels,
		prometheus.Labels{"hostname": common.GetHostname()},
	)
	slog.Debug(fmt.Sprintf("Adding metric %s to registry %T\n", name, prometheus.DefaultRegisterer))
	// }
}

func (mc *MetriclyCollector) UpdateMetric(name string, value float64, labels []string) {
	mc.Mutex.Lock()
	defer mc.Mutex.Unlock()

	if len(labels) > 0 {
		// Only Network and Disk collectors use labels while updating metrics
		// because multiple metrics get reported for same resource

		// Create a key from the metric name and labels
		name = fmt.Sprintf("%s|%s", name, strings.Join(labels, "|"))
	}

	// prepend exporter name to every metric name
	name = fmt.Sprintf("metricly_%s", name)

	// if _, exists := mc.metrics[name]; exists {
	// labels = append(labels, common.GetHostname())

	mc.Data[name] = metricData{
		Value:  value,
		Labels: labels,
	}
	// }

}
