package main

import (
	"log"
	"net/http"
	"time"
)

type Application struct {
	server *http.Server
}

func newApplicatoin() Application {
	server := &http.Server{
		Addr:              ":8000",
		ReadHeaderTimeout: 2 * time.Second,
	}
	return Application{
		server: server,
	}
}

func (app *Application) serve() {
	log.Printf("serving on http://localhost%s\n", app.server.Addr)
	wsChat := app.runWsChat()
	app.setRoutes(wsChat)
	if err := app.server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func (app *Application) runWsChat() *WsChat {
	chat := newWsChat()
	go chat.handle()
	return chat
}
