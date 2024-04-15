package main

import (
	"log"

	"github.com/ProstoyVadila/golang_ws_app/internal/chat"
	"github.com/ProstoyVadila/golang_ws_app/internal/gopool"
	"github.com/mailru/easygo/netpoll"
)

type Application struct {
	config *Config
	pool   *gopool.Pool
	chat   *chat.Chat
	poller netpoll.Poller
	exit   chan struct{}
}

func NewApplication(config *Config) *Application {
	// Initialize netpoll instance. We will use it to be noticed about incoming
	// events from listener of user connections.
	poller, err := netpoll.New(nil)
	if err != nil {
		log.Fatal(err)
	}

	// Make pool of X size, Y sized work queue and one pre-spawned
	// goroutine.
	pool := gopool.New(config.workers, config.queue, 1)
	chat := chat.New(pool)
	exit := make(chan struct{})
	return &Application{
		config: config,
		pool:   pool,
		chat:   chat,
		poller: poller,
		exit:   exit,
	}
}

func (app *Application) run() {
	log.Println("starting server...")
	run_ws(app)
	<-app.exit
}
