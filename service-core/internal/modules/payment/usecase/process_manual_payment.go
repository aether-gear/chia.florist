package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	apperrors "service-core/internal/common/errors"
	inventoryRepo "service-core/internal/modules/inventory/repository"
	orderDomain "service-core/internal/modules/order/domain"
	orderRepo "service-core/internal/modules/order/repository"
	"service-core/internal/modules/payment/domain"
	paymentDomain "service-core/internal/modules/payment/domain"
	paymentRepo "service-core/internal/modules/payment/repository"
	transaction "service-core/internal/shared/transaction"

	"github.com/google/uuid"
)

type ProcessManualPaymentUsecase struct {
	paymentRepo      paymentRepo.PaymentRepository
	paymentAccRepo   paymentRepo.PaymentAccountRepository
	paymentEventRepo paymentRepo.PaymentEventRepository
	orderRepo        orderRepo.OrderRepository
	orderItemRepo    orderRepo.OrderItemRepository
	inventoryRepo    inventoryRepo.InventoryRepository
	transactor       transaction.Transactor
	executor         transaction.Executor
}

func NewProcessManualPaymentUsecase(
	paymentRepo paymentRepo.PaymentRepository,
	paymentAccRepo paymentRepo.PaymentAccountRepository,
	paymentEventRepo paymentRepo.PaymentEventRepository,
	orderRepo orderRepo.OrderRepository,
	orderItemRepo orderRepo.OrderItemRepository,
	inventoryRepo inventoryRepo.InventoryRepository,
	transactor transaction.Transactor,
	executor transaction.Executor,
) *ProcessManualPaymentUsecase {
	return &ProcessManualPaymentUsecase{
		paymentRepo:      paymentRepo,
		paymentAccRepo:   paymentAccRepo,
		paymentEventRepo: paymentEventRepo,
		orderRepo:        orderRepo,
		orderItemRepo:    orderItemRepo,
		inventoryRepo:    inventoryRepo,
		transactor:       transactor,
		executor:         executor,
	}
}

type ProcessManualPaymentInput struct {
	PaymentID uuid.UUID
	Action    string
}

func (u *ProcessManualPaymentUsecase) Execute(
	ctx context.Context,
	input ProcessManualPaymentInput,
) error {
	err := u.transactor.WithinTransaction(ctx, func(exec transaction.Executor) error {
		payment, err := u.paymentRepo.
			GetByID(ctx, exec, input.PaymentID)
		if err != nil {
			return fmt.Errorf("failed to retrieve payment: %w", err)
		}

		if payment == nil {
			return apperrors.NewNotFound("payment not found")
		}

		if payment.Provider !=
			string(domain.PaymentProviderManual) {

			return apperrors.NewBadRequest("only manual payments can be processed manually")
		}

		if payment.Status !=
			paymentDomain.PaymentStatusPending {

			return apperrors.NewConflict("payment is not pending")
		}

		// Derive payment and order status transitions
		// from manual action input
		//
		// This ensures a strict coupling between
		// payment status and order status: they must always
		// transition together within the same transactional boundary
		var (
			newPaymentStatus paymentDomain.PaymentStatus
			newOrderStatus   orderDomain.OrderStatus
			action           string
		)

		switch input.Action {
		case "confirm":
			newPaymentStatus = paymentDomain.PaymentStatusPaid
			newOrderStatus = orderDomain.OrderStatusConfirmed
			action = "commit"

		case "reject", "cancel":
			newPaymentStatus = paymentDomain.PaymentStatusFailed
			newOrderStatus = orderDomain.OrderStatusCancelled
			action = "release"

		default:
			return apperrors.NewBadRequest(fmt.Sprintf("invalid manual action: %s", input.Action))
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
				if err := u.inventoryRepo.Commit(
					ctx,
					exec,
					item.ProductID,
					item.ShopID,
					item.Quantity,
				); err != nil {
					return fmt.Errorf("failed to commit inventory for product %s: %w", item.ProductID, err)
				}

			case "release":
				if err := u.inventoryRepo.Release(
					ctx,
					exec,
					item.ProductID,
					item.ShopID,
					item.Quantity,
				); err != nil {
					return fmt.Errorf("failed to release inventory for product %s: %w", item.ProductID, err)
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
			"status":       string(newPaymentStatus),
			"confirmed_by": "manual_admin",
			"action":       input.Action,
			"gross_amount": payment.Amount,
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
