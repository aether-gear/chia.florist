package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	apperrors "service-core/internal/common/errors"
	paymentgateway "service-core/internal/infra/payment-gateway"
	inventoryRepo "service-core/internal/modules/inventory/repository"
	orderDomain "service-core/internal/modules/order/domain"
	orderRepo "service-core/internal/modules/order/repository"
	paymentDomain "service-core/internal/modules/payment/domain"
	paymentRepo "service-core/internal/modules/payment/repository"
	transaction "service-core/internal/shared/transaction"

	"github.com/google/uuid"
)

type ProcessPaymentWebhookUsecase struct {
	paymentRepo      paymentRepo.PaymentRepository
	paymentAccRepo   paymentRepo.PaymentAccountRepository
	paymentEventRepo paymentRepo.PaymentEventRepository
	orderRepo        orderRepo.OrderRepository
	orderItemRepo    orderRepo.OrderItemRepository
	inventoryRepo    inventoryRepo.InventoryRepository
	paymentGateway   paymentgateway.Provider
	transactor       transaction.Transactor
	executor         transaction.Executor
}

func NewProcessPaymentWebhookUsecase(
	paymentRepo paymentRepo.PaymentRepository,
	paymentAccRepo paymentRepo.PaymentAccountRepository,
	paymentEventRepo paymentRepo.PaymentEventRepository,
	orderRepo orderRepo.OrderRepository,
	orderItemRepo orderRepo.OrderItemRepository,
	inventoryRepo inventoryRepo.InventoryRepository,
	paymentGateway paymentgateway.Provider,
	transactor transaction.Transactor,
	executor transaction.Executor,
) *ProcessPaymentWebhookUsecase {
	return &ProcessPaymentWebhookUsecase{
		paymentRepo:      paymentRepo,
		paymentAccRepo:   paymentAccRepo,
		paymentEventRepo: paymentEventRepo,
		orderRepo:        orderRepo,
		orderItemRepo:    orderItemRepo,
		inventoryRepo:    inventoryRepo,
		paymentGateway:   paymentGateway,
		transactor:       transactor,
		executor:         executor,
	}
}

type ProcessPaymentWebhookInput struct {
	Payload map[string]any
}

func (u *ProcessPaymentWebhookUsecase) Execute(
	ctx context.Context,
	input ProcessPaymentWebhookInput,
) error {
	notifResult, err := u.paymentGateway.
		ParseNotification(ctx, input.Payload)
	if err != nil {
		return fmt.Errorf("failed to parse gateway notification: %w", err)
	}

	orderID, err := uuid.Parse(notifResult.GatewayOrderID)
	if err != nil {
		return apperrors.NewBadRequest(fmt.Sprintf("invalid order ID in gateway response: %s",
			notifResult.GatewayOrderID))
	}

	err = u.transactor.WithinTransaction(ctx, func(exec transaction.Executor) error {
		payment, err := u.paymentRepo.
			GetByOrderID(ctx, exec, orderID)
		if err != nil {
			return fmt.Errorf("failed to retrieve payment: %w", err)
		}

		if payment == nil {
			return apperrors.NewNotFound("payment not found for order")
		}

		if payment.Status != paymentDomain.PaymentStatusPending {
			return nil
		}

		// Derive payment and order status transitions
		// from webhook input
		//
		// This ensures a strict coupling between
		// payment status and order status: they must always
		// transition together within the same transactional boundary
		var (
			newPaymentStatus paymentDomain.PaymentStatus
			newOrderStatus   orderDomain.OrderStatus
			action           string
		)

		switch notifResult.Status {
		case paymentgateway.NotificationStatusSettlement:
			newPaymentStatus = paymentDomain.PaymentStatusPaid
			newOrderStatus = orderDomain.OrderStatusConfirmed
			action = "commit"

		case paymentgateway.NotificationStatusExpire:
			newPaymentStatus = paymentDomain.PaymentStatusExpired
			newOrderStatus = orderDomain.OrderStatusCancelled
			action = "release"

		case paymentgateway.NotificationStatusCancel:
			newPaymentStatus = paymentDomain.PaymentStatusCancelled
			newOrderStatus = orderDomain.OrderStatusCancelled
			action = "release"

		case paymentgateway.NotificationStatusDeny:
			newPaymentStatus = paymentDomain.PaymentStatusFailed
			newOrderStatus = orderDomain.OrderStatusCancelled
			action = "release"

		default:
			return nil
		}

		if err := u.paymentRepo.UpdateStatus(
			ctx,
			exec,
			payment.ID,
			newPaymentStatus,
		); err != nil {
			return fmt.Errorf("failed to update payment status: %w", err)
		}

		if err := u.orderRepo.UpdateStatus(
			ctx,
			exec,
			payment.OrderID,
			newOrderStatus,
		); err != nil {
			return fmt.Errorf("failed to update order status: %w", err)
		}

		if action == "commit" ||
			action == "release" {

			orderItems, err := u.orderItemRepo.
				ListByOrderID(ctx, exec, payment.OrderID)
			if err != nil {
				return fmt.Errorf("failed to list order items: %w", err)
			}

			// Apply inventory state transition based
			// on payment outcome
			//
			// This step ensures stock consistency after
			// payment resolution:
			//   - commit   → finalize reserved stock
			// 				  (reduce available inventory)
			//   - release  → rollback reserved stock
			// 				  back to available pool
			for _, item := range orderItems {
				switch action {
				case "commit":
					if err := u.inventoryRepo.
						Commit(
							ctx,
							exec,
							item.ProductID,
							item.ShopID,
							item.Quantity,
						); err != nil {
						return fmt.Errorf("failed to commit inventory for product %s: %w",
							item.ProductID, err)
					}

				case "release":
					if err := u.inventoryRepo.
						Release(
							ctx,
							exec,
							item.ProductID,
							item.ShopID,
							item.Quantity,
						); err != nil {
						return fmt.Errorf("failed to release inventory for product %s: %w",
							item.ProductID, err)
					}
				}
			}
		}

		// Adjust payment account load tracking
		// after successful resolution
		//
		// A successful or failed payment resolution
		// reduces the active load previously reserved
		// during checkout
		if payment.PaymentAccountID != nil {
			if err := u.paymentAccRepo.
				DecrementLoad(ctx, exec, *payment.PaymentAccountID); err != nil {

				return fmt.Errorf("failed to decrement payment account load: %w", err)
			}
		}

		// Emit payment event for audit
		// and downstream processing
		//
		// This event acts as a durable audit log
		// of the final payment resolution state
		payloadBytes, err := json.Marshal(map[string]any{
			"status":        string(newPaymentStatus),
			"raw_status":    notifResult.RawStatus,
			"fraud_status":  notifResult.FraudStatus,
			"gross_amount":  notifResult.GrossAmount,
			"gateway_tx_id": notifResult.GatewayTransactionID,
		})
		if err != nil {
			return fmt.Errorf("failed to marshal payment event payload: %w", err)
		}

		paymentEvent := paymentDomain.PaymentEvent{
			ID:        uuid.New(),
			PaymentID: payment.ID,
			EventName: string(newPaymentStatus),
			Payload:   payloadBytes,
			CreatedAt: time.Now(),
		}

		if err := u.paymentEventRepo.
			Create(ctx, exec, paymentEvent); err != nil {
			return fmt.Errorf("failed to create payment event: %w", err)
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
