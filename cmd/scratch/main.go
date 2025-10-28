package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Luc1808/TaskAPI/internal/repository"
	"github.com/Luc1808/TaskAPI/internal/repository/postgres"
	"github.com/Luc1808/TaskAPI/pkg/models"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

func mustEnv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		log.Fatalf("missing env %s", k)
	}
	return v
}

func main() {
	// Example DSN; prefer DATABASE_URL or build from parts.
	// e.g. DATABASE_URL="postgres://user:pass@localhost:5432/taskapi?sslmode=disable"
	dsn := mustEnv("DATABASE_URL")

	db := sqlx.MustConnect("pgx", dsn)
	defer db.Close()

	repo := postgres.NewTaskRepo(db)

	ctx := context.Background()

	// 1) Create
	due := time.Now().Add(48 * time.Hour)
	t := &models.Task{
		Title:       "Phase 4 – wire repo",
		Description: "Implement domain + repo + adapter",
		Status:      models.StatusTodo,
		DueAt:       &due,
	}
	created, err := repo.Create(ctx, t)
	if err != nil {
		log.Fatal("create:", err)
	}
	fmt.Println("Created ID:", created.ID)

	// 2) GetByID
	got, err := repo.GetByID(ctx, created.ID)
	if err != nil {
		log.Fatal("get:", err)
	}
	fmt.Println("Fetched title:", got.Title)

	// 3) Update (status)
	got.Status = models.StatusInProgress
	updated, err := repo.Update(ctx, got)
	if err != nil {
		log.Fatal("update:", err)
	}
	fmt.Println("Updated status:", updated.Status)

	// 4) List (simple filter + pagination)
	list, err := repo.List(ctx, repository.ListFilter{
		Status: &updated.Status,
	}, repository.Pagination{Limit: 10, Offset: 0})
	if err != nil {
		log.Fatal("list:", err)
	}
	fmt.Println("List count:", len(list))

	// 5) Delete
	if err := repo.Delete(ctx, created.ID); err != nil {
		log.Fatal("delete:", err)
	}
	fmt.Println("Deleted:", created.ID)
}
