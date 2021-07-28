package transport

import "github.com/ulexxander/transport-madness/services"

type Publisher interface {
	PublishUserCreated(user *services.User)
	PublishMessageCreated(msg *services.Message)
}
