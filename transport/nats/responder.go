package nats

import (
	"log"

	"github.com/nats-io/nats.go"
	"github.com/ulexxander/transport-madness/services"
)

type Responder struct {
	Conn            *nats.Conn
	UsersService    *services.UsersService
	MessagesService *services.MessagesService
	Log             *log.Logger
}

func (rs *Responder) Setup() {
	rs.Conn.Subscribe("users_all", func(msg *nats.Msg) {
		data := rs.UsersService.UsersAll()
		rs.respondData(msg, data)
	})

	rs.Conn.Subscribe("user_create", func(msg *nats.Msg) {
		var input services.UserCreateInput
		if err := messagePayload(msg, &input); err != nil {
			rs.respondError(msg, err)
			return
		}
		data, err := rs.UsersService.CreateUser(input)
		if err != nil {
			rs.respondError(msg, err)
			return
		}
		rs.respondData(msg, data)
	})

	rs.Conn.Subscribe("messages_pagination", func(msg *nats.Msg) {
		var input services.MessagesPaginationInput
		if err := messagePayload(msg, &input); err != nil {
			rs.respondError(msg, err)
			return
		}
		data, err := rs.MessagesService.MessagesPagination(input)
		if err != nil {
			rs.respondError(msg, err)
			return
		}
		rs.respondData(msg, data)
	})

	rs.Conn.Subscribe("message_create", func(msg *nats.Msg) {
		var input services.MessageCreateInput
		if err := messagePayload(msg, &input); err != nil {
			rs.respondError(msg, err)
			return
		}
		data, err := rs.MessagesService.CreateMessage(input)
		if err != nil {
			rs.respondError(msg, err)
			return
		}
		rs.respondData(msg, data)
	})
}
