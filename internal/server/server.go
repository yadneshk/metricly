package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	v1 "metricly/api/v1"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"metricly/config"
)

func StartMetriclyServer(ctx context.Context, conf *config.Config) {

	mux := http.NewServeMux()
	v1.HandleRoutes(mux, conf)

	metricsURL := fmt.Sprintf("%s:%s", conf.Server.Address, conf.Server.Port)

	server := &http.Server{
		Addr:    metricsURL,
		Handler: mux,
	}
	slog.Info(fmt.Sprintf("Starting to host metrics on %s ...", metricsURL))

	errChan := make(chan error)
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errChan <- err
		}
	}()

	slog.Info("Started Metricly...")

	select {
	case err := <-errChan:
		slog.Error(fmt.Sprintf("failed to listen and serve: %v", err))
	case <-signalChan:
		slog.Error("Received shutdown signal...")

		// derived ctx from
		ctx, cancel := context.WithDeadline(ctx, time.Now().Add(5*time.Second))
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			slog.Error(fmt.Sprintf("failed to shutdown gracefully: %v", err))
		}
	}
}
