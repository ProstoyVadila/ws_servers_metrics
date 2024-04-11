# WebSocket servers comparison

Rust server is based on Rocket framework.
Python app is based on Tornado.

## Metrics

**WS_CONNECTIONS** - amount of ws connections at the moment

**ALL_WS_CONNECTIONS_TOTAL** - amount of all new ws connections

**WS_BROADCAST_DURATION_SECONDS** - duration of broadcasting message to all connections in seconds

**WS_MESSAGE_HANDLING_DURATION_SECONDS** - duration of handling message in seconds

**WS_CLOSED_CONNS_ERRORS_TOTAL** - amount of closed websocket connections

## Dashboard

grafana - http://localhost:3000
dashboard - [here](http://localhost:3000/d/ee75b6b8-f1c6-4ef1-9d39-fe50cc55a274/websocket-server3a-rust-vs-python?orgId=1&refresh=5s)

## Quick Start

### Local Load

1. **Install k6**

On mac:

```
brew install k6
```

More about [k6 installation](https://k6.io/docs/get-started/installation/)

2. **raise containers**

```
docker compose up -d
```

3. **run k6 load**

```
k6 run k6/scripts/both.js
```

### Container Load

1. **Uncomment k6 container in compose.yaml**
2. **raise all containers**

```
docker compose up -d
```

### Cleanup

```
docker compose down
```

### Load Options

Only on rust ws app:

```
k6 run k6/scripts/rust.js
```

Only on python ws app:

```
k6 run k6/scripts/python.js
```

with console debug output

```
k6 run -v ...
```

## Ideas for Improvement

Feel free to add more scenarios or change an existing one.

### Get additional metrics from k6 to Prometheus

https://k6.io/docs/results-output/real-time/prometheus-remote-write/
