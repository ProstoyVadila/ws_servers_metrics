package main

import (
	"log"
	"net"
	"time"

	"github.com/ProstoyVadila/golang_ws_app/internal/gopool"
	"github.com/gobwas/ws"
	"github.com/mailru/easygo/netpoll"
)

// deadliner is a wrapper around net.Conn that sets read/write deadlines before
// every Read() or Write() call.
type deadliner struct {
	net.Conn
	t time.Duration
}

func (d deadliner) Write(p []byte) (int, error) {
	if err := d.Conn.SetWriteDeadline(time.Now().Add(d.t)); err != nil {
		return 0, err
	}
	return d.Conn.Write(p)
}

func (d deadliner) Read(p []byte) (int, error) {
	if err := d.Conn.SetReadDeadline(time.Now().Add(d.t)); err != nil {
		return 0, err
	}
	return d.Conn.Read(p)
}

func nameConn(conn net.Conn) string {
	return conn.LocalAddr().String() + " > " + conn.RemoteAddr().String()
}

func run_ws(app *Application) {
	// handle_ws is a new incoming connection handler.
	// It upgrades TCP connection to WebSocket, registers netpoll listener on
	// it and stores it as a chat user in Chat instance.
	//
	// We will call it below within accept() loop.
	handle_ws := func(conn net.Conn) {
		// NOTE: we wrap conn here to show that ws could work with any kind of
		// io.ReadWriter.
		safeConn := deadliner{conn, app.config.ioTimeout}

		// Zero-copy upgrade to WebSocket connection.
		hs, err := ws.Upgrade(safeConn)
		if err != nil {
			log.Printf("%s: upgrade error: %v", nameConn(conn), err)
			conn.Close()
			return
		}
		log.Printf("%s: established websocket connection: %+v", nameConn(conn), hs)

		// Register incoming user in chat.
		user := app.chat.Register(safeConn)

		// Create netpoll event descriptor for conn.
		// We want to handle only read events of it.
		desc := netpoll.Must(netpoll.HandleRead(conn))

		// Subscribe to events about conn.
		app.poller.Start(desc, func(event netpoll.Event) {
			if event&(netpoll.EventReadHup|netpoll.EventHup) != 0 {
				// When ReadHup or Hup received, this mean that client has
				// closed at least write end of the connection or connections
				// itself. So we want to stop receive events about such conn
				// and remove it from the chat registry.
				app.poller.Stop(desc)
				app.chat.Remove(user)
				return
			}

			// Here we can read some new message from connection.
			// We can not read it right here in callback, because then we will
			// block the poller's inner loop.
			// We do not want to spawn a new goroutine to read single message.
			// But we want to reuse previously spawned goroutine.
			app.pool.Schedule(func() {
				if err := user.Receive(); err != nil {
					// When receive failed, we can only disconnect broken
					// connection and stop to receive events about it.
					app.poller.Stop(desc)
					app.chat.Remove(user)
				}
			})
		})
	}

	// Create incoming connections listener.
	listener, err := net.Listen("tcp", app.config.addr)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("websocket is listening on %s", listener.Addr().String())

	// Create netpoll descriptor for the listener.
	// We use OneShot here to manually resume events stream when we want to.
	acceptDesc := netpoll.Must(netpoll.HandleListener(
		listener, netpoll.EventRead|netpoll.EventOneShot,
	))

	// accept is a channel to signal about next incoming connection Accept()
	// results.
	accept := make(chan error, 1)

	// Subscribe to events about listener.
	app.poller.Start(acceptDesc, func(event netpoll.Event) {
		// We do not want to accept incoming connection when goroutine pool is
		// busy. So if there are no free goroutines during 1ms we want to
		// cooldown the server and do not receive connection for some short
		// time.
		err := app.pool.ScheduleTimeout(time.Millisecond, func() {
			conn, err := listener.Accept()
			if err != nil {
				accept <- err
				return
			}

			accept <- nil
			handle_ws(conn)
		})

		if err == nil {
			err = <-accept
		}

		if err != nil {
			if err != gopool.ErrScheduleTimeout {
				goto cooldown
			}
			if ne, ok := err.(net.Error); ok && ne.Timeout() {
				goto cooldown
			}

		cooldown:
			delay := 5 * time.Millisecond
			log.Printf("accept error: %v; retrying in %s", err, delay)
			time.Sleep(delay)
		}

		app.poller.Resume(acceptDesc)
	})
}
