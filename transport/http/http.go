package http

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/pkg/errors"
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
		users := rs.UsersService.UsersAll()
		rs.respondData(w, users)
	})

	rs.Mux.HandleFunc("/api/users/create", func(w http.ResponseWriter, r *http.Request) {
		var input services.UserCreateInput
		if err := rs.readBody(r, &input); err != nil {
			rs.respondError(w, err)
			return
		}
		user, err := rs.UsersService.CreateUser(input)
		if err != nil {
			rs.respondError(w, err)
			return
		}
		rs.respondData(w, user)
	})
}

func (rs *Responder) readBody(r *http.Request, out interface{}) error {
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(out); err != nil {
		return errors.Wrap(err, "invalid body")
	}
	return nil
}

type responseSuccess struct {
	Data interface{}
}

type responseError struct {
	// TODO: error codes and http response codes
	Error string
}

func (r *Responder) respondData(w http.ResponseWriter, data interface{}) {
	res := responseSuccess{
		Data: data,
	}
	r.respond(w, res)
}

func (r *Responder) respondError(w http.ResponseWriter, err error) {
	res := responseError{
		Error: err.Error(),
	}
	r.respond(w, res)
}

func (r *Responder) respond(w http.ResponseWriter, payload interface{}) {
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		r.Log.Println("error when encoding http response:", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
