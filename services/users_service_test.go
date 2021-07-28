package services

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/ulexxander/transport-madness/testutils"
)

func TestUsersService_CreateUser(t *testing.T) {
	r := require.New(t)

	publisher := testutils.Publisher{}
	us := NewUsersService(&publisher)

	users := us.UsersAll()
	r.Len(users, 0)

	user, err := us.UserByUsername(UserByUsernameInput{})
	r.ErrorIs(err, ErrUsernameEmpty)
	r.Nil(user)

	user, err = us.CreateUser(UserCreateInput{})
	r.ErrorIs(err, ErrUsernameEmpty)
	r.Nil(user)

	user, err = us.UserByUsername(UserByUsernameInput{"alex"})
	r.Error(err)
	r.Nil(user)

	user, err = us.CreateUser(UserCreateInput{"alex"})
	r.Nil(err)
	r.Equal("alex", user.Username)

	user, err = us.UserByUsername(UserByUsernameInput{"alex"})
	r.Nil(err)
	r.Equal("alex", user.Username)

	users = us.UsersAll()
	r.Len(users, 1)
}
