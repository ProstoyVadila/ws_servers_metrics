package main

import "net/http"

func (app *Application) setRoutes(wsChat *WsChat) {
	http.HandleFunc("GET /", HomeHandler)
	http.HandleFunc("GET /healthcheck", HealthCheckHandler)
	http.HandleFunc("GET /ws", func(w http.ResponseWriter, r *http.Request) {
		initWebsocket(wsChat, w, r)
	})
}
