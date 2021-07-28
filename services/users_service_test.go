package services

import (
	"testing"
)

func TestUsersService_CreateUser(t *testing.T) {
	us := NewUsersService()

	user, err := us.UserByUsername("alex")
	if err == nil {
		t.Fatalf("expected to get error, got nil")
	}
	if user != nil {
		t.Fatalf("expected not to get user, got: %v", user)
	}

	user = us.CreateUser("alex")
	if user.Username != "alex" {
		t.Fatalf("expected created user to have username alex got: %v", user.Username)
	}

	user, err = us.UserByUsername("alex")
	if err != nil {
		t.Fatalf("unexpected error when getting alex: %s", err)
	}
	if user.Username != "alex" {
		t.Fatalf("expected fetched user to have username alex got: %v", user.Username)
	}
}
