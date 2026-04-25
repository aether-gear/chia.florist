package persistence

import (
	"context"
	"errors"
	"fmt"
	"time"

	database "service-core/internal/infra/db"
	"service-core/internal/modules/payment/domain"
	"service-core/internal/modules/payment/repository"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type paymentAccountRepositoryImpl struct {
	db *pgxpool.Pool
}

func NewPaymentAccountRepository(conn *database.Connection) repository.PaymentAccountRepository {
	return &paymentAccountRepositoryImpl{
		db: conn.Pool,
	}
}

func (r *paymentAccountRepositoryImpl) Save(acc domain.PaymentAccount) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

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

	_, err := r.db.Exec(ctx, query,
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

func (r *paymentAccountRepositoryImpl) GetByID(paymentID uuid.UUID) (*domain.PaymentAccount, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

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

	err := r.db.QueryRow(ctx, query, paymentID).Scan(
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

func (r *paymentAccountRepositoryImpl) AcquireLeastLoaded(methodID uuid.UUID) (*domain.PaymentAccount, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx acquire least loaded payment account failed: %w", err)
	}
	defer tx.Rollback(ctx)

	query := `
		SELECT id, method_id, account_name, account_number,
			phone_number, qr_string, is_active,
			current_load, last_used_at, created_at, updated_at
		FROM payment_accounts
		WHERE method_id = $1 AND is_active = true AND deleted_at IS NULL
		ORDER BY current_load ASC
		LIMIT 1
		FOR UPDATE SKIP LOCKED
	`

	var acc domain.PaymentAccount

	err = tx.QueryRow(ctx, query, methodID).Scan(
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

	updateQuery := `
		UPDATE payment_accounts
		SET current_load = current_load + 1,
		    last_used_at = NOW(),
		    updated_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
		RETURNING current_load, last_used_at, updated_at
	`

	err = tx.QueryRow(ctx, updateQuery, acc.ID).Scan(
		&acc.CurrentLoad,
		&acc.LastUsedAt,
		&acc.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("udpate payment account current load failed: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit tx acquire least loaded payment account failed: %w", err)
	}

	return &acc, nil
}

func (r *paymentAccountRepositoryImpl) IncrementLoad(accountID uuid.UUID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		UPDATE payment_accounts
		SET current_load = current_load + 1
			last_used = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`

	_, err := r.db.Exec(ctx, query, accountID)
	if err != nil {
		return fmt.Errorf("update payment account current load failed: %w", err)
	}

	return nil
}

func (r *paymentAccountRepositoryImpl) DecrementLoad(accountID uuid.UUID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		UPDATE payment_accounts
		SET current_load = GREATEST(current_load - 1, 0)
		WHERE id = $1 AND deleted_at IS NULL
	`

	_, err := r.db.Exec(ctx, query, accountID)
	if err != nil {
		return fmt.Errorf("update payment account current load failed: %w", err)
	}

	return nil
}

func (r *paymentAccountRepositoryImpl) ListByMethodID(methodID uuid.UUID) ([]domain.PaymentAccount, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT id, method_id, account_name, account_number,
		       phone_number, qr_string, is_active,
		       current_load, last_used_at, created_at
		FROM payment_accounts
		WHERE method_id = $1 AND deleted_at IS NULL
		ORDER BY created_at ASC
	`

	rows, err := r.db.Query(ctx, query, methodID)
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

func (r *paymentAccountRepositoryImpl) ListAll() ([]domain.PaymentAccount, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT id, method_id, account_name, account_number,
		       phone_number, qr_string, is_active,
		       current_load, last_used_at, created_at
		FROM payment_accounts
		WHERE deleted_at IS NULL
		ORDER BY created_at ASC
	`

	rows, err := r.db.Query(ctx, query)
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
