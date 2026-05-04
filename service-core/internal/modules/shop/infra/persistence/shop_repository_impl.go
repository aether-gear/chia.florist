package persistence

import (
	"context"
	"errors"
	"fmt"
	database "service-core/internal/infra/db"
	"service-core/internal/modules/shop/domain"
	"service-core/internal/modules/shop/repository"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type shopRepositoryImpl struct {
	db *pgxpool.Pool
}

func NewShopRepositoryImpl(conn *database.Connection) repository.ShopRepository {
	return &shopRepositoryImpl{
		db: conn.Pool,
	}
}

func (r *shopRepositoryImpl) GetByID(shopID uuid.UUID) (*domain.Shop, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

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

	err := r.db.QueryRow(ctx, query, shopID).Scan(
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

func (r *shopRepositoryImpl) Create(shop domain.Shop) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

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

	_, err := r.db.Exec(ctx, query,
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
