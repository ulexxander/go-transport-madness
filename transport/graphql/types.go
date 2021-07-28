package graphql

import (
	"github.com/graph-gophers/graphql-go"
	"github.com/ulexxander/transport-madness/services"
)

type User struct {
	services.User
	createdAt graphql.Time
}

func (u User) CreatedAt() graphql.Time {
	return u.createdAt
}

func ConvertUsers(su []services.User) []User {
	var cu []User
	for _, u := range su {
		cu = append(cu, User{u, graphql.Time{Time: u.CreatedAt}})
	}
	return cu
}
