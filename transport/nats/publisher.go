package nats

import (
	"log"

	"github.com/nats-io/nats.go"
	"github.com/ulexxander/transport-madness/models"
)

type Publisher struct {
	Conn *nats.Conn
	Log  *log.Logger
}

func (p *Publisher) PublishUserCreated(user *models.User) {
	p.publish("user_created", user)
}

func (p *Publisher) PublishMessageCreated(msg *models.Message) {
	p.publish("message_created", msg)
}
