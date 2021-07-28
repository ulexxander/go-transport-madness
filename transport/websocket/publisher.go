package websocket

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/ulexxander/transport-madness/services"
)

type Publisher struct {
	Mux *http.ServeMux
	Log *log.Logger

	mu         sync.Mutex
	lastConnID int64
	connsByID  map[int64]*websocket.Conn
}

func (p *Publisher) Setup() {
	p.initFields()
	upgrader := websocket.Upgrader{}
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	p.Mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			p.Log.Println("websocket connection upgrade failed:", err)
			return
		}

		p.mu.Lock()
		connID := p.lastConnID
		p.connsByID[connID] = conn
		p.lastConnID++
		p.mu.Unlock()

		p.Log.Printf("incoming websocket connection with id %d\n", connID)

		if err := writeEvent(conn, "connected", eventConnected{
			ConnID: connID,
		}); err != nil {
			p.connRemove(connID)
		}
	})
}

func (p *Publisher) initFields() {
	p.connsByID = make(map[int64]*websocket.Conn)
}

func (p *Publisher) PublishUserCreated(user *services.User) {
	p.broadcast("user_created", user)
}

func (p *Publisher) PublishMessageCreated(msg *services.Message) {
	p.broadcast("message_created", msg)
}

func (p *Publisher) connRemove(connID int64) {
	p.mu.Lock()
	conn := p.connsByID[connID]
	delete(p.connsByID, connID)
	p.mu.Unlock()

	if err := conn.Close(); err != nil {
		p.Log.Printf("could not close websocket connection with id %d\n", connID)
		return
	}
	p.Log.Printf("removed websocket connection with id %d\n", connID)
}
