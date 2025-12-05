package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	_ "github.com/lib/pq" // postgres driver

	"go-shop-app-backend/internal/infra/config"
)

const (
	pingTimeout      = 3 * time.Second
	maxPingAttempts  = 5
	pingRetryBackoff = 2 * time.Second
)

func NewPostgres(cfg *config.Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.DBDSN)
	if err != nil {
		return nil, fmt.Errorf("cannot open db: %w", err)
	}

	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(30 * time.Minute)

	if err := pingWithRetry(db); err != nil {
		return nil, fmt.Errorf("db ping failed after retries: %w", err)
	}

	return db, nil
}

func pingWithRetry(db *sql.DB) error {
	var lastErr error

	for attempt := 1; attempt <= maxPingAttempts; attempt++ {
		ctx, cancel := context.WithTimeout(context.Background(), pingTimeout)

		err := db.PingContext(ctx)
		cancel()

		if err == nil {
			return nil
		}

		lastErr = err
		fmt.Printf("db ping attempt %d/%d failed: %v\n", attempt, maxPingAttempts, err)

		time.Sleep(pingRetryBackoff)
	}

	return lastErr
}

func HealthCheck(db *sql.DB) error {
	ctx, cancel := context.WithTimeout(context.Background(), pingTimeout)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return errors.New("database unreachable")
	}

	return nil
}
