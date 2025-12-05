package domain

import "time"

type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusPaid      OrderStatus = "paid"
	OrderStatusCancelled OrderStatus = "cancelled"
)

type Order struct {
	ID         int64
	UserID     int64
	Status     OrderStatus
	TotalPrice int64
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type OrderItem struct {
	ID         int64
	OrderID    int64
	ProductID  int64
	Quantity   int64
	UnitPrice  int64
	TotalPrice int64
}
