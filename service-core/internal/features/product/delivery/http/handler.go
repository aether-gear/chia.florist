package http

import (
	"encoding/json"
	"net/http"
	"service-core/internal/features/product/repository"
	application "service-core/internal/features/product/usecase"
	"strconv"
	"strings"
)

type ProductHandler struct {
	findUsecase *application.FindProductsUsecase
	getUsecase  *application.GetProductUsecase
}

func NewProductHandler(
	find *application.FindProductsUsecase,
	get *application.GetProductUsecase,
) *ProductHandler {
	return &ProductHandler{
		findUsecase: find,
		getUsecase:  get,
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

	product, err := h.getUsecase.Execute(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	result := ToDetailResponse(product)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
