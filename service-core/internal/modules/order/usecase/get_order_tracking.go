package usecase

import (
	"context"
	"fmt"
	"sort"
	"time"

	apperrors "service-core/internal/common/errors"
	shipping "service-core/internal/infra/shipping"
	addressRepo "service-core/internal/modules/address/repository"
	"service-core/internal/modules/order/repository"
	shipmentDomain "service-core/internal/modules/shipment/domain"
	shipmentRepo "service-core/internal/modules/shipment/repository"
	transaction "service-core/internal/shared/transaction"

	"github.com/google/uuid"
)

type GetOrderTrackingUsecase struct {
	executor          transaction.Executor
	orderRepo         repository.OrderRepository
	shipmentRepo      shipmentRepo.ShipmentRepository
	shipmentEventRepo shipmentRepo.ShipmentEventRepository
	logisticsProvider shipping.LogisticsProvider
	addressRepo       addressRepo.CustomerAddressRepository
}

func NewGetOrderTrackingUsecase(
	executor transaction.Executor,
	orderRepo repository.OrderRepository,
	shipmentRepo shipmentRepo.ShipmentRepository,
	shipmentEventRepo shipmentRepo.ShipmentEventRepository,
	logisticsProvider shipping.LogisticsProvider,
	addressRepo addressRepo.CustomerAddressRepository,
) *GetOrderTrackingUsecase {
	return &GetOrderTrackingUsecase{
		executor:          executor,
		orderRepo:         orderRepo,
		shipmentRepo:      shipmentRepo,
		shipmentEventRepo: shipmentEventRepo,
		logisticsProvider: logisticsProvider,
		addressRepo:       addressRepo,
	}
}

type GetOrderTrackingInput struct {
	OrderID    uuid.UUID
	CustomerID uuid.UUID
}

type TrackingTimelineEvent struct {
	Status      string    `json:"status"`
	Description string    `json:"description"`
	Location    string    `json:"location"`
	Timestamp   time.Time `json:"timestamp"`
}

type GetOrderTrackingResult struct {
	OrderID        uuid.UUID               `json:"order_id"`
	ShipmentID     uuid.UUID               `json:"shipment_id"`
	Courier        string                  `json:"courier"`
	TrackingNumber *string                 `json:"tracking_number,omitempty"`
	Timeline       []TrackingTimelineEvent `json:"timeline"`
}

func (u *GetOrderTrackingUsecase) Execute(
	ctx context.Context,
	input GetOrderTrackingInput,
) (*GetOrderTrackingResult, error) {
	order, err := u.orderRepo.GetByID(ctx, u.executor, input.OrderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %w", err)
	}
	if order == nil {
		return nil, apperrors.NewNotFound("order not found")
	}

	if order.CustomerID != input.CustomerID {
		return nil, apperrors.NewUnauthorized("not authorized")
	}

	shipment, err := u.shipmentRepo.GetByOrderID(ctx, u.executor, order.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get shipment: %w", err)
	}
	if shipment == nil {
		return nil, apperrors.NewNotFound("shipment not found")
	}

	internalEvents, err := u.shipmentEventRepo.
		ListByShipmentID(ctx, u.executor, shipment.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to list internal shipment events: %w", err)
	}

	var timeline []TrackingTimelineEvent
	for _, e := range internalEvents {
		timeline = append(timeline, TrackingTimelineEvent{
			Status:      e.Status,
			Description: e.Description,
			Location:    e.Location,
			Timestamp:   e.Timestamp,
		})
	}

	if shipment.FulfillmentMethod == shipmentDomain.FulfillmentMethodCourier &&
		shipment.TrackingNumber != nil && *shipment.TrackingNumber != "" &&
		u.logisticsProvider != nil {

		var lastPhone *string
		customerAddr, err := u.addressRepo.
			GetByID(ctx, u.executor, order.AddressID)
		if err == nil &&
			customerAddr != nil &&
			customerAddr.Phone != nil {

			lastPhone = customerAddr.Phone
		}

		trackInput := shipping.TrackShipmentInput{
			Courier:        shipment.Courier,
			TrackingNumber: *shipment.TrackingNumber,
			LastPhone:      lastPhone,
		}

		externalEvents, err := u.logisticsProvider.TrackShipment(ctx, trackInput)
		if err != nil {
			// Fail-safe: log warning/ignore external tracking errors
			// so customer can still see internal timeline
			//
			// System will not fail the request if external provider is down
		} else {
			for _, e := range externalEvents {
				timeline = append(timeline, TrackingTimelineEvent{
					Status:      e.Status,
					Description: e.Description,
					Location:    e.Location,
					Timestamp:   e.Timestamp,
				})
			}
		}
	}

	sort.Slice(timeline, func(i, j int) bool {
		return timeline[i].Timestamp.Before(timeline[j].Timestamp)
	})

	result := GetOrderTrackingResult{
		OrderID:        order.ID,
		ShipmentID:     shipment.ID,
		Courier:        shipment.Courier,
		TrackingNumber: shipment.TrackingNumber,
		Timeline:       timeline,
	}

	return &result, nil
}
