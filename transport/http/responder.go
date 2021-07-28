package http

import (
	"log"
	"net/http"

	"github.com/ulexxander/transport-madness/services"
)

type Responder struct {
	Mux             *http.ServeMux
	UsersService    *services.UsersService
	MessagesService *services.MessagesService
	Log             *log.Logger
}

func (rs *Responder) Setup() {
	rs.Mux.HandleFunc("/api/users/all", func(w http.ResponseWriter, r *http.Request) {
		data := rs.UsersService.UsersAll()
		rs.respondData(w, data)
	})

	rs.Mux.HandleFunc("/api/users/create", func(w http.ResponseWriter, r *http.Request) {
		var input services.UserCreateInput
		if err := requestBody(r, &input); err != nil {
			rs.respondError(w, err)
			return
		}
		data, err := rs.UsersService.CreateUser(input)
		if err != nil {
			rs.respondError(w, err)
			return
		}
		rs.respondData(w, data)
	})

	rs.Mux.HandleFunc("/api/messages/pagination", func(w http.ResponseWriter, r *http.Request) {
		page, err := queryInt(r, "page")
		if err != nil {
			rs.respondError(w, err)
			return
		}
		pageSize, err := queryInt(r, "pageSize")
		if err != nil {
			rs.respondError(w, err)
			return
		}
		data, err := rs.MessagesService.MessagesPagination(services.MessagesPaginationInput{
			Page:     page,
			PageSize: pageSize,
		})
		if err != nil {
			rs.respondError(w, err)
			return
		}
		rs.respondData(w, data)
	})

	rs.Mux.HandleFunc("/api/messages/create", func(w http.ResponseWriter, r *http.Request) {
		var input services.MessageCreateInput
		if err := requestBody(r, &input); err != nil {
			rs.respondError(w, err)
			return
		}
		data, err := rs.MessagesService.CreateMessage(input)
		if err != nil {
			rs.respondError(w, err)
			return
		}
		rs.respondData(w, data)
	})
}
