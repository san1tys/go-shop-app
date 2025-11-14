package products

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(input CreateProductInput) (Product, error) {
	return s.repo.Create(input)
}

func (s *Service) GetAll() ([]Product, error) {
	return s.repo.GetAll()
}

func (s *Service) GetByID(id int64) (Product, error) {
	return s.repo.GetByID(id)
}

func (s *Service) Update(id int64, input UpdateProductInput) (Product, error) {
	return s.repo.Update(id, input)
}

func (s *Service) Delete(id int64) error {
	return s.repo.Delete(id)
}
