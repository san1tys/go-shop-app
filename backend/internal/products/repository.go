package products
// repository.go
import "context"

type Repository interface {
	Create(ctx context.Context, input CreateProductInput) (*Product, error)
	GetAll(ctx context.Context) ([]*Product, error)
	GetByID(ctx context.Context, id int64) (*Product, error)
	Update(ctx context.Context, id int64, input UpdateProductInput) (*Product, error)
	Delete(ctx context.Context, id int64) error
}
