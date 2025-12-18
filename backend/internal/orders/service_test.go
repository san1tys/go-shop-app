package orders

import (
	"context"
	"errors"
	"testing"

	"go-shop-app-backend/internal/domain"
)

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

func TestService_CreateOrder_Validation(t *testing.T) {
	repo := &mockOrderRepo{}
	svc := NewService(repo, nil) 

	tests := []struct {
		name    string
		userID  int64
		input   CreateOrderInput
		wantErr bool
	}{
		{
			name:   "invalid user id",
			userID: 0,
			input: CreateOrderInput{
				Items: []CreateOrderItemInput{{ProductID: 1, Quantity: 1, UnitPrice: 100}},
			},
			wantErr: true,
		},
		{
			name:    "no items",
			userID:  1,
			input:   CreateOrderInput{},
			wantErr: true,
		},
		{
			name:   "invalid product id",
			userID: 1,
			input: CreateOrderInput{
				Items: []CreateOrderItemInput{{ProductID: 0, Quantity: 1, UnitPrice: 100}},
			},
			wantErr: true,
		},
		{
			name:   "invalid quantity",
			userID: 1,
			input: CreateOrderInput{
				Items: []CreateOrderItemInput{{ProductID: 1, Quantity: 0, UnitPrice: 100}},
			},
			wantErr: true,
		},
		{
			name:   "invalid unit price",
			userID: 1,
			input: CreateOrderInput{
				Items: []CreateOrderItemInput{{ProductID: 1, Quantity: 1, UnitPrice: 0}},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := svc.CreateOrder(context.Background(), tt.userID, tt.input)
			if tt.wantErr && err == nil {
				t.Fatalf("expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.wantErr && !domain.IsValidationError(err) {
				t.Fatalf("expected validation error, got %v", err)
			}
		})
	}
}

func TestService_CreateOrder_Success(t *testing.T) {
	var capturedTotal int64
	repo := &mockOrderRepo{
		createOrderFn: func(ctx context.Context, userID int64, totalPrice int64) (*Order, error) {
			capturedTotal = totalPrice
			return &Order{
				ID:         1,
				UserID:     userID,
				Status:     OrderStatusPending,
				TotalPrice: totalPrice,
			}, nil
		},
		addOrderItemsFn: func(ctx context.Context, orderID int64, items []CreateOrderItemInput) ([]OrderItem, error) {
			if orderID != 1 {
				return nil, errors.New("unexpected order id")
			}
			result := make([]OrderItem, len(items))
			for i, it := range items {
				result[i] = OrderItem{
					ID:         int64(i + 1),
					OrderID:    orderID,
					ProductID:  it.ProductID,
					Quantity:   it.Quantity,
					UnitPrice:  it.UnitPrice,
					TotalPrice: it.UnitPrice * it.Quantity,
				}
			}
			return result, nil
		},
		getByIDFn: func(ctx context.Context, id int64) (*Order, []OrderItem, error) {
			return nil, nil, errors.New("not used")
		},
		listByUserFn: func(ctx context.Context, userID int64, limit, offset int) ([]*Order, error) {
			return nil, errors.New("not used")
		},
		updateStatusFn: func(ctx context.Context, id int64, status OrderStatus) error {
			return errors.New("not used")
		},
	}

	svc := NewService(repo, nil)

	input := CreateOrderInput{
		Items: []CreateOrderItemInput{
			{ProductID: 1, Quantity: 2, UnitPrice: 100},
			{ProductID: 2, Quantity: 1, UnitPrice: 50},
		},
	}

	order, items, err := svc.CreateOrder(context.Background(), 10, input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedTotal := int64(2*100 + 1*50)
	if capturedTotal != expectedTotal {
		t.Fatalf("expected total %d, got %d", expectedTotal, capturedTotal)
	}
	if order.TotalPrice != expectedTotal {
		t.Fatalf("order.TotalPrice = %d, want %d", order.TotalPrice, expectedTotal)
	}
	if len(items) != len(input.Items) {
		t.Fatalf("expected %d items, got %d", len(input.Items), len(items))
	}
}

func TestService_ListByUser_Validation(t *testing.T) {
	repo := &mockOrderRepo{
		listByUserFn: func(ctx context.Context, userID int64, limit, offset int) ([]*Order, error) {
			return []*Order{}, nil
		},
		createOrderFn: func(ctx context.Context, userID int64, totalPrice int64) (*Order, error) {
			return nil, errors.New("not used")
		},
		addOrderItemsFn: func(ctx context.Context, orderID int64, items []CreateOrderItemInput) ([]OrderItem, error) {
			return nil, errors.New("not used")
		},
		getByIDFn: func(ctx context.Context, id int64) (*Order, []OrderItem, error) {
			return nil, nil, errors.New("not used")
		},
		updateStatusFn: func(ctx context.Context, id int64, status OrderStatus) error {
			return errors.New("not used")
		},
	}
	svc := NewService(repo, nil)

	_, err := svc.ListByUser(context.Background(), 0, 1, 10)
	if err == nil || !domain.IsValidationError(err) {
		t.Fatalf("expected validation error for invalid user id, got %v", err)
	}

	_, err = svc.ListByUser(context.Background(), 1, 1, 101)
	if err == nil || !domain.IsValidationError(err) {
		t.Fatalf("expected validation error for too big pageSize, got %v", err)
	}
}

func TestService_Cancel_Validation(t *testing.T) {
	repo := &mockOrderRepo{
		updateStatusFn: func(ctx context.Context, id int64, status OrderStatus) error {
			return nil
		},
		createOrderFn: func(ctx context.Context, userID int64, totalPrice int64) (*Order, error) {
			return nil, errors.New("not used")
		},
		addOrderItemsFn: func(ctx context.Context, orderID int64, items []CreateOrderItemInput) ([]OrderItem, error) {
			return nil, errors.New("not used")
		},
		getByIDFn: func(ctx context.Context, id int64) (*Order, []OrderItem, error) {
			return nil, nil, errors.New("not used")
		},
		listByUserFn: func(ctx context.Context, userID int64, limit, offset int) ([]*Order, error) {
			return nil, errors.New("not used")
		},
	}
	svc := NewService(repo, nil)

	if err := svc.Cancel(context.Background(), 0); err == nil || !domain.IsValidationError(err) {
		t.Fatalf("expected validation error for invalid id, got %v", err)
	}
}
