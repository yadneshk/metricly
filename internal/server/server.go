package server

import (
	"fmt"
	v1 "metricly/api/v1"
	"net/http"
	"time"
)

func StartServer(address string, port string, interval time.Duration) error {

	mux := http.NewServeMux()
	v1.HandleRoutes(mux)

	metricsURL := fmt.Sprintf("%s:%s", address, port)

	server := &http.Server{
		Addr:    metricsURL,
		Handler: mux,
	}

	return server.ListenAndServe()
}
