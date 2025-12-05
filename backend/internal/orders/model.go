package orders

import "time"

type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusPaid      OrderStatus = "paid"
	OrderStatusCancelled OrderStatus = "cancelled"
)

type OrderItem struct {
	ID         int64 `json:"id"`
	OrderID    int64 `json:"order_id"`
	ProductID  int64 `json:"product_id"`
	Quantity   int64 `json:"quantity"`
	UnitPrice  int64 `json:"unit_price"`
	TotalPrice int64 `json:"total_price"`
}

type Order struct {
	ID         int64       `json:"id"`
	UserID     int64       `json:"user_id"`
	Status     OrderStatus `json:"status"`
	TotalPrice int64       `json:"total_price"`
	Items      []OrderItem `json:"items,omitempty"`
	CreatedAt  time.Time   `json:"created_at"`
	UpdatedAt  time.Time   `json:"updated_at"`
}

type CreateOrderItemInput struct {
	ProductID int64 `json:"product_id"`
	Quantity  int64 `json:"quantity"`
	UnitPrice int64 `json:"unit_price"`
}

type CreateOrderInput struct {
	Items []CreateOrderItemInput `json:"items"`
}
