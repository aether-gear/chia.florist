package http

import "service-core/internal/modules/product/repository"

func ToListResponse(p repository.ProductWithInventory) ProductOverviewResponse {
	return ProductOverviewResponse{
		ID:            p.Product.ID,
		SKU:           p.Product.SKU,
		Name:          p.Product.Name,
		Status:        ProductStatusDTO(p.Product.Status),
		Price:         p.Product.Price,
		Stock:         p.Inventory.Stock,
		ReservedStock: p.Inventory.ReservedStock,
	}
}

func ToDetailResponse(p repository.ProductWithInventory) ProductDetailResponse {
	return ProductDetailResponse{
		ID:            p.Product.ID,
		SKU:           p.Product.SKU,
		Name:          p.Product.Name,
		Description:   p.Product.Description,
		Status:        ProductStatusDTO(p.Product.Status),
		Price:         p.Product.Price,
		Weight:        p.Product.Weight,
		Stock:         p.Inventory.Stock,
		ReservedStock: p.Inventory.ReservedStock,
		CreatedAt:     p.Product.CreatedAt,
		UpdatedAt:     p.Product.UpdatedAt,
		ArchivedAt:    p.Product.ArchivedAt,
	}
}
