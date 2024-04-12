# WebSocket servers under high load

Rust ws server is based on Rocket framework. \
Python ws server is based on Tornado.

It's not a comparison of languages or frameworks. The main idea was to take real world applications (simplified versions of course) and find out how they will behave under high load and handle almost instantaneous sessions.

## Metrics

### My custom metrics

**WS_CONNECTIONS** [_Gauge_] - amount of ws connections at the moment

**ALL_WS_CONNECTIONS_TOTAL** [_Counter_] - all ws connections

**WS_CLOSED_CONNS_ERRORS_TOTAL** [_Counter_] - amount of closed websocket connections errors

**WS_BROADCAST_DURATION_SECONDS** [_Histogram_] - duration of broadcasting message to all connections in seconds

**WS_MESSAGE_HANDLING_DURATION_SECONDS** [_Histogram_] - duration of handling message with json de/serialization and fields validation in seconds

### Other default metrics

**k6 metrics** – [here](https://k6.io/docs/using-k6/metrics/reference/)

**cadvisor metrics** - [here](https://github.com/google/cadvisor/blob/master/docs/storage/prometheus.md)

## Dashboard

You can check metrics here after setting up docker containers

**grafana** - http://localhost:3000

### **My custom dashboard** - [here](http://localhost:3000/d/ee75b6b8-f1c6-4ef1-9d39-fe50cc55a274/websocket-server3a-rust-vs-python?orgId=1&refresh=5s)

## Quick Start

### Local Load

1. **Install k6**

On mac:

```
brew install k6
```

More about [k6 installation](https://k6.io/docs/get-started/installation/)

2. **Set up containers**

```
docker compose up -d
```

3. **run k6 load**

```
K6_PROMETHEUS_RW_SERVER_URL=http://localhost:9090/api/v1/write \
k6 run -o experimental-prometheus-rw k6/scripts/both.js
```

### Container Load

1. **Uncomment k6 container manifest in compose.yaml**
2. **set up all containers**

```
docker compose up -d
```

### Cleanup

```
docker compose down
```

### Load Options

Only rust ws app:

```
k6 run k6/scripts/rust.js
```

Only python ws app:

```
k6 run k6/scripts/python.js
```

Add logs from console debug output

```
k6 run -v ...
```

### k6 metrics

To push metrics from k6 to prometheus you need to specify
**K6_PROMETHEUS_RW_SERVER_URL** env variable and add ` -o experimental-prometheus-rw` flag to `k6 run`. More [here](https://k6.io/docs/results-output/real-time/prometheus-remote-write/)

## Ideas for Improvement

Feel free to add more scenarios and metrics.

## Issues

It seems values of CPU usage by container from cadvisor in grafana are lower than should be. Values from Docker Desktop are more correct.
