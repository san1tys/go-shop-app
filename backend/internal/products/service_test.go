package products

import (
	"context"
	"testing"
	"time"

	
)

// Mock Repository
type mockProductRepo struct {
	createFn func(ctx context.Context, input CreateProductInput) (*Product, error)
	getAllFn func(ctx context.Context, limit, offset int) ([]*Product, error)
	getByIDFn func(ctx context.Context, id int64) (*Product, error)
	updateFn func(ctx context.Context, id int64, input UpdateProductInput) (*Product, error)
	deleteFn func(ctx context.Context, id int64) error
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

// ---- Тесты ----
func TestService_Create_Success(t *testing.T) {
	repo := &mockProductRepo{
		createFn: func(ctx context.Context, input CreateProductInput) (*Product, error) {
			return &Product{ID: 1, Name: input.Name, Price: input.Price, Stock: input.Stock, CreatedAt: time.Now(), UpdatedAt: time.Now()}, nil
		},
	}
	svc := NewService(repo)
	input := CreateProductInput{Name: "Test", Price: 100, Stock: 10}
	p, err := svc.Create(context.Background(), input)
	if err != nil {
		t.Fatal(err)
	}
	if p.Name != "Test" {
		t.Fatalf("expected name Test, got %s", p.Name)
	}
}
