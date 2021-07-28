package graphql

import (
	"github.com/graph-gophers/graphql-go"
	"github.com/ulexxander/transport-madness/models"
	"github.com/ulexxander/transport-madness/services"
)

type User struct {
	Username  string
	CreatedAt graphql.Time
}

func ConvertUser(su models.User) User {
	return User{
		Username:  su.Username,
		CreatedAt: graphql.Time{Time: su.CreatedAt},
	}
}

func ConvertUsers(su []models.User) []User {
	var cu []User
	for _, u := range su {
		cu = append(cu, ConvertUser(u))
	}
	return cu
}

type Message struct {
	SenderUsername string
	Content        string
	CreatedAt      graphql.Time
}

func ConvertMessage(sm models.Message) Message {
	return Message{
		SenderUsername: sm.SenderUsername,
		Content:        sm.Content,
		CreatedAt:      graphql.Time{Time: sm.CreatedAt},
	}
}

func ConvertMessages(sm []models.Message) []Message {
	var cm []Message
	for _, m := range sm {
		cm = append(cm, ConvertMessage(m))
	}
	return cm
}

type MessagePaginationArgs struct {
	Input struct {
		Page     int32
		PageSize int32
	}
}

func (mpa *MessagePaginationArgs) Convert() services.MessagesPaginationInput {
	return services.MessagesPaginationInput{
		Page:     int(mpa.Input.Page),
		PageSize: int(mpa.Input.PageSize),
	}
}
