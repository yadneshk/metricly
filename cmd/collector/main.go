package main

import (
	"context"
	"log"

	"metricly/config"
	"metricly/internal/pollster"
	"metricly/internal/server"

	"github.com/prometheus/client_golang/prometheus"
)

func main() {
	// Load configuration (hardcoded defaults for simplicity)
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config file %v", err)
	}

	cc := pollster.CreateMetricCollector()
	prometheus.MustRegister(cc)

	// Context for clean shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start metrics collection
	server.StartMetricsCollection(ctx, config.CollectionInterval, cc)

	err = server.StartServer(config)

	if err != nil {
		log.Fatalf("Error starting server %v", err)
	}

}
