package services

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/ulexxander/transport-madness/models"
	"github.com/ulexxander/transport-madness/testutils"
)

func TestMessageService_CreateMessage(t *testing.T) {
	r := require.New(t)

	pub := testutils.Publisher{R: r}
	us := NewUsersService(&pub)
	us.CreateUser(UserCreateInput{"alex"})

	ms := NewMessagesService(us, &pub)

	msg, err := ms.CreateMessage(MessageCreateInput{SenderUsername: "", Content: "abc"})
	r.ErrorIs(err, ErrUsernameEmpty)
	r.Nil(msg)
	pub.NoMessages()

	msg, err = ms.CreateMessage(MessageCreateInput{SenderUsername: "abc", Content: ""})
	r.ErrorIs(err, ErrContentEmpty)
	r.Nil(msg)
	pub.NoMessages()

	msg, err = ms.CreateMessage(MessageCreateInput{"noname", "hello i do not exist :("})
	r.Error(err)
	r.Nil(msg)
	pub.NoMessages()

	firstMsgContent := "first message"
	msg, err = ms.CreateMessage(MessageCreateInput{"alex", firstMsgContent})
	r.NoError(err)
	r.Equal(firstMsgContent, msg.Content)
	pub.LastMessageEqual(msg)
	pub.LenMessages(1)

	msgPage, err := ms.MessagesPagination(MessagesPaginationInput{-1, 5})
	r.ErrorIs(err, ErrPaginationRangeInvalid)
	r.Nil(msgPage)

	msgPage, err = ms.MessagesPagination(MessagesPaginationInput{0, 0})
	r.ErrorIs(err, ErrPaginationRangeInvalid)
	r.Nil(msgPage)

	msgPage, err = ms.MessagesPagination(MessagesPaginationInput{0, 5})
	r.NoError(err)
	r.Equal([]models.Message{*msg}, msgPage)

	us.CreateUser(UserCreateInput{"spammer"})
	for i := 0; i < 12; i++ {
		msg, err := ms.CreateMessage(MessageCreateInput{"spammer", fmt.Sprintf("spam %d", i)})
		r.NoError(err)
		pub.LastMessageEqual(msg)
		pub.LenMessages(i + 2)
	}

	msgPage, err = ms.MessagesPagination(MessagesPaginationInput{0, 5})
	r.NoError(err)
	r.Len(msgPage, 5)
	r.Equal("spam 11", msgPage[0].Content)

	msgPage, err = ms.MessagesPagination(MessagesPaginationInput{1, 5})
	r.NoError(err)
	r.Len(msgPage, 5)
	r.Equal("spam 6", msgPage[0].Content)

	msgPage, err = ms.MessagesPagination(MessagesPaginationInput{2, 5})
	r.NoError(err)
	r.Len(msgPage, 3)
	r.Equal(firstMsgContent, msgPage[2].Content)
}
