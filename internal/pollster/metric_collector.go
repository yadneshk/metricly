package pollster

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

type metric struct {
	description string
	value       float64
	labels      map[string]string
}

type MetricCollector struct {
	metrics map[string]metric
	mu      sync.RWMutex
}

func CreateMetricCollector() *MetricCollector {
	return &MetricCollector{
		metrics: make(map[string]metric),
	}
}

func (cc *MetricCollector) Describe(ch chan<- *prometheus.Desc) {
	cc.mu.RLock()
	defer cc.mu.RUnlock()

	for metricName, metricData := range cc.metrics {
		ch <- prometheus.NewDesc(
			metricName,
			metricData.description,
			nil,
			metricData.labels,
		)
	}
}

func (cc *MetricCollector) Collect(ch chan<- prometheus.Metric) {
	cc.mu.RLock()
	defer cc.mu.RUnlock()

	for metricName, metricData := range cc.metrics {
		metricDescriptor := prometheus.NewDesc(
			metricName,
			metricData.description,
			nil,
			metricData.labels,
		)
		ch <- prometheus.MustNewConstMetric(
			metricDescriptor,
			prometheus.GaugeValue,
			metricData.value,
		)
	}
}

func (cc *MetricCollector) UpdateMetric(name string, value float64, description string, labels map[string]string) {
	cc.mu.Lock()
	defer cc.mu.Unlock()
	cc.metrics[name] = metric{
		description: description,
		value:       value,
		labels:      labels,
	}
}
