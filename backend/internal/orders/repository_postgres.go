package orders

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"go-shop-app-backend/internal/domain"
)

type postgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) Repository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) CreateOrder(ctx context.Context, userID int64, totalPrice int64) (*Order, error) {
	const query = `
        INSERT INTO orders (user_id, status, total_price)
        VALUES ($1, $2, $3)
        RETURNING id, user_id, status, total_price, created_at, updated_at
    `

	var o Order
	err := r.db.QueryRowContext(
		ctx,
		query,
		userID,
		OrderStatusPending,
		totalPrice,
	).Scan(
		&o.ID,
		&o.UserID,
		&o.Status,
		&o.TotalPrice,
		&o.CreatedAt,
		&o.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("insert order: %w", err)
	}

	return &o, nil
}

func (r *postgresRepository) AddOrderItems(ctx context.Context, orderID int64, items []CreateOrderItemInput) ([]OrderItem, error) {
	const query = `
        INSERT INTO order_items (order_id, product_id, quantity, unit_price, total_price)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id, order_id, product_id, quantity, unit_price, total_price
    `

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	var result []OrderItem

	for _, it := range items {
		var row OrderItem
		total := it.UnitPrice * it.Quantity

		err := tx.QueryRowContext(
			ctx,
			query,
			orderID,
			it.ProductID,
			it.Quantity,
			it.UnitPrice,
			total,
		).Scan(
			&row.ID,
			&row.OrderID,
			&row.ProductID,
			&row.Quantity,
			&row.UnitPrice,
			&row.TotalPrice,
		)
		if err != nil {
			return nil, fmt.Errorf("insert order item: %w", err)
		}

		result = append(result, row)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit order items tx: %w", err)
	}

	return result, nil
}

func (r *postgresRepository) GetByID(ctx context.Context, id int64) (*Order, []OrderItem, error) {
	const orderQuery = `
        SELECT id, user_id, status, total_price, created_at, updated_at
        FROM orders
        WHERE id = $1
    `

	const itemsQuery = `
        SELECT id, order_id, product_id, quantity, unit_price, total_price
        FROM order_items
        WHERE order_id = $1
    `

	var o Order
	err := r.db.QueryRowContext(ctx, orderQuery, id).Scan(
		&o.ID,
		&o.UserID,
		&o.Status,
		&o.TotalPrice,
		&o.CreatedAt,
		&o.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil, domain.ErrNotFound
		}
		return nil, nil, fmt.Errorf("get order by id: %w", err)
	}

	rows, err := r.db.QueryContext(ctx, itemsQuery, id)
	if err != nil {
		return nil, nil, fmt.Errorf("query order items: %w", err)
	}
	defer rows.Close()

	var items []OrderItem
	for rows.Next() {
		var it OrderItem
		if err := rows.Scan(
			&it.ID,
			&it.OrderID,
			&it.ProductID,
			&it.Quantity,
			&it.UnitPrice,
			&it.TotalPrice,
		); err != nil {
			return nil, nil, fmt.Errorf("scan order item: %w", err)
		}
		items = append(items, it)
	}

	if err := rows.Err(); err != nil {
		return nil, nil, fmt.Errorf("rows error: %w", err)
	}

	return &o, items, nil
}

func (r *postgresRepository) ListByUser(ctx context.Context, userID int64) ([]*Order, error) {
	const query = `
        SELECT id, user_id, status, total_price, created_at, updated_at
        FROM orders
        WHERE user_id = $1
        ORDER BY created_at DESC
    `

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("query orders by user: %w", err)
	}
	defer rows.Close()

	var orders []*Order
	for rows.Next() {
		var o Order
		if err := rows.Scan(
			&o.ID,
			&o.UserID,
			&o.Status,
			&o.TotalPrice,
			&o.CreatedAt,
			&o.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan order: %w", err)
		}
		orders = append(orders, &o)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return orders, nil
}

func (r *postgresRepository) UpdateStatus(ctx context.Context, id int64, status OrderStatus) error {
	const query = `
        UPDATE orders
        SET status = $1, updated_at = now()
        WHERE id = $2
    `

	res, err := r.db.ExecContext(ctx, query, status, id)
	if err != nil {
		return fmt.Errorf("update order status: %w", err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("update order status rows affected: %w", err)
	}

	if affected == 0 {
		return domain.ErrNotFound
	}

	return nil
}
