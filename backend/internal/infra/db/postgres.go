package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq" // postgres driver

	"go-shop-app-backend/internal/infra/config"
)

func NewPostgres(cfg *config.Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.DBDSN)
	if err != nil {
		return nil, fmt.Errorf("cannot open db: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("cannot ping db: %w", err)
	}

	return db, nil
}
