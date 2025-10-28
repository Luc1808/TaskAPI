package models

import (
	"errors"
	"fmt"
	"time"
)

type TaskStatus string

const (
	StatusTodo       TaskStatus = "todo"
	StatusInProgress TaskStatus = "in_progress"
	StatusDone       TaskStatus = "done"
)

type Task struct {
	ID          string     `db:"id" json:"id"`
	Title       string     `db:"title" json:"title"`
	Description string     `db:"description" json:"description"`
	Status      TaskStatus `db:"status" json:"status"`
	DueAt       *time.Time `db:"due_at" json:"due_at"`
	CreatedAt   time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at" json:"updated_at"`
}

// Errors to be used everywhere
var (
	ErrNotFound   = errors.New("task not found")
	ErrValidation = errors.New("validation error")
)

func (t *Task) Validate() error {
	if len(t.Title) == 0 {
		return fmt.Errorf("%w: title is required", ErrValidation)
	}
	switch t.Status {
	case StatusTodo, StatusInProgress, StatusDone:
	default:
		return fmt.Errorf("%w: invalid status %q", ErrValidation, t.Status)
	}
	return nil
}
