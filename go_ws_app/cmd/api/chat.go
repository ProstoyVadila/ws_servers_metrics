package main

type WsChat struct {
	clients map[*WsClient]bool
	// inbound messages
	broadcast chan []byte
	// register new clients
	register chan *WsClient
	// remove clients
	remove chan *WsClient
}

func newWsChat() *WsChat {
	return &WsChat{
		broadcast: make(chan []byte),
		register:  make(chan *WsClient),
		remove:    make(chan *WsClient),
		clients:   make(map[*WsClient]bool),
	}
}

func (ws *WsChat) handle() {
	for {
		select {
		case client := <-ws.register:
			ws.clients[client] = true
		case client := <-ws.remove:
			if _, exists := ws.clients[client]; exists {
				delete(ws.clients, client)
				close(client.sendCh)
			}
		case message := <-ws.broadcast:
			for client := range ws.clients {
				select {
				case client.sendCh <- message:
				default:
					close(client.sendCh)
					delete(ws.clients, client)
				}
			}
		}
	}
}
