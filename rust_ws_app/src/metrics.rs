use rocket_prometheus::{
    self, 
    prometheus::{Histogram, HistogramOpts, IntCounter, IntGauge}, 
    PrometheusMetrics
};
use once_cell::sync::Lazy;
use log;


pub static WS_CONNECTIONS: Lazy<IntGauge> = Lazy::new(|| {
    IntGauge::new("ws_connections", "an amount of ws connections to the server")
        .expect("Cannot create ws_connections metric")
});
pub static ALL_WS_CONNECTIONS_TOTAL: Lazy<IntCounter> = Lazy::new(|| {
    IntCounter::new("all_ws_connections_total", "a counter of all new connections to the server")
        .expect("Cannot create new_connections_total metric")
});
pub static WS_CONN_CLOSED_ERRORS_TOTAL: Lazy<IntCounter> = Lazy::new(|| {
    IntCounter::new("ws_conn_closed_erros_total", "an amount of closed websocket connections")
        .expect("Cannot create ws_conn_closed_erros_total")
});

pub static BUCKETS: [f64; 14] = [0.0001, 0.0005, 0.0025, 0.005, 0.0075, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1.0, 2.0, 5.0];

pub static WS_BROADCAST_DURATION_SECONDS: Lazy<Histogram> = Lazy::new(|| {
    Histogram::with_opts(
        HistogramOpts::new(
            "ws_broadcast_duration_seconds",
            "a duration of broadcasting message to all connections in seconds"
        ).buckets(BUCKETS.to_vec())
    ).expect("Cannot create ws_broadcast_duration_seconds")
});
pub static WS_MESSAGE_HANDLING_DURATION_SECONDS: Lazy<Histogram> = Lazy::new(|| {
    Histogram::with_opts(
        HistogramOpts::new(
            "ws_message_handling_duration_seconds", 
            "a duration of handling message in seconds"
        ).buckets(BUCKETS.to_vec())
    ).expect("Cannot create ws_message_handling_duration_seconds")
});


pub fn get_prometheus() -> PrometheusMetrics {
    log::info!("Setting up prometheus metrics");
    let prom = PrometheusMetrics::new();
    prom.registry().register(Box::new(WS_CONNECTIONS.clone()))
        .expect("Cannot register ws_connections metric");
    prom.registry().register(Box::new(ALL_WS_CONNECTIONS_TOTAL.clone()))
        .expect("Cannot register ws_connections_total metric");
    prom.registry().register(Box::new(WS_CONN_CLOSED_ERRORS_TOTAL.clone()))
        .expect("Cannot register ws_conn_closed_erros_total");
    prom.registry().register(Box::new(WS_BROADCAST_DURATION_SECONDS.clone()))
        .expect("Cannot register ws_broadcast_duration_seconds");
    prom.registry().register(Box::new(WS_MESSAGE_HANDLING_DURATION_SECONDS.clone()))
        .expect("Cannot register ws_message_handling_duration_seconds");
    prom
}
