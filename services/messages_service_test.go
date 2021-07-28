package services

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMessageService_CreateMessage(t *testing.T) {
	r := require.New(t)

	us := NewUsersService()
	us.CreateUser(UserCreateInput{"alex"})

	ms := NewMessagesService(us)

	msg, err := ms.CreateMessage(CreateMessageInput{SenderUsername: "", Content: "abc"})
	r.ErrorIs(err, ErrUsernameEmpty)
	r.Nil(msg)

	msg, err = ms.CreateMessage(CreateMessageInput{SenderUsername: "abc", Content: ""})
	r.ErrorIs(err, ErrContentEmpty)
	r.Nil(msg)

	msg, err = ms.CreateMessage(CreateMessageInput{"noname", "hello i do not exist :("})
	r.Error(err)
	r.Nil(msg)

	firstMsgContent := "first message"
	msg, err = ms.CreateMessage(CreateMessageInput{"alex", firstMsgContent})
	r.NoError(err)
	r.Equal(firstMsgContent, msg.Content)

	msgPage := ms.MessagesPage(MessagesPageInput{0, 5})
	r.Equal([]Message{*msg}, msgPage)

	us.CreateUser(UserCreateInput{"spammer"})
	for i := 0; i < 12; i++ {
		_, err := ms.CreateMessage(CreateMessageInput{"spammer", fmt.Sprintf("spam %d", i)})
		r.NoError(err)
	}

	msgPage = ms.MessagesPage(MessagesPageInput{0, 5})
	r.Len(msgPage, 5)
	r.Equal("spam 11", msgPage[0].Content)

	msgPage = ms.MessagesPage(MessagesPageInput{1, 5})
	r.Len(msgPage, 5)
	r.Equal("spam 6", msgPage[0].Content)

	msgPage = ms.MessagesPage(MessagesPageInput{2, 5})
	r.Len(msgPage, 3)
	r.Equal(firstMsgContent, msgPage[2].Content)
}
