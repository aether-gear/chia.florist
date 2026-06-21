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

type paymentAccountRepositoryImpl struct{}

func NewPaymentAccountRepository() repository.PaymentAccountRepository {
	return &paymentAccountRepositoryImpl{}
}

func (r *paymentAccountRepositoryImpl) Save(
	ctx context.Context,
	exec transaction.Executor,
	acc domain.PaymentAccount,
) error {
	query := `
		INSERT INTO payment_accounts (
			id,
			method_id,
			account_name, 
			account_number,
			phone_number,
			qr_string,
			is_active,
			current_load,
			last_used_at,
			created_at
		) VALUES (
		 	$1,$2,$3,$4,$5,$6,$7,$8,$9,$10
		)
	`

	_, err := exec.Exec(ctx, query,
		acc.ID,
		acc.MethodID,
		acc.AccountName,
		acc.AccountNumber,
		acc.PhoneNumber,
		acc.QRString,
		acc.IsActive,
		acc.CurrentLoad,
		acc.LastUsedAt,
		acc.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("insert payment account failed: %w", err)
	}

	return nil
}

func (r *paymentAccountRepositoryImpl) GetByID(
	ctx context.Context,
	exec transaction.Executor,
	paymentID uuid.UUID,
) (*domain.PaymentAccount, error) {
	query := `
		SELECT 
			id,
			method_id,
			account_name,
			account_number,
			phone_number,
			qr_string,
			is_active,
			current_load,
			last_used_at,
			created_at,
			updated_at
		FROM payment_accounts
		WHERE id = $1 AND deleted_at IS NULL
	`

	var acc domain.PaymentAccount
	err := exec.QueryRow(ctx, query, paymentID).Scan(
		&acc.ID,
		&acc.MethodID,
		&acc.AccountName,
		&acc.AccountNumber,
		&acc.PhoneNumber,
		&acc.QRString,
		&acc.IsActive,
		&acc.CurrentLoad,
		&acc.LastUsedAt,
		&acc.CreatedAt,
		&acc.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("query payment account by id failed: %w", err)
	}

	return &acc, nil
}

func (r *paymentAccountRepositoryImpl) RetrieveLeastLoaded(
	ctx context.Context,
	exec transaction.Executor,
	methodID uuid.UUID,
) (*domain.PaymentAccount, error) {
	query := `
		SELECT
			id,
			method_id,
			account_name,
			account_number,
			phone_number,
			qr_string,
			is_active,
			current_load,
			last_used_at,
			created_at,
			updated_at
		FROM payment_accounts
		WHERE method_id = $1 AND is_active = true AND deleted_at IS NULL
		ORDER BY current_load ASC
		LIMIT 1
		FOR UPDATE SKIP LOCKED
	`

	var acc domain.PaymentAccount

	err := exec.QueryRow(ctx, query, methodID).Scan(
		&acc.ID,
		&acc.MethodID,
		&acc.AccountName,
		&acc.AccountNumber,
		&acc.PhoneNumber,
		&acc.QRString,
		&acc.IsActive,
		&acc.CurrentLoad,
		&acc.LastUsedAt,
		&acc.CreatedAt,
		&acc.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("query payment account least loaded failed: %w", err)
	}

	return &acc, nil
}

func (r *paymentAccountRepositoryImpl) IncrementLoad(
	ctx context.Context,
	exec transaction.Executor,
	accountID uuid.UUID,
) error {
	query := `
		UPDATE payment_accounts
		SET
			current_load = current_load + 1,
			last_used_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`

	_, err := exec.Exec(ctx, query, accountID)
	if err != nil {
		return fmt.Errorf("update payment account current load failed: %w", err)
	}

	return nil
}

func (r *paymentAccountRepositoryImpl) DecrementLoad(
	ctx context.Context,
	exec transaction.Executor,
	accountID uuid.UUID,
) error {
	query := `
		UPDATE payment_accounts
		SET current_load = GREATEST(current_load - 1, 0)
		WHERE id = $1 AND deleted_at IS NULL
	`

	_, err := exec.Exec(ctx, query, accountID)
	if err != nil {
		return fmt.Errorf("update payment account current load failed: %w", err)
	}

	return nil
}

func (r *paymentAccountRepositoryImpl) ListByMethodID(
	ctx context.Context,
	exec transaction.Executor,
	methodID uuid.UUID,
) ([]domain.PaymentAccount, error) {
	query := `
		SELECT 
			id,
			method_id,
			account_name,
			account_number,
			phone_number,
			qr_string,
			is_active,
			current_load,
			last_used_at,
			created_at
		FROM payment_accounts
		WHERE method_id = $1 AND deleted_at IS NULL
		ORDER BY created_at ASC
	`

	rows, err := exec.Query(ctx, query, methodID)
	if err != nil {
		return nil, fmt.Errorf("query payment accounts by method id failed: %w", err)
	}
	defer rows.Close()

	var result []domain.PaymentAccount

	for rows.Next() {
		var acc domain.PaymentAccount

		err := rows.Scan(
			&acc.ID,
			&acc.MethodID,
			&acc.AccountName,
			&acc.AccountNumber,
			&acc.PhoneNumber,
			&acc.QRString,
			&acc.IsActive,
			&acc.CurrentLoad,
			&acc.LastUsedAt,
			&acc.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("mapping payment account model to domain failed: %w", err)
		}

		result = append(result, acc)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate payment accounts failed: %w", err)
	}

	return result, nil
}

func (r *paymentAccountRepositoryImpl) ListAll(
	ctx context.Context,
	exec transaction.Executor,
) ([]domain.PaymentAccount, error) {
	query := `
		SELECT
			id,
			method_id,
			account_name,
			account_number,
			phone_number,
			qr_string,
			is_active,
			current_load,
			last_used_at,
			created_at
		FROM payment_accounts
		WHERE deleted_at IS NULL
		ORDER BY created_at ASC
	`

	rows, err := exec.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query payment accounts failed: %w", err)
	}
	defer rows.Close()

	var result []domain.PaymentAccount
	for rows.Next() {
		var acc domain.PaymentAccount

		err := rows.Scan(
			&acc.ID,
			&acc.MethodID,
			&acc.AccountName,
			&acc.AccountNumber,
			&acc.PhoneNumber,
			&acc.QRString,
			&acc.IsActive,
			&acc.CurrentLoad,
			&acc.LastUsedAt,
			&acc.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("mapping payment account model to domain failed: %w", err)
		}

		result = append(result, acc)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate payment accounts failed: %w", err)
	}

	return result, nil
}
