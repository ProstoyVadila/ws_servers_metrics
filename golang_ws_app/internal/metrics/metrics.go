package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	BUCKETS        = []float64{0.0001, 0.0005, 0.0025, 0.005, 0.0075, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1.0, 2.0, 5.0}
	WS_CONNECTIONS = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "ws_connestions",
		Help: "ws connections at the moment",
	})
	ALL_WS_CONNECTIONS_TOTAL = promauto.NewCounter(prometheus.CounterOpts{
		Name: "all_ws_connections_total",
		Help: "all new ws connections",
	})
	WS_CONN_CLOSED_ERRORS_TOTAL = promauto.NewCounter(prometheus.CounterOpts{
		Name: "ws_conn_closed_errors_total",
		Help: "an amount of closed websocket connections",
	})
	WS_BROADCAST_DURATION_SECONDS = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "ws_broadcast_duration_seconds",
		Help:    "a duration of broadcasting message to all connections in seconds",
		Buckets: BUCKETS,
	})
	WS_MESSAGE_HANDLING_DURATION_SECONDS = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "ws_message_handling_duration_seconds",
		Help:    "Response latency in seconds",
		Buckets: BUCKETS,
	})
)
