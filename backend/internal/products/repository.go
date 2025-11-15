package products

type Repository interface {
	Create(p CreateProductInput) (Product, error)
	GetAll() ([]Product, error)
	GetByID(id int64) (Product, error)
	Update(id int64, input UpdateProductInput) (Product, error)
	Delete(id int64) error
}
