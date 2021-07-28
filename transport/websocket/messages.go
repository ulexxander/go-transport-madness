package websocket

import (
	"encoding/json"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

type eventConnected struct {
	ConnID int64
}

func write(conn *websocket.Conn, payload interface{}) error {
	encoded, err := json.Marshal(payload)
	if err != nil {
		return errors.Wrap(err, "failed to encode message")
	}

	w, err := conn.NextWriter(websocket.TextMessage)
	if err != nil {
		return errors.Wrap(err, "failed to aquire next writer")
	}
	defer w.Close()

	if _, err := w.Write(encoded); err != nil {
		return errors.Wrap(err, "failed to write message")
	}

	return nil
}

type messageOut struct {
	Event string
	Data  interface{}
}

func writeEvent(conn *websocket.Conn, event string, data interface{}) error {
	res := messageOut{
		Event: event,
		Data:  data,
	}
	return write(conn, res)
}

func (p *Publisher) broadcast(event string, data interface{}) {
	for connID, conn := range p.connsByID {
		if err := writeEvent(conn, event, data); err != nil {
			p.Log.Printf("failed to broadcast event %s to websocket connection with ID %d\n", event, connID)
			continue
		}
	}
}
