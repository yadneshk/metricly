package v1

import (
	"net/http"

	"metricly/config"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// MetricsHandler serves the Prometheus metrics endpoint.
func MetricsHandler(conf *config.Config) http.Handler {

	// no need to pass registry to handler since all metrics are added to global registry
	return promhttp.Handler()

}
