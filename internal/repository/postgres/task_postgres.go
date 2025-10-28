package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/Luc1808/TaskAPI/internal/repository"
	"github.com/Luc1808/TaskAPI/pkg/models"
	"github.com/jmoiron/sqlx"
)

type TaskRepo struct {
	db *sqlx.DB
}

func NewTaskRepo(db *sqlx.DB) *TaskRepo {
	return &TaskRepo{db: db}
}

func (r *TaskRepo) Create(ctx context.Context, t *models.Task) (*models.Task, error) {
	if err := t.Validate(); err != nil {
		return nil, err
	}

	const q = `
		INSERT INTO public.tasks (title, description, status, due_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at;
		`
	if err := r.db.QueryRowContext(ctx, q,
		t.Title, t.Description, t.Status,
		t.DueAt).Scan(&t.ID, &t.CreatedAt, &t.UpdatedAt); err != nil {
		return nil, err
	}

	return t, nil
}

func (r *TaskRepo) GetByID(ctx context.Context, id string) (*models.Task, error) {
	const q = `
		SELECT id, title, description, status, due_at, created_at, updated_at
		FROM public.tasks
		WHERE id = $1;
		`
	var out models.Task
	if err := r.db.GetContext(ctx, &out, q, id); err != nil {
		// In case there's no rows, it could return "no rows"
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNotFound
		}
		return nil, err
	}

	return &out, nil
}

func (r *TaskRepo) List(ctx context.Context, f repository.ListFilter, p repository.Pagination) ([]models.Task, error) {
	base := `
	SELECT id, title, description, status, due_at, created_at, updated_at
	FROM public.tasks
	`

	where := []string{"1=1"}
	args := []any{}
	arg := 1

	if f.Status != nil {
		where = append(where, fmt.Sprintf("status = $%d", arg))
		args = append(args, *f.Status)
		arg++
	}
	if f.Search != "" {
		where = append(where, fmt.Sprintf("(title ILIKE $%d OR description ILIKE $%d)", arg, arg))
		args = append(args, "%"+f.Search+"%")
		arg++
	}

	order := "ORDER BY created_at"
	limit := 20
	if p.Limit > 0 {
		limit = p.Limit
	}

	query := fmt.Sprintf("%s WHERE %s %s LIMIT %d OFFSET %d;",
		base, strings.Join(where, " AND "), order, limit, p.Offset)

	out := []models.Task{}
	if err := r.db.SelectContext(ctx, &out, query, args...); err != nil {
		return nil, err
	}

	return out, nil
}

func (r *TaskRepo) Update(ctx context.Context, t *models.Task) (*models.Task, error) {
	if err := t.Validate(); err != nil {
		return nil, err
	}

	const q = `
		UPDATE public.tasks
		SET title = $1,
		description = $2,
		status = $3,
		due_at = $4,
		updated_at = now()
		WHERE id = $5
		RETURNING created_at, updated_at;
		`
	var createdAt, updatedAt = t.CreatedAt, t.UpdatedAt
	if err := r.db.QueryRowxContext(ctx, q, t.Title, t.Description, t.Status, t.DueAt, t.ID).Scan(&createdAt, &updatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNotFound
		}
		return nil, err
	}
	t.CreatedAt = createdAt
	t.UpdatedAt = updatedAt

	return t, nil
}

func (r *TaskRepo) Delete(ctx context.Context, id string) error {
	const q = `DELETE FROM public.tasks WHERE id = $1;`
	res, err := r.db.ExecContext(ctx, q, id)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return models.ErrNotFound
	}

	return nil
}
