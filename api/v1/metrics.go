package v1

import (
	"net/http"
	"time"

	"metricly/config"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// MetricsHandler serves the Prometheus metrics endpoint.
func MetricsHandler(conf *config.Config) http.Handler {

	// no need to pass registry to handler since all metrics are added to global registry
	handler := promhttp.Handler()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		start := time.Now()

		handler.ServeHTTP(w, r)

		logAPIRequests(r, time.Since(start).Milliseconds(), http.StatusOK)

	})

}
