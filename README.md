# **Metricly**

**Metricly** is a Go-based metrics collection framework designed to gather and expose system-level metrics, such as CPU usage, memory usage, disk I/O, network statistics, and more. These metrics are collected and exposed in Prometheus-compatible format for easy monitoring and visualization.

![Sample Metricly dashboard](doc/metricly_dashboard.png)
![Sample Network dashboard](doc/metricly_network_dashboard.png)
![Sample Disk dashboard](doc/metricly_disk_dashboard.png)
---

## **Features**
- **Extensible Collectors**:
  - CPU Usage
  - Disk I/O and Space Usage
  - Network Throughput
  - Memory Usage
- **Prometheus Integration**:
  - Exposes metrics in a format compatible with `Prometheus`.
- **Configurable**:
  - Customize collection intervals and server settings via YAML configuration or environment variables.
- **Podman-Compatible**:
  - Run Metricly as a containerized service.
- **Custom Logging**:
  - Logs incoming and outgoing API requests with support for multiple log levels (INFO, DEBUG, ERROR).
- **Metrics Visualization**
  - Provides an inbuilt `Grafana` dashboard to visualize all metrics.
---

## **Getting Started**

### **Prerequisites**
- **Go 1.23+**
- **Podman** (for containerized deployment)
- **Prometheus** (for metrics scraping)

---

### **Installation**

#### Clone the repository:
```bash
$ git clone https://github.com/yadneshk/metricly.git
$ cd metricly
```

#### Build the binary:
```bash
$ go build -o metricly cmd/collector/main.go
$ ./metricly --config /path/to/config.yaml
```
OR
#### Use Makefile to build
```bash
$ make build
$ ./build/metricly --config /path/to/config.yaml
```
OR
#### Use Makefile to build and run
```bash
$ make run CONFIG_FILE=/path/to/config.yaml
```
---

### **Usage**

#### **Configuration**
Metricly uses a YAML configuration file to specify server settings and collection intervals. It also supports overriding configuration with environment variables. The default config file is `/etc/metricly/config.yaml`.

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

**Setting configurations through environment variables:**

| **Env Variables**   |  **Default Values**     | **Description**             |
|-----------------------|-----------------------|-----------------------------|
| `SERVER_ADDRESS`      |  `0.0.0.0`            | Address for server          |
| `SERVER_PORT`         |   `8080`              | Listen, serve on this port  |
| `PROMETHEUS_ADDRESS`  |   `0.0.0.0`           | Prometheus IP address       |
| `PROMETHEUS_PORT`     |    `9090`             | Prometheus serving port     |
| `COLLECTION_INTERVAL` |    `10s`              | Collect metrics after interval |
| `DEBUG`               |    `true`             | Log level                   |
| `HOSTNAME`            |                       | If empty, `os.Hostname()`   |
| `PROC_CPU_STAT`       |    `/proc/stat`       | Source for CPU metrics      |
| `PROC_MEMORY_INFO`    | `/proc/meminfo`       | Source for Memory metrics   |
| `PROC_DISK_STATS`     |  `/proc/diskstats`    | Source for Disk I/O usage   |
| `PROC_DISK_MOUNTS`    |  `/proc/mounts`       | Source for Disk Space Usage |
| `PROC_NETWORK_DEV`    |  `/proc/net/dev`      | Source for Network metrics  |

---

#### **Podman Compose Deployment**
Users of Fedora 36 and later can use the dnf package manager to install podman-compose like so:
```bash
$ sudo dnf install podman-compose
```
```bash
$ podman-compose --version
podman-compose version 1.2.0
podman version 5.3.1
```

Running `run_compose_up` deploys `Metricly`, `Prometheus` and `Grafana`
```bash
$ make run_compose_up
...
✔ Container metricly_prometheus  Started 
✔ Container metricly_metricly    Started
✔ Container metricly_grafana     Started
```

Check containers
```bash
$ podman ps
CONTAINER ID   IMAGE                    COMMAND                  CREATED         STATUS         PORTS     NAMES
b6f86fba674e   metricly-metricly        "./metricly --config…"   3 minutes ago   Up 3 minutes             metricly_metricly
7f97e165d514   grafana/grafana          "/run.sh"                4 minutes ago   Up 3 minutes             metricly_grafana
5d24779a1af3   prom/prometheus:latest   "/bin/prometheus --c…"   4 minutes ago   Up 3 minutes             metricly_prometheus
```

`Grafana` dashboard would be available at `http://127.0.0.1:3000`

Use `admin/admin` to login into dashboard. `Metricly` dashboard comes preloaded only in this deployment approach.


To destroy all containers
```bash
$ make run_compose_down
```

**Run with Podman:**
```bash
$ podman build -t metricly .
$ podman run --rm -d -p 8080:8080 --name metricly \
-v ./config/config.yaml:/etc/metricly/config.yaml:ro \
-e HOSTNAME=${HOSTNAME} \
--health-cmd "/root/healthcheck metricly" \
localhost/metricly:latest
```

Prebuilt images with latest commits are pushed to [Quay](https://quay.io/repository/yadneshk/metricly?tab=tags).

---

### **Metrics Exposed**

| **Metric Name**                   | **Description**                        | **Unit**   | **Labels** |
|-----------------------------------|----------------------------------------|------------|------------|
| `cpu_total`                       | Total CPU usage                        |  percent   | `hostname` |
| `cpu_system`                      | Total system CPU usage                 |  percent   | `hostname` |
| `cpu_user`                        | Total user CPU usage                   |  percent   | `hostname` |
| `cpu_steal`                       | Total steal                            |  percent   | `hostname` |
| `memory_total_bytes`              | Total memory                           |  bytes     | `hostname` |
| `memory_available_bytes`          | Total available memory                 |  bytes     | `hostname` |
| `memory_free_bytes`               | Free memory                            |  bytes     | `hostname` |
| `memory_hugepages_free`           | Free hugepages                         |  count     | `hostname` |
| `memory_hugepages_total`          | Total hugepages                        |  count     | `hostname` |
| `memory_hugepages_rsvd`           | Reserved hugepages                     |  count     | `hostname` |
| `memory_hugepages_surp`           | Surplus hugepages                      |  count     | `hostname` |
| `network_rx_bytes`                | Bytes received                         |  bytes/s   | `interface`, `hostname` |
| `network_tx_bytes`                | Bytes transmitted                      |  bytes/s   |  `interface`, `hostname` |
| `network_rx_packets`              | Packets received                       |  packets/s | `interface`, `hostname` |
| `network_tx_packets`              | Packets transmitted                    |  packets/s | `interface`, `hostname` |
| `network_rx_drops`                | Packets droppped while receiving       | packets/s  | `interface`, `hostname` |
| `network_tx_drops`                | Packets droppped while transmitting    | packets/s  | `interface`, `hostname` |
| `network_rx_errors`               | Malformed packets while receiving      | packets/s  | `interface`, `hostname` |
| `network_tx_errors`               | Malformed packets while transmitting   | packets/s  | `interface`, `hostname` |
| `disk_available_bytes`            | Available Disk space                   | bytes      | `interface`, `hostname` |
| `disk_total_bytes`                | Total Disk Space                       | bytes      | `interface`, `hostname` |
| `disk_usage_percentage`           | Disk Usage                             | percent    | `interface`, `hostname` |
| `disk_used_bytes`                 | Disk Usage                             | bytes      | `interface`, `hostname` |
| `disk_io_in_progress`             | Current disk IO operations in progress | count      | `interface`, `hostname` |
| `disk_io_time_spent_seconds`      | Time spent on IO operations in seconds | milliseconds | `interface`, `hostname` |
| `disk_read_throughput_bytes`      | Disk read throughput in bytes          | bytes      | `interface`, `hostname` |
| `disk_write_throughput_bytes`     | Disk write throughput in bytes         | bytes      | `interface`, `hostname` |
| `disk_reads_completed_total`      | Total disk reads completed             | bytes      | `interface`, `hostname` |
| `disk_writes_completed_total`     | Total disk writes completed            | bytes      | `interface`, `hostname` |
| `disk_weighted_io_time_seconds`   | Weighted time spent on IO in seconds   | milliseconds | `interface`, `hostname` |

---

### APIs Exposed ###

The Metricly exporter provides the following API endpoints:
1. Metrics
    - Path: `/api/v1/metrics`
    - Method: `GET`
    - Description: Returns the metrics collected by the exporter in Prometheus format.
    - Example: Get all metrics exposed
      ```bash
      $ curl http://localhost:8080/api/v1/metrics
      metricly_cpu_steal{hostname="fedora"} 0
      metricly_cpu_system{hostname="fedora"} 4.76
      metricly_cpu_total{hostname="fedora"} 13.22
      metricly_cpu_user{hostname="fedora"} 7.41
      metricly_disk_available_bytes{hostname="fedora",mount_point="/"} 4.66716291072e+11
      ```

2. Query
    - Path: `/api/v1/query`
    - Method: `GET`
    - Description: Provides a simple health check to verify the exporter is running.
    - Request Parameters
        | **Parameter** | **Type**  | **Required**  | **Description** |  **Example Value**   |
        |---------------|-----------|---------------|-----------------|----------------------|
        |    `metric`   | `string`  |    Yes        |   Metric Name   | `metricly_cpu_total` |
        |   `timestamp` | `unix_timestamp`  |    No         |   Metric at specific timestamp   | `2024-11-29T12:25:27Z` |
    - Example: Find current CPU utilization
      ```bash
      $ curl http://localhost:8080/api/v1/query?metric=metricly_cpu_total
      {
        "status": "success",
        "data": {
          "resultType": "vector",
          "result": [
            {
              "metric": {
                "__name__": "metricly_cpu_total",
                "hostname": "fedora",
                "instance": "127.0.0.1:8080",
                "job": "metricly"
              },
              "value": [
                1732883589.795,
                "17.94"
              ]
            }
          ]
        }
      }
      ```
    - Example: Find CPU utilization at a specific timestamp
      ```bash
      $ curl http://localhost:8080/api/v1/query?metric=metricly_cpu_total&timestamp=2024-11-29T12:25:27Z
      {
        "status": "success",
        "data": {
          "resultType": "vector",
          "result": [
            {
              "metric": {
                "__name__": "metricly_cpu_total",
                "hostname": "fedora",
                "instance": "127.0.0.1:8080",
                "job": "metricly"
              },
              "value": [
                1732883127,
                "28.41"
              ]
            }
          ]
        }
      }    
      ```

3. Query Range
    - Path: `/api/v1/query_range`
    - Method: `GET`
    - Description: Query metrics over a specific range of time.
    - Request Parameters
        | **Parameter** | **Type**  | **Required**  | **Description** |  **Example Value**   |
        |---------------|-----------|---------------|-----------------|----------------------|
        |  `metric`     | `string`  |    Yes        |   Metric Name   | `metricly_cpu_total` |
        |   `start`     | `unix_timestamp`  |    Yes        |   Start timestamp   | `2024-11-29T12:25:27Z` |
        |   `end`       | `unix_timestamp`  |    Yes        |   End timestamp   | `2024-11-29T12:25:27Z` |    
        |   `step`      | `seconds`  |    Yes        |   Query resolution step width   | `15s` |    
    - Example: Find CPU utilization in 15 minute window by specifying `start` and `end` timestamps
      ```bash
      $ curl http://127.0.0.1:8080/api/v1/query_range?metric=metricly_cpu_total&start=2024-11-29T12:25:00Z&end=2024-11-29T12:40:00Z&step=15s
      {
        "status": "success",
        "data": {
          "resultType": "matrix",
          "result": [
            {
              "metric": {
                "__name__": "metricly_cpu_total",
                "hostname": "fedora",
                "instance": "127.0.0.1:8080",
                "job": "metricly"
              },
              "values": [
                [
                  1732883100,
                  "13.1"
                ],
                [
                  1732883115,
                  "14.84"
                ],
                [
                  1732883130,
                  "28.41"
                ],
                .
                .
                .
                [
                  1732884000,
                  "14.84"
                ]
              ]
            }
          ]
        }
      }
      ```

4. Aggregate
    - Path: `/api/v1/aggregate`
    - Method: `GET`
    - Description: Aggregate metrics over time using mathematical operations.
    - Request Parameters
        | **Parameter** | **Type**  | **Required**  | **Description** |  **Example Value**   |
        |---------------|-----------|---------------|-----------------|----------------------|
        |  `metric`     | `string`  |    Yes        |   Metric Name   | `metricly_cpu_total` |
        |  `operation`  | `string`  |    Yes        |   Aggregation operation   | `avg`, `max`,`min` |
        |   `window`    | `duration` |    Yes        |   End timestamp   | `d`,`h`,`m`, `ms`, `s`, `w`, `y` |
    - Example: Find average CPU utilization in the last 2 hours.
      ```bash
      $ curl http://127.0.0.1:8080/api/v1/aggregate?metric=metricly_cpu_total&operation=avg&window=2h
      {
        "status": "success",
        "data": {
          "resultType": "vector",
          "result": [
            {
              "metric": {
                "hostname": "fedora",
                "instance": "127.0.0.1:8080",
                "job": "metricly"
              },
              "value": [
                1732886126.191,
                "17.223917274939144"
              ]
            }
          ]
        }
      }      
      ```
    - Example: Find the maximum CPU utilization in the last 2 hours.
      ```bash
      $ curl http://127.0.0.1:8080/api/v1/aggregate?metric=metricly_cpu_total&operation=max&window=2h
      {
        "status": "success",
        "data": {
          "resultType": "vector",
          "result": [
            {
              "metric": {
                "hostname": "fedora",
                "instance": "127.0.0.1:8080",
                "job": "metricly"
              },
              "value": [
                1732886765.71,
                "84.05"
              ]
            }
          ]
        }
      }
      ```      

### **Development**

#### **Testing**
Run unit tests:
```bash
$ go test ./...
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

### **Contact**
For questions or support, please create an issue or contact the maintainer at [yadnesh45@gmail.com](mailto:yadnesh45@gmail.com).

---
