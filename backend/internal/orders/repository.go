package orders

import "context"

type Repository interface {
	CreateOrder(ctx context.Context, userID int64, totalPrice int64) (*Order, error)
	AddOrderItems(ctx context.Context, orderID int64, items []CreateOrderItemInput) ([]OrderItem, error)
	GetByID(ctx context.Context, id int64) (*Order, []OrderItem, error)
	ListByUser(ctx context.Context, userID int64, limit, offset int) ([]*Order, error)
	UpdateStatus(ctx context.Context, id int64, status OrderStatus) error
}
