package server

import (
	"fmt"
	"log"
	"metricly/internal/pollster"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func StartServer(port int, cc *pollster.MetricCollector, interval time.Duration) error {
	prometheus.MustRegister(cc)

	go func() {
		for {
			pollster.ReportCpuUsage(cc)
			time.Sleep(interval * time.Second)
		}
	}()

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte(
			`<html>
			<head><title>Metricly Exporter</title></head>
			<body>
			<h1>Metricly Exporter</h1>
			<p><a href='/metrics'>Metrics</a></p>
			</body>
			</html>`))
		if err != nil {
			log.Fatal(err)
		}
	})
	mux.Handle("/metrics", promhttp.Handler())

	metricsURL := fmt.Sprintf("%s:%d", "", port)

	server := &http.Server{
		Addr:    metricsURL,
		Handler: mux,
	}

	return server.ListenAndServe()
}
