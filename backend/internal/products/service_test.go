package products

import (
	"context"
	"errors"
	"testing"

	"go-shop-app-backend/internal/domain"
)

type mockProductRepo struct {
	createFn  func(ctx context.Context, input CreateProductInput) (*Product, error)
	getAllFn  func(ctx context.Context, limit, offset int) ([]*Product, error)
	getByIDFn func(ctx context.Context, id int64) (*Product, error)
	updateFn  func(ctx context.Context, id int64, input UpdateProductInput) (*Product, error)
	deleteFn  func(ctx context.Context, id int64) error
}

func (m *mockProductRepo) Create(ctx context.Context, input CreateProductInput) (*Product, error) {
	return m.createFn(ctx, input)
}

func (m *mockProductRepo) GetAll(ctx context.Context, limit, offset int) ([]*Product, error) {
	return m.getAllFn(ctx, limit, offset)
}

func (m *mockProductRepo) GetByID(ctx context.Context, id int64) (*Product, error) {
	return m.getByIDFn(ctx, id)
}

func (m *mockProductRepo) Update(ctx context.Context, id int64, input UpdateProductInput) (*Product, error) {
	return m.updateFn(ctx, id, input)
}

func (m *mockProductRepo) Delete(ctx context.Context, id int64) error {
	return m.deleteFn(ctx, id)
}

func TestService_Create_Validation(t *testing.T) {
	repo := &mockProductRepo{
		createFn: func(ctx context.Context, input CreateProductInput) (*Product, error) {
			return &Product{ID: 1, Name: input.Name, Price: input.Price, Stock: input.Stock}, nil
		},
		getAllFn: func(ctx context.Context, limit, offset int) ([]*Product, error) {
			return nil, nil
		},
		getByIDFn: func(ctx context.Context, id int64) (*Product, error) {
			return nil, nil
		},
		updateFn: func(ctx context.Context, id int64, input UpdateProductInput) (*Product, error) {
			return nil, nil
		},
		deleteFn: func(ctx context.Context, id int64) error {
			return nil
		},
	}

	svc := NewService(repo)

	tests := []struct {
		name    string
		input   CreateProductInput
		wantErr bool
	}{
		{
			name: "empty name",
			input: CreateProductInput{
				Name:  "",
				Price: 100,
				Stock: 10,
			},
			wantErr: true,
		},
		{
			name: "non-positive price",
			input: CreateProductInput{
				Name:  "Product",
				Price: 0,
				Stock: 10,
			},
			wantErr: true,
		},
		{
			name: "negative stock",
			input: CreateProductInput{
				Name:  "Product",
				Price: 100,
				Stock: -1,
			},
			wantErr: true,
		},
		{
			name: "ok",
			input: CreateProductInput{
				Name:  "Product",
				Price: 100,
				Stock: 10,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, err := svc.Create(context.Background(), tt.input)
			if tt.wantErr && err == nil {
				t.Fatalf("expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !tt.wantErr && p == nil {
				t.Fatalf("expected product, got nil")
			}
		})
	}
}

func TestService_GetAll_Validation(t *testing.T) {
	repo := &mockProductRepo{
		getAllFn: func(ctx context.Context, limit, offset int) ([]*Product, error) {
			return []*Product{}, nil
		},
		createFn: func(ctx context.Context, input CreateProductInput) (*Product, error) {
			return nil, nil
		},
		getByIDFn: func(ctx context.Context, id int64) (*Product, error) {
			return nil, nil
		},
		updateFn: func(ctx context.Context, id int64, input UpdateProductInput) (*Product, error) {
			return nil, nil
		},
		deleteFn: func(ctx context.Context, id int64) error {
			return nil
		},
	}

	svc := NewService(repo)

	_, err := svc.GetAll(context.Background(), 1, 101)
	if err == nil || !domain.IsValidationError(err) {
		t.Fatalf("expected validation error for too big pageSize, got %v", err)
	}
}

func TestService_GetByID(t *testing.T) {
	product := &Product{ID: 1, Name: "P1", Price: 100, Stock: 10}

	repo := &mockProductRepo{
		getByIDFn: func(ctx context.Context, id int64) (*Product, error) {
			if id == product.ID {
				return product, nil
			}
			return nil, domain.ErrNotFound
		},
		createFn: func(ctx context.Context, input CreateProductInput) (*Product, error) {
			return nil, nil
		},
		getAllFn: func(ctx context.Context, limit, offset int) ([]*Product, error) {
			return nil, nil
		},
		updateFn: func(ctx context.Context, id int64, input UpdateProductInput) (*Product, error) {
			return nil, nil
		},
		deleteFn: func(ctx context.Context, id int64) error {
			return nil
		},
	}

	svc := NewService(repo)

	_, err := svc.GetByID(context.Background(), 0)
	if err == nil || !domain.IsValidationError(err) {
		t.Fatalf("expected validation error for id <= 0, got %v", err)
	}

	p, err := svc.GetByID(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.ID != product.ID {
		t.Fatalf("expected id %d, got %d", product.ID, p.ID)
	}

	_, err = svc.GetByID(context.Background(), 2)
	if err == nil || !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("expected ErrNotFound for unknown id, got %v", err)
	}
}

func TestService_Update_Validation(t *testing.T) {
	repo := &mockProductRepo{
		updateFn: func(ctx context.Context, id int64, input UpdateProductInput) (*Product, error) {
			return &Product{ID: id}, nil
		},
		createFn: func(ctx context.Context, input CreateProductInput) (*Product, error) {
			return nil, nil
		},
		getAllFn: func(ctx context.Context, limit, offset int) ([]*Product, error) {
			return nil, nil
		},
		getByIDFn: func(ctx context.Context, id int64) (*Product, error) {
			return nil, nil
		},
		deleteFn: func(ctx context.Context, id int64) error {
			return nil
		},
	}

	svc := NewService(repo)

	_, err := svc.Update(context.Background(), 0, UpdateProductInput{})
	if err == nil || !domain.IsValidationError(err) {
		t.Fatalf("expected validation error for id <= 0, got %v", err)
	}

	price := int64(0)
	_, err = svc.Update(context.Background(), 1, UpdateProductInput{Price: &price})
	if err == nil || !domain.IsValidationError(err) {
		t.Fatalf("expected validation error for non-positive price, got %v", err)
	}

	stock := int64(-1)
	_, err = svc.Update(context.Background(), 1, UpdateProductInput{Stock: &stock})
	if err == nil || !domain.IsValidationError(err) {
		t.Fatalf("expected validation error for negative stock, got %v", err)
	}
}

func TestService_Delete_Validation(t *testing.T) {
	repo := &mockProductRepo{
		deleteFn: func(ctx context.Context, id int64) error {
			if id == 1 {
				return nil
			}
			return domain.ErrNotFound
		},
		createFn: func(ctx context.Context, input CreateProductInput) (*Product, error) {
			return nil, nil
		},
		getAllFn: func(ctx context.Context, limit, offset int) ([]*Product, error) {
			return nil, nil
		},
		getByIDFn: func(ctx context.Context, id int64) (*Product, error) {
			return nil, nil
		},
		updateFn: func(ctx context.Context, id int64, input UpdateProductInput) (*Product, error) {
			return nil, nil
		},
	}

	svc := NewService(repo)

	if err := svc.Delete(context.Background(), 0); err == nil || !domain.IsValidationError(err) {
		t.Fatalf("expected validation error for id <= 0, got %v", err)
	}

	if err := svc.Delete(context.Background(), 2); err == nil || !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("expected ErrNotFound for unknown id, got %v", err)
	}

	if err := svc.Delete(context.Background(), 1); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}


