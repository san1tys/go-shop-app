package products
// repository_postgres.go
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

func (r *postgresRepository) Create(ctx context.Context, input CreateProductInput) (*Product, error) {
	const query = `
        INSERT INTO products (name, description, price, stock)
        VALUES ($1, $2, $3, $4)
        RETURNING id, name, description, price, stock, created_at, updated_at
    `

	var p Product
	err := r.db.QueryRowContext(
		ctx,
		query,
		input.Name,
		input.Description,
		input.Price,
		input.Stock,
	).Scan(
		&p.ID,
		&p.Name,
		&p.Description,
		&p.Price,
		&p.Stock,
		&p.CreatedAt,
		&p.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("insert product: %w", err)
	}

	return &p, nil
}

func (r *postgresRepository) GetAll(ctx context.Context) ([]*Product, error) {
	const query = `
        SELECT id, name, description, price, stock, created_at, updated_at
        FROM products
        ORDER BY id
    `

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query products: %w", err)
	}
	defer rows.Close()

	var products []*Product
	for rows.Next() {
		var p Product
		if err := rows.Scan(
			&p.ID,
			&p.Name,
			&p.Description,
			&p.Price,
			&p.Stock,
			&p.CreatedAt,
			&p.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan product: %w", err)
		}
		products = append(products, &p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return products, nil
}

func (r *postgresRepository) GetByID(ctx context.Context, id int64) (*Product, error) {
	const query = `
        SELECT id, name, description, price, stock, created_at, updated_at
        FROM products
        WHERE id = $1
    `

	var p Product
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&p.ID,
		&p.Name,
		&p.Description,
		&p.Price,
		&p.Stock,
		&p.CreatedAt,
		&p.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("get product by id: %w", err)
	}

	return &p, nil
}

func (r *postgresRepository) Update(ctx context.Context, id int64, input UpdateProductInput) (*Product, error) {
	current, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if input.Name != nil {
		current.Name = *input.Name
	}
	if input.Description != nil {
		current.Description = *input.Description
	}
	if input.Price != nil {
		current.Price = *input.Price
	}
	if input.Stock != nil {
		current.Stock = *input.Stock
	}

	const query = `
        UPDATE products
        SET name = $1,
            description = $2,
            price = $3,
            stock = $4,
            updated_at = now()
        WHERE id = $5
    `

	res, err := r.db.ExecContext(
		ctx,
		query,
		current.Name,
		current.Description,
		current.Price,
		current.Stock,
		id,
	)
	if err != nil {
		return nil, fmt.Errorf("update product: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("update product rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return nil, domain.ErrNotFound
	}

	// updated_at обновляется в БД, можно при желании перечитать запись
	// но чаще в CRUD это не критично. Если хочешь точно свежее updated_at —
	// можно повторно вызвать GetByID(ctx, id) и вернуть её.
	return current, nil
}

func (r *postgresRepository) Delete(ctx context.Context, id int64) error {
	const query = `DELETE FROM products WHERE id = $1`

	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete product: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("delete product rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return domain.ErrNotFound
	}

	return nil
}
