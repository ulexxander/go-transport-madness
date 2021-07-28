package services

import (
	"errors"
	"time"

	"github.com/ulexxander/transport-madness/models"
	"github.com/ulexxander/transport-madness/transport"
)

var (
	ErrContentEmpty           = errors.New("content cannot be empty")
	ErrPaginationRangeInvalid = errors.New("invalid pagination range")
)

type MessagesService struct {
	usersService *UsersService
	publisher    transport.Publisher
	messages     []models.Message
}

func NewMessagesService(usersService *UsersService, publisher transport.Publisher) *MessagesService {
	return &MessagesService{
		usersService: usersService,
		publisher:    publisher,
	}
}

type MessagesPaginationInput struct {
	Page     int
	PageSize int
}

func (mpi *MessagesPaginationInput) Validate() error {
	if mpi.Page < 0 {
		return ErrPaginationRangeInvalid
	}
	if mpi.PageSize < 1 {
		return ErrPaginationRangeInvalid
	}
	return nil
}

func (ms *MessagesService) MessagesPagination(input MessagesPaginationInput) ([]models.Message, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	start := input.Page * input.PageSize
	stop := (input.Page + 1) * input.PageSize
	last := len(ms.messages)

	if start >= last {
		start = last
	}

	if stop >= last {
		stop = last
	}

	page := ms.messages[start:stop]
	if page == nil {
		return []models.Message{}, nil
	}
	return page, nil
}

type MessageCreateInput struct {
	SenderUsername string
	Content        string
}

func (cmi *MessageCreateInput) Validate() error {
	if cmi.SenderUsername == "" {
		return ErrUsernameEmpty
	}
	if cmi.Content == "" {
		return ErrContentEmpty
	}
	return nil
}

func (ms *MessagesService) CreateMessage(input MessageCreateInput) (*models.Message, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	user, err := ms.usersService.UserByUsername(UserByUsernameInput{input.SenderUsername})
	if err != nil {
		return nil, err
	}

	message := models.Message{
		SenderUsername: user.Username,
		Content:        input.Content,
		CreatedAt:      time.Now(),
	}

	ms.publisher.PublishMessageCreated(&message)

	ms.messages = append([]models.Message{message}, ms.messages...)
	return &message, nil
}
