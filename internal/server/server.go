package server

import (
	"fmt"
	"log/slog"
	v1 "metricly/api/v1"
	"net/http"

	"metricly/config"
)

func StartServer(conf *config.Config) error {

	mux := http.NewServeMux()
	v1.HandleRoutes(mux, conf)

	metricsURL := fmt.Sprintf("%s:%s", conf.Server.Address, conf.Server.Port)

	server := &http.Server{
		Addr:    metricsURL,
		Handler: mux,
	}
	slog.Info("Start hosting metrics...")
	return server.ListenAndServe()
}
