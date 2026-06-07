package persistence

import (
	"context"
	"errors"
	"fmt"

	"service-core/internal/modules/merchant/domain"
	"service-core/internal/modules/merchant/repository"
	transaction "service-core/internal/shared/transaction"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type merchantRepositoryImpl struct{}

func NewMerchantRepositoryImpl() repository.MerchantRepository {
	return &merchantRepositoryImpl{}
}

func (r *merchantRepositoryImpl) GetByAccountID(
	ctx context.Context,
	exec transaction.Executor,
	accountID uuid.UUID,
) (*domain.Merchant, error) {
	query := `
		SELECT
			id,
			account_id,
			created_at,
			updated_at
		FROM
			merchants
		WHERE
			account_id = $1 AND deleted_at IS NOT NULL
		LIMIT 1
	`

	var merchant domain.Merchant
	err := exec.QueryRow(ctx, query, accountID).Scan(
		&merchant.ID,
		&merchant.AccountID,
		&merchant.CreatedAt,
		&merchant.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("query merchant by account id failed: %w", err)
	}

	return &merchant, nil
}
