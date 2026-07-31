package usecase

import (
	"context"
	"fmt"
	"strconv"
	"time"

	apperrors "service-core/internal/common/errors"
	applogger "service-core/internal/common/logger"
	shipping "service-core/internal/infra/shipping"
	addressRepo "service-core/internal/modules/address/repository"
	"service-core/internal/modules/order/domain"
	"service-core/internal/modules/order/repository"
	shipmentDomain "service-core/internal/modules/shipment/domain"
	shipmentRepo "service-core/internal/modules/shipment/repository"
	transaction "service-core/internal/shared/transaction"

	"github.com/google/uuid"
)

// DEFAULT_SHIPMENT_WEIGHT_GRAMS is a placeholder weight (in grams) used when
// creating a Komerce shipment order. Replace with actual product weight once
// weight tracking is added to the inventory / product schema.
const DEFAULT_SHIPMENT_WEIGHT_GRAMS = 1000

// DEFAULT_SHIPMENT_ITEM_QTY is a placeholder item quantity used for the
// Komerce order creation call. Refactor once per-item weight is available.
const DEFAULT_SHIPMENT_ITEM_QTY = 1

type UpdateOrderStatusUsecase struct {
	executor        transaction.Executor
	transactor      transaction.Transactor
	orderRepo       repository.OrderRepository
	orderItemRepo   repository.OrderItemRepository
	shipmentRepo    shipmentRepo.ShipmentRepository
	addressRepo     addressRepo.CustomerAddressRepository
	shopAddressRepo addressRepo.ShopAddressRepository
	logistics       shipping.LogisticsProvider
	auditLogger     applogger.AuditLogger
}

func NewUpdateOrderStatusUsecase(
	executor transaction.Executor,
	transactor transaction.Transactor,
	orderRepo repository.OrderRepository,
	orderItemRepo repository.OrderItemRepository,
	shipmentRepo shipmentRepo.ShipmentRepository,
	addressRepo addressRepo.CustomerAddressRepository,
	shopAddressRepo addressRepo.ShopAddressRepository,
	logistics shipping.LogisticsProvider,
	auditLogger applogger.AuditLogger,
) *UpdateOrderStatusUsecase {
	return &UpdateOrderStatusUsecase{
		executor:        executor,
		transactor:      transactor,
		orderRepo:       orderRepo,
		orderItemRepo:   orderItemRepo,
		shipmentRepo:    shipmentRepo,
		addressRepo:     addressRepo,
		shopAddressRepo: shopAddressRepo,
		logistics:       logistics,
		auditLogger:     auditLogger,
	}
}

type UpdateOrderStatusInput struct {
	OrderID uuid.UUID
	Status  domain.OrderStatus

	// TrackingNumber is an optional override used when the server is running
	// in manual logistics mode. Automated providers (e.g. Komerce) ignore it.
	TrackingNumber *string

	// FulfillmentMethod is an optional override. If not provided, it defaults
	// to "courier".
	FulfillmentMethod *string
}

type UpdateOrderStatusResult struct {
	Order    domain.Order
	Shipment *shipmentDomain.Shipment
}

func (u *UpdateOrderStatusUsecase) Execute(
	ctx context.Context,
	input UpdateOrderStatusInput,
) (res *UpdateOrderStatusResult, err error) {
	var oldStatus string
	defer func() {
		if err != nil {
			u.auditLogger.Log(ctx, applogger.AuditEvent{
				Category:   "user_action",
				Action:     "update_order_status",
				Resource:   "order",
				ResourceID: input.OrderID.String(),
				Outcome:    applogger.OutcomeFailure,
				Metadata: map[string]any{
					"error":      err.Error(),
					"old_status": oldStatus,
					"new_status": string(input.Status),
				},
			})
		} else {
			u.auditLogger.Log(ctx, applogger.AuditEvent{
				Category:   "user_action",
				Action:     "update_order_status",
				Resource:   "order",
				ResourceID: input.OrderID.String(),
				Outcome:    applogger.OutcomeSuccess,
				Metadata: map[string]any{
					"old_status": oldStatus,
					"new_status": string(input.Status),
				},
			})
		}
	}()

	order, err := u.orderRepo.GetByID(ctx, u.executor,
		input.OrderID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %w", err)
	}
	if order == nil {
		return nil, apperrors.NewNotFound("order not found")
	}

	oldStatus = string(order.Status)

	if errStatus := order.UpdateStatus(input.Status); errStatus != nil {
		return nil, apperrors.NewInvalidInput(errStatus.Error())
	}

	// Simple transitions (no side-effects)
	if input.Status != domain.OrderStatusShipped {
		err = u.transactor.WithinTransaction(ctx, func(exec transaction.Executor) error {
			return u.orderRepo.UpdateStatus(ctx, exec,
				order.ID,
				input.Status,
			)
		})
		if err != nil {
			return nil, fmt.Errorf("failed to update order status: %w", err)
		}

		result := UpdateOrderStatusResult{Order: *order}
		return &result, nil
	}

	// processing → shipped
	items, err := u.orderItemRepo.ListByOrderID(ctx,
		u.executor,
		order.ID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to list order items: %w", err)
	}
	if len(items) == 0 {
		return nil, apperrors.NewInvalidInput("order has no items")
	}

	// All items within a single order share
	// the same courier selection (captured at checkout per shop).
	//
	// Use the first item as the reference.
	first := items[0]

	method := shipmentDomain.FulfillmentMethodCourier
	if input.FulfillmentMethod != nil &&
		*input.FulfillmentMethod != "" {
		method = shipmentDomain.FulfillmentMethod(*input.FulfillmentMethod)
	}

	var courierCode, courierService string
	if method == shipmentDomain.FulfillmentMethodCourier {
		if first.CourierCode != nil {
			courierCode = *first.CourierCode
		}

		if first.CourierService != nil {
			courierService = *first.CourierService
		}

		if courierCode == "" || courierService == "" {
			return nil, apperrors.NewInvalidInput("order items have no courier information")
		}
	} else {
		courierCode = string(shipmentDomain.FulfillmentMethodSelfDelivery)
		courierService = string(shipmentDomain.FulfillmentMethodSelfDelivery)
	}

	// Resolve customer destination district ID
	// from the order's address.
	customerAddr, err := u.addressRepo.GetByID(ctx, u.executor,
		order.AddressID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get customer address: %w", err)
	}
	if customerAddr == nil {
		return nil, apperrors.NewNotFound("customer address not found")
	}

	destAreaID, err := strconv.Atoi(customerAddr.Detail.DistrictID)
	if err != nil {
		return nil, fmt.Errorf("invalid destination district ID %q: %w", customerAddr.Detail.DistrictID, err)
	}

	// Resolve shop origin district ID
	// from the shop's default address.
	shopAddr, err := u.shopAddressRepo.GetDefaultByShopID(ctx, u.executor,
		first.ShopID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get shop address: %w", err)
	}
	if shopAddr == nil {
		return nil, apperrors.NewNotFound("shop address not found")
	}

	originAreaID, err := strconv.Atoi(shopAddr.Detail.DistrictID)
	if err != nil {
		return nil, fmt.Errorf("invalid origin district ID %q: %w", shopAddr.Detail.DistrictID, err)
	}

	var trackingNumber *string
	if method == shipmentDomain.FulfillmentMethodCourier {
		// Build a human-readable item name
		// for the Komerce order.
		itemName := first.ProductName
		if len(items) > 1 {
			itemName = fmt.Sprintf("%s (+%d more)", first.ProductName, len(items)-1)
		}

		orderInput := shipping.CreateOrderInput{
			OriginAreaID:         originAreaID,
			DestinationAreaID:    destAreaID,
			CourierCode:          courierCode,
			CourierService:       courierService,
			Weight:               DEFAULT_SHIPMENT_WEIGHT_GRAMS,
			UniqueOrderID:        order.Number,
			ItemName:             itemName,
			ItemPrice:            order.Subtotal,
			ItemQty:              DEFAULT_SHIPMENT_ITEM_QTY,
			ShipperName:          shopAddr.Label,
			ShipperPhone:         derefPhone(shopAddr.Phone),
			ShipperAddress:       shopAddr.Detail.FullAddress,
			ReceiverName:         customerAddr.ReceiverName,
			ReceiverPhone:        derefPhone(customerAddr.Phone),
			ReceiverAddress:      customerAddr.Detail.FullAddress,
			ManualTrackingNumber: input.TrackingNumber,
		}

		// Cost is pulled from the order's ShippingFee
		// stored at checkout (OQ1).
		//
		// Weight uses the placeholder constant above
		// until product weight is tracked.
		komerceResult, err := u.logistics.CreateOrder(ctx, orderInput)
		if err != nil {
			return nil, fmt.Errorf("failed to create Komerce shipment order: %w", err)
		}

		tracking := komerceResult.TrackingNumber
		trackingNumber = &tracking
	}

	now := time.Now()
	shipment := shipmentDomain.Shipment{
		ID:                uuid.New(),
		OrderID:           order.ID,
		Status:            shipmentDomain.ShipmentStatusCreated,
		FulfillmentMethod: method,
		TrackingNumber:    trackingNumber,
		Courier:           courierCode,
		Service:           courierService,
		Cost:              order.ShippingFee,
		Weight:            DEFAULT_SHIPMENT_WEIGHT_GRAMS,
		OriginID:          shopAddr.Detail.DistrictID,
		DestinationID:     customerAddr.Detail.DistrictID,
		CreatedAt:         now,
	}

	if err := shipment.Validate(); err != nil {
		return nil, apperrors.NewInvalidInput(err.Error())
	}

	var created *shipmentDomain.Shipment
	err = u.transactor.WithinTransaction(ctx, func(exec transaction.Executor) error {
		if err := u.shipmentRepo.Create(ctx, exec,
			shipment,
		); err != nil {
			return fmt.Errorf("failed to persist shipment: %w", err)
		}

		if err := u.orderRepo.UpdateStatus(ctx, exec,
			order.ID,
			input.Status,
		); err != nil {
			return fmt.Errorf("failed to update order status: %w", err)
		}

		created = &shipment
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &UpdateOrderStatusResult{
		Order:    *order,
		Shipment: created,
	}, nil
}

// derefPhone safely dereferences a nullable phone pointer, returning an
// empty string when nil. Komerce accepts empty phone numbers gracefully.
func derefPhone(phone *string) string {
	if phone == nil {
		return ""
	}
	return *phone
}
