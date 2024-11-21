package v1

import (
	"log"
	"net/http"
	"metricly/config"
)

func HandleRoutes(mux *http.ServeMux, conf *config.Config) {
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte(
			`<html>
			<head><title>Metricly Exporter</title></head>
			<body>
			<h1>Metricly Exporter</h1>
			<p><a href='/api/v1/metrics'>Metrics</a></p>
			</body>
			</html>`))
		if err != nil {
			log.Fatal(err)
		}
	})

	mux.HandleFunc("/api/v1/metrics", MetricsHandler(conf))

	mux.HandleFunc("/api/v1/query", PrometheusQueryHandler(conf))

	mux.HandleFunc("/api/v1/query_range", prometheusQueryRangeHandler(conf))
}
