package repository

import (
	"context"

	"github.com/Luc1808/TaskAPI/pkg/models"
)

type ListFilter struct {
	Status *models.TaskStatus
	Search string
}

type Pagination struct {
	Limit  int
	Offset int
}

type TaskRepository interface {
	Create(ctx context.Context, t *models.Task) (*models.Task, error)
	GetByID(ctx context.Context, id string) (*models.Task, error)
	List(ctx context.Context, f ListFilter, p Pagination) ([]models.Task, error)
	Update(ctx context.Context, t *models.Task) (*models.Task, error)
	Delete(ctx context.Context, id string) error
}
