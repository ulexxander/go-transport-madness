package http

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

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
}

func requestBody(r *http.Request, out interface{}) error {
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(out); err != nil {
		return errors.Wrap(err, "invalid body")
	}
	return nil
}

func queryInt(r *http.Request, key string) (int, error) {
	valStr := r.URL.Query().Get(key)
	if valStr == "" {
		return 0, errors.Errorf("query parameter %s is missing", key)
	}
	valInt, err := strconv.Atoi(valStr)
	if err != nil {
		return 0, err
	}
	return valInt, nil
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
