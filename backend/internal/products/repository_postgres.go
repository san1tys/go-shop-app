package products

import (
	"database/sql"
)

type postgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) Repository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) Create(input CreateProductInput) (Product, error) {
	query := `
        INSERT INTO products (name, description, price, stock)
        VALUES ($1, $2, $3, $4)
        RETURNING id, name, description, price, stock, created_at, updated_at
    `

	var p Product
	err := r.db.QueryRow(
		query,
		input.Name, input.Description, input.Price, input.Stock,
	).Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.Stock, &p.CreatedAt, &p.UpdatedAt)

	return p, err
}

func (r *postgresRepository) GetAll() ([]Product, error) {
	rows, err := r.db.Query(`SELECT id, name, description, price, stock, created_at, updated_at FROM products`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var p Product
		err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.Stock, &p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}

func (r *postgresRepository) GetByID(id int64) (Product, error) {
	var p Product
	err := r.db.QueryRow(
		`SELECT id, name, description, price, stock, created_at, updated_at FROM products WHERE id = $1`,
		id,
	).Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.Stock, &p.CreatedAt, &p.UpdatedAt)

	return p, err
}

func (r *postgresRepository) Update(id int64, input UpdateProductInput) (Product, error) {
	// simplified dynamic update
	p, err := r.GetByID(id)
	if err != nil {
		return p, err
	}

	if input.Name != nil {
		p.Name = *input.Name
	}
	if input.Description != nil {
		p.Description = *input.Description
	}
	if input.Price != nil {
		p.Price = *input.Price
	}
	if input.Stock != nil {
		p.Stock = *input.Stock
	}

	_, err = r.db.Exec(`
        UPDATE products
        SET name=$1, description=$2, price=$3, stock=$4, updated_at=now()
        WHERE id = $5
    `, p.Name, p.Description, p.Price, p.Stock, id)

	return p, err
}

func (r *postgresRepository) Delete(id int64) error {
	_, err := r.db.Exec(`DELETE FROM products WHERE id = $1`, id)
	return err
}
