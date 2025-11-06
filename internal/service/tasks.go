package service

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/Luc1808/TaskAPI/internal/repository"
	"github.com/Luc1808/TaskAPI/pkg/models"
	"github.com/google/uuid"
)

var (
	ErrInvalidTitle  = errors.New("title is required and must be <= 140 characters")
	ErrInvalidStatus = errors.New("status is invalid")
	ErrNotFound      = errors.New("task not found")
)

var allowedStatus = map[string]bool{
	"todo":        true,
	"in_progress": true,
	"done":        true,
}

type CreateTaskInput struct {
	Title       string
	Description string
	Status      string
}

type UpdateTaskInput struct {
	Title       *string
	Description *string
	Status      *string
}

type ListOptions struct {
	Status   string
	Search   string
	Page     string
	PageSize string
}

type TaskService struct {
	repo repository.TaskRepository
}

func NewTaskService(r repository.TaskRepository) *TaskService {
	return &TaskService{
		repo: r,
	}
}

func validateTitle(t string) error {
	trimmed := strings.TrimSpace(t)
	if trimmed == "" {
		return ErrInvalidTitle
	}
	if len(trimmed) > 140 {
		return ErrInvalidTitle
	}
	return nil
}

func validateStatus(s string) error {
	if s == "" {
		return nil
	}
	if !allowedStatus[s] {
		return ErrInvalidStatus
	}
	return nil
}

func (s *TaskService) CreateTask(ctx context.Context, in CreateTaskInput) (*models.Task, error) {
	if err := validateTitle(in.Title); err != nil {
		return &models.Task{}, err
	}

	status := in.Status
	if status == "" {
		status = "todo"
	}
	if err := validateStatus(status); err != nil {
		return &models.Task{}, err
	}

	now := time.Now().UTC()

	task := &models.Task{
		ID:          uuid.NewString(),
		Title:       strings.TrimSpace(in.Title),
		Description: in.Description,
		Status:      models.TaskStatus(status),
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	return s.repo.Create(ctx, task)
}

func (s *TaskService) GetTask(ctx context.Context, id string) (*models.Task, error) {
	t, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return &models.Task{}, err
		}

		return &models.Task{}, err
	}
	return t, nil
}

func (s *TaskService) ListTasks(ctx context.Context, in ListOptions) ([]models.Task, error) {
	page := parsePositiveInt(in.Page, 1)
	size := parsePositiveInt(in.PageSize, 20)
	offset := (page - 1) * size

	var statusPtr *models.TaskStatus
	if in.Status != "" {
		if err := validateStatus(in.Status); err != nil {
			return nil, err
		}
		st := models.TaskStatus(in.Status)
		statusPtr = &st
	}

	repoFilter := repository.ListFilter{
		Status: statusPtr,
		Search: in.Search,
	}

	repoPagination := repository.Pagination{
		Limit:  size,
		Offset: offset,
	}

	tasks, err := s.repo.List(ctx, repoFilter, repoPagination)
	if err != nil {
		return nil, err
	}

	return tasks, err
	// return s.repo.List(ctx, filter, pagination)
}

func (s *TaskService) UpdateTask(ctx context.Context, id string, in UpdateTaskInput) (models.Task, error) {
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return models.Task{}, err
		}
		return models.Task{}, err
	}

	if in.Title != nil {
		if err := validateTitle(*in.Title); err != nil {
			return models.Task{}, err
		}
		existing.Title = strings.TrimSpace(*in.Title)
	}
	if in.Description != nil {
		existing.Description = *in.Description
	}
	if in.Status != nil {
		if err := validateStatus(*in.Status); err != nil {
			return models.Task{}, err
		}
		if *in.Status != "" {
			existing.Status = models.TaskStatus(*in.Status)
		}
	}

	existing.UpdatedAt = time.Now().UTC()

	updated, err := s.repo.Update(ctx, existing)
	if err != nil {
		return models.Task{}, err
	}

	return *updated, nil
}

func (s *TaskService) DeleteTask(ctx context.Context, id string) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return ErrNotFound
		}
		return err
	}

	return nil
}

func parsePositiveInt(s string, def int) int {
	if s == "" {
		return def
	}
	n, err := strconv.Atoi(s)
	if err != nil || n <= 0 {
		return def
	}
	return n
}
