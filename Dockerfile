# Dockerfile
FROM golang:1.23.3 AS builder
WORKDIR /app

# Copy and build the Go project
COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o metricly cmd/collector/main.go

# # Use a minimal base image for the final image
FROM quay.io/jitesoft/alpine:3.20.3
COPY --from=builder /app/metricly /bin/metricly
WORKDIR /metricly
COPY ./config/healthcheck .

RUN mkdir /etc/metricly
RUN chown -R nobody:nobody /etc/metricly /metricly

# Expose the port
EXPOSE 8080
USER nobody

# Run the metrics collector
ENTRYPOINT ["/bin/metricly"]

# Default agrs
CMD ["--config", "/etc/metricly/config.yaml"]
