package services

import (
	"errors"
	"time"
)

var (
	ErrContentEmpty     = errors.New("content cannot be empty")
	ErrPageRangeInvalid = errors.New("invalid page range")
)

type Message struct {
	SenderUsername string
	Content        string
	CreatedAt      time.Time
}

type MessagesService struct {
	usersService *UsersService
	messages     []Message
}

func NewMessagesService(usersService *UsersService) *MessagesService {
	return &MessagesService{
		usersService: usersService,
	}
}

type MessagesPageInput struct {
	Page     int
	PageSize int
}

func (mpi *MessagesPageInput) Validate() error {
	if mpi.Page < 0 {
		return ErrPageRangeInvalid
	}
	if mpi.PageSize < 1 {
		return ErrPageRangeInvalid
	}
	return nil
}

func (ms *MessagesService) MessagesPage(input MessagesPageInput) []Message {
	start := input.Page * input.PageSize
	stop := (input.Page + 1) * input.PageSize
	last := len(ms.messages)

	if start >= last {
		start = last
	}

	if stop >= last {
		stop = last
	}

	return ms.messages[start:stop]
}

type CreateMessageInput struct {
	SenderUsername string
	Content        string
}

func (cmi *CreateMessageInput) Validate() error {
	if cmi.SenderUsername == "" {
		return ErrUsernameEmpty
	}
	if cmi.Content == "" {
		return ErrContentEmpty
	}
	return nil
}

func (ms *MessagesService) CreateMessage(input CreateMessageInput) (*Message, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	user, err := ms.usersService.UserByUsername(UserByUsernameInput{input.SenderUsername})
	if err != nil {
		return nil, err
	}

	message := Message{
		SenderUsername: user.Username,
		Content:        input.Content,
		CreatedAt:      time.Now(),
	}

	ms.messages = append([]Message{message}, ms.messages...)
	return &message, nil
}
