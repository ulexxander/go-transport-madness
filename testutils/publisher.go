package testutils

import (
	"github.com/stretchr/testify/require"
	"github.com/ulexxander/transport-madness/models"
)

type Publisher struct {
	R               *require.Assertions
	UsersCreated    []models.User
	MessagesCreated []models.Message
}

func (p *Publisher) PublishUserCreated(user *models.User) {
	p.UsersCreated = append(p.UsersCreated, *user)
}

func (p *Publisher) PublishMessageCreated(msg *models.Message) {
	p.MessagesCreated = append(p.MessagesCreated, *msg)
}

func (p *Publisher) NoUsers() {
	p.R.Len(p.UsersCreated, 0)
}

func (p *Publisher) NoMessages() {
	p.R.Len(p.MessagesCreated, 0)
}

func (p *Publisher) LastUserEqual(user *models.User) {
	p.R.Equal(user, p.UsersCreated[len(p.UsersCreated)-1])
}

func (p *Publisher) LastMessageEqual(msg *models.Message) {
	p.R.Equal(msg, &p.MessagesCreated[len(p.MessagesCreated)-1])
}

func (p *Publisher) LenUsers(count int) {
	p.R.Len(p.UsersCreated, count)
}

func (p *Publisher) LenMessages(count int) {
	p.R.Len(p.MessagesCreated, count)
}
