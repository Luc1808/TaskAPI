package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Luc1808/TaskAPI/internal/api/middleware"
	"github.com/Luc1808/TaskAPI/internal/service"
	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
)

func NewRouter(taskSvc *service.TaskService) http.Handler {
	r := chi.NewRouter()

	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)
	r.Use(middleware.RequestID())

	h := NewTaskHandler(taskSvc)

	r.Get("/healthz", h.HealthHandler)

	r.Route("/tasks", func(tr chi.Router) {
		tr.Get("/", h.ListTasks)
		tr.Post("/", h.CreateTask)

		tr.Route("/{id}", func(ir chi.Router) {
			ir.Get("/", h.GetTask)
			ir.Put("/", h.UpdateTask)
			ir.Delete("/", h.DeleteTask)
		})
	})

	return r
}

func decodeJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	r.Body = http.MaxBytesReader(w, r.Body, 1_000_000)

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(dst); err != nil {
		return service.WrapValidation(errors.New("body must contain a single JSON object"))
	}

	return nil
}
