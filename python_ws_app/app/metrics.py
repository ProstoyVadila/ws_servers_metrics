from prometheus_client import Gauge, Counter, Histogram

buckets = [0.0001, 0.0005, 0.0025, 0.005, 0.0075, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1.0, 2.0, 5.0]

WS_CONNECTIONS = Gauge("ws_connections", "amount of ws connections at the moment")
ALL_WS_CONNECTIONS_TOTAL = Counter("all_ws_connections_total", "amount of all new ws connections")

WS_BROADCAST_DURATION_SECONDS = Histogram(
    "ws_broadcast_duration_seconds", 
    "a duration of broadcasting message to all connections in seconds",
    buckets=buckets
    )
WS_MESSAGE_HANDLING_DURATION_SECONDS = Histogram(
    "ws_message_handling_duration_seconds", 
    "a duration of handling message in seconds",
    buckets=buckets
    )

