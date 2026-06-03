package http

import (
	"net/http"

	apperrors "service-core/internal/common/errors"
	apphttp "service-core/internal/common/http"
	"service-core/internal/modules/courier/usecase"

	"github.com/google/uuid"
)

type CourierHandler struct {
	configureShopCourier *usecase.ConfigureShopCourierUsecase
}

func NewCourierHandler(
	configureShopCourier *usecase.ConfigureShopCourierUsecase,
) *CourierHandler {
	return &CourierHandler{
		configureShopCourier: configureShopCourier,
	}
}

func (h *CourierHandler) ConfigureCourierShop(w http.ResponseWriter, r *http.Request) error {
	var req configureCourierShopRequest

	if err := apphttp.DecodeJSON(r, &req); err != nil {
		return apperrors.NewBadRequest("invalid body request")
	}

	shopID, err := uuid.Parse(req.ShopID)
	if err != nil {
		return apperrors.NewBadRequest("invalid shop id")
	}

	inputs := make([]usecase.ConfigureShopCourierInput, 0, len(req.Couriers))

	for _, courier := range req.Couriers {
		inputs = append(inputs, usecase.ConfigureShopCourierInput{
			Code:   courier.Code,
			Active: courier.IsEnabled,
		})
	}

	err = h.configureShopCourier.Execute(r.Context(), shopID, inputs)
	if err != nil {
		return err
	}

	response := map[string]string{
		"message": "courier shops successfully configured",
	}

	apphttp.WriteJSON(w, http.StatusOK, response)
	return nil
}
