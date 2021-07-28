package websocket

import (
	"encoding/json"

	"github.com/gorilla/websocket"
)

type eventConnected struct {
	ConnID int64
}

func (p *Publisher) writeEventConnected(conn *websocket.Conn, connID int64) {
	p.writeEvent(conn, "connected", eventConnected{
		ConnID: connID,
	})
}

type messageSuccess struct {
	Event string
	Data  interface{}
}

func (p *Publisher) writeEvent(conn *websocket.Conn, event string, data interface{}) {
	res := messageSuccess{
		Event: event,
		Data:  data,
	}
	p.write(conn, res)
}

func (p *Publisher) write(conn *websocket.Conn, payload interface{}) {
	encoded, err := json.Marshal(payload)
	if err != nil {
		p.Log.Println("error when encoding websocket message:", err)
		return
	}

	w, err := conn.NextWriter(websocket.TextMessage)
	if err != nil {
		p.Log.Println("error when aquiring websocket next writer:", err)
		return
	}
	defer w.Close()

	if _, err := w.Write(encoded); err != nil {
		p.Log.Println("error when writing websocket message:", err)
		return
	}
}
