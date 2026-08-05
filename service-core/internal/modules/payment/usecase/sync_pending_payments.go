package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	apperrors "service-core/internal/common/errors"
	applogger "service-core/internal/common/logger"
	paymentgateway "service-core/internal/infra/payment-gateway"
	inventoryDomain "service-core/internal/modules/inventory/domain"
	inventoryRepo "service-core/internal/modules/inventory/repository"
	orderDomain "service-core/internal/modules/order/domain"
	orderRepo "service-core/internal/modules/order/repository"
	paymentDomain "service-core/internal/modules/payment/domain"
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
	transactor     transaction.Transactor
	orderRepo      orderRepo.OrderRepository
	orderItemRepo  orderRepo.OrderItemRepository
	inventoryRepo  inventoryRepo.InventoryRepository
}

func NewSyncPendingPaymentsUsecase(
	paymentRepo paymentRepo.PaymentRepository,
	paymentGateway paymentgateway.Provider,
	processWebhook *ProcessPaymentWebhookUsecase,
	executor transaction.Executor,
	logger applogger.Logger,
	lookbackWindow time.Duration,
	transactor transaction.Transactor,
	orderRepo orderRepo.OrderRepository,
	orderItemRepo orderRepo.OrderItemRepository,
	inventoryRepo inventoryRepo.InventoryRepository,
) *SyncPendingPaymentsUsecase {
	return &SyncPendingPaymentsUsecase{
		paymentRepo:    paymentRepo,
		paymentGateway: paymentGateway,
		processWebhook: processWebhook,
		executor:       executor,
		logger:         logger,
		lookbackWindow: lookbackWindow,
		transactor:     transactor,
		orderRepo:      orderRepo,
		orderItemRepo:  orderItemRepo,
		inventoryRepo:  inventoryRepo,
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

	payments, err := u.paymentRepo.ListPendingGateway(ctx, u.executor,
		since,
	)
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
			u.logger.Error(ctx, "gateway payment is missing provider order id",
				applogger.Field{Key: "payment_id", Value: payment.ID.String()},
			)
			continue
		}

		gatewayOrderID := *payment.ProviderOrderID
		result, err := u.paymentGateway.GetTransactionStatus(ctx, gatewayOrderID)
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
			if payment.ExpiresAt != nil && time.Now().UTC().After(*payment.ExpiresAt) {
				msg = "payment has expired locally, cancelling at gateway and expiring locally"
				u.logger.Info(ctx, msg,
					applogger.Field{Key: "payment_id", Value: payment.ID.String()},
					applogger.Field{Key: "gateway_order_id", Value: gatewayOrderID},
				)

				// Best effort cancel at gateway
				if err := u.paymentGateway.CancelTransaction(ctx, gatewayOrderID); err != nil {
					u.logger.Warn(ctx, "failed to cancel transaction at gateway, proceeding with local expiry",
						applogger.Field{Key: "payment_id", Value: payment.ID.String()},
						applogger.Field{Key: "gateway_order_id", Value: gatewayOrderID},
						applogger.Field{Key: "error", Value: err.Error()},
					)
				}

				if err := u.expirePaymentLocally(ctx, payment); err != nil {
					u.logger.Error(ctx, "failed to locally expire payment",
						applogger.Field{Key: "payment_id", Value: payment.ID.String()},
						applogger.Field{Key: "error", Value: err.Error()},
					)
				} else {
					resolved++
				}
			}
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

func (u *SyncPendingPaymentsUsecase) expirePaymentLocally(
	ctx context.Context,
	payment paymentDomain.Payment,
) error {
	return u.transactor.WithinTransaction(ctx, func(exec transaction.Executor) error {
		order, err := u.orderRepo.GetByID(ctx, exec, payment.OrderID)
		if err != nil {
			return fmt.Errorf("failed to retrieve order: %w", err)
		}
		if order == nil {
			return fmt.Errorf("order not found: %s", payment.OrderID)
		}

		if err := order.UpdateStatus(orderDomain.OrderStatusExpired); err != nil {
			return fmt.Errorf("invalid order status transition: %w", err)
		}

		if err := u.paymentRepo.UpdateStatus(ctx, exec,
			payment.ID,
			paymentDomain.PaymentStatusExpired,
		); err != nil {
			return fmt.Errorf("failed to update payment status: %w", err)
		}

		if err := u.orderRepo.UpdateStatus(ctx, exec,
			payment.OrderID,
			orderDomain.OrderStatusExpired,
		); err != nil {
			return fmt.Errorf("failed to update order status: %w", err)
		}

		orderItems, err := u.orderItemRepo.ListByOrderID(ctx, exec,
			payment.OrderID,
		)
		if err != nil {
			return fmt.Errorf("failed to list order items: %w", err)
		}

		for _, item := range orderItems {
			if err := u.inventoryRepo.Release(ctx, exec,
				item.ProductID,
				item.ShopID,
				item.Quantity,
			); err != nil {
				if errors.Is(err, inventoryDomain.ErrInsufficientReserved) ||
					errors.Is(err, apperrors.ErrNotFound) {

					msg := "inventory anomaly during payment reconciliation expiry: reserved stock insufficient or missing"
					u.logger.Warn(ctx, msg,
						applogger.Field{Key: "payment_id", Value: payment.ID.String()},
						applogger.Field{Key: "order_id", Value: payment.OrderID.String()},
						applogger.Field{Key: "product_id", Value: item.ProductID.String()},
						applogger.Field{Key: "shop_id", Value: item.ShopID.String()},
						applogger.Field{Key: "requested_qty", Value: item.Quantity},
						applogger.Field{Key: "reason", Value: err.Error()},
					)
					continue
				}
				return fmt.Errorf("failed to release inventory for product %s: %w", item.ProductID, err)
			}
		}

		return nil
	})
}
