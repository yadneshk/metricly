package v1

import (
	"net/http"
	"time"

	"metricly/config"
	"metricly/internal/pollster"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// MetricsHandler serves the Prometheus metrics endpoint.
func MetricsHandler(conf *config.Config) http.HandlerFunc {

	cc := pollster.CreateMetricCollector()
	prometheus.MustRegister(cc)

	pollster.RegisterCPUMetrics(cc)
	pollster.RegisterMemoryMetrics(cc)
	pollster.RegisterNetworkMetrics(cc)

	go func() {
		for {
			pollster.ReportCpuUsage(cc, conf)
		}
	}()

	go func() {
		for {
			pollster.ReportMemoryUsage(cc)
			time.Sleep(conf.CollectionInterval * time.Second)
		}
	}()

	go func() {
		for {
			pollster.ReportNetworkUsage(cc)
			time.Sleep(conf.CollectionInterval * time.Second)
		}
	}()

	handler := promhttp.Handler()

	return func(w http.ResponseWriter, r *http.Request) {
		handler.ServeHTTP(w, r)
	}
}
