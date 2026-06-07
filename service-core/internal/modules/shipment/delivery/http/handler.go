package http

import (
	"net/http"

	apperrors "service-core/internal/common/errors"
	apphttp "service-core/internal/common/http"
	"service-core/internal/modules/shipment/usecase"
)

type ShipmentHandler struct {
	estimateShippOpts *usecase.EstimateShippingOptionsUsecase
}

func NewShipmentHandler(
	estimateShippOpts *usecase.EstimateShippingOptionsUsecase,
) *ShipmentHandler {
	return &ShipmentHandler{
		estimateShippOpts: estimateShippOpts,
	}
}

func (h *ShipmentHandler) EstimateShippingOptions(w http.ResponseWriter, r *http.Request) error {
	var req estimateShippingOptionsRequest

	if err := apphttp.DecodeJSON(r, &req); err != nil {
		return apperrors.NewBadRequest("invalid body request")
	}

	if req.Weight <= 0 {
		return apperrors.NewBadRequest("invalid weight")
	}
	if req.Origin <= 0 {
		return apperrors.NewBadRequest("invalid origin")
	}
	if req.Destination <= 0 {
		return apperrors.NewBadRequest("invalid destination")
	}

	input := usecase.EstimateShippingOptionsInput{
		Origin:      req.Origin,
		Destination: req.Destination,
		Weight:      req.Weight,
		Couriers:    req.Couriers,
		PriceFilter: req.PriceFilter,
	}

	results, err := h.estimateShippOpts.Execute(r.Context(), input)
	if err != nil {
		return err
	}

	couriers := make([]estimateShippingOptionsResponse, 0, len(results))

	for _, result := range results {
		option := estimateShippingOptionsResponse{
			Name:        result.Name,
			Code:        result.Code,
			Service:     result.Service,
			Description: result.Description,
			Cost:        result.Cost,
			Etd:         result.Etd,
		}

		couriers = append(couriers, option)
	}

	response := map[string]any{
		"couriers": couriers,
	}

	apphttp.WriteJSON(w, http.StatusOK, response)
	return nil
}
