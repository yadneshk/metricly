# Dockerfile
FROM golang:1.23.3 AS builder
WORKDIR /app

# Copy and build the Go project
COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux GOARCH=x86_64 go build -o metricly cmd/collector/main.go

# # Use a minimal base image for the final image
FROM quay.io/jitesoft/alpine:3.20.3
WORKDIR /root
COPY --from=builder /app/metricly .
COPY ./config/healthcheck .
RUN chmod +x /root/healthcheck
RUN mkdir /etc/metricly
# Expose the port
EXPOSE 8080

# Run the metrics collector
ENTRYPOINT ["./metricly"]

# Default agrs
CMD ["--config", "/etc/metricly/config.yaml"]
