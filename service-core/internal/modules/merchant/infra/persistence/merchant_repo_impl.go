package persistence

import (
	"context"
	"errors"
	"fmt"

	database "service-core/internal/infra/db"
	"service-core/internal/modules/merchant/domain"
	"service-core/internal/modules/merchant/repository"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type merchantRepositoryImpl struct {
	db *pgxpool.Pool
}

func NewMerchantRepositoryImpl(conn *database.Connection) repository.MerchantRepository {
	return &merchantRepositoryImpl{
		db: conn.Pool,
	}
}

func (r *merchantRepositoryImpl) GetByAccountID(ctx context.Context, accountID uuid.UUID) (*domain.Merchant, error) {
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
	err := r.db.QueryRow(ctx, query, accountID).Scan(
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
