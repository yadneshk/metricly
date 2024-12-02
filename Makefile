# Variables
APP_NAME := metricly
BUILD_DIR := build
SRC_DIR := ./cmd/collector/
BIN := $(BUILD_DIR)/$(APP_NAME)
COMPOSE_FILE := docker-compose.yml
CONFIG_FILE ?= ./config/config.yaml # Default config file location

# Build the Go binary
.PHONY: build
build:
	@echo "Building the Go binary..."
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 go build -o $(BIN) $(SRC_DIR)

# Run Podman Compose to deploy the containers
.PHONY: run
run: build
	@echo "Running Metricly with config: $(CONFIG_FILE)"
	./$(BUILD_DIR)/$(APP_NAME) --config $(CONFIG_FILE)

# Run all test cases
.PHONY: tests
tests:
	@echo "Running tests"
	go test ./... -v -coverprofile=coverage.out

.PHONY: run_container
run_container:
	@echo "Building & spawning Metricly container..."
	podman build . -t metricly
	podman run -d --rm --replace --name metricly \
	-p 8080:8080 \
	-v ./config/config.yaml:/etc/metricly/config.yaml:ro,z \
	-v /:/host/root:ro,slave \
	--health-cmd "/root/healthcheck metricly" \
	-e HOSTNAME=${HOSTNAME} \
	-e PROC_CPU_STAT=/host/root/proc/stat \
	-e PROC_MEMORY_INFO=/host/root/proc/meminfo \
	-e PROC_NETWORK_DEV=/host/root/proc/net/dev \
	-e PROC_DISK_MOUNTS=/host/root/proc/mounts \
	-e PROC_DISK_STATS=/host/root/proc/diskstats \
	localhost/metricly:latest

# Run Podman Compose to deploy the containers
.PHONY: run_compose_up
run_compose_up:
	@echo "Deploying containers using Podman Compose..."
	podman-compose -f $(COMPOSE_FILE) up --build -d --no-cache

# Stop and clean up the containers
.PHONY: run_compose_down
run_compose_down:
	@echo "Stopping and removing containers..."
	podman-compose -f $(COMPOSE_FILE) down

# Show help
.PHONY: help
help:
	@echo "Usage:"
	@echo "  make build            Build the Go binary."
	@echo "  make run              Run Metricly Exporter"
	@echo "  make run_compose_up   Deploy Metricly, Prometheus and Grafana"
	@echo "  make run_compose_down Prune all containers from compose up"
	@echo "  make help             Show this help message."
