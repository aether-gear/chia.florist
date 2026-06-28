package usecase

import (
	"context"
	"fmt"

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
	executor      transaction.Executor
	orderRepo     repository.OrderRepository
	orderItemRepo repository.OrderItemRepository
	paymentRepo   paymentRepo.PaymentRepository
	shipmentRepo  shipmentRepo.ShipmentRepository
}

func NewGetOrderUsecase(
	executor transaction.Executor,
	orderRepo repository.OrderRepository,
	orderItemRepo repository.OrderItemRepository,
	paymentRepo paymentRepo.PaymentRepository,
	shipmentRepo shipmentRepo.ShipmentRepository,
) *GetOrderUsecase {
	return &GetOrderUsecase{
		executor:      executor,
		orderRepo:     orderRepo,
		orderItemRepo: orderItemRepo,
		paymentRepo:   paymentRepo,
		shipmentRepo:  shipmentRepo,
	}
}

type GetOrderInput struct {
	OrderID uuid.UUID

	// CustomerID, when set, enforces that the order must belong to this customer.
	// Use for customer-facing endpoints. Leave nil for admin endpoints.
	CustomerID *uuid.UUID
}

type GetOrderResult struct {
	Order    domain.Order
	Items    []domain.OrderItem
	Payment  *paymentDomain.Payment
	Shipment *shipmentDomain.Shipment
}

func (u *GetOrderUsecase) Execute(
	ctx context.Context,
	input GetOrderInput,
) (*GetOrderResult, error) {
	order, err := u.orderRepo.GetByID(ctx, u.executor, input.OrderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %w", err)
	}
	if order == nil {
		return nil, nil
	}

	// Enforce ownership for customer path.
	if input.CustomerID != nil && order.CustomerID != *input.CustomerID {
		return nil, nil
	}

	items, err := u.orderItemRepo.ListByOrderID(ctx, u.executor, order.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to list order items: %w", err)
	}
	if items == nil {
		items = []domain.OrderItem{}
	}

	payment, err := u.paymentRepo.GetByOrderID(ctx, u.executor, order.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment: %w", err)
	}

	shipment, err := u.shipmentRepo.GetByOrderID(ctx, u.executor, order.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get shipment: %w", err)
	}

	return &GetOrderResult{
		Order:    *order,
		Items:    items,
		Payment:  payment,
		Shipment: shipment,
	}, nil
}
