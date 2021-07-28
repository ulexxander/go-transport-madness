package services

import "time"

type Message struct {
	SenderUsername string
	Content        string
	CreatedAt      time.Time
}

type MessageService struct {
	usersService *UsersService
	messages     []Message
}

func NewMessageService(usersService *UsersService) *MessageService {
	return &MessageService{
		usersService: usersService,
	}
}

func (ms *MessageService) MessagesPage(page int, pageSize int) []Message {
	start := page * pageSize
	stop := (page + 1) * pageSize
	last := len(ms.messages)

	if start >= last {
		start = last
	}

	if stop >= last {
		stop = last
	}

	return ms.messages[start:stop]
}

func (ms *MessageService) CreateMessage(senderUsername, content string) (*Message, error) {
	user, err := ms.usersService.UserByUsername(senderUsername)
	if err != nil {
		return nil, err
	}

	message := Message{
		SenderUsername: user.Username,
		Content:        content,
		CreatedAt:      time.Now(),
	}

	ms.messages = append([]Message{message}, ms.messages...)
	return &message, nil
}
