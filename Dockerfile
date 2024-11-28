# Dockerfile
FROM golang:1.23.3 AS builder
WORKDIR /app

# Copy and build the Go project
COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o metricly cmd/collector/main.go

# # Use a minimal base image for the final image
FROM alpine:latest
WORKDIR /root
COPY --from=builder /app/metricly .
RUN mkdir /etc/metricly
# Expose the port
EXPOSE 8080

# Run the metrics collector
ENTRYPOINT ["./metricly"]

# Default agrs
CMD ["--config", "/etc/metricly/config.yaml"]
