package usecase

import (
	"context"
	"fmt"
	"time"

	applogger "service-core/internal/common/logger"
	paymentgateway "service-core/internal/infra/payment-gateway"
	paymentRepo "service-core/internal/modules/payment/repository"
	transaction "service-core/internal/shared/transaction"
)

// SyncPendingPaymentsUsecase is the background reconciliation
// job usecase.
//
// It scans for gateway payments that are still 'pending'
// within a look-back window and queries Midtrans directly
// for their current status.
//
// If the status has changed, it drives the payment through
// the existing webhook processing usecase — including the
// idempotency gate — so no double-processing can
// occur even if a late webhook also arrives.
type SyncPendingPaymentsUsecase struct {
	paymentRepo    paymentRepo.PaymentRepository
	paymentGateway paymentgateway.Provider
	processWebhook *ProcessPaymentWebhookUsecase
	executor       transaction.Executor
	logger         applogger.Logger
	lookbackWindow time.Duration
}

func NewSyncPendingPaymentsUsecase(
	paymentRepo paymentRepo.PaymentRepository,
	paymentGateway paymentgateway.Provider,
	processWebhook *ProcessPaymentWebhookUsecase,
	executor transaction.Executor,
	logger applogger.Logger,
	lookbackWindow time.Duration,
) *SyncPendingPaymentsUsecase {
	return &SyncPendingPaymentsUsecase{
		paymentRepo:    paymentRepo,
		paymentGateway: paymentGateway,
		processWebhook: processWebhook,
		executor:       executor,
		logger:         logger,
		lookbackWindow: lookbackWindow,
	}
}

// It fetches all pending gateway payments within
// the look-back window, queries Midtrans for each
// one's current status, and feeds any resolved
// status back through ProcessPaymentWebhookUsecase
// as a synthetic payload.
//
// Per-payment errors are logged and skipped —
// one failure must never prevent the remaining payments
// from being reconciled.
func (u *SyncPendingPaymentsUsecase) Execute(ctx context.Context) {
	since := time.Now().UTC().Add(-u.lookbackWindow)
	var msg string

	payments, err := u.paymentRepo.
		ListPendingGateway(ctx, u.executor, since)
	if err != nil {
		msg = "failed to list pending gateway payments"
		u.logger.Error(ctx, msg,
			applogger.Field{Key: "error", Value: err.Error()},
		)
		return
	}

	if len(payments) == 0 {
		msg = "no pending gateway payments in window"
		u.logger.Info(ctx, msg)
		return
	}

	msg = "payment sync: starting reconciliation cycle"
	u.logger.Info(ctx, msg,
		applogger.Field{Key: "count", Value: len(payments)},
		applogger.Field{Key: "since", Value: since.Format(time.RFC3339)},
	)

	resolved := 0
	for _, payment := range payments {
		if payment.ProviderOrderID == nil {
			continue
		}

		gatewayOrderID := *payment.ProviderOrderID

		result, err := u.paymentGateway.
			GetTransactionStatus(ctx, gatewayOrderID)
		if err != nil {
			msg = "failed to fetch transaction status"
			u.logger.Error(ctx, msg,
				applogger.Field{Key: "payment_id", Value: payment.ID.String()},
				applogger.Field{Key: "gateway_order_id", Value: gatewayOrderID},
				applogger.Field{Key: "error", Value: err.Error()},
			)
			continue
		}

		// Skip if Midtrans still reports pending — nothing to do yet.
		if result.Status == paymentgateway.NotificationStatusPending {
			continue
		}

		// Build a synthetic webhook payload that matches what Midtrans
		// would have sent. ProcessPaymentWebhookUsecase.Execute uses only
		// order_id and transaction_status from the raw payload before
		// calling ParseNotification, which itself re-fetches from Midtrans.
		syntheticPayload := map[string]any{
			"order_id":           result.GatewayOrderID,
			"transaction_status": result.RawStatus,
		}

		if err := u.processWebhook.Execute(ctx, ProcessPaymentWebhookInput{
			Payload: syntheticPayload,
		}); err != nil {
			msg = "failed to process reconciled payment"
			u.logger.Error(ctx, msg,
				applogger.Field{Key: "payment_id", Value: payment.ID.String()},
				applogger.Field{Key: "gateway_order_id", Value: gatewayOrderID},
				applogger.Field{Key: "gateway_status", Value: string(result.Status)},
				applogger.Field{Key: "error", Value: err.Error()},
			)
			continue
		}

		msg = "successfully reconciled payment"
		u.logger.Info(ctx, msg,
			applogger.Field{Key: "payment_id", Value: payment.ID.String()},
			applogger.Field{Key: "gateway_order_id", Value: gatewayOrderID},
			applogger.Field{Key: "gateway_status", Value: string(result.Status)},
		)
		resolved++
	}

	msg = "reconciliation cycle complete"
	u.logger.Info(ctx, msg,
		applogger.Field{Key: "total", Value: len(payments)},
		applogger.Field{Key: "resolved", Value: resolved},
		applogger.Field{Key: "skipped", Value: fmt.Sprintf("%d", len(payments)-resolved)},
	)
}
