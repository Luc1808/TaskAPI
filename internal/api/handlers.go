package api

import (
	"encoding/json"
	"net/http"
	"time"
)

type healthReponse struct {
	Status string    `json:"status"`
	Time   time.Time `json:"time"`
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content Type", "application/json")
	resp := healthReponse{
		Status: "ok",
		Time:   time.Now().UTC(),
	}
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}
