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
	getShop    *usecase.GetShopUsecase
	createShop *usecase.CreateShopUsecase
}

func NewAddressHandler(
	getShop *usecase.GetShopUsecase,
	createShop *usecase.CreateShopUsecase,
) *ShopHandler {
	return &ShopHandler{
		getShop:    getShop,
		createShop: createShop,
	}
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
