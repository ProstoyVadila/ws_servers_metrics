package main

import (
	"io"
	"log"
	"net"
	"net/http"
)

func upstream(name, network, addr string) http.Handler {
	if conn, err := net.Dial(network, addr); err != nil {
		log.Printf("warning: test upstream %q error: %v", name, err)
	} else {
		log.Printf("upstream %q ok", name)
		conn.Close()
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		peer, err := net.Dial(network, addr)
		if err != nil {
			log.Printf("dial upstream error: %v", err)
			w.WriteHeader(502)
			return
		}
		if err := r.Write(peer); err != nil {
			log.Printf("write request to upstream error: %v", err)
			w.WriteHeader(502)
			return
		}
		hj, ok := w.(http.Hijacker)
		if !ok {
			w.WriteHeader(500)
			return
		}
		conn, _, err := hj.Hijack()
		if err != nil {
			w.WriteHeader(500)
			return
		}
		log.Printf(
			"serving %s < %s <~> %s > %s",
			peer.RemoteAddr(), peer.LocalAddr(), conn.RemoteAddr(), conn.LocalAddr(),
		)

		go func() {
			defer peer.Close()
			defer conn.Close()
			io.Copy(peer, conn)
		}()
		go func() {
			defer peer.Close()
			defer conn.Close()
			io.Copy(conn, peer)
		}()
	})
}
