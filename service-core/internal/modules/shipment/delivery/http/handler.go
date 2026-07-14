package http

import (
	"net/http"

	apperrors "service-core/internal/common/errors"
	apphttp "service-core/internal/common/http"
	"service-core/internal/modules/shipment/domain"
	"service-core/internal/modules/shipment/usecase"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type ShipmentHandler struct {
	estimateShippOpts    *usecase.EstimateShippingOptionsUsecase
	updateShipmentStatus *usecase.UpdateShipmentStatusUsecase
	updateShipment       *usecase.UpdateShipmentUsecase
}

func NewShipmentHandler(
	estimateShippOpts *usecase.EstimateShippingOptionsUsecase,
	updateShipmentStatus *usecase.UpdateShipmentStatusUsecase,
	updateShipment *usecase.UpdateShipmentUsecase,
) *ShipmentHandler {
	return &ShipmentHandler{
		estimateShippOpts:    estimateShippOpts,
		updateShipmentStatus: updateShipmentStatus,
		updateShipment:       updateShipment,
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

func (h *ShipmentHandler) UpdateShipmentStatus(w http.ResponseWriter, r *http.Request) error {
	shipmentIDStr := chi.URLParam(r, "shipmentID")
	shipmentID, err := uuid.Parse(shipmentIDStr)
	if err != nil {
		return apperrors.NewBadRequest("invalid shipment id")
	}

	var req updateShipmentStatusRequest
	if err := apphttp.DecodeJSON(r, &req); err != nil {
		return apperrors.NewBadRequest("invalid request body")
	}

	if req.Status == "" {
		return apperrors.NewBadRequest("status is required")
	}

	input := usecase.UpdateShipmentStatusInput{
		ShipmentID:  shipmentID,
		Status:      domain.ShipmentStatus(req.Status),
		Description: req.Description,
		Location:    req.Location,
	}

	res, err := h.updateShipmentStatus.Execute(r.Context(), input)
	if err != nil {
		return err
	}

	resp := buildShipmentResponse(res.Shipment)

	apphttp.WriteJSON(w, http.StatusOK, resp)
	return nil
}

func (h *ShipmentHandler) UpdateShipment(w http.ResponseWriter, r *http.Request) error {
	shipmentIDStr := chi.URLParam(r, "shipmentID")
	shipmentID, err := uuid.Parse(shipmentIDStr)
	if err != nil {
		return apperrors.NewBadRequest("invalid shipment id")
	}

	var req updateShipmentRequest
	if err := apphttp.DecodeJSON(r, &req); err != nil {
		return apperrors.NewBadRequest("invalid request body")
	}

	input := usecase.UpdateShipmentInput{
		ShipmentID:     shipmentID,
		TrackingNumber: req.TrackingNumber,
		Courier:        req.Courier,
		Service:        req.Service,
	}

	res, err := h.updateShipment.Execute(r.Context(), input)
	if err != nil {
		return err
	}

	resp := buildShipmentResponse(res.Shipment)
	apphttp.WriteJSON(w, http.StatusOK, resp)

	return nil
}

func buildShipmentResponse(s *domain.Shipment) shipmentResponse {
	return shipmentResponse{
		ID:                s.ID.String(),
		OrderID:           s.OrderID.String(),
		Status:            string(s.Status),
		FulfillmentMethod: string(s.FulfillmentMethod),
		TrackingNumber:    s.TrackingNumber,
		Courier:           s.Courier,
		Service:           s.Service,
		Cost:              s.Cost,
		Weight:            s.Weight,
		CreatedAt:         s.CreatedAt,
	}
}
