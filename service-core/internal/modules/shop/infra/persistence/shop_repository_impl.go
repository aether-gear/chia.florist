package persistence

import (
	"context"
	"errors"
	"fmt"
	"service-core/internal/modules/shop/domain"
	"service-core/internal/modules/shop/repository"
	transaction "service-core/internal/shared/transaction"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type shopRepositoryImpl struct{}

func NewShopRepositoryImpl() repository.ShopRepository {
	return &shopRepositoryImpl{}
}

func (r *shopRepositoryImpl) GetByID(
	ctx context.Context,
	exec transaction.Executor,
	shopID uuid.UUID,
) (*domain.Shop, error) {
	query := `
		SELECT
			id,
			name,
			slug,
			description,
			is_active,
			created_at,
			updated_at
		FROM shops
		WHERE id = $1
		LIMIT 1
	`

	var s domain.Shop
	err := exec.QueryRow(ctx, query, shopID).Scan(
		&s.ID,
		&s.Name,
		&s.Slug,
		&s.Description,
		&s.IsActive,
		&s.CreatedAt,
		&s.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("query shop by id failed: %w", err)
	}

	return &s, nil
}

func (r *shopRepositoryImpl) Create(
	ctx context.Context,
	exec transaction.Executor,
	shop domain.Shop,
) error {
	query := `
		INSERT INTO shops (
			id,
			name,
			slug,
			description,
			is_active,
			created_at
		) VALUES ($1,$2,$3,$4,$5,$6)
	`

	_, err := exec.Exec(ctx, query,
		shop.ID,
		shop.Name,
		shop.Slug,
		shop.Description,
		shop.IsActive,
		shop.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("insert shop failed: %w", err)
	}
	return nil
}
