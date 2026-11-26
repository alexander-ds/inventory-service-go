package inventory

import "context"

// Service contiene la lógica de negocio. Hoy delega en Repository; las reglas
// de stock y movimientos se agregarán aquí, no en los handlers.
type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) ListProducts(ctx context.Context) ([]Product, error) {
	return s.repo.ListProducts(ctx)
}

func (s *Service) GetProduct(ctx context.Context, id string) (*Product, error) {
	return s.repo.GetProduct(ctx, id)
}

func (s *Service) CreateProduct(ctx context.Context, input ProductInput) (*Product, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	return s.repo.CreateProduct(ctx, input)
}

func (s *Service) UpdateProduct(ctx context.Context, id string, input ProductInput) (*Product, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	return s.repo.UpdateProduct(ctx, id, input)
}

func (s *Service) DeleteProduct(ctx context.Context, id string) error {
	return s.repo.DeleteProduct(ctx, id)
}
