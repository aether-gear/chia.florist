package http

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"service-core/internal/modules/product/domain"
	"service-core/internal/modules/product/repository"
	"service-core/internal/modules/product/usecase"

	"service-core/internal/common/errors"
	apphttp "service-core/internal/common/http"

	"github.com/google/uuid"
)

type ProductHandler struct {
	findProducts  *usecase.FindProductsUsecase
	getProduct    *usecase.GetProductUsecase
	createProduct *usecase.CreateProductUsecase
}

func NewProductHandler(
	findProducts *usecase.FindProductsUsecase,
	getProduct *usecase.GetProductUsecase,
	createProduct *usecase.CreateProductUsecase,
) *ProductHandler {
	return &ProductHandler{
		findProducts:  findProducts,
		getProduct:    getProduct,
		createProduct: createProduct,
	}
}

func (h *ProductHandler) FindProducts(w http.ResponseWriter, r *http.Request) error {
	query := r.URL.Query()

	page, _ := strconv.Atoi(query.Get("page"))
	limit, _ := strconv.Atoi(query.Get("limit"))

	name := query.Get("name")
	id := query.Get("id")

	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	params := repository.FindProductParams{
		Page:  page,
		Limit: limit,
	}

	if name != "" {
		params.Name = &name
	}
	if id != "" {
		params.ID = &id
	}

	products, total, err := h.findProducts.Execute(params)
	if err != nil {
		return err
	}

	results := make([]ProductOverviewResponse, 0, len(products))
	for _, p := range products {
		result := ProductOverviewResponse{
			ID:            p.Product.ID,
			SKU:           p.Product.SKU,
			Name:          p.Product.Name,
			Slug:          p.Product.Slug,
			Status:        ProductStatusDTO(p.Product.Status),
			Price:         p.Product.Price,
			TotalStock:    p.Inventory.TotalStock,
			ReservedStock: p.Inventory.ReservedStock,
		}

		results = append(results, result)
	}

	response := map[string]interface{}{
		"products": results,
		"page":     page,
		"limit":    limit,
		"total":    total,
	}

	apphttp.WriteJSON(w, http.StatusOK, response)
	return nil
}

func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) error {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 || parts[2] == "" {
		return errors.ErrBadRequest
	}

	id := parts[2]
	if id == "" {
		return errors.ErrBadRequest
	}

	parsedID, err := uuid.Parse(id)
	if err != nil {
		return errors.ErrBadRequest
	}

	product, err := h.getProduct.Execute(parsedID)
	if err != nil {
		return err
	}
	if product == nil {
		return errors.ErrNotFound
	}

	inventories := make([]ProductInventoryView, 0, len(product.ShopInventories))
	for _, inventory := range product.ShopInventories {
		inventories = append(inventories, ProductInventoryView{
			ID:            inventory.ID,
			ShopID:        inventory.ShopID,
			TotalStock:    inventory.TotalStock,
			ReservedStock: inventory.ReservedStock,
			Available:     inventory.Available(),
		})
	}

	response := ProductDetailResponse{
		ID:            product.Product.ID,
		SKU:           product.Product.SKU,
		Name:          product.Product.Name,
		Slug:          product.Product.Slug,
		Description:   product.Product.Description,
		Status:        ProductStatusDTO(product.Product.Status),
		Price:         product.Product.Price,
		Weight:        product.Product.Weight,
		TotalStock:    product.Inventory.TotalStock,
		ReservedStock: product.Inventory.ReservedStock,
		Inventories:   inventories,
		CreatedAt:     product.Product.CreatedAt,
		UpdatedAt:     product.Product.UpdatedAt,
		ArchivedAt:    product.Product.ArchivedAt,
	}

	apphttp.WriteJSON(w, http.StatusOK, response)
	return nil
}

func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) error {
	var req CreateProductRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return errors.ErrBadRequest
	}

	if req.Name == "" || req.SKU == "" {
		return errors.ErrBadRequest
	}

	if req.Price < 0 {
		return errors.ErrBadRequest
	}

	if !req.Status.isStatusValid() {
		return errors.ErrBadRequest
	}

	err := h.createProduct.Execute(usecase.CreateProductInput{
		SKU:         req.SKU,
		Name:        req.Name,
		Description: req.Description,
		Status:      domain.ProductStatus(req.Status),
		Price:       req.Price,
		Weight:      req.Weight,
	})
	if err != nil {
		return err
	}

	response := map[string]string{
		"message": "product successfully created",
	}

	apphttp.WriteJSON(w, http.StatusOK, response)
	return nil
}

func (s ProductStatusDTO) isStatusValid() bool {
	switch s {
	case ProductStatusActive,
		ProductStatusInactive,
		ProductStatusArchived:
		return true
	default:
		return false
	}
}
