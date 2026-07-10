package persistence

import (
	"context"
	"errors"
	"fmt"

	"service-core/internal/modules/payment/domain"
	"service-core/internal/modules/payment/repository"
	transaction "service-core/internal/shared/transaction"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type paymentMethodRepositoryImpl struct{}

func NewPaymentMethodRepository() repository.PaymentMethodRepository {
	return &paymentMethodRepositoryImpl{}
}

func (r *paymentMethodRepositoryImpl) Save(
	ctx context.Context,
	exec transaction.Executor,
	method domain.PaymentMethod,
) error {
	query := `
		INSERT INTO payment_methods (
			id,
			name,
			code,
			type,
			is_active,
			description,
			fee_type,
			fee_amount,
			fee_rate,
			created_at,
			updated_at
		) VALUES (
		 	$1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11
		)
		ON CONFLICT (id) DO UPDATE SET
			name = EXCLUDED.name,
			code = EXCLUDED.code,
			type = EXCLUDED.type,
			is_active = EXCLUDED.is_active,
			description = EXCLUDED.description,
			fee_type = EXCLUDED.fee_type,
			fee_amount = EXCLUDED.fee_amount,
			fee_rate = EXCLUDED.fee_rate,
			updated_at = NOW()
	`

	_, err := exec.Exec(ctx, query,
		method.ID,
		method.Name,
		method.Code,
		method.Type,
		method.IsActive,
		method.Description,
		method.FeeType,
		method.FeeFixed,
		method.FeePercentage,
		method.CreatedAt,
		method.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("save payment method failed: %w", err)
	}

	return nil
}

func (r *paymentMethodRepositoryImpl) FindByName(
	ctx context.Context,
	exec transaction.Executor,
	name string,
) (*domain.PaymentMethod, error) {
	query := `
		SELECT
			id,
			name,
			code,
			type,
			is_active,
			description,
			created_at,
			updated_at
		FROM payment_methods
		WHERE name LIKE $1 || '%'
		LIMIT 1
	`

	var method domain.PaymentMethod
	err := exec.QueryRow(ctx, query, name).Scan(
		&method.ID,
		&method.Name,
		&method.Code,
		&method.Type,
		&method.IsActive,
		&method.Description,
		&method.CreatedAt,
		&method.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("query payment method by name failed: %w", err)
	}

	return &method, nil
}

func (r *paymentMethodRepositoryImpl) GetByID(
	ctx context.Context,
	exec transaction.Executor,
	paymentID uuid.UUID,
) (*domain.PaymentMethod, error) {
	query := `
		SELECT 
			id,
			name,
			code,
			type,
			is_active,
			description,
			created_at,
			updated_at
		FROM payment_methods
		WHERE id = $1
		LIMIT 1
	`

	var method domain.PaymentMethod
	err := exec.QueryRow(ctx, query, paymentID).Scan(
		&method.ID,
		&method.Name,
		&method.Code,
		&method.Type,
		&method.IsActive,
		&method.Description,
		&method.CreatedAt,
		&method.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("query payment method by id failed: %w", err)
	}

	return &method, nil
}

func (r *paymentMethodRepositoryImpl) ListAll(
	ctx context.Context,
	exec transaction.Executor,
) ([]domain.PaymentMethod, error) {
	query := `
		SELECT 
			id,
			name,
			code,
			type,
			is_active,
			description,
			created_at,
			updated_at
		FROM payment_methods
		WHERE deleted_at IS NULL
	`

	rows, err := exec.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query payment methods failed: %w", err)
	}
	defer rows.Close()

	var result []domain.PaymentMethod
	for rows.Next() {
		var row domain.PaymentMethod

		err := rows.Scan(
			&row.ID,
			&row.Name,
			&row.Code,
			&row.Type,
			&row.IsActive,
			&row.Description,
			&row.CreatedAt,
			&row.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("mapping payment method model to domain failed: %w", err)
		}

		result = append(result, row)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate payment methods failed: %w", err)
	}

	return result, nil
}
