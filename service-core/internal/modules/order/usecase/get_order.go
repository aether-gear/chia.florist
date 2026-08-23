package usecase

import (
	"context"
	"fmt"

	cartDomain "service-core/internal/modules/cart/domain"
	"service-core/internal/modules/order/domain"
	"service-core/internal/modules/order/repository"
	paymentDomain "service-core/internal/modules/payment/domain"
	paymentRepo "service-core/internal/modules/payment/repository"
	shipmentDomain "service-core/internal/modules/shipment/domain"
	shipmentRepo "service-core/internal/modules/shipment/repository"
	transaction "service-core/internal/shared/transaction"

	"github.com/google/uuid"
)

type GetOrderUsecase struct {
	executor               transaction.Executor
	orderRepo              repository.OrderRepository
	orderItemRepo          repository.OrderItemRepository
	customDesignRepo       repository.OrderItemCustomDesignRepository
	paymentRepo            paymentRepo.PaymentRepository
	paymentChannelDataRepo paymentRepo.PaymentChannelDataRepository
	shipmentRepo           shipmentRepo.ShipmentRepository
	shipmentEventRepo      shipmentRepo.ShipmentEventRepository
}

func NewGetOrderUsecase(
	executor transaction.Executor,
	orderRepo repository.OrderRepository,
	orderItemRepo repository.OrderItemRepository,
	customDesignRepo repository.OrderItemCustomDesignRepository,
	paymentRepo paymentRepo.PaymentRepository,
	paymentChannelDataRepo paymentRepo.PaymentChannelDataRepository,
	shipmentRepo shipmentRepo.ShipmentRepository,
	shipmentEventRepo shipmentRepo.ShipmentEventRepository,
) *GetOrderUsecase {
	return &GetOrderUsecase{
		executor:               executor,
		orderRepo:              orderRepo,
		orderItemRepo:          orderItemRepo,
		customDesignRepo:       customDesignRepo,
		paymentRepo:            paymentRepo,
		paymentChannelDataRepo: paymentChannelDataRepo,
		shipmentRepo:           shipmentRepo,
		shipmentEventRepo:      shipmentEventRepo,
	}
}

type GetOrderInput struct {
	OrderID uuid.UUID

	// CustomerID, when set, enforces that the order
	// must belong to this customer.
	//
	// Use for customer-facing endpoints.
	// Leave nil for admin endpoints.
	CustomerID *uuid.UUID
}

type GetOrderResult struct {
	Order         domain.Order
	Items         []domain.OrderItem
	CustomDesigns map[uuid.UUID]domain.OrderItemCustomDesign
	Payment       *paymentDomain.Payment
	ChannelData   *paymentDomain.PaymentChannelData
	Shipment      *shipmentDomain.Shipment
	Shipments     []shipmentDomain.Shipment
}

func (u *GetOrderUsecase) Execute(
	ctx context.Context,
	input GetOrderInput,
) (*GetOrderResult, error) {
	order, err := u.orderRepo.GetByID(ctx, u.executor,
		input.OrderID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %w", err)
	}
	if order == nil {
		return nil, nil
	}

	// Enforce ownership for customer path.
	if input.CustomerID != nil &&
		order.CustomerID != *input.CustomerID {
		return nil, nil
	}

	items, err := u.orderItemRepo.ListByOrderID(ctx, u.executor,
		order.ID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to list order items: %w", err)
	}
	if items == nil {
		items = []domain.OrderItem{}
	}

	var customDesignsMap map[uuid.UUID]domain.OrderItemCustomDesign
	var customItemIDs []uuid.UUID
	for _, itm := range items {
		if itm.ProductVariantType == cartDomain.ProductVariantTypeCustom || itm.ProductID == nil {
			customItemIDs = append(customItemIDs, itm.ID)
		}
	}
	if len(customItemIDs) > 0 && u.customDesignRepo != nil {
		customDesignsMap, _ = u.customDesignRepo.ListByOrderItemIDs(ctx, u.executor, customItemIDs)
	}

	payment, err := u.paymentRepo.GetByOrderID(ctx, u.executor,
		order.ID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment: %w", err)
	}

	var channelData *paymentDomain.PaymentChannelData
	if payment != nil {
		cd, err := u.paymentChannelDataRepo.GetByPaymentID(ctx, u.executor,
			payment.ID,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to get payment channel data: %w", err)
		}
		channelData = cd
	}

	shipments, err := u.shipmentRepo.ListByOrderID(ctx, u.executor,
		order.ID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to list shipments: %w", err)
	}

	for i := range shipments {
		events, err := u.shipmentEventRepo.ListByShipmentID(ctx, u.executor,
			shipments[i].ID,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to list shipment events: %w", err)
		}
		if events == nil {
			events = []shipmentDomain.ShipmentEvent{}
		}

		shipments[i].Events = events

		// Associate item IDs
		var itemIDs []uuid.UUID
		for _, itm := range items {
			if itm.ShipmentID != nil && *itm.ShipmentID == shipments[i].ID {
				itemIDs = append(itemIDs, itm.ID)
			}
		}
		shipments[i].ItemIDs = itemIDs
	}

	var firstShipment *shipmentDomain.Shipment
	if len(shipments) > 0 {
		firstShipment = &shipments[0]
	}

	return &GetOrderResult{
		Order:         *order,
		Items:         items,
		CustomDesigns: customDesignsMap,
		Payment:       payment,
		ChannelData:   channelData,
		Shipment:      firstShipment,
		Shipments:     shipments,
	}, nil
}

