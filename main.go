package main

import (
	"io"
	"log"
	"net"
	stdHTTP "net/http"
	"os"
	"os/signal"
	"syscall"

	stdNats "github.com/nats-io/nats.go"
	"github.com/pkg/errors"
	"github.com/ulexxander/transport-madness/services"
	"github.com/ulexxander/transport-madness/transport/graphql"
	"github.com/ulexxander/transport-madness/transport/grpc"
	"github.com/ulexxander/transport-madness/transport/http"
	"github.com/ulexxander/transport-madness/transport/nats"
	"github.com/ulexxander/transport-madness/transport/websocket"
	stdGRPC "google.golang.org/grpc"
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

	natsURL := "nats://localhost:4222"
	log.Println("connecting to nats server", natsURL)
	natsConn, err := stdNats.Connect(natsURL)
	if err != nil {
		return errors.Wrap(err, "could not connect to nats")
	}

	natsPublisher := nats.Publisher{
		Conn: natsConn,
		Log:  log,
	}

	usersService := services.NewUsersService(&websocketPublisher)
	messagesService := services.NewMessagesService(usersService, &natsPublisher)

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

	natsResponder := nats.Responder{
		Conn:            natsConn,
		UsersService:    usersService,
		MessagesService: messagesService,
		Log:             log,
	}
	natsResponder.Setup()

	grpcAddr := ":4008"
	log.Println("creating tcp listener for grpc server on addr", grpcAddr)
	grpcListener, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		return errors.Wrap(err, "could not create listener for grpc server")
	}
	grpcServer := stdGRPC.NewServer()

	go func() {
		if err := grpcServer.Serve(grpcListener); err != nil {
			log.Println("failed to serve grpc", err)
		}
	}()

	grpcResponder := grpc.Responder{
		Server:         grpcServer,
		UsersServer:    &grpc.UsersServer{UsersService: usersService},
		MessagesServer: &grpc.MessagesServer{MessagesService: messagesService},
	}
	grpcResponder.Setup()

	go func() {
		log.Println("starting http server on", httpAddr)
		if err := httpServ.ListenAndServe(); err != nil {
			log.Println("failed to listen and serve http", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	log.Println("received interrupt from the os")
	return nil
}
