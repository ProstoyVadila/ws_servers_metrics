package main

import "flag"

type Config struct {
	Addr       string
	WsChatAddr string
}

func NewConfig() *Config {
	addr := flag.String("listen", "127.0.0.1:8888", "port to listen")
	chatAddr := flag.String("chat_addr", "127.0.0.1:8000", "chat tcp addr to proxy pass")

	flag.Parse()
	return &Config{
		Addr:       *addr,
		WsChatAddr: *chatAddr,
	}
}
