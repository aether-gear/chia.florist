package repository

import (
	"service-core/internal/modules/product/domain"

	"github.com/google/uuid"
)

type FindProductParams struct {
	Page  int
	Limit int
	ID    *string
	Name  *string
}

type ProductWithInventory struct {
	Product   domain.Product
	Inventory struct {
		Stock         int
		ReservedStock int
	}
}

type ProductRepository interface {
	FindProducts(params FindProductParams) ([]ProductWithInventory, int, error)
	GetByID(id uuid.UUID) (*ProductWithInventory, error)
	FindByIDs(IDs []uuid.UUID) ([]ProductWithInventory, error)
	CreateProduct(product *domain.Product) error
	CreateInventory(inventory *domain.Inventory) error
}
