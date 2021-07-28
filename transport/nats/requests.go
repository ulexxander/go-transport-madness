package nats

import (
	"encoding/json"

	"github.com/nats-io/nats.go"
	"github.com/pkg/errors"
)

type responseSuccess struct {
	Data interface{}
}

type responseError struct {
	// TODO: error codes
	Error string
}

func (r *Responder) respondData(msg *nats.Msg, data interface{}) {
	res := responseSuccess{
		Data: data,
	}
	r.respond(msg, res)
}

func (r *Responder) respondError(msg *nats.Msg, err error) {
	res := responseError{
		Error: err.Error(),
	}
	r.respond(msg, res)
}

func (r *Responder) respond(msg *nats.Msg, payload interface{}) {
	encoded, err := json.Marshal(payload)
	if err != nil {
		r.Log.Println("error when encoding nats response:", err)
		return
	}

	if err := msg.Respond(encoded); err != nil {
		r.Log.Println("error when responding to nats request:", err)
	}
}

func messagePayload(msg *nats.Msg, out interface{}) error {
	if err := json.Unmarshal(msg.Data, out); err != nil {
		return errors.Wrap(err, "invalid payload")
	}
	return nil
}
