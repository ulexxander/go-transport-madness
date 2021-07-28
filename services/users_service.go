package services

import (
	"time"

	"github.com/pkg/errors"
)

var (
	ErrUsernameEmpty = errors.New("username cannot be empty")
)

type User struct {
	Username  string
	CreatedAt time.Time
}

type UsersService struct {
	usersByUsername map[string]User
}

func NewUsersService() *UsersService {
	return &UsersService{
		usersByUsername: make(map[string]User),
	}
}

func (us *UsersService) UsersAll() []User {
	users := []User{}
	for _, user := range us.usersByUsername {
		users = append(users, user)
	}
	return users
}

type UserByUsernameInput struct {
	Username string
}

func (ubui *UserByUsernameInput) Validate() error {
	if ubui.Username == "" {
		return ErrUsernameEmpty
	}
	return nil
}

func (us *UsersService) UserByUsername(input UserByUsernameInput) (*User, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	user, ok := us.usersByUsername[input.Username]
	if !ok {
		return nil, errors.Errorf("user with username %s does not exist", input.Username)
	}

	return &user, nil
}

type UserCreateInput struct {
	Username string
}

func (uci *UserCreateInput) Validate() error {
	if uci.Username == "" {
		return ErrUsernameEmpty
	}
	return nil
}

func (us *UsersService) CreateUser(input UserCreateInput) (*User, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	user := User{
		Username:  input.Username,
		CreatedAt: time.Now(),
	}

	us.usersByUsername[input.Username] = user
	return &user, nil
}
