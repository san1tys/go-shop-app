package users

import "context"

type Repository interface {
	Create(ctx context.Context, email, name, passwordHash, role string) (*UserWithPassword, error)
	GetByEmail(ctx context.Context, email string) (*UserWithPassword, error)
	GetByID(ctx context.Context, id int64) (*UserWithPassword, error)
}
