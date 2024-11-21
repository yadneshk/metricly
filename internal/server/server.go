package server

import (
	"fmt"
	v1 "metricly/api/v1"
	"net/http"
	"time"
)

func StartServer(port int, interval time.Duration) error {

	mux := http.NewServeMux()
	v1.HandleRoutes(mux)

	metricsURL := fmt.Sprintf("%s:%d", "", port)

	server := &http.Server{
		Addr:    metricsURL,
		Handler: mux,
	}

	return server.ListenAndServe()
}
