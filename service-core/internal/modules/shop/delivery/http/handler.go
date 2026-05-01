package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"service-core/internal/modules/shop/usecase"

	"service-core/internal/common/errors"
	apphttp "service-core/internal/common/http"
)

type ShopHandler struct {
	createShop *usecase.CreateShopUsecase
}

func NewAddressHandler(
	createShop *usecase.CreateShopUsecase,
) *ShopHandler {
	return &ShopHandler{
		createShop: createShop,
	}
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
