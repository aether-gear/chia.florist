package http

import (
	"encoding/json"
	"net/http"

	"service-core/internal/modules/shipment/usecase"

	"service-core/internal/common/errors"
	apphttp "service-core/internal/common/http"
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
	var req EstimateShippingOptionsRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return errors.ErrBadRequest
	}

	if req.Weight <= 0 {
		return errors.ErrBadRequest
	}

	if req.Origin <= 0 || req.Destination <= 0 {
		return errors.ErrBadRequest
	}

	input := usecase.EstimateShippingOptionsInput{
		Origin:      req.Origin,
		Destination: req.Destination,
		Weight:      req.Weight,
		Couriers:    req.Couriers,
		PriceFilter: req.PriceFilter,
	}

	results, err := h.estimateShippOpts.Execute(input)
	if err != nil {
		return err
	}

	couriers := make([]EstimateShippingOptionsResponse, 0, len(results))

	for _, result := range results {
		option := EstimateShippingOptionsResponse{
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
