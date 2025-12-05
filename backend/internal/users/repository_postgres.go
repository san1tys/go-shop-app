package users

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"go-shop-app-backend/internal/domain"
)

type postgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) Repository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) Create(ctx context.Context, email, name, passwordHash, role string) (*UserWithPassword, error) {
	const query = `
        INSERT INTO users (email, name, password_hash, role)
        VALUES ($1, $2, $3, $4)
        RETURNING id, email, name, password_hash, role, created_at, updated_at
    `

	var u UserWithPassword
	err := r.db.QueryRowContext(
		ctx,
		query,
		email,
		name,
		passwordHash,
		role,
	).Scan(
		&u.ID,
		&u.Email,
		&u.Name,
		&u.PasswordHash,
		&u.Role,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
	if err != nil {
		// здесь можно разобрать ошибку на уникальность email по коду PG,
		// но для простоты пока оставим как есть
		return nil, fmt.Errorf("insert user: %w", err)
	}

	return &u, nil
}

func (r *postgresRepository) GetByEmail(ctx context.Context, email string) (*UserWithPassword, error) {
	const query = `
        SELECT id, email, name, password_hash, role, created_at, updated_at
        FROM users
        WHERE email = $1
    `

	var u UserWithPassword
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&u.ID,
		&u.Email,
		&u.Name,
		&u.PasswordHash,
		&u.Role,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("get user by email: %w", err)
	}

	return &u, nil
}

func (r *postgresRepository) GetByID(ctx context.Context, id int64) (*UserWithPassword, error) {
	const query = `
        SELECT id, email, name, password_hash, role, created_at, updated_at
        FROM users
        WHERE id = $1
    `

	var u UserWithPassword
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&u.ID,
		&u.Email,
		&u.Name,
		&u.PasswordHash,
		&u.Role,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("get user by id: %w", err)
	}

	return &u, nil
}
