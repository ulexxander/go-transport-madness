package main

import (
	"log"
	stdHTTP "net/http"
	"os"

	"github.com/ulexxander/transport-madness/services"
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

	log.Println("listening on", httpAddr)
	return httpServ.ListenAndServe()
}
