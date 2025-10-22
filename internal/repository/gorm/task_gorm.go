package gorm

import (
	"context"
	"errors"
	"time"

	"github.com/Luc1808/TaskAPI/pkg/models"
	"github.com/prometheus/common/model"
	"gorm.io/gorm"
)

type TaskRepo struct {
	db *gorm.DB
}

func NewTaskRepo(db *gorm.DB) *TaskRepo {
	return &TaskRepo{db: db}
}

type TaskRow struct {
	ID          int64      `gorm:"column:id;primaryKey;autoIncrement"`
	Title       string     `gorm:"column:title;type:text;not null"`
	Description string     `gorm:"column:description;type:text;not null;default:''"`
	Status      string     `gorm:"column:status;type:text;not null"`
	Priority    int        `gorm:"column:priority;not null;default:0"`
	DueAt       *time.Time `gorm:"column:due_at"`
	CreatedAt   time.Time  `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time  `gorm:"column:updated_at;autoUpdateTime"`
}

func (TaskRow) TableName() string { return "public.tasks" }

func toRow(t *models.Task) *TaskRow {
	return &TaskRow{
		ID:          t.ID,
		Title:       t.Title,
		Description: t.Description,
		Status:      string(t.Status),
		Priority:    t.Priority,
		DueAt:       t.DueAt,
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
	}
}

func toDomain(r *TaskRow) *models.Task {
	return &models.Task{
		ID:          r.ID,
		Title:       r.Title,
		Description: r.Description,
		Status:      models.TaskStatus(r.Status),
		Priority:    r.Priority,
		DueAt:       r.DueAt,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
	}
}

func (r *TaskRepo) Create(ctx context.Context, t *models.Task) (*models.Task, error) {
	if err := t.Validate(); err != nil {
		return nil, err
	}

	row := toRow(t)
	if err := r.db.WithContext(ctx).Error; err != nil {
		return nil, err
	}

	return toDomain(row), nil
}

func (r *TaskRepo) GetByID(ctx context.Context, id int64) (*models.Task, error) {
	var row TaskRow
	err := r.db.WithContext(ctx).First(&row, "id=?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, models.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return toDomain(&row), nil
}
