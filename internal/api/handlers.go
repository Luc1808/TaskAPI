package api

import (
	"net/http"
	"time"

	"github.com/Luc1808/TaskAPI/internal/service"
	"github.com/go-chi/chi/v5"
)

type TaskHandler struct {
	svc *service.TaskService
}

type healthReponse struct {
	Status string    `json:"status"`
	Time   time.Time `json:"time"`
}

func NewTaskHandler(svc *service.TaskService) *TaskHandler {
	return &TaskHandler{svc: svc}
}

func (h *TaskHandler) HealthHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"status": "ok",
	})
}

func (h *TaskHandler) ListTasks(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")

	pageStr := r.URL.Query().Get("page")
	sizeStr := r.URL.Query().Get("page_size")

	result, err := h.svc.ListTasks(r.Context(), service.ListOptions{
		Status:   status,
		Page:     pageStr,
		PageSize: sizeStr,
	})
	if err != nil {
		writeError(w, err)
	}

	writeJSON(w, http.StatusOK, result)
}

func (h *TaskHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	task, err := h.svc.GetTask(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, task)
}

func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var req service.CreateTaskInput
	if err := decodeJSON(w, r, &req); err != nil {
		writeError(w, err)
		return
	}

	newTask, err := h.svc.CreateTask(r.Context(), service.CreateTaskInput{
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
	})
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, newTask)
}

func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req service.UpdateTaskInput
	if err := decodeJSON(w, r, &req); err != nil {
		writeError(w, err)
		return
	}

	updated, err := h.svc.UpdateTask(r.Context(), id, service.UpdateTaskInput{
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
	})
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, updated)
}

func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := h.svc.DeleteTask(r.Context(), id); err != nil {
		writeError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
