package orders

import (
	"context"
	"fmt"

	"go-shop-app-backend/internal/domain"
)

type Service interface {
	CreateOrder(ctx context.Context, userID int64, input CreateOrderInput) (*Order, []OrderItem, error)
	GetByID(ctx context.Context, id int64) (*Order, []OrderItem, error)
	ListByUser(ctx context.Context, userID int64, page, pageSize int) ([]*Order, error)
	Cancel(ctx context.Context, id int64) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreateOrder(ctx context.Context, userID int64, input CreateOrderInput) (*Order, []OrderItem, error) {
	if userID <= 0 {
		return nil, nil, domain.NewValidationError("user_id is required")
	}
	if len(input.Items) == 0 {
		return nil, nil, domain.NewValidationError("at least one item is required")
	}

	var total int64
	for _, it := range input.Items {
		if it.ProductID <= 0 {
			return nil, nil, domain.NewValidationError("product_id must be positive")
		}
		if it.Quantity <= 0 {
			return nil, nil, domain.NewValidationError("quantity must be positive")
		}
		if it.UnitPrice <= 0 {
			return nil, nil, domain.NewValidationError("unit_price must be positive")
		}
		total += it.UnitPrice * it.Quantity
	}

	order, err := s.repo.CreateOrder(ctx, userID, total)
	if err != nil {
		return nil, nil, fmt.Errorf("create order: %w", err)
	}

	items, err := s.repo.AddOrderItems(ctx, order.ID, input.Items)
	if err != nil {
		return nil, nil, fmt.Errorf("add order items: %w", err)
	}

	order.Items = items

	return order, items, nil
}

func (s *service) GetByID(ctx context.Context, id int64) (*Order, []OrderItem, error) {
	if id <= 0 {
		return nil, nil, domain.NewValidationError("invalid id")
	}

	order, items, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, nil, err
	}

	order.Items = items

	return order, items, nil
}

func (s *service) ListByUser(ctx context.Context, userID int64, page, pageSize int) ([]*Order, error) {
	if userID <= 0 {
		return nil, domain.NewValidationError("invalid user_id")
	}

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		return nil, domain.NewValidationError("pageSize must be less than or equal to 100")
	}

	offset := (page - 1) * pageSize

	return s.repo.ListByUser(ctx, userID, pageSize, offset)
}

func (s *service) Cancel(ctx context.Context, id int64) error {
	if id <= 0 {
		return domain.NewValidationError("invalid id")
	}

	if err := s.repo.UpdateStatus(ctx, id, OrderStatusCancelled); err != nil {
		return fmt.Errorf("cancel order: %w", err)
	}

	return nil
}
