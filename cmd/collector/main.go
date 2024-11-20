package main

import (
	"log"
	"time"

	"metricly/internal/pollster"
	"metricly/internal/server"
)

const (
	defaultCollectionInterval = 10 * time.Second
	defaultPort               = 8080
)

func main() {
	// Load configuration (hardcoded defaults for simplicity)
	collectionInterval := defaultCollectionInterval
	port := defaultPort

	// cc := pollster.CreateMetricCollector()
	cc := pollster.CreateMetricCollector()
	err := server.StartServer(port, cc, collectionInterval)

	if err != nil {
		log.Fatalf("Error starting server %v", err)
	}

}
