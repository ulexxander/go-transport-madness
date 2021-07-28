package main

import (
	"io"
	"log"
	stdHTTP "net/http"
	"os"

	"github.com/pkg/errors"
	"github.com/ulexxander/transport-madness/services"
	"github.com/ulexxander/transport-madness/transport/graphql"
	"github.com/ulexxander/transport-madness/transport/http"
)

func main() {
	log := log.New(os.Stdout, "", log.LstdFlags)

	if err := run(log); err != nil {
		log.Fatalln("service died", err)
	}
}

func run(log *log.Logger) error {
	log.Println("starting service")

	usersService := services.NewUsersService()
	messagesService := services.NewMessagesService(usersService)
	var _ = messagesService

	httpAddr := ":4007"
	httpMux := stdHTTP.NewServeMux()
	httpServ := stdHTTP.Server{
		Addr:    httpAddr,
		Handler: httpMux,
	}

	httpResponder := http.Responder{
		Mux:             httpMux,
		UsersService:    usersService,
		MessagesService: messagesService,
		Log:             log,
	}
	httpResponder.Setup()

	graphqlSchemaFile, err := os.Open("./transport/graphql/schema.graphql")
	if err != nil {
		return errors.Wrap(err, "could not open graphql schema file")
	}
	graphqlSchemaBytes, err := io.ReadAll(graphqlSchemaFile)
	if err != nil {
		return errors.Wrap(err, "could not read graphql schema file")
	}
	graphqlQueryResolver := graphql.Query{
		UsersService:    usersService,
		MessagesService: messagesService,
	}
	graphqlResponder := graphql.Responder{
		Mux:           httpMux,
		Schema:        string(graphqlSchemaBytes),
		QueryResolver: &graphqlQueryResolver,
		Log:           log,
	}
	graphqlResponder.Setup()

	log.Println("listening on", httpAddr)
	return httpServ.ListenAndServe()
}
