package main

import (
	"flag"
	"log"
	"time"
)

type Config struct {
	ioTimeout         time.Duration
	addr              string
	metricsServerAddr string
	workers           int
	queue             int
}

func NewConfig() *Config {
	addr := flag.String("listen", "localhost:8000", "address to bind to")
	metricsServerAddr := flag.String("metrics-server", "localhost:8001", "address for pprof http")
	workers := flag.Int("workers", 256, "max workers count")
	queue := flag.Int("queue", 16, "workers task queue size")
	ioTimeout := flag.Duration("io_timeout", time.Millisecond*100, "i/o operations timeout")

	log.Println("seting up config...")
	flag.Parse()
	return &Config{
		addr:              *addr,
		metricsServerAddr: *metricsServerAddr,
		workers:           *workers,
		queue:             *queue,
		ioTimeout:         *ioTimeout,
	}
}
