package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Luc1808/TaskAPI/internal/repository"
	"github.com/Luc1808/TaskAPI/internal/repository/gorm"
	"github.com/Luc1808/TaskAPI/pkg/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func mustEnv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		log.Fatalf("missing env %s", k)
	}
	return v
}

func main() {
	dsn := mustEnv("DATABASE_URL")

	// Minimal GORM setup (no AutoMigrate here; DB is managed by migrations)
	gdb, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		log.Fatal("gorm open:", err)
	}

	repo := postgresgorm.NewTaskRepo(gdb)
	ctx := context.Background()

	// Create
	due := time.Now().Add(24 * time.Hour)
	t := &models.Task{
		Title:       "GORM adapter test",
		Description: "Create → Read → Update → Delete",
		Status:      models.StatusTodo,
		DueAt:       &due,
	}
	created, err := repo.Create(ctx, t)
	if err != nil {
		log.Fatal("create:", err)
	}
	fmt.Println("Created:", created.ID)

	// Get
	got, err := repo.GetByID(ctx, created.ID)
	if err != nil {
		log.Fatal("get:", err)
	}
	fmt.Println("Fetched:", got.Title, got.Status)

	// Update
	got.Status = models.StatusInProgress
	updated, err := repo.Update(ctx, got)
	if err != nil {
		log.Fatal("update:", err)
	}
	fmt.Println("Updated status:", updated.Status, "UpdatedAt:", updated.UpdatedAt)

	// List
	list, err := repo.List(ctx, repository.ListFilter{Status: &updated.Status}, repository.Pagination{Limit: 10})
	if err != nil {
		log.Fatal("list:", err)
	}
	fmt.Println("List count:", len(list))

	// Delete
	if err := repo.Delete(ctx, created.ID); err != nil {
		log.Fatal("delete:", err)
	}
	fmt.Println("Deleted:", created.ID)
}
