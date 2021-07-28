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
