package http

import (
	"encoding/json"
	"net/http"

	"service-core/internal/modules/shipment/usecase"

	"service-core/internal/common/errors"
	apphttp "service-core/internal/common/http"
)

type ShipmentHandler struct {
	estimateShippingCost *usecase.EstimateShippingCostUsecase
}

func NewShipmentHandler(
	estimateShippingCost *usecase.EstimateShippingCostUsecase,
) *ShipmentHandler {
	return &ShipmentHandler{
		estimateShippingCost: estimateShippingCost,
	}
}

func (h *ShipmentHandler) EstimateShippingCost(w http.ResponseWriter, r *http.Request) error {
	var req EstimateShippingCostRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return errors.ErrBadRequest
	}

	if req.Weight <= 0 {
		return errors.ErrBadRequest
	}

	if req.Origin <= 0 || req.Destination <= 0 {
		return errors.ErrBadRequest
	}

	input := usecase.EstimateShippingCostInput{
		Origin:      req.Origin,
		Destination: req.Destination,
		Weight:      req.Weight,
		Couriers:    req.Couriers,
		PriceFilter: req.PriceFilter,
	}

	results, err := h.estimateShippingCost.Execute(input)
	if err != nil {
		return err
	}

	couriers := make([]EstimateShippingCostResponse, 0, len(results))

	for _, result := range results {
		option := EstimateShippingCostResponse{
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
