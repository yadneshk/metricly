package main

import (
	"log"
	"time"

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

	err := server.StartServer(port, collectionInterval)

	if err != nil {
		log.Fatalf("Error starting server %v", err)
	}

}
