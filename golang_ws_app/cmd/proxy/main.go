package main

import (
	"log"
	"net/http"
)

func main() {
	config := NewConfig()

	http.HandleFunc("GET /", HomeHandler)
	http.HandleFunc("GET /healthcheck", HealthCheckHandler)
	http.Handle("GET /ws", upstream("chat", "tcp", config.WsChatAddr))
	// http.Handle("GET /metrics", upstream("chat_metrics", "tcp", config.WsChatAddr))
	// http.Handle("GET /ws", )

	log.Printf("proxy is listening on %q", config.Addr)
	http.ListenAndServe(config.Addr, nil)
}
