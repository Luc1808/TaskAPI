package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Luc1808/TaskAPI/internal/repository"
	"github.com/Luc1808/TaskAPI/pkg/models"
)

type fakeTaskRepo struct {
	store map[string]models.Task
}

func newFakeTaskRepo() *fakeTaskRepo {
	return &fakeTaskRepo{
		store: make(map[string]models.Task),
	}
}

func (f *fakeTaskRepo) Create(ctx context.Context, t *models.Task) (*models.Task, error) {
	copy := *t
	f.store[t.ID] = copy
	return &copy, nil
}

func (f *fakeTaskRepo) GetByID(ctx context.Context, id string) (*models.Task, error) {
	t, ok := f.store[id]
	if !ok {
		return &models.Task{}, models.ErrNotFound
	}
	copy := t
	return &copy, nil
}

func (f *fakeTaskRepo) List(ctx context.Context, filter repository.ListFilter, pagination repository.Pagination) ([]models.Task, error) {
	out := make([]models.Task, 0, len(f.store))
	for _, v := range f.store {
		out = append(out, v)
	}
	return out, nil
}

func (f *fakeTaskRepo) Update(ctx context.Context, t *models.Task) (*models.Task, error) {
	_, ok := f.store[t.ID]
	if !ok {
		return nil, models.ErrNotFound
	}

	copy := *t
	f.store[t.ID] = copy

	return &copy, nil
}

func (f *fakeTaskRepo) Delete(ctx context.Context, id string) error {
	if _, ok := f.store[id]; !ok {
		return models.ErrNotFound
	}
	delete(f.store, id)
	return nil
}

// --- TESTS ---

func testCreateTask_RejectsEmptyTitle(t *testing.T) {
	repo := newFakeTaskRepo()
	svc := NewTaskService(repo)

	_, err := svc.CreateTask(context.Background(), CreateTaskInput{
		Title: "",
	})
	if !errors.Is(err, ErrInvalidTitle) {
		t.Fatalf("expected ErrInvalidtitle, got %v", err)
	}
}

func TestCreateTask_EnforcesMaxTitleLength(t *testing.T) {
	repo := newFakeTaskRepo()
	svc := NewTaskService(repo)

	longTitle := make([]byte, 141)
	for i := 0; i < 141; i++ {
		longTitle[i] = 'a'
	}

	_, err := svc.CreateTask(context.Background(), CreateTaskInput{
		Title: string(longTitle),
	})
	if !errors.Is(err, ErrInvalidTitle) {
		t.Fatalf("expected ErrInvalidTitle for long title, got %v", err)
	}
}

func TestCreateTask_DefaultStatusTodo(t *testing.T) {
	repo := newFakeTaskRepo()
	svc := NewTaskService(repo)

	task, err := svc.CreateTask(context.Background(), CreateTaskInput{
		Title: "Buy milk",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if task.Status != "todo" {
		t.Fatalf("expected default status 'todo', got %q", task.Status)
	}
}

func TestUpdateTask_RefreshesUpdatedAt(t *testing.T) {
	repo := newFakeTaskRepo()
	svc := NewTaskService(repo)

	created, err := svc.CreateTask(context.Background(), CreateTaskInput{
		Title:  "Original title",
		Status: "todo",
	})
	if err != nil {
		t.Fatalf("create err: %v", err)
	}

	before := created.UpdatedAt

	time.Sleep(10 * time.Millisecond)

	newTitle := "Changed title"
	updated, err := svc.UpdateTask(context.Background(), created.ID, UpdateTaskInput{
		Title: &newTitle,
	})
	if err != nil {
		t.Fatalf("update err: %v", err)
	}

	if !updated.UpdatedAt.After(before) {
		t.Fatalf("expected UpdatedAt to be refreshed. before=%v after=%v", before, updated.UpdatedAt)
	}

	if updated.Title != "Changed title" {
		t.Fatalf("expected title to be updated, got %q", updated.Title)
	}
}
