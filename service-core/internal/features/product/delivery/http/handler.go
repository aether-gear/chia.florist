package http

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"service-core/internal/features/product/domain"
	"service-core/internal/features/product/repository"
	"service-core/internal/features/product/usecase"
	"strconv"
	"strings"

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

func (h *ProductHandler) FindProducts(w http.ResponseWriter, r *http.Request) {
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	responses := make([]ProductOverviewResponse, 0, len(products))
	for _, p := range products {
		responses = append(responses, ToListResponse(p))
	}

	response := map[string]interface{}{
		"products": responses,
		"page":     page,
		"limit":    limit,
		"total":    total,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 || parts[2] == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	id := parts[2]
	if id == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	parsedID, err := uuid.Parse(id)
	if err != nil {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}

	product, err := h.getUsecase.Execute(parsedID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if product == nil {
		http.Error(w, "not found", http.StatusNotFound)
	}

	result := ToDetailResponse(*product)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var req CreateProductRequest

	body, _ := io.ReadAll(r.Body)

	r.Body = io.NopCloser(bytes.NewBuffer(body))

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		http.Error(w, "missing name", http.StatusBadRequest)
		return
	}
	if req.SKU == "" {
		http.Error(w, "missing sku", http.StatusBadRequest)
		return
	}
	if req.Price < 0 {
		http.Error(w, "price cannot be negative", http.StatusBadRequest)
		return
	}
	if req.InitialStock < 0 {
		http.Error(w, "stock cannot be negative", http.StatusBadRequest)
		return
	}
	if !req.Status.isStatusValid() {
		http.Error(w, "invalid status", http.StatusBadRequest)
		return
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "product successfully created",
	})
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
