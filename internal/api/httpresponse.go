package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Luc1808/TaskAPI/internal/service"
)

type envelope struct {
	Data  any    `json:"data"`
	Error string `json:"error"`
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	res := envelope{
		Data:  data,
		Error: "",
	}

	_ = json.NewEncoder(w).Encode(res)
}

func writeError(w http.ResponseWriter, err error) {
	status := http.StatusInternalServerError
	msg := "internal error"

	if (errors.Is(err, service.ErrInvalidStatus)) || (errors.Is(err, service.ErrInvalidTitle)) {
		status = http.StatusBadRequest
		msg = err.Error()
	} else if errors.Is(err, service.ErrNotFound) {
		status = http.StatusNotFound
		msg = err.Error()
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	res := envelope{
		Data:  nil,
		Error: msg,
	}

	_ = json.NewEncoder(w).Encode(res)
}
