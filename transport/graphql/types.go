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

func ConvertUser(su services.User) User {
	return User{su, graphql.Time{Time: su.CreatedAt}}
}

func ConvertUsers(su []services.User) []User {
	var cu []User
	for _, u := range su {
		cu = append(cu, ConvertUser(u))
	}
	return cu
}
