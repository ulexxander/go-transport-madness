package main

import (
	"io"
	"log"
	stdHTTP "net/http"
	"os"

	stdNats "github.com/nats-io/nats.go"
	"github.com/pkg/errors"
	"github.com/ulexxander/transport-madness/services"
	"github.com/ulexxander/transport-madness/transport/graphql"
	"github.com/ulexxander/transport-madness/transport/http"
	"github.com/ulexxander/transport-madness/transport/nats"
	"github.com/ulexxander/transport-madness/transport/websocket"
)

func main() {
	log := log.New(os.Stdout, "", log.LstdFlags)

	if err := run(log); err != nil {
		log.Fatalln("service died", err)
	}
}

func run(log *log.Logger) error {
	log.Println("starting service")

	httpAddr := ":4007"
	httpMux := stdHTTP.NewServeMux()
	httpServ := stdHTTP.Server{
		Addr:    httpAddr,
		Handler: httpMux,
	}

	websocketPublisher := websocket.Publisher{
		Mux: httpMux,
		Log: log,
	}
	websocketPublisher.Setup()

	usersService := services.NewUsersService(&websocketPublisher)
	messagesService := services.NewMessagesService(usersService, &websocketPublisher)

	httpResponder := http.Responder{
		Mux:             httpMux,
		UsersService:    usersService,
		MessagesService: messagesService,
		Log:             log,
	}
	httpResponder.Setup()

	graphqlSchemaFilepath := "./transport/graphql/schema.graphql"
	log.Println("opening graphql schema file", graphqlSchemaFilepath)
	graphqlSchemaFile, err := os.Open(graphqlSchemaFilepath)
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

	natsURL := "nats://localhost:4222"
	log.Println("connecting to nats server", natsURL)
	natsConn, err := stdNats.Connect(natsURL)
	if err != nil {
		return errors.Wrap(err, "could not connect to nats")
	}
	natsResponder := nats.Responder{
		Conn:            natsConn,
		UsersService:    usersService,
		MessagesService: messagesService,
		Log:             log,
	}
	natsResponder.Setup()

	log.Println("starting http server on", httpAddr)
	return httpServ.ListenAndServe()
}
