package v1

import (
	"net/http"
	"time"

	"metricly/internal/pollster"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// MetricsHandler serves the Prometheus metrics endpoint.
func MetricsHandler() http.HandlerFunc {

	cc := pollster.CreateMetricCollector()
	prometheus.MustRegister(cc)

	go func() {
		for {
			pollster.ReportCpuUsage(cc)
			time.Sleep(10 * time.Second)
		}
	}()

	handler := promhttp.Handler()

	return func(w http.ResponseWriter, r *http.Request) {
		handler.ServeHTTP(w, r)
	}
}
