package products

import (
	"context"
	"errors"
	"fmt"

	"go-shop-app-backend/internal/domain"
)

type Service interface {
	Create(ctx context.Context, input CreateProductInput) (*Product, error)
	GetAll(ctx context.Context, page, pageSize int) ([]*Product, error)
	GetByID(ctx context.Context, id int64) (*Product, error)
	Update(ctx context.Context, id int64, input UpdateProductInput) (*Product, error)
	Delete(ctx context.Context, id int64) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(ctx context.Context, input CreateProductInput) (*Product, error) {
	if input.Name == "" {
		return nil, domain.NewValidationError("name is required")
	}
	if input.Price <= 0 {
		return nil, domain.NewValidationError("price must be greater than 0")
	}
	if input.Stock < 0 {
		return nil, domain.NewValidationError("stock cannot be negative")
	}

	product, err := s.repo.Create(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("create product: %w", err)
	}

	return product, nil
}

func (s *service) GetAll(ctx context.Context, page, pageSize int) ([]*Product, error) {
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

	products, err := s.repo.GetAll(ctx, pageSize, offset)
	if err != nil {
		return nil, fmt.Errorf("get all products: %w", err)
	}

	return products, nil
}

func (s *service) GetByID(ctx context.Context, id int64) (*Product, error) {
	if id <= 0 {
		return nil, domain.NewValidationError("invalid id")
	}

	product, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("get product by id: %w", err)
	}

	return product, nil
}

func (s *service) Update(ctx context.Context, id int64, input UpdateProductInput) (*Product, error) {
	if id <= 0 {
		return nil, domain.NewValidationError("invalid id")
	}

	if input.Price != nil && *input.Price <= 0 {
		return nil, domain.NewValidationError("price must be greater than 0")
	}
	if input.Stock != nil && *input.Stock < 0 {
		return nil, domain.NewValidationError("stock cannot be negative")
	}

	product, err := s.repo.Update(ctx, id, input)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("update product: %w", err)
	}

	return product, nil
}

func (s *service) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return domain.NewValidationError("invalid id")
	}

	err := s.repo.Delete(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return domain.ErrNotFound
		}
		return fmt.Errorf("delete product: %w", err)
	}

	return nil
}
