package repository

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func InitDB() (*sql.DB, error) {
	// Validate env variables
	required := []string{"DB_HOST", "DB_PORT", "DB_USER", "DB_PASSWORD", "DB_NAME", "DB_SSLMODE"}
	for _, k := range required {
		if os.Getenv(k) == "" {
			return nil, fmt.Errorf("missing required env var %s", k)
		}
	}

	// Build DSN for pgx (database/sql)
	dsn := fmt.Sprintf(
		"user=%s password=%s host=%s port=%s dbname=%s sslmode=%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_SSLMODE"),
	)

	// Open pool and set limits
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("db open faile: %w", err)
	}
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(30 * time.Minute)

	// Ping with small retry window
	if err := pingWithRetry(db, 5, 600*time.Millisecond); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("db ping failed: %w", err)
	}

	return db, nil
}

func pingWithRetry(db *sql.DB, attempts int, backoff time.Duration) error {
	var last error
	for range attempts {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		last = db.PingContext(ctx)
		cancel()
		if last == nil {
			return nil
		}
		time.Sleep(backoff)
	}
	return last
}
