package chat

import (
	"bytes"
	"encoding/json"
	"log"
	"net"
	"sort"
	"sync"
	"time"

	"github.com/ProstoyVadila/golang_ws_app/internal/gopool"
	"github.com/ProstoyVadila/golang_ws_app/internal/models"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

type Chat struct {
	mutex sync.RWMutex
	idx   uint
	users []*User
	ns    map[string]*User
	pool  *gopool.Pool
	out   chan []byte
}

func New(pool *gopool.Pool) *Chat {
	chat := &Chat{
		pool: pool,
		ns:   make(map[string]*User),
		out:  make(chan []byte),
	}
	go chat.writer()

	return chat
}

func (c *Chat) Register(conn net.Conn) *User {
	log.Println("Adding new user")
	user := NewUser(c, conn)
	c.mutex.Lock()
	{
		user.id = c.idx
		user.name = "rand name" // TODO: ???

		c.users = append(c.users, user)
		c.idx++
	}
	c.mutex.Unlock()

	return user
}

func (c *Chat) Remove(user *User) {
	log.Printf("Removing user %d\n", user.id)
	c.mutex.Lock()
	removed := c.remove(user)
	c.mutex.Unlock()

	if !removed {
		log.Printf("cannot remove user %d %s\n", user.id, user.name)
	}
}

func (c *Chat) remove(user *User) bool {
	if _, exists := c.ns[user.name]; !exists {
		log.Printf("cannot find a user %d %s\n", user.id, user.name)
		return false
	}

	delete(c.ns, user.name)

	i := sort.Search(len(c.users), func(i int) bool {
		return c.users[i].id >= user.id
	})
	if i >= len(c.users) {
		panic("chat: inconsistent state")
	}

	without := make([]*User, len(c.users)-1)
	copy(without[:i], c.users[:i])
	copy(without[i:], c.users[i+1:])
	c.users = without

	return true
}

func (c *Chat) Broadcast(wsMsg *models.WsMessage) error {
	log.Println("broadcast message")

	var buf bytes.Buffer

	w := wsutil.NewWriter(&buf, ws.StateServerSide, ws.OpText)
	encoder := json.NewEncoder(w)

	if err := encoder.Encode(wsMsg); err != nil {
		return err
	}
	if err := w.Flush(); err != nil {
		return err
	}

	c.out <- buf.Bytes()
	return nil
}

// writer writes broadcast messages from chat.out channel
func (c *Chat) writer() {
	for b := range c.out {
		c.mutex.RLock()
		users := c.users
		c.mutex.RUnlock()

		for _, user := range users {
			user := user
			c.pool.Schedule(func() {
				if err := user.writeRaw(b); err != nil {
					log.Printf("cannot send message to user %d %s %v", user.id, user.name, err)
				}
			})
		}
	}
}

func timestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}
