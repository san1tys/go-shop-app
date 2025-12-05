package users

import (
	"context"
	"fmt"
	"strings"

	"go-shop-app-backend/internal/domain"
	"go-shop-app-backend/internal/infra/auth"
	"go-shop-app-backend/pkg/utils"
)

type Service interface {
	Register(ctx context.Context, input RegisterInput) (*AuthResponse, error)
	Login(ctx context.Context, input LoginInput) (*AuthResponse, error)
	GetByID(ctx context.Context, id int64) (*User, error)
}

type service struct {
	repo       Repository
	jwtManager *auth.Manager
}

func NewService(repo Repository, jwtManager *auth.Manager) Service {
	return &service{
		repo:       repo,
		jwtManager: jwtManager,
	}
}

func (s *service) Register(ctx context.Context, input RegisterInput) (*AuthResponse, error) {
	email := strings.TrimSpace(strings.ToLower(input.Email))
	name := strings.TrimSpace(input.Name)

	if email == "" {
		return nil, domain.NewValidationError("email is required")
	}
	if name == "" {
		return nil, domain.NewValidationError("name is required")
	}
	if len(input.Password) < 6 {
		return nil, domain.NewValidationError("password must be at least 6 characters")
	}

	if existing, err := s.repo.GetByEmail(ctx, email); err == nil && existing != nil {
		return nil, domain.NewValidationError("email is already in use")
	}

	hash, err := utils.HashPassword(input.Password)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	u, err := s.repo.Create(ctx, email, name, hash, string(domain.UserRoleUser))
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	token, err := s.jwtManager.GenerateToken(u.ID, u.Role)
	if err != nil {
		return nil, fmt.Errorf("generate token: %w", err)
	}

	resp := &AuthResponse{
		Token: token,
		User: User{
			ID:        u.ID,
			Email:     u.Email,
			Name:      u.Name,
			Role:      u.Role,
			CreatedAt: u.CreatedAt,
			UpdatedAt: u.UpdatedAt,
		},
	}

	return resp, nil
}

func (s *service) Login(ctx context.Context, input LoginInput) (*AuthResponse, error) {
	email := strings.TrimSpace(strings.ToLower(input.Email))

	if email == "" {
		return nil, domain.NewValidationError("email is required")
	}
	if input.Password == "" {
		return nil, domain.NewValidationError("password is required")
	}

	u, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		if err == domain.ErrNotFound {
			return nil, domain.NewValidationError("invalid email or password")
		}
		return nil, fmt.Errorf("get user by email: %w", err)
	}

	if err := utils.CheckPassword(u.PasswordHash, input.Password); err != nil {
		return nil, domain.NewValidationError("invalid email or password")
	}

	token, err := s.jwtManager.GenerateToken(u.ID, u.Role)
	if err != nil {
		return nil, fmt.Errorf("generate token: %w", err)
	}

	resp := &AuthResponse{
		Token: token,
		User: User{
			ID:        u.ID,
			Email:     u.Email,
			Name:      u.Name,
			Role:      u.Role,
			CreatedAt: u.CreatedAt,
			UpdatedAt: u.UpdatedAt,
		},
	}

	return resp, nil
}

func (s *service) GetByID(ctx context.Context, id int64) (*User, error) {
	if id <= 0 {
		return nil, domain.NewValidationError("invalid id")
	}

	u, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &User{
		ID:        u.ID,
		Email:     u.Email,
		Name:      u.Name,
		Role:      u.Role,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}, nil
}
