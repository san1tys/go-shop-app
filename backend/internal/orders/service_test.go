package orders

import (
	"context"
	
	"testing"
	"time"

	
)

// Mock Repository
type mockOrderRepo struct {
	createOrderFn   func(ctx context.Context, userID int64, totalPrice int64) (*Order, error)
	addOrderItemsFn func(ctx context.Context, orderID int64, items []CreateOrderItemInput) ([]OrderItem, error)
	getByIDFn       func(ctx context.Context, id int64) (*Order, []OrderItem, error)
	listByUserFn    func(ctx context.Context, userID int64, limit, offset int) ([]*Order, error)
	updateStatusFn  func(ctx context.Context, id int64, status OrderStatus) error
}

func (m *mockOrderRepo) CreateOrder(ctx context.Context, userID int64, totalPrice int64) (*Order, error) {
	return m.createOrderFn(ctx, userID, totalPrice)
}

func (m *mockOrderRepo) AddOrderItems(ctx context.Context, orderID int64, items []CreateOrderItemInput) ([]OrderItem, error) {
	return m.addOrderItemsFn(ctx, orderID, items)
}

func (m *mockOrderRepo) GetByID(ctx context.Context, id int64) (*Order, []OrderItem, error) {
	return m.getByIDFn(ctx, id)
}

func (m *mockOrderRepo) ListByUser(ctx context.Context, userID int64, limit, offset int) ([]*Order, error) {
	return m.listByUserFn(ctx, userID, limit, offset)
}

func (m *mockOrderRepo) UpdateStatus(ctx context.Context, id int64, status OrderStatus) error {
	return m.updateStatusFn(ctx, id, status)
}

// ---- Тесты ----
func TestService_CreateOrder_Success(t *testing.T) {
	repo := &mockOrderRepo{
		createOrderFn: func(ctx context.Context, userID int64, totalPrice int64) (*Order, error) {
			return &Order{ID: 1, UserID: userID, TotalPrice: totalPrice, CreatedAt: time.Now(), UpdatedAt: time.Now()}, nil
		},
		addOrderItemsFn: func(ctx context.Context, orderID int64, items []CreateOrderItemInput) ([]OrderItem, error) {
			var result []OrderItem
			for i, it := range items {
				result = append(result, OrderItem{
					ID:         int64(i + 1),
					OrderID:    orderID,
					ProductID:  it.ProductID,
					Quantity:   it.Quantity,
					UnitPrice:  it.UnitPrice,
					TotalPrice: it.Quantity * it.UnitPrice,
				})
			}
			return result, nil
		},
	}
	svc := NewService(repo, nil)
	input := CreateOrderInput{Items: []CreateOrderItemInput{{ProductID: 1, Quantity: 2, UnitPrice: 100}}}

	order, items, err := svc.CreateOrder(context.Background(), 1, input)
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != len(input.Items) {
		t.Fatalf("expected %d items, got %d", len(input.Items), len(items))
	}
	if order.TotalPrice != 200 {
		t.Fatalf("expected total 200, got %d", order.TotalPrice)
	}
}
