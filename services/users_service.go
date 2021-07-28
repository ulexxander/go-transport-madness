package services

import (
	"time"

	"github.com/pkg/errors"
)

type User struct {
	Username  string
	CreatedAt time.Time
}

type UsersService struct {
	usersByUsername map[string]*User
}

func NewUsersService() *UsersService {
	return &UsersService{
		usersByUsername: make(map[string]*User),
	}
}

func (us *UsersService) UserByUsername(username string) (*User, error) {
	user, ok := us.usersByUsername[username]
	if !ok {
		return nil, errors.Errorf("user with username %s does not exist", username)
	}
	return user, nil
}

func (us *UsersService) CreateUser(username string) *User {
	user := User{
		Username:  username,
		CreatedAt: time.Now(),
	}

	us.usersByUsername[username] = &user
	return &user
}
