package postgresgorm

import (
	"context"
	"errors"
	"time"

	"github.com/Luc1808/TaskAPI/internal/repository"
	"github.com/Luc1808/TaskAPI/pkg/models"
	"gorm.io/gorm"
)

type TaskRepo struct {
	db *gorm.DB
}

func NewTaskRepo(db *gorm.DB) *TaskRepo {
	return &TaskRepo{db: db}
}

type TaskRow struct {
	ID          string     `gorm:"column:id;type:uuid;default:gen_random_uuid();primaryKey"`
	Title       string     `gorm:"column:title;type:text;not null"`
	Description string     `gorm:"column:description;type:text;not null;default:''"`
	Status      string     `gorm:"column:status;type:text;not null"`
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
	if err := r.db.WithContext(ctx).Create(row).Error; err != nil {
		return nil, err
	}

	return toDomain(row), nil
}

func (r *TaskRepo) GetByID(ctx context.Context, id string) (*models.Task, error) {
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

func (r *TaskRepo) List(ctx context.Context, f repository.ListFilter, p repository.Pagination) ([]models.Task, error) {
	q := r.db.WithContext(ctx).Model(&TaskRow{})

	if f.Status != nil {
		q = q.Where("status = ?", *f.Status)
	}
	if f.Search != "" {
		like := "%" + f.Search + "%"
		q = q.Where("(title ILIKE ? OR description ILIKE ?)", like, like)
	}

	limit := 50
	if p.Limit > 0 {
		limit = p.Limit
	}

	var rows []TaskRow
	if err := q.Order("created_at DESC").Limit(limit).Offset(p.Offset).Find(&rows).Error; err != nil {
		return nil, err
	}

	out := make([]models.Task, len(rows))
	for i := range rows {
		out[i] = *toDomain(&rows[i])
	}
	return out, nil
}

func (r *TaskRepo) Update(ctx context.Context, t *models.Task) (*models.Task, error) {
	if err := t.Validate(); err != nil {
		return nil, err
	}

	data := map[string]any{
		"title":       t.Title,
		"description": t.Description,
		"status":      string(t.Status),
		"due_at":      t.DueAt,
	}

	tx := r.db.WithContext(ctx).Model(&TaskRow{}).Where("id = ?", t.ID).Updates(data)
	if tx.Error != nil {
		return nil, tx.Error
	}
	if tx.RowsAffected == 0 {
		return nil, models.ErrNotFound
	}

	return r.GetByID(ctx, t.ID)
}

func (r *TaskRepo) Delete(ctx context.Context, id string) error {
	tx := r.db.WithContext(ctx).Where("id = ?", id).Delete(&TaskRow{})
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return models.ErrNotFound
	}
	return nil
}
