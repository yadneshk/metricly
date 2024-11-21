package main

import (
	"log"

	"metricly/config"
	"metricly/internal/server"
)

func main() {
	// Load configuration (hardcoded defaults for simplicity)
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config file %v", err)
	}

	err = server.StartServer(config.Server.Address, config.Server.Port, config.CollectionInterval)

	if err != nil {
		log.Fatalf("Error starting server %v", err)
	}

}
