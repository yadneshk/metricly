package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"metricly/config"
	collector "metricly/internal/collector"
	"metricly/internal/server"

	"github.com/prometheus/client_golang/prometheus"
)

func main() {
	// Load configuration
	configPath := flag.String("config", "", "configuration file path")
	flag.Parse()
	config, err := config.LoadConfig(configPath)
	if err != nil {
		slog.Error(fmt.Sprintf("Error loading config file %v", err))
	}

	if config.Debug {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}

	cc := collector.CreateMetricCollector()
	prometheus.MustRegister(cc)

	// Context for clean shutdown (parent ctx)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start metrics collection before starting server
	server.StartMetricsCollection(ctx, config.CollectionInterval, cc)

	// Start metricly metrics hosting server
	server.StartMetriclyServer(ctx, config)

}
