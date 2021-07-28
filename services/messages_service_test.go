package services

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMessageService_CreateMessage(t *testing.T) {
	r := require.New(t)

	us := NewUsersService()
	us.CreateUser("alex")

	ms := NewMessageService(us)

	msg, err := ms.CreateMessage("noname", "hello i do not exist :(")
	r.Nil(msg)
	r.Error(err)

	firstMsgContent := "first message"
	msg, err = ms.CreateMessage("alex", firstMsgContent)
	r.NoError(err)
	r.Equal(firstMsgContent, msg.Content)

	msgPage := ms.MessagesPage(0, 5)
	r.Equal([]Message{*msg}, msgPage)

	us.CreateUser("spammer")
	for i := 0; i < 12; i++ {
		_, err := ms.CreateMessage("spammer", fmt.Sprintf("spam %d", i))
		r.NoError(err)
	}

	msgPage = ms.MessagesPage(0, 5)
	r.Len(msgPage, 5)
	r.Equal("spam 11", msgPage[0].Content)

	msgPage = ms.MessagesPage(1, 5)
	r.Len(msgPage, 5)
	r.Equal("spam 6", msgPage[0].Content)

	msgPage = ms.MessagesPage(2, 5)
	r.Len(msgPage, 3)
	r.Equal(firstMsgContent, msgPage[2].Content)
}
