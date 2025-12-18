package users

import (
	"context"
	"testing"
	"time"

	"go-shop-app-backend/internal/domain"
	"go-shop-app-backend/internal/infra/auth"
	"go-shop-app-backend/pkg/utils"
)

// Mock implementation of UserRepository
type mockUserRepo struct {
	getByEmailFn func(ctx context.Context, email string) (*UserWithPassword, error)
	createFn     func(ctx context.Context, email, name, hash, role string) (*UserWithPassword, error)
	getByIDFn    func(ctx context.Context, id int64) (*UserWithPassword, error)
}

func (m *mockUserRepo) GetByEmail(ctx context.Context, email string) (*UserWithPassword, error) {
	return m.getByEmailFn(ctx, email)
}

func (m *mockUserRepo) Create(ctx context.Context, email, name, hash, role string) (*UserWithPassword, error) {
	return m.createFn(ctx, email, name, hash, role)
}

func (m *mockUserRepo) GetByID(ctx context.Context, id int64) (*UserWithPassword, error) {
	return m.getByIDFn(ctx, id)
}

func TestService_Register_Success(t *testing.T) {
	repo := &mockUserRepo{
		getByEmailFn: func(ctx context.Context, email string) (*UserWithPassword, error) {
			return nil, domain.ErrNotFound
		},
		createFn: func(ctx context.Context, email, name, hash, role string) (*UserWithPassword, error) {
			return &UserWithPassword{
				User: User{
					ID:        1,
					Email:     email,
					Name:      name,
					Role:      role,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				PasswordHash: hash,
			}, nil
		},
	}

	jwtManager := auth.NewManager("testsecret", time.Hour)
	svc := NewService(repo, jwtManager)

	input := RegisterInput{
		Email:    "test@example.com",
		Name:     "Test User",
		Password: "password",
	}

	resp, err := svc.Register(context.Background(), input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.User.Email != input.Email {
		t.Fatalf("expected email %s, got %s", input.Email, resp.User.Email)
	}
}

func TestService_Login_Success(t *testing.T) {
	password := "password"
	hash, _ := utils.HashPassword(password)

	repo := &mockUserRepo{
		getByEmailFn: func(ctx context.Context, email string) (*UserWithPassword, error) {
			return &UserWithPassword{
				User: User{
					ID:        1,
					Email:     email,
					Name:      "Test User",
					Role:      "user",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				PasswordHash: hash,
			}, nil
		},
	}

	jwtManager := auth.NewManager("testsecret", time.Hour)
	svc := NewService(repo, jwtManager)

	input := LoginInput{
		Email:    "test@example.com",
		Password: password,
	}

	resp, err := svc.Login(context.Background(), input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.User.Email != input.Email {
		t.Fatalf("expected email %s, got %s", input.Email, resp.User.Email)
	}
}

func TestService_Register_Validation(t *testing.T) {
	repo := &mockUserRepo{
		getByEmailFn: func(ctx context.Context, email string) (*UserWithPassword, error) {
			return nil, domain.ErrNotFound
		},
		createFn: func(ctx context.Context, email, name, hash, role string) (*UserWithPassword, error) {
			return &UserWithPassword{
				User: User{
					ID:        1,
					Email:     email,
					Name:      name,
					Role:      role,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				PasswordHash: hash,
			}, nil
		},
	}

	jwtManager := auth.NewManager("testsecret", time.Hour)
	svc := NewService(repo, jwtManager)

	tests := []struct {
		name    string
		input   RegisterInput
		wantErr bool
	}{
		{"empty email", RegisterInput{Email: "", Name: "Name", Password: "password"}, true},
		{"empty name", RegisterInput{Email: "a@b.com", Name: "", Password: "password"}, true},
		{"short password", RegisterInput{Email: "a@b.com", Name: "Name", Password: "123"}, true},
		{"valid input", RegisterInput{Email: "a@b.com", Name: "Name", Password: "password"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := svc.Register(context.Background(), tt.input)
			if tt.wantErr && err == nil {
				t.Fatalf("expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}
