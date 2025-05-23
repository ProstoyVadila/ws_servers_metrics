package main

import (
	"bytes"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

type WsClient struct {
	chat   *WsChat
	conn   *websocket.Conn
	sendCh chan []byte
}

func newWsClient(chat *WsChat, conn *websocket.Conn, sendCh chan []byte) *WsClient {
	return &WsClient{
		chat:   chat,
		conn:   conn,
		sendCh: sendCh,
	}
}

func (c *WsClient) readPump() {
	defer func() {
		c.chat.remove <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetWriteDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(appData string) error {
		c.conn.SetWriteDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(
				err,
				websocket.CloseGoingAway,
				websocket.CloseAbnormalClosure,
			) {
				log.Printf("got UnexpectedCloseError: %v\n", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		c.chat.broadcast <- message
	}
}

func (c *WsClient) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.sendCh:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
			}
			ws, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				log.Printf("cannot send message: %s\n", err)
				return
			}
			ws.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.sendCh)
			for i := 0; i < n; i++ {
				ws.Write(newline)
				ws.Write(<-c.sendCh)
			}

			if err := ws.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
