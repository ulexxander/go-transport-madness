package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/pkg/errors"
)

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
