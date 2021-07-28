package websocket

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
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

		p.Log.Println("incoming websocket connection id:", connID)

		// TODO: error handling in future
		p.writeEventConnected(conn, connID)
	})
}

func (p *Publisher) initFields() {
	p.connsByID = make(map[int64]*websocket.Conn)
}
