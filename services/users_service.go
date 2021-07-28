package services

import (
	"time"

	"github.com/pkg/errors"
	"github.com/ulexxander/transport-madness/models"
	"github.com/ulexxander/transport-madness/transport"
)

var (
	ErrUsernameEmpty = errors.New("username cannot be empty")
)

type UsersService struct {
	publisher       transport.Publisher
	usersByUsername map[string]models.User
}

func NewUsersService(publisher transport.Publisher) *UsersService {
	return &UsersService{
		usersByUsername: make(map[string]models.User),
		publisher:       publisher,
	}
}

func (us *UsersService) UsersAll() []models.User {
	users := []models.User{}
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

func (us *UsersService) UserByUsername(input UserByUsernameInput) (*models.User, error) {
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

func (us *UsersService) CreateUser(input UserCreateInput) (*models.User, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	user := models.User{
		Username:  input.Username,
		CreatedAt: time.Now(),
	}

	us.usersByUsername[input.Username] = user
	us.publisher.PublishUserCreated(&user)
	return &user, nil
}
