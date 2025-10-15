package api

import (
	"net/http"

	"github.com/Luc1808/TaskAPI/internal/api/middleware"
	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
)

func NewRouter() http.Handler {
	r := chi.NewRouter()

	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)
	r.Use(middleware.RequestID())

	r.Get("/healthz", HealthHandler)

	return r
}
