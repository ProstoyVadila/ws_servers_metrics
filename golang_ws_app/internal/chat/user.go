package chat

import (
	"encoding/json"
	"io"
	"log"
	"strconv"
	"sync"

	"github.com/ProstoyVadila/golang_ws_app/internal/models"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

type User struct {
	mutex sync.Mutex
	conn  io.ReadWriteCloser
	id    uint
	name  string
	chat  *Chat
}

func NewUser(chat *Chat, conn io.ReadWriteCloser) *User {
	return &User{
		chat: chat,
		conn: conn,
	}
}

// Receive reads next message from user's underlying connection.
// It blocks until full message received.
func (u *User) Receive() error {
	req, err := u.readRequest()
	if err != nil {
		log.Printf("cannot read request %v\n", err)
		log.Printf("closing connection for user %d %s", u.id, u.name)
		// TODO: perhaps it's better not to close conn after invalid message
		u.conn.Close()
	}

	if req == nil {
		// Handled some control message.
		return nil
	}
	log.Printf("Got new message %v", req.Show())

	if req.UserId == "" {
		req.UserId = "new user id"
	}
	if req.Data == "" {
		req.Data = strconv.Itoa(int(timestamp()))
	}

	switch req.ActionType {
	case models.BROADCAST:
		u.chat.Broadcast(req)
	case models.DIRECT:
		u.WriteDirect(req)
	case models.PING:
		//
	case models.PONG:
		//
	default:
		u.WriteDirect(req)
	}
	return nil
}

// func (u *User) writeErrorTo(req *models.WsMessage, err models.MessageParams) error {
// 	return u.write(models.WsError{
// 		UserId: req.UserId,
// 		Error:  err,
// 	})
// }

//	func (u *User) writeNotice(method string, params models.MessageParams) error {
//		return u.write(models.WsMessage{
//			ActionType: method,
//			Params:     params,
//		})
//	}

func (u *User) WriteDirect(req *models.WsMessage) {
	log.Println("writing direct message")
	if err := u.write(*req); err != nil {
		log.Printf("cannot write direct message: %v\n", err)
	}
}

func (u *User) write(data interface{}) error {
	w := wsutil.NewWriter(u.conn, ws.StateServerSide, ws.OpText)
	encoder := json.NewEncoder(w)

	u.mutex.Lock()
	defer u.mutex.Unlock()

	if err := encoder.Encode(data); err != nil {
		return err
	}

	if err := w.Flush(); err != nil {
		return err
	}
	return nil
}

func (u *User) writeRaw(p []byte) error {
	u.mutex.Lock()
	defer u.mutex.Unlock()

	_, err := u.conn.Write(p)

	return err
}

// readRequests reads json-rpc request from connection.
// It takes io mutex.
func (u *User) readRequest() (*models.WsMessage, error) {
	u.mutex.Lock()
	defer u.mutex.Unlock()

	wsHeader, reader, err := wsutil.NextReader(u.conn, ws.StateServerSide)
	if err != nil {
		log.Printf("cannot get reader %s\n", err)
		return nil, err
	}

	if wsHeader.OpCode.IsControl() {
		return nil, wsutil.ControlFrameHandler(u.conn, ws.StateServerSide)(wsHeader, reader)
	}

	req := &models.WsMessage{}
	decoder := json.NewDecoder(reader)
	if err := decoder.Decode(req); err != nil {
		log.Printf("cannot decode a request from user %d %s\n", u.id, u.name)
		return nil, err
	}
	return req, nil
}
