package repository

import "service-core/internal/modules/product/domain"

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
