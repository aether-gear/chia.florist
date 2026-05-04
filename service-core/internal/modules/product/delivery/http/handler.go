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
	findUsecase   *usecase.FindProductsUsecase
	getUsecase    *usecase.GetProductUsecase
	createUsecase *usecase.CreateProductUsecase
}

func NewProductHandler(
	find *usecase.FindProductsUsecase,
	get *usecase.GetProductUsecase,
	create *usecase.CreateProductUsecase,
) *ProductHandler {
	return &ProductHandler{
		findUsecase:   find,
		getUsecase:    get,
		createUsecase: create,
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

	products, total, err := h.findUsecase.Execute(params)
	if err != nil {
		return err
	}

	results := make([]ProductOverviewResponse, 0, len(products))
	for _, p := range products {
		results = append(results, ToListResponse(p))
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

	product, err := h.getUsecase.Execute(parsedID)
	if err != nil {
		return err
	}

	if product == nil {
		return errors.ErrNotFound
	}

	response := ToDetailResponse(*product)

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

	if req.Price < 0 || req.InitialStock < 0 {
		return errors.ErrBadRequest
	}

	if !req.Status.isStatusValid() {
		return errors.ErrBadRequest
	}

	err := h.createUsecase.Execute(usecase.CreateProductInput{
		SKU:          req.SKU,
		Name:         req.Name,
		Description:  req.Description,
		Status:       domain.ProductStatus(req.Status),
		Price:        req.Price,
		Weight:       req.Weight,
		InitialStock: req.InitialStock,
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
