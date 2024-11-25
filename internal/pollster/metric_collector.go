package pollster

import (
	"fmt"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

type metricData struct {
	description string
	value       float64
	labels      []string
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

	for name, data := range mc.data {
		if desc, exists := mc.metrics[name]; exists {
			ch <- prometheus.MustNewConstMetric(
				desc,
				prometheus.GaugeValue,
				data.value,
				data.labels...,
			)
		}
	}
}

func (mc *MetriclyCollector) addMetric(name string, description string, labels []string) {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()

	// prepend exporter name to every metric name
	name = fmt.Sprintf("metricly_%s", name)

	if _, exists := mc.metrics[name]; !exists {
		mc.metrics[name] = prometheus.NewDesc(
			name,
			description,
			labels,
			nil,
		)
	}
}

func (mc *MetriclyCollector) updateMetric(name string, value float64, labels []string) {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()

	// prepend exporter name to every metric name
	name = fmt.Sprintf("metricly_%s", name)

	if _, exists := mc.metrics[name]; exists {
		mc.data[name] = metricData{
			value:  value,
			labels: labels,
		}
	}

}
