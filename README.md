Here’s a well-structured GitHub README file for the **Metricly** project:

---

# **Metricly**

**Metricly** is a Go-based metrics collection framework designed to gather and expose system-level metrics, such as CPU usage, memory usage, disk I/O, network statistics, and more. These metrics are collected and exposed in Prometheus-compatible format for easy monitoring and visualization.

---

## **Features**
- **Extensible Collectors**:
  - CPU Usage
  - Disk I/O and Space Usage
  - Network Throughput
  - Memory Usage
- **Prometheus Integration**:
  - Exposes metrics in a format compatible with Prometheus.
- **Configurable**:
  - Customize collection intervals and server settings via YAML configuration or environment variables.
- **Docker-Compatible**:
  - Run Metricly as a containerized service.
- **Lightweight**:
  - Minimal resource overhead, suitable for production environments.
- **Custom Logging**:
  - Logs incoming and outgoing API requests with support for multiple log levels (INFO, DEBUG, ERROR).

---

## **Getting Started**

### **Prerequisites**
- **Go 1.22+**
- **Docker** (for containerized deployment)
- **Prometheus** (for metrics scraping)

---

### **Installation**

#### Clone the repository:
```bash
git clone https://github.com/your-username/metricly.git
cd metricly
```

#### Build the binary:
```bash
go build -o metricly cmd/collector/main.go
./metricly --config /path/to/config.yaml
```
OR
#### Use Makefile to build
```bash
make build
./build/metricly --config /path/to/config.yaml
```
OR
#### Use Makefile to build and run
```bash
make run CONFIG_FILE=/path/to/config.yaml
```
---

### **Usage**

#### **Configuration**
Metricly uses a YAML configuration file to specify server settings and collection intervals. It also supports overriding configuration with environment variables. The default config file is `./config/config.yaml`

**Sample `config.yaml`:**
```yaml
server:
  address: "0.0.0.0"
  port: "8080"
prometheus:
  address: "0.0.0.0"
  port: "9090"
interval: 10s
debug: false
```

**Environment Variables:**
- `SERVER_ADDRESS`
- `SERVER_PORT`
- `PROMETHEUS_ADDRESS`
- `PROMETHEUS_PORT`
- `COLLECTION_INTERVAL`
- `DEBUG`

---

#### **Docker/Podman Compose Deployment**
For `podman compose` make sure you have installed `podman-compose` using `pip`
```
pip install podman-compose==1.2.0
```

Run `podman compose up`
```
make run_compose_up
```

To destroy all containers
```
make run_compose_down
```

**Run with Docker:**
```bash
docker build -t metricly .
docker run -p 8080:8080 -v /proc:/host/proc:ro --cap-add SYS_ADMIN metricly
```

---

### **Metrics Exposed**

| **Metric Name**      | **Description**                           | **Labels**             |
|-----------------------|-------------------------------------------|------------------------|
| `cpu_total`           | Total CPU usage as a percentage           | `hostname`            |
| `disk_read_bytes`     | Total bytes read from disk per device      | `device`, `hostname`  |
| `disk_write_bytes`    | Total bytes written to disk per device     | `device`, `hostname`  |
| `network_rx_bytes`    | Total received bytes on a network device   | `interface`, `hostname` |
| `network_tx_bytes`    | Total transmitted bytes on a network device| `interface`, `hostname` |

---

### **Development**

#### **Testing**
Run unit tests:
```bash
go test ./...
```

#### **Logging**
Metricly uses Go’s `log/slog` library for structured logging. Customize log levels by modifying the configuration.

---

### **Contributing**
Contributions are welcome! If you find a bug or want to add a feature:
1. Fork the repository.
2. Create a new branch (`git checkout -b feature-name`).
3. Commit your changes (`git commit -m "Add feature"`).
4. Push to your branch (`git push origin feature-name`).
5. Open a Pull Request.

---

### **License**
This project is licensed under the MIT License. See the `LICENSE` file for details.

---

### **Contact**
For questions or support, please create an issue or contact the maintainer at [yadnesh45@gmail.com](mailto:yadnesh45@gmail.com).

---

With this README, users will have a clear understanding of Metricly’s functionality and how to get started with it.
