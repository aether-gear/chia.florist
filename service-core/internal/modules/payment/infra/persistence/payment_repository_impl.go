package persistence

import (
	"context"
	"errors"
	"fmt"
	"time"

	"service-core/internal/modules/payment/domain"
	"service-core/internal/modules/payment/repository"
	transaction "service-core/internal/shared/transaction"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type paymentRepositoryImpl struct{}

func NewPaymentRepositoryImpl() repository.PaymentRepository {
	return &paymentRepositoryImpl{}
}

func (r *paymentRepositoryImpl) GetByID(
	ctx context.Context,
	exec transaction.Executor,
	id uuid.UUID,
) (*domain.Payment, error) {
	query := `
		SELECT
			id,
			order_id,
			method_id,
			payment_account_id,
			provider,
			provider_payment_id,
			provider_order_id,
			amount,
			status,
			expires_at,
			created_at,
			updated_at
		FROM payments
		WHERE id = $1
	`

	payment, err := r.
		scanPayment(exec.QueryRow(ctx, query, id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("query payment by id failed: %w", err)
	}

	return payment, nil
}

func (r *paymentRepositoryImpl) GetByOrderID(
	ctx context.Context,
	exec transaction.Executor,
	orderID uuid.UUID,
) (*domain.Payment, error) {
	query := `
		SELECT
			id,
			order_id,
			method_id,
			payment_account_id,
			provider,
			provider_payment_id,
			provider_order_id,
			amount,
			status,
			expires_at,
			created_at,
			updated_at
		FROM payments
		WHERE order_id = $1
		ORDER BY created_at DESC
		LIMIT 1
	`

	payment, err := r.scanPayment(
		exec.QueryRow(ctx, query, orderID),
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("query payment by order id failed: %w", err)
	}

	return payment, nil
}

func (r *paymentRepositoryImpl) UpdateStatus(
	ctx context.Context,
	exec transaction.Executor,
	id uuid.UUID,
	status domain.PaymentStatus,
) error {
	query := `
		UPDATE payments
		SET
			status = $2,
			updated_at = NOW()
		WHERE id = $1
	`

	_, err := exec.
		Exec(ctx, query, id, status)

	return err
}

func (r *paymentRepositoryImpl) Save(
	ctx context.Context,
	exec transaction.Executor,
	payment domain.Payment,
) error {
	query := `
		INSERT INTO payments (
			id,
			order_id,
			method_id,
			payment_account_id,
			provider,
			provider_payment_id,
			provider_order_id,
			amount,
			status,
			expires_at,
			created_at,
			updated_at
		)
		VALUES (
			$1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12
		)
		ON CONFLICT (id)
		DO UPDATE SET
			order_id = EXCLUDED.order_id,
			method_id = EXCLUDED.method_id,
			payment_account_id = EXCLUDED.payment_account_id,
			provider = EXCLUDED.provider,
			provider_payment_id = EXCLUDED.provider_payment_id,
			provider_order_id = EXCLUDED.provider_order_id,
			amount = EXCLUDED.amount,
			status = EXCLUDED.status,
			expires_at = EXCLUDED.expires_at,
			updated_at = EXCLUDED.updated_at
	`

	_, err := exec.Exec(ctx, query,
		payment.ID,
		payment.OrderID,
		payment.MethodID,
		payment.PaymentAccountID,
		payment.Provider,
		payment.ProviderPaymentID,
		payment.ProviderOrderID,
		payment.Amount,
		payment.Status,
		payment.ExpiresAt,
		payment.CreatedAt,
		payment.UpdatedAt,
	)

	return err
}

func (r *paymentRepositoryImpl) scanPayment(
	row pgx.Row,
) (*domain.Payment, error) {
	var payment domain.Payment
	err := row.Scan(
		&payment.ID,
		&payment.OrderID,
		&payment.MethodID,
		&payment.PaymentAccountID,
		&payment.Provider,
		&payment.ProviderPaymentID,
		&payment.ProviderOrderID,
		&payment.Amount,
		&payment.Status,
		&payment.ExpiresAt,
		&payment.CreatedAt,
		&payment.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("mapping payment model to domain failed: %w", err)
	}

	return &payment, nil
}

func (r *paymentRepositoryImpl) ListByOrderIDs(
	ctx context.Context,
	exec transaction.Executor,
	orderIDs []uuid.UUID,
) ([]domain.Payment, error) {
	if len(orderIDs) == 0 {
		return []domain.Payment{}, nil
	}

	query := `
		SELECT DISTINCT ON (order_id)
			id,
			order_id,
			method_id,
			payment_account_id,
			provider,
			provider_payment_id,
			provider_order_id,
			amount,
			status,
			expires_at,
			created_at,
			updated_at
		FROM payments
		WHERE order_id = ANY($1::uuid[])
		ORDER BY order_id, created_at DESC
	`

	orderIDStrings := make([]string, len(orderIDs))
	for i, id := range orderIDs {
		orderIDStrings[i] = id.String()
	}

	rows, err := exec.Query(ctx, query, orderIDStrings)
	if err != nil {
		return nil, fmt.Errorf("query payments by order ids failed: %w", err)
	}
	defer rows.Close()

	payments, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.Payment, error) {
		var p domain.Payment
		err := row.Scan(
			&p.ID,
			&p.OrderID,
			&p.MethodID,
			&p.PaymentAccountID,
			&p.Provider,
			&p.ProviderPaymentID,
			&p.ProviderOrderID,
			&p.Amount,
			&p.Status,
			&p.ExpiresAt,
			&p.CreatedAt,
			&p.UpdatedAt,
		)
		return p, err
	})
	if err != nil {
		return nil, fmt.Errorf("scan payments failed: %w", err)
	}

	return payments, nil
}

func (r *paymentRepositoryImpl) ListPendingGateway(
	ctx context.Context,
	exec transaction.Executor,
	since time.Time,
) ([]domain.Payment, error) {
	query := `
		SELECT
			id,
			order_id,
			method_id,
			payment_account_id,
			provider,
			provider_payment_id,
			provider_order_id,
			amount,
			status,
			expires_at,
			created_at,
			updated_at
		FROM payments
		WHERE status = 'pending'
		  AND provider = 'gateway'
		  AND provider_order_id IS NOT NULL
		  AND created_at >= $1
		ORDER BY created_at ASC
	`

	rows, err := exec.Query(ctx, query, since)
	if err != nil {
		return nil, fmt.Errorf("query pending gateway payments failed: %w", err)
	}
	defer rows.Close()

	payments, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.Payment, error) {
		var p domain.Payment
		err := row.Scan(
			&p.ID,
			&p.OrderID,
			&p.MethodID,
			&p.PaymentAccountID,
			&p.Provider,
			&p.ProviderPaymentID,
			&p.ProviderOrderID,
			&p.Amount,
			&p.Status,
			&p.ExpiresAt,
			&p.CreatedAt,
			&p.UpdatedAt,
		)
		return p, err
	})
	if err != nil {
		return nil, fmt.Errorf("scan pending gateway payments failed: %w", err)
	}

	return payments, nil
}
