package main

import (
	"log"
	"net/http"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("go_ws_app"))
}

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ok"))
}

func initWebsocket(wsChat *WsChat, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	sendCh := make(chan []byte, 256)
	client := newWsClient(wsChat, conn, sendCh)
	client.chat.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
}
