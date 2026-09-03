package usecase

import (
	"context"
	"fmt"
	"log"
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
	trackingCache     *shipping.TrackingCache
}

func NewGetOrderTrackingUsecase(
	executor transaction.Executor,
	orderRepo repository.OrderRepository,
	shipmentRepo shipmentRepo.ShipmentRepository,
	shipmentEventRepo shipmentRepo.ShipmentEventRepository,
	logisticsProvider shipping.LogisticsProvider,
	addressRepo addressRepo.CustomerAddressRepository,
	trackingCache *shipping.TrackingCache,
) *GetOrderTrackingUsecase {
	if trackingCache == nil {
		trackingCache = shipping.NewTrackingCache(shipping.DefaultTrackingCacheTTL)
	}
	return &GetOrderTrackingUsecase{
		executor:          executor,
		orderRepo:         orderRepo,
		shipmentRepo:      shipmentRepo,
		shipmentEventRepo: shipmentEventRepo,
		logisticsProvider: logisticsProvider,
		addressRepo:       addressRepo,
		trackingCache:     trackingCache,
	}
}

type GetOrderTrackingInput struct {
	OrderID    uuid.UUID
	CustomerID uuid.UUID
	ShipmentID *uuid.UUID
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
	Warning        *string                 `json:"warning,omitempty"`
	Timeline       []TrackingTimelineEvent `json:"timeline"`
}

func (u *GetOrderTrackingUsecase) Execute(
	ctx context.Context,
	input GetOrderTrackingInput,
) (*GetOrderTrackingResult, error) {
	order, err := u.orderRepo.GetByID(ctx, u.executor,
		input.OrderID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %w", err)
	}
	if order == nil {
		return nil, apperrors.NewNotFound("order not found")
	}

	if input.CustomerID != uuid.Nil && order.CustomerID != input.CustomerID {
		return nil, apperrors.NewUnauthorized("not authorized")
	}

	var shipment *shipmentDomain.Shipment
	if input.ShipmentID != nil {
		s, err := u.shipmentRepo.GetByID(ctx, u.executor, *input.ShipmentID)
		if err != nil {
			return nil, fmt.Errorf("failed to get shipment: %w", err)
		}
		if s != nil && s.OrderID == order.ID {
			shipment = s
		}
	}

	if shipment == nil {
		shipments, err := u.shipmentRepo.ListByOrderID(ctx, u.executor, order.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to list shipments: %w", err)
		}
		if len(shipments) > 0 {
			shipment = &shipments[0]
		}
	}

	if shipment == nil {
		return nil, apperrors.NewNotFound("shipment not found")
	}

	internalEvents, err := u.shipmentEventRepo.ListByShipmentID(ctx, u.executor,
		shipment.ID,
	)
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

	var warningMsg *string

	if shipment.FulfillmentMethod == shipmentDomain.FulfillmentMethodCourier &&
		shipment.TrackingNumber != nil && *shipment.TrackingNumber != "" &&
		u.logisticsProvider != nil {

		var externalEvents []shipping.TrackingEvent

		// Check in-memory TTL cache to shield external provider from rate limits (429)
		if cachedEvents, hit := u.trackingCache.Get(shipment.Courier, *shipment.TrackingNumber); hit {
			externalEvents = cachedEvents
		} else {
			var lastPhone *string
			customerAddr, err := u.addressRepo.GetByID(ctx, u.executor,
				order.AddressID,
			)
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

			fetchedEvents, err := u.logisticsProvider.TrackShipment(ctx, trackInput)
			if err != nil {
				msg := fmt.Sprintf("Live courier tracking is unavailable for %s (AWB: %s). Manual status updates are required.", shipment.Courier, *shipment.TrackingNumber)
				warningMsg = &msg
				log.Printf("[GetOrderTracking Warning] %s (detail: %v)", msg, err)
				// Fail-safe: if external provider returns rate limit (429) or error, try returning stale cached events if available
				if staleEvents, staleHit := u.trackingCache.GetStale(shipment.Courier, *shipment.TrackingNumber); staleHit {
					externalEvents = staleEvents
				}
			} else {
				externalEvents = fetchedEvents
				u.trackingCache.Set(shipment.Courier, *shipment.TrackingNumber, fetchedEvents)
			}
		}

		for _, e := range externalEvents {
			timeline = append(timeline, TrackingTimelineEvent{
				Status:      e.Status,
				Description: e.Description,
				Location:    e.Location,
				Timestamp:   e.Timestamp,
			})
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
		Warning:        warningMsg,
		Timeline:       timeline,
	}

	return &result, nil
}
