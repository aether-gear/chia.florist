package usecase

import (
	"context"
	"fmt"

	apperrors "service-core/internal/common/errors"
	paymentgateway "service-core/internal/infra/payment-gateway"
	orderRepo "service-core/internal/modules/order/repository"
	"service-core/internal/modules/payment/domain"
	"service-core/internal/modules/payment/repository"
	transaction "service-core/internal/shared/transaction"

	"github.com/google/uuid"
)

// CheckPaymentStatusUsecase is the customer-
// triggered payment sync.
//
// When a customer's order appears stuck as 'pending'
// after they have paid, they can call this usecase to immediately
// query Midtrans for the current status and resolve the payment —
// without waiting for the next background reconciliation tick.
type CheckPaymentStatusUsecase struct {
	orderRepo      orderRepo.OrderRepository
	repository     repository.PaymentRepository
	paymentGateway paymentgateway.Provider
	processWebhook *ProcessPaymentWebhookUsecase
	executor       transaction.Executor
}

func NewCheckPaymentStatusUsecase(
	orderRepo orderRepo.OrderRepository,
	repository repository.PaymentRepository,
	paymentGateway paymentgateway.Provider,
	processWebhook *ProcessPaymentWebhookUsecase,
	executor transaction.Executor,
) *CheckPaymentStatusUsecase {
	return &CheckPaymentStatusUsecase{
		orderRepo:      orderRepo,
		repository:     repository,
		paymentGateway: paymentGateway,
		processWebhook: processWebhook,
		executor:       executor,
	}
}

type CheckPaymentStatusInput struct {
	OrderID    uuid.UUID
	CustomerID uuid.UUID
}

type CheckPaymentStatusResult struct {
	// Status is the payment status after the check
	// (may be unchanged if Midtrans still reports pending).
	Status domain.PaymentStatus

	// Synced is true when the status was resolved
	// from Midtrans during this call
	// (i.e. it was pending before and is now terminal).
	Synced bool
}

func (u *CheckPaymentStatusUsecase) Execute(
	ctx context.Context,
	input CheckPaymentStatusInput,
) (*CheckPaymentStatusResult, error) {
	order, err := u.orderRepo.
		GetByID(ctx, u.executor, input.OrderID)
	if err != nil {
		return nil, fmt.Errorf("check payment status: retrieve order: %w", err)
	}

	if order == nil {
		return nil, apperrors.NewNotFound("order not found")
	}

	if order.CustomerID != input.CustomerID {
		return nil, apperrors.NewNotFound("order not found")
	}

	payment, err := u.repository.
		GetByOrderID(ctx, u.executor, input.OrderID)
	if err != nil {
		return nil, fmt.Errorf("check payment status: retrieve payment: %w", err)
	}

	if payment == nil {
		return nil, apperrors.NewNotFound("payment not found")
	}

	// Only gateway payments with a ProviderOrderID can be synced.
	if payment.Provider != "gateway" ||
		payment.ProviderOrderID == nil {
		return &CheckPaymentStatusResult{Status: payment.Status, Synced: false}, nil
	}

	// Only pending payments need a sync — return early for any terminal state.
	if payment.Status != domain.PaymentStatusPending {
		return &CheckPaymentStatusResult{Status: payment.Status, Synced: false}, nil
	}

	result, err := u.paymentGateway.
		GetTransactionStatus(ctx, *payment.ProviderOrderID)
	if err != nil {
		return nil, fmt.Errorf("check payment status: gateway status check failed: %w", err)
	}

	if result.Status == paymentgateway.NotificationStatusPending {
		return &CheckPaymentStatusResult{
			Status: payment.Status,
			Synced: false,
		}, nil
	}

	// Status has resolved — drive it through the standard
	// webhook processing pipeline.
	//
	// This naturally passes through the Option 2 idempotency gate,
	// so if the webhook also arrives concurrently, one is a safe no-op.
	syntheticPayload := map[string]any{
		"order_id":           result.GatewayOrderID,
		"transaction_status": result.RawStatus,
	}

	if err := u.processWebhook.Execute(ctx, ProcessPaymentWebhookInput{
		Payload: syntheticPayload,
	}); err != nil {
		return nil, fmt.Errorf("check payment status: process resolved status: %w", err)
	}

	resolvedStatus := mapGatewayStatus(result.Status)
	res := CheckPaymentStatusResult{
		Status: resolvedStatus,
		Synced: true,
	}

	return &res, nil
}

// mapGatewayStatus converts a gateway NotificationStatus to the domain
// PaymentStatus that the processWebhook usecase would have set.
func mapGatewayStatus(s paymentgateway.NotificationStatus) domain.PaymentStatus {
	switch s {
	case paymentgateway.NotificationStatusSettlement:
		return domain.PaymentStatusPaid
	case paymentgateway.NotificationStatusExpire:
		return domain.PaymentStatusExpired
	case paymentgateway.NotificationStatusCancel, paymentgateway.NotificationStatusDeny:
		return domain.PaymentStatusCancelled
	default:
		return domain.PaymentStatusPending
	}
}
