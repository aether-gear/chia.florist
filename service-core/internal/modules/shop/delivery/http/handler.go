package http

import (
	"net/http"
	"strconv"

	apperrors "service-core/internal/common/errors"
	apphttp "service-core/internal/common/http"
	authzSvc "service-core/internal/modules/authorization/infra/service"
	"service-core/internal/modules/shop/usecase"
)

type ShopHandler struct {
	findShops  *usecase.FindShopsUsecase
	getShop    *usecase.GetShopUsecase
	createShop *usecase.CreateShopUsecase
}

func NewAddressHandler(
	findShops *usecase.FindShopsUsecase,
	getShop *usecase.GetShopUsecase,
	createShop *usecase.CreateShopUsecase,
) *ShopHandler {
	return &ShopHandler{
		findShops:  findShops,
		getShop:    getShop,
		createShop: createShop,
	}
}

func (h *ShopHandler) FindShops(w http.ResponseWriter, r *http.Request) error {
	page := apphttp.QueryIntDefault(r, "page", 1)
	if page <= 0 {
		page = 1
	}
	limit := apphttp.QueryIntDefault(r, "limit", 10)
	if limit <= 0 {
		limit = 10
	}

	name := apphttp.Query(r, "name")
	id := apphttp.Query(r, "id")
	sort := apphttp.Query(r, "sort")

	input := usecase.FindShopsInput{
		Page:  page,
		Limit: limit,
		Sort:  sort,
	}
	if name != "" {
		input.Name = &name
	}
	if id != "" {
		input.ID = &id
	}

	shops, total, err := h.findShops.Execute(r.Context(), input)
	if err != nil {
		return err
	}

	var shopsResponse []getShopResponse
	for _, shop := range shops {
		s := getShopResponse{
			ID:          shop.ID,
			Name:        shop.Name,
			Slug:        shop.Slug,
			Description: shop.Description,
			IsActive:    shop.IsActive,
			CreatedAt:   shop.CreatedAt,
			UpdatedAt:   shop.UpdatedAt,
		}

		shopsResponse = append(shopsResponse, s)
	}

	response := map[string]any{
		"shops": shopsResponse,
		"page":  page,
		"limit": limit,
		"total": total,
	}

	apphttp.WriteJSON(w, http.StatusOK, response)
	return nil
}

func (h *ShopHandler) GetShopByID(w http.ResponseWriter, r *http.Request) error {
	id, err := apphttp.ParamUUID(r, "id")
	if err != nil {
		return apperrors.NewBadRequest("invalid shop id")
	}

	result, err := h.getShop.GetByID(r.Context(), id)
	if err != nil {
		return err
	}

	response := map[string]getShopResponse{
		"shop": {
			ID:          result.ID,
			Name:        result.Name,
			Slug:        result.Slug,
			Description: result.Description,
			IsActive:    result.IsActive,
			CreatedAt:   result.CreatedAt,
			UpdatedAt:   result.UpdatedAt,
		},
	}

	apphttp.WriteJSON(w, http.StatusOK, response)
	return nil
}

func (h *ShopHandler) CreateShop(w http.ResponseWriter, r *http.Request) error {
	actor, ok := authzSvc.GetActor(r.Context())
	if !ok {
		return apperrors.NewUnauthorized("authentication required")
	}

	var req createShopRequest
	if err := apphttp.DecodeJSON(r, &req); err != nil {
		return apperrors.NewBadRequest("invalid request body")
	}

	if req.Name == "" {
		return apperrors.NewBadRequest("invalid name")
	}

	var parsedIsActive bool
	parsedIsActive, err := strconv.ParseBool(req.IsActive)
	if err != nil {
		return apperrors.NewBadRequest("invalid active status")
	}

	input := usecase.CreateShopInput{
		Name:        req.Name,
		Description: req.Description,
		IsActive:    parsedIsActive,
	}

	err = h.createShop.Execute(
		r.Context(),
		*actor,
		input,
	)
	if err != nil {
		return err
	}

	response := map[string]string{
		"message": "shop successfully created",
	}

	apphttp.WriteJSON(w, http.StatusOK, response)
	return nil
}
