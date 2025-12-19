package users

import (
	"context"
	"errors"
	"testing"
	"time"

	"go-shop-app-backend/internal/domain"
	"go-shop-app-backend/internal/infra/auth"
	"go-shop-app-backend/pkg/utils"
)

type mockUserRepo struct {
	createFn     func(ctx context.Context, email, name, passwordHash, role string) (*UserWithPassword, error)
	getByEmailFn func(ctx context.Context, email string) (*UserWithPassword, error)
	getByIDFn    func(ctx context.Context, id int64) (*UserWithPassword, error)
}

func (m *mockUserRepo) Create(ctx context.Context, email, name, passwordHash, role string) (*UserWithPassword, error) {
	return m.createFn(ctx, email, name, passwordHash, role)
}

func (m *mockUserRepo) GetByEmail(ctx context.Context, email string) (*UserWithPassword, error) {
	return m.getByEmailFn(ctx, email)
}

func (m *mockUserRepo) GetByID(ctx context.Context, id int64) (*UserWithPassword, error) {
	return m.getByIDFn(ctx, id)
}

func newTestJWTManager() *auth.Manager {
	// короткий TTL для тестов
	return auth.NewManager("test-secret", time.Minute)
}

func TestService_Register_Validation(t *testing.T) {
	repo := &mockUserRepo{
		getByEmailFn: func(ctx context.Context, email string) (*UserWithPassword, error) {
			return nil, domain.ErrNotFound
		},
		createFn: func(ctx context.Context, email, name, passwordHash, role string) (*UserWithPassword, error) {
			return &UserWithPassword{
				User: User{
					ID:    1,
					Email: email,
					Name:  name,
					Role:  role,
				},
				PasswordHash: passwordHash,
			}, nil
		},
		getByIDFn: func(ctx context.Context, id int64) (*UserWithPassword, error) {
			return nil, domain.ErrNotFound
		},
	}

	svc := NewService(repo, newTestJWTManager())

	tests := []struct {
		name    string
		input   RegisterInput
		wantErr bool
	}{
		{
			name: "empty email",
			input: RegisterInput{
				Email:    "",
				Name:     "Test",
				Password: "123456",
			},
			wantErr: true,
		},
		{
			name: "empty name",
			input: RegisterInput{
				Email:    "test@example.com",
				Name:     "",
				Password: "123456",
			},
			wantErr: true,
		},
		{
			name: "short password",
			input: RegisterInput{
				Email:    "test@example.com",
				Name:     "Test",
				Password: "123",
			},
			wantErr: true,
		},
		{
			name: "ok",
			input: RegisterInput{
				Email:    "test@example.com",
				Name:     "Test",
				Password: "123456",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := svc.Register(context.Background(), tt.input)
			if tt.wantErr && err == nil {
				t.Fatalf("expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !tt.wantErr {
				if resp == nil || resp.Token == "" {
					t.Fatalf("expected non-nil response with token")
				}
				if resp.User.Email != "test@example.com" {
					t.Fatalf("unexpected user email: %s", resp.User.Email)
				}
			}
		})
	}
}

func TestService_Login(t *testing.T) {
	hashed, err := utils.HashPassword("password")
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}

	user := &UserWithPassword{
		User: User{
			ID:    1,
			Email: "user@example.com",
			Name:  "User",
			Role:  string(domain.UserRoleUser),
		},
		PasswordHash: hashed,
	}

	repo := &mockUserRepo{
		getByEmailFn: func(ctx context.Context, email string) (*UserWithPassword, error) {
			if email == user.Email {
				return user, nil
			}
			return nil, domain.ErrNotFound
		},
		createFn: func(ctx context.Context, email, name, passwordHash, role string) (*UserWithPassword, error) {
			return nil, errors.New("not implemented")
		},
		getByIDFn: func(ctx context.Context, id int64) (*UserWithPassword, error) {
			return nil, domain.ErrNotFound
		},
	}

	svc := NewService(repo, newTestJWTManager())

	t.Run("invalid email", func(t *testing.T) {
		_, err := svc.Login(context.Background(), LoginInput{
			Email:    "",
			Password: "password",
		})
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
		if !domain.IsValidationError(err) {
			t.Fatalf("expected validation error, got %v", err)
		}
	})

	t.Run("invalid password", func(t *testing.T) {
		_, err := svc.Login(context.Background(), LoginInput{
			Email:    "user@example.com",
			Password: "",
		})
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
		if !domain.IsValidationError(err) {
			t.Fatalf("expected validation error, got %v", err)
		}
	})

	t.Run("wrong password", func(t *testing.T) {
		_, err := svc.Login(context.Background(), LoginInput{
			Email:    "user@example.com",
			Password: "wrong",
		})
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
		if !domain.IsValidationError(err) {
			t.Fatalf("expected validation error, got %v", err)
		}
	})

	t.Run("success", func(t *testing.T) {
		resp, err := svc.Login(context.Background(), LoginInput{
			Email:    "user@example.com",
			Password: "password",
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.Token == "" {
			t.Fatalf("expected non-empty token")
		}
		if resp.User.ID != user.ID {
			t.Fatalf("unexpected user id: %d", resp.User.ID)
		}
	})
}
