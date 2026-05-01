package persistence

import (
	"context"
	"fmt"
	database "service-core/internal/infra/db"
	"service-core/internal/modules/shop/domain"
	"service-core/internal/modules/shop/repository"

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

func (r *shopRepositoryImpl) Create(shop domain.Shop) error {
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

	_, err := r.db.Exec(context.Background(), query,
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
