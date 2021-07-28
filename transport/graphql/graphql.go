package graphql

import (
	"context"
	"log"
	"net/http"

	"github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"
	"github.com/ulexxander/transport-madness/services"
)

type Responder struct {
	Mux           *http.ServeMux
	Schema        string
	QueryResolver *Query
	Log           *log.Logger
}

func (rs *Responder) Setup() {
	logger := Logger{StdLog: rs.Log}
	graphql.Logger(&logger)

	schema := graphql.MustParseSchema(rs.Schema, rs.QueryResolver, graphql.UseFieldResolvers())
	rs.Mux.Handle("/graphql", &relay.Handler{Schema: schema})
}

type Logger struct {
	StdLog *log.Logger
}

func (l *Logger) LogPanic(ctx context.Context, value interface{}) {
	l.StdLog.Println("panic during graphql query execution:", value)
}

type Query struct {
	UsersService    *services.UsersService
	MessagesService *services.MessagesService
}

func (q *Query) UsersAll() ([]User, error) {
	data := q.UsersService.UsersAll()
	return ConvertUsers(data), nil
}

func (q *Query) UserCreate(args struct {
	Input services.UserCreateInput
}) (*User, error) {
	data, err := q.UsersService.CreateUser(args.Input)
	if err != nil {
		return nil, err
	}
	u := ConvertUser(*data)
	return &u, nil
}

func (q *Query) MessagesPagination(args MessagePaginationArgs) ([]Message, error) {
	data, err := q.MessagesService.MessagesPagination(args.Convert())
	if err != nil {
		return nil, err
	}
	return ConvertMessages(data), nil
}
