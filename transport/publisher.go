package transport

import "github.com/ulexxander/transport-madness/models"

type Publisher interface {
	PublishUserCreated(user *models.User)
	PublishMessageCreated(msg *models.Message)
}
