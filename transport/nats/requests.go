package nats

import (
	"encoding/json"

	"github.com/nats-io/nats.go"
	"github.com/pkg/errors"
)

type reqData struct {
	Data interface{}
}

type reqError struct {
	// TODO: error codes
	Error string
}

func (rs *Responder) respondData(msg *nats.Msg, data interface{}) {
	rs.respond(msg, reqData{
		Data: data,
	})
}

func (rs *Responder) respondError(msg *nats.Msg, err error) {
	rs.respond(msg, reqError{
		Error: err.Error(),
	})
}

func (rs *Responder) respond(msg *nats.Msg, payload interface{}) {
	encoded, err := json.Marshal(payload)
	if err != nil {
		rs.Log.Println("error when encoding nats response:", err)
		return
	}

	if err := msg.Respond(encoded); err != nil {
		rs.Log.Println("error when responding to nats request:", err)
		return
	}
}

func requestPayload(msg *nats.Msg, out interface{}) error {
	if err := json.Unmarshal(msg.Data, out); err != nil {
		return errors.Wrap(err, "invalid payload")
	}
	return nil
}
