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
	"service-core/internal/modules/order/domain"
	"service-core/internal/modules/order/repository"
	productDomain "service-core/internal/modules/product/domain"
	productRepo "service-core/internal/modules/product/repository"
	shipmentDomain "service-core/internal/modules/shipment/domain"
	shipmentRepo "service-core/internal/modules/shipment/repository"
	transaction "service-core/internal/shared/transaction"

	"github.com/google/uuid"
)

type DispatchShopShipmentUsecase struct {
	executor        transaction.Executor
	transactor      transaction.Transactor
	orderRepo       repository.OrderRepository
	orderItemRepo   repository.OrderItemRepository
	productRepo     productRepo.ProductRepository
	shipmentRepo    shipmentRepo.ShipmentRepository
	addressRepo     addressRepo.CustomerAddressRepository
	shopAddressRepo addressRepo.ShopAddressRepository
	logistics       shipping.LogisticsProvider
	auditLogger     applogger.AuditLogger
}

func NewDispatchShopShipmentUsecase(
	executor transaction.Executor,
	transactor transaction.Transactor,
	orderRepo repository.OrderRepository,
	orderItemRepo repository.OrderItemRepository,
	productRepo productRepo.ProductRepository,
	shipmentRepo shipmentRepo.ShipmentRepository,
	addressRepo addressRepo.CustomerAddressRepository,
	shopAddressRepo addressRepo.ShopAddressRepository,
	logistics shipping.LogisticsProvider,
	auditLogger applogger.AuditLogger,
) *DispatchShopShipmentUsecase {
	return &DispatchShopShipmentUsecase{
		executor:        executor,
		transactor:      transactor,
		orderRepo:       orderRepo,
		orderItemRepo:   orderItemRepo,
		productRepo:     productRepo,
		shipmentRepo:    shipmentRepo,
		addressRepo:     addressRepo,
		shopAddressRepo: shopAddressRepo,
		logistics:       logistics,
		auditLogger:     auditLogger,
	}
}

type DispatchShopShipmentInput struct {
	OrderID           uuid.UUID
	ShopID            uuid.UUID
	FulfillmentMethod string
	Courier           string
	Service           string
	TrackingNumber    *string
	ItemIDs           []uuid.UUID
}

type DispatchShopShipmentResult struct {
	Order           domain.Order
	Shipment        shipmentDomain.Shipment
	AllItemsShipped bool
}

func (u *DispatchShopShipmentUsecase) Execute(
	ctx context.Context,
	input DispatchShopShipmentInput,
) (res *DispatchShopShipmentResult, err error) {
	defer func() {
		if err != nil {
			u.auditLogger.Log(ctx, applogger.AuditEvent{
				Category:   "user_action",
				Action:     "dispatch_shop_shipment",
				Resource:   "order",
				ResourceID: input.OrderID.String(),
				Outcome:    applogger.OutcomeFailure,
				Metadata: map[string]any{
					"error":   err.Error(),
					"shop_id": input.ShopID.String(),
				},
			})
		} else {
			u.auditLogger.Log(ctx, applogger.AuditEvent{
				Category:   "user_action",
				Action:     "dispatch_shop_shipment",
				Resource:   "order",
				ResourceID: input.OrderID.String(),
				Outcome:    applogger.OutcomeSuccess,
				Metadata: map[string]any{
					"shop_id":     input.ShopID.String(),
					"shipment_id": res.Shipment.ID.String(),
					"all_shipped": res.AllItemsShipped,
				},
			})
		}
	}()

	if len(input.ItemIDs) == 0 {
		return nil, apperrors.NewInvalidInput("each shipment must contain at least one order item")
	}

	order, err := u.orderRepo.GetByID(ctx, u.executor, input.OrderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %w", err)
	}
	if order == nil {
		return nil, apperrors.NewNotFound("order not found")
	}

	if order.Status != domain.OrderStatusProcessing && order.Status != domain.OrderStatusConfirmed {
		return nil, apperrors.NewInvalidInput(fmt.Sprintf("cannot dispatch shipment for order in %s status", order.Status))
	}

	// Fetch all order items
	allItems, err := u.orderItemRepo.ListByOrderID(ctx, u.executor, order.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to list order items: %w", err)
	}
	if len(allItems) == 0 {
		return nil, apperrors.NewInvalidInput("order has no items")
	}

	itemMap := make(map[uuid.UUID]domain.OrderItem, len(allItems))
	for _, item := range allItems {
		itemMap[item.ID] = item
	}

	var shipmentItems []domain.OrderItem
	for _, itemID := range input.ItemIDs {
		item, ok := itemMap[itemID]
		if !ok {
			return nil, apperrors.NewInvalidInput(fmt.Sprintf("order item %s does not belong to this order", itemID))
		}
		if item.ShopID != input.ShopID {
			return nil, apperrors.NewForbidden(fmt.Sprintf("order item %s does not belong to shop %s", itemID, input.ShopID))
		}
		if item.ShipmentID != nil {
			return nil, apperrors.NewInvalidInput(fmt.Sprintf("order item %s is already assigned to a shipment", itemID))
		}
		shipmentItems = append(shipmentItems, item)
	}

	first := shipmentItems[0]

	method := shipmentDomain.FulfillmentMethodCourier
	if input.FulfillmentMethod != "" {
		method = shipmentDomain.FulfillmentMethod(input.FulfillmentMethod)
	}

	var courierCode, courierService string
	if method == shipmentDomain.FulfillmentMethodCourier {
		courierCode = input.Courier
		courierService = input.Service
		if courierCode == "" && first.CourierCode != nil {
			courierCode = *first.CourierCode
		}
		if courierService == "" && first.CourierService != nil {
			courierService = *first.CourierService
		}
		if courierCode == "" || courierService == "" {
			return nil, apperrors.NewInvalidInput("shipment is missing courier information")
		}
	} else {
		courierCode = string(shipmentDomain.FulfillmentMethodSelfDelivery)
		courierService = string(shipmentDomain.FulfillmentMethodSelfDelivery)
	}

	customerAddr, err := u.addressRepo.GetByID(ctx, u.executor, order.AddressID)
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

	shopAddr, err := u.shopAddressRepo.GetDefaultByShopID(ctx, u.executor, input.ShopID)
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

	// Fetch product details for accurate weight calculation
	var productIDs []uuid.UUID
	for _, item := range shipmentItems {
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

	var totalWeightGrams int
	var totalItemQty int
	var subtotal int64
	for _, item := range shipmentItems {
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
		totalWeightGrams += weight * qty
		totalItemQty += qty
		subtotal += item.Subtotal
	}

	var createdLogisticsOrderNo string
	rollbackLogistics := func() {
		if createdLogisticsOrderNo != "" {
			_ = u.logistics.CancelOrder(ctx, createdLogisticsOrderNo)
		}
	}

	trackingNumber := input.TrackingNumber
	if method == shipmentDomain.FulfillmentMethodCourier && (trackingNumber == nil || *trackingNumber == "") {
		itemName := first.ProductName
		if len(shipmentItems) > 1 {
			itemName = fmt.Sprintf("%s (+%d more)", first.ProductName, len(shipmentItems)-1)
		}

		uniqueOrderID := fmt.Sprintf("%s-SH-%s", order.Number, input.ShopID.String()[:8])

		orderInput := shipping.CreateOrderInput{
			OriginAreaID:         originAreaID,
			DestinationAreaID:    destAreaID,
			CourierCode:          courierCode,
			CourierService:       courierService,
			Weight:               totalWeightGrams,
			UniqueOrderID:        uniqueOrderID,
			ItemName:             itemName,
			ItemPrice:            subtotal,
			ItemQty:              totalItemQty,
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
			return nil, fmt.Errorf("failed to create shipment order: %w", err)
		}
		if komerceResult != nil && komerceResult.KomerceOrderNo != "" {
			createdLogisticsOrderNo = komerceResult.KomerceOrderNo
		}

		if komerceResult != nil && komerceResult.TrackingNumber != "" {
			tracking := komerceResult.TrackingNumber
			trackingNumber = &tracking
		}
	}

	shipmentCost := first.ShippingFee

	now := appclock.Now()
	shipment := shipmentDomain.Shipment{
		ID:                uuid.New(),
		OrderID:           order.ID,
		Status:            shipmentDomain.ShipmentStatusCreated,
		FulfillmentMethod: method,
		TrackingNumber:    trackingNumber,
		Courier:           courierCode,
		Service:           courierService,
		Cost:              shipmentCost,
		Weight:            totalWeightGrams,
		OriginID:          shopAddr.Detail.DistrictID,
		DestinationID:     customerAddr.Detail.DistrictID,
		CreatedAt:         now,
		ItemIDs:           input.ItemIDs,
	}

	if err := shipment.Validate(); err != nil {
		rollbackLogistics()
		return nil, apperrors.NewInvalidInput(err.Error())
	}

	// Check whether all items across all shops in this order will now be shipped
	unshippedCount := 0
	for _, item := range allItems {
		if item.ShipmentID == nil {
			unshippedCount++
		}
	}
	allItemsShipped := unshippedCount <= len(input.ItemIDs)

	err = u.transactor.WithinTransaction(ctx, func(exec transaction.Executor) error {
		if err := u.shipmentRepo.Create(ctx, exec, shipment); err != nil {
			return fmt.Errorf("failed to persist shipment: %w", err)
		}
		if err := u.orderItemRepo.AssignShipment(ctx, exec, shipment.ID, input.ItemIDs); err != nil {
			return fmt.Errorf("failed to link items to shipment: %w", err)
		}
		if allItemsShipped {
			if err := u.orderRepo.UpdateStatus(ctx, exec, order.ID, domain.OrderStatusShipped); err != nil {
				return fmt.Errorf("failed to update order status to shipped: %w", err)
			}
			order.Status = domain.OrderStatusShipped
		}
		return nil
	})
	if err != nil {
		rollbackLogistics()
		return nil, err
	}

	return &DispatchShopShipmentResult{
		Order:           *order,
		Shipment:        shipment,
		AllItemsShipped: allItemsShipped,
	}, nil
}
