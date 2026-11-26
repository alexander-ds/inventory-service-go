package inventory

import "context"

// Repository abstrae la persistencia para poder mockearla en los tests.
type Repository interface {
	ListProducts(ctx context.Context) ([]Product, error)
	GetProduct(ctx context.Context, id string) (*Product, error)
	CreateProduct(ctx context.Context, input ProductInput) (*Product, error)
	UpdateProduct(ctx context.Context, id string, input ProductInput) (*Product, error)
	DeleteProduct(ctx context.Context, id string) error
}
