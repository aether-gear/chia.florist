package usecase

import (
	"context"
	"fmt"
	"strconv"

	appclock "service-core/internal/common/clock"
	apperrors "service-core/internal/common/errors"
	applogger "service-core/internal/common/logger"
	shipping "service-core/internal/infra/shipping"
	addressRepo "service-core/internal/modules/address/repository"
	inventoryRepo "service-core/internal/modules/inventory/repository"
	"service-core/internal/modules/order/domain"
	"service-core/internal/modules/order/repository"
	paymentDomain "service-core/internal/modules/payment/domain"
	paymentRepo "service-core/internal/modules/payment/repository"
	productDomain "service-core/internal/modules/product/domain"
	productRepo "service-core/internal/modules/product/repository"
	shipmentDomain "service-core/internal/modules/shipment/domain"
	shipmentRepo "service-core/internal/modules/shipment/repository"
	transaction "service-core/internal/shared/transaction"

	"github.com/google/uuid"
)

// DEFAULT_SHIPMENT_WEIGHT_GRAMS is a fallback weight (in grams) used when
// product weight is not specified on the product schema.
const DEFAULT_SHIPMENT_WEIGHT_GRAMS = 1000

// DEFAULT_SHIPMENT_ITEM_QTY is a placeholder item quantity used for the
// Komerce order creation call.
const DEFAULT_SHIPMENT_ITEM_QTY = 1

type UpdateOrderStatusUsecase struct {
	executor        transaction.Executor
	transactor      transaction.Transactor
	orderRepo       repository.OrderRepository
	orderItemRepo   repository.OrderItemRepository
	inventoryRepo   inventoryRepo.InventoryRepository
	paymentRepo     paymentRepo.PaymentRepository
	productRepo     productRepo.ProductRepository
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
	inventoryRepo inventoryRepo.InventoryRepository,
	paymentRepo paymentRepo.PaymentRepository,
	productRepo productRepo.ProductRepository,
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
		inventoryRepo:   inventoryRepo,
		paymentRepo:     paymentRepo,
		productRepo:     productRepo,
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

	// When confirming an unconfirmed order,
	// invoke the domain Confirm method to stamp
	//
	// ConfirmedAt and calculate the 3-day handling SLA expiration
	// (HandlingExpiresAt).
	if input.Status == domain.OrderStatusConfirmed && order.ConfirmedAt == nil {
		if errStatus := order.Confirm(appclock.Now(), domain.DefaultHandlingSLAWindow); errStatus != nil {
			return nil, apperrors.NewInvalidInput(errStatus.Error())
		}
	} else {
		if errStatus := order.UpdateStatus(input.Status); errStatus != nil {
			return nil, apperrors.NewInvalidInput(errStatus.Error())
		}
	}

	switch input.Status {
	case domain.OrderStatusConfirmed:
		payment, err := u.paymentRepo.GetByOrderID(ctx, u.executor, order.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get payment: %w", err)
		}
		if payment == nil || payment.Status != paymentDomain.PaymentStatusPaid {
			return nil, apperrors.NewInvalidInput("cannot confirm order without confirmed payment")
		}

		items, err := u.orderItemRepo.ListByOrderID(ctx, u.executor, order.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to list order items: %w", err)
		}

		err = u.transactor.WithinTransaction(ctx, func(exec transaction.Executor) error {
			for _, item := range items {
				if item.ProductID == nil {
					continue
				}
				if err := u.inventoryRepo.Commit(ctx, exec, *item.ProductID, item.ShopID, item.Quantity); err != nil {
					return fmt.Errorf("failed to commit inventory for product %s: %w", item.ProductID, err)
				}
			}
			if err := u.orderRepo.UpdateStatusWithSLA(ctx, exec, order.ID, input.Status, order.ConfirmedAt, order.HandlingExpiresAt); err != nil {
				return fmt.Errorf("failed to update order status: %w", err)
			}
			return nil
		})
		if err != nil {
			return nil, err
		}

		result := UpdateOrderStatusResult{Order: *order}
		return &result, nil

	case domain.OrderStatusProcessing:
		payment, err := u.paymentRepo.GetByOrderID(ctx, u.executor, order.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get payment: %w", err)
		}
		if payment == nil || payment.Status != paymentDomain.PaymentStatusPaid {
			return nil, apperrors.NewInvalidInput("cannot move order to processing without confirmed payment")
		}

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

	case domain.OrderStatusDelivered:
		shipment, err := u.shipmentRepo.GetByOrderID(ctx, u.executor, order.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get shipment: %w", err)
		}

		err = u.transactor.WithinTransaction(ctx, func(exec transaction.Executor) error {
			if shipment != nil {
				if err := shipment.UpdateStatus(shipmentDomain.ShipmentStatusDelivered); err != nil {
					return fmt.Errorf("failed to update shipment status: %w", err)
				}
				if err := u.shipmentRepo.Update(ctx, exec, *shipment); err != nil {
					return fmt.Errorf("failed to persist shipment: %w", err)
				}
			}
			return u.orderRepo.UpdateStatus(ctx, exec, order.ID, input.Status)
		})
		if err != nil {
			return nil, err
		}

		result := UpdateOrderStatusResult{Order: *order, Shipment: shipment}
		return &result, nil

	case domain.OrderStatusCancelled:
		items, err := u.orderItemRepo.ListByOrderID(ctx, u.executor, order.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to list order items: %w", err)
		}

		err = u.transactor.WithinTransaction(ctx, func(exec transaction.Executor) error {
			for _, item := range items {
				if item.ProductID == nil {
					continue
				}
				if err := u.inventoryRepo.Release(ctx, exec, *item.ProductID, item.ShopID, item.Quantity); err != nil {
					return fmt.Errorf("failed to release inventory for product %s: %w", item.ProductID, err)
				}
			}
			return u.orderRepo.UpdateStatus(ctx, exec, order.ID, input.Status)
		})
		if err != nil {
			return nil, err
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

	// Resolve customer destination district ID from the order's address
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

	// Fetch product details for all items to calculate accurate shipping weights
	var productIDs []uuid.UUID
	for _, item := range items {
		if item.ProductID != nil {
			productIDs = append(productIDs, *item.ProductID)
		}
	}
	products, err := u.productRepo.FindByIDs(ctx, u.executor, productIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to load products: %w", err)
	}
	productMap := make(map[uuid.UUID]productDomain.Product, len(products))
	for _, p := range products {
		productMap[p.ID] = p
	}

	// Group order items by ShopID for per-shop shipment processing
	type shopGroup struct {
		shopID uuid.UUID
		items  []domain.OrderItem
	}
	var groups []shopGroup
	groupMap := make(map[uuid.UUID]int)
	for _, item := range items {
		idx, exists := groupMap[item.ShopID]
		if !exists {
			groupMap[item.ShopID] = len(groups)
			groups = append(groups, shopGroup{
				shopID: item.ShopID,
				items:  []domain.OrderItem{item},
			})
		} else {
			groups[idx].items = append(groups[idx].items, item)
		}
	}

	method := shipmentDomain.FulfillmentMethodCourier
	if input.FulfillmentMethod != nil && *input.FulfillmentMethod != "" {
		method = shipmentDomain.FulfillmentMethod(*input.FulfillmentMethod)
	}

	now := appclock.Now()
	var shipmentsToCreate []shipmentDomain.Shipment

	for idx, group := range groups {
		first := group.items[0]

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

		shopAddr, err := u.shopAddressRepo.GetDefaultByShopID(ctx, u.executor,
			group.shopID,
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

		// Calculate total weight and item quantity for this shop's shipment
		var shopTotalWeightGrams int
		var shopTotalItemQty int
		var shopSubtotal int64
		for _, item := range group.items {
			qty := item.Quantity
			if qty <= 0 {
				qty = 1
			}
			weight := DEFAULT_SHIPMENT_WEIGHT_GRAMS
			if item.ProductID != nil {
				if p, ok := productMap[*item.ProductID]; ok && p.Weight != nil && *p.Weight > 0 {
					weight = int(*p.Weight)
				}
			}
			shopTotalWeightGrams += weight * qty
			shopTotalItemQty += qty
			shopSubtotal += item.Subtotal
		}

		var trackingNumber *string
		if method == shipmentDomain.FulfillmentMethodCourier {
			itemName := first.ProductName
			if len(group.items) > 1 {
				itemName = fmt.Sprintf("%s (+%d more)", first.ProductName, len(group.items)-1)
			}

			uniqueOrderID := order.Number
			if len(groups) > 1 {
				uniqueOrderID = fmt.Sprintf("%s-%d", order.Number, idx+1)
			}

			orderInput := shipping.CreateOrderInput{
				OriginAreaID:         originAreaID,
				DestinationAreaID:    destAreaID,
				CourierCode:          courierCode,
				CourierService:       courierService,
				Weight:               shopTotalWeightGrams,
				UniqueOrderID:        uniqueOrderID,
				ItemName:             itemName,
				ItemPrice:            shopSubtotal,
				ItemQty:              shopTotalItemQty,
				ShipperName:          shopAddr.Label,
				ShipperPhone:         derefPhone(shopAddr.Phone),
				ShipperAddress:       shopAddr.Detail.FullAddress,
				ReceiverName:         customerAddr.ReceiverName,
				ReceiverPhone:        derefPhone(customerAddr.Phone),
				ReceiverAddress:      customerAddr.Detail.FullAddress,
				ManualTrackingNumber: input.TrackingNumber,
			}

			komerceResult, err := u.logistics.CreateOrder(ctx, orderInput)
			if err != nil {
				return nil, fmt.Errorf("failed to create Komerce shipment order: %w", err)
			}

			tracking := komerceResult.TrackingNumber
			trackingNumber = &tracking
		}

		shipmentCost := first.ShippingFee
		if shipmentCost == 0 && len(groups) == 1 {
			shipmentCost = order.ShippingFee
		}

		shipment := shipmentDomain.Shipment{
			ID:                uuid.New(),
			OrderID:           order.ID,
			Status:            shipmentDomain.ShipmentStatusCreated,
			FulfillmentMethod: method,
			TrackingNumber:    trackingNumber,
			Courier:           courierCode,
			Service:           courierService,
			Cost:              shipmentCost,
			Weight:            shopTotalWeightGrams,
			OriginID:          shopAddr.Detail.DistrictID,
			DestinationID:     customerAddr.Detail.DistrictID,
			CreatedAt:         now,
		}

		if err := shipment.Validate(); err != nil {
			return nil, apperrors.NewInvalidInput(err.Error())
		}

		shipmentsToCreate = append(shipmentsToCreate, shipment)
	}

	var firstCreated *shipmentDomain.Shipment
	err = u.transactor.WithinTransaction(ctx, func(exec transaction.Executor) error {
		for i, shipment := range shipmentsToCreate {
			if err := u.shipmentRepo.Create(ctx, exec, shipment); err != nil {
				return fmt.Errorf("failed to persist shipment: %w", err)
			}
			if i == 0 {
				firstCreated = &shipmentsToCreate[0]
			}
		}

		if err := u.orderRepo.UpdateStatus(ctx, exec,
			order.ID,
			input.Status,
		); err != nil {
			return fmt.Errorf("failed to update order status: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &UpdateOrderStatusResult{
		Order:    *order,
		Shipment: firstCreated,
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
