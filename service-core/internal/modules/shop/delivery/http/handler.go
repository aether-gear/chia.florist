package http

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"service-core/internal/modules/shop/usecase"

	"service-core/internal/common/errors"
	apphttp "service-core/internal/common/http"

	"github.com/google/uuid"
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
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 || parts[2] == "" {
		return errors.ErrBadRequest
	}

	id := parts[2]
	if id == "" {
		return errors.ErrBadRequest
	}
	parsedShopID, err := uuid.Parse(id)
	if err != nil {
		return errors.ErrBadRequest
	}

	result, err := h.getShop.GetByID(parsedShopID)
	if err != nil {
		return err
	}

	response := map[string]GetShopResponse{
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
	var req CreateShopRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return errors.ErrBadRequest
	}

	if req.Name == "" {
		return errors.ErrBadRequest
	}

	var parsedBool bool
	parsedBool, err := strconv.ParseBool(req.IsActive)
	if err != nil {
		return errors.ErrBadRequest
	}

	shopInput := usecase.CreateShopInput{
		Name:        req.Name,
		Description: req.Description,
		IsActive:    parsedBool,
	}

	err = h.createShop.Execute(shopInput)
	if err != nil {
		return err
	}

	response := map[string]string{
		"message": "shop successfully created",
	}

	apphttp.WriteJSON(w, http.StatusOK, response)
	return nil
}
