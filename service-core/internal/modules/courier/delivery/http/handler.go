package http

import (
	"encoding/json"
	"net/http"

	"service-core/internal/common/errors"
	apphttp "service-core/internal/common/http"

	"service-core/internal/modules/courier/usecase"

	"github.com/google/uuid"
)

type CourierHandler struct {
	configureCourierShop *usecase.ConfigureShopCourierUsecase
}

func NewCourierHandler(
	configureCourierShop *usecase.ConfigureShopCourierUsecase,
) *CourierHandler {
	return &CourierHandler{
		configureCourierShop: configureCourierShop,
	}
}

func (h *CourierHandler) ConfigureCourierShop(w http.ResponseWriter, r *http.Request) error {
	var req ConfigureCourierShopRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return errors.ErrBadRequest
	}

	shopID, err := uuid.Parse(req.ShopID)
	if err != nil {
		return errors.ErrBadRequest
	}

	inputs := make([]usecase.ConfigureShopCourierInput, 0, len(req.Couriers))

	for _, courier := range req.Couriers {
		inputs = append(inputs, usecase.ConfigureShopCourierInput{
			Code:   courier.Code,
			Active: courier.IsEnabled,
		})
	}

	err = h.configureCourierShop.Execute(shopID, inputs)
	if err != nil {
		return err
	}

	response := map[string]string{
		"message": "courier shops successfully configured",
	}

	apphttp.WriteJSON(w, http.StatusOK, response)
	return nil
}
