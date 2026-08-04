package usecase

import (
	"context"
	"fmt"
	"sync"
	"time"

	applogger "service-core/internal/common/logger"
	paymentgateway "service-core/internal/infra/payment-gateway"
	inventoryRepo "service-core/internal/modules/inventory/repository"
	orderDomain "service-core/internal/modules/order/domain"
	orderRepo "service-core/internal/modules/order/repository"
	paymentDomain "service-core/internal/modules/payment/domain"
	paymentRepo "service-core/internal/modules/payment/repository"
	transaction "service-core/internal/shared/transaction"
)

type ExpirePastDuePaymentsUsecase struct {
	paymentRepo    paymentRepo.PaymentRepository
	paymentGateway paymentgateway.Provider
	executor       transaction.Executor
	transactor     transaction.Transactor
	orderRepo      orderRepo.OrderRepository
	orderItemRepo  orderRepo.OrderItemRepository
	inventoryRepo  inventoryRepo.InventoryRepository
	logger         applogger.Logger
	batchSize      int
	concurrency    int
}

func NewExpirePastDuePaymentsUsecase(
	paymentRepo paymentRepo.PaymentRepository,
	paymentGateway paymentgateway.Provider,
	executor transaction.Executor,
	transactor transaction.Transactor,
	orderRepo orderRepo.OrderRepository,
	orderItemRepo orderRepo.OrderItemRepository,
	inventoryRepo inventoryRepo.InventoryRepository,
	logger applogger.Logger,
	batchSize int,
	concurrency int,
) *ExpirePastDuePaymentsUsecase {
	if batchSize <= 0 {
		batchSize = 100
	}

	if concurrency <= 0 {
		concurrency = 5
	}

	return &ExpirePastDuePaymentsUsecase{
		paymentRepo:    paymentRepo,
		paymentGateway: paymentGateway,
		executor:       executor,
		transactor:     transactor,
		orderRepo:      orderRepo,
		orderItemRepo:  orderItemRepo,
		inventoryRepo:  inventoryRepo,
		logger:         logger,
		batchSize:      batchSize,
		concurrency:    concurrency,
	}
}

func (u *ExpirePastDuePaymentsUsecase) Execute(ctx context.Context) {
	now := time.Now().UTC()
	payments, err := u.paymentRepo.ListPastDuePending(ctx, u.executor, now,
		u.batchSize,
	)
	if err != nil {
		u.logger.Error(ctx, "failed to list past-due pending payments",
			applogger.Field{Key: "error", Value: err.Error()},
		)
		return
	}

	if len(payments) == 0 {
		return
	}

	u.logger.Info(ctx, "payment expiry job: processing past-due payments",
		applogger.Field{Key: "count", Value: len(payments)},
	)

	paymentChan := make(chan paymentDomain.Payment, len(payments))
	for _, p := range payments {
		paymentChan <- p
	}
	close(paymentChan)

	numWorkers := min(len(payments), u.concurrency)

	var (
		successCount int64
		wg           sync.WaitGroup
		mu           sync.Mutex
	)

	for range numWorkers {
		wg.Add(1)

		go func() {
			defer wg.Done()
			for payment := range paymentChan {
				if err := u.expireSinglePayment(ctx, payment); err != nil {
					u.logger.Error(ctx, "failed to expire payment",
						applogger.Field{Key: "payment_id", Value: payment.ID.String()},
						applogger.Field{Key: "order_id", Value: payment.OrderID.String()},
						applogger.Field{Key: "error", Value: err.Error()},
					)
				} else {
					mu.Lock()
					successCount++
					mu.Unlock()
				}
			}
		}()
	}

	wg.Wait()

	u.logger.Info(ctx, "payment expiry job: processing complete",
		applogger.Field{Key: "total", Value: len(payments)},
		applogger.Field{Key: "expired_count", Value: successCount},
	)
}

func (u *ExpirePastDuePaymentsUsecase) expireSinglePayment(
	ctx context.Context,
	payment paymentDomain.Payment,
) error {
	// Best-effort gateway cancellation
	if payment.Provider == "gateway" &&
		payment.ProviderOrderID != nil &&
		*payment.ProviderOrderID != "" {

		if err := u.paymentGateway.CancelTransaction(ctx, *payment.ProviderOrderID); err != nil {
			u.logger.Warn(ctx, "failed to cancel transaction at gateway during expiry, proceeding locally",
				applogger.Field{Key: "payment_id", Value: payment.ID.String()},
				applogger.Field{Key: "gateway_order_id", Value: *payment.ProviderOrderID},
				applogger.Field{Key: "error", Value: err.Error()},
			)
		}
	}

	return u.transactor.WithinTransaction(ctx, func(exec transaction.Executor) error {
		order, err := u.orderRepo.GetByID(ctx, exec, payment.OrderID)
		if err != nil {
			return fmt.Errorf("failed to retrieve order: %w", err)
		}
		if order == nil {
			return fmt.Errorf("order not found: %s", payment.OrderID)
		}

		if err := order.UpdateStatus(orderDomain.OrderStatusCancelled); err != nil {
			return fmt.Errorf("invalid order status transition: %w", err)
		}

		if err := u.paymentRepo.UpdateStatus(ctx, exec,
			payment.ID,
			paymentDomain.PaymentStatusExpired,
		); err != nil {
			return fmt.Errorf("failed to update payment status to expired: %w", err)
		}

		if err := u.orderRepo.UpdateStatus(ctx, exec,
			payment.OrderID,
			orderDomain.OrderStatusCancelled,
		); err != nil {
			return fmt.Errorf("failed to update order status to cancelled: %w", err)
		}

		orderItems, err := u.orderItemRepo.ListByOrderID(ctx, exec, payment.OrderID)
		if err != nil {
			return fmt.Errorf("failed to list order items: %w", err)
		}

		for _, item := range orderItems {
			if err := u.inventoryRepo.Release(ctx, exec,
				item.ProductID,
				item.ShopID,
				item.Quantity,
			); err != nil {
				return fmt.Errorf("failed to release inventory for product %s: %w", item.ProductID, err)
			}
		}

		return nil
	})
}
