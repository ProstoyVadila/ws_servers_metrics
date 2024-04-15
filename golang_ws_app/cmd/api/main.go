package main

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	config := NewConfig()

	// set metrics
	go func() {
		log.Printf("starting http server for metrics on %s\n", config.metricsServerAddr)
		http.Handle("GET /metrics", promhttp.Handler())
		http.ListenAndServe(config.metricsServerAddr, nil)
	}()

	app := NewApplication(config)
	app.run()
}
