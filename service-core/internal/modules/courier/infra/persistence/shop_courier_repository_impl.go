package persistence

import (
	"context"
	"fmt"
	"time"

	database "service-core/internal/infra/db"

	"service-core/internal/modules/courier/domain"
	"service-core/internal/modules/courier/repository"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type shopCourierRepositoryImpl struct {
	db *pgxpool.Pool
}

func NewShopCourierRepositoryImpl(conn *database.Connection) repository.ShopCourierRepository {
	return &shopCourierRepositoryImpl{
		db: conn.Pool,
	}
}

func (r *shopCourierRepositoryImpl) GetByShopID(shopID uuid.UUID) ([]domain.ShopCourier, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT
			shop_id,
			code,
			active
		FROM shop_couriers
		WHERE shop_id = $1
	`

	rows, err := r.db.Query(ctx, query, shopID)
	if err != nil {
		return nil, fmt.Errorf("query shop courier by shop id failed: %w", err)
	}
	defer rows.Close()

	var shopCouriers []domain.ShopCourier
	for rows.Next() {
		var sC domain.ShopCourier

		err := rows.Scan(
			&sC.ShopID,
			&sC.Code,
			&sC.Active,
		)

		if err != nil {
			return nil, fmt.Errorf("mapping shop courier model to domain failed: %w", err)
		}

		shopCouriers = append(shopCouriers, sC)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate shop couriers failed: %w", err)
	}

	return shopCouriers, nil
}

func (r *shopCourierRepositoryImpl) SaveShopCouriers(shopCouriers []domain.ShopCourier) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if len(shopCouriers) == 0 {
		return nil
	}

	shopIDs := make([]uuid.UUID, len(shopCouriers))
	codes := make([]string, len(shopCouriers))
	actives := make([]bool, len(shopCouriers))

	for i, sc := range shopCouriers {
		shopIDs[i] = sc.ShopID
		codes[i] = sc.Code
		actives[i] = sc.Active
	}

	query := `
		INSERT INTO shop_couriers (
			shop_id,
			code,
			active
		)
		SELECT *
		FROM UNNEST(
			$1::uuid[],
			$2::text[],
			$3::boolean[]
		)
		ON CONFLICT (shop_id, code)
		DO UPDATE SET
			active = EXCLUDED.active
	`

	_, err := r.db.Exec(
		ctx,
		query,
		shopIDs,
		codes,
		actives,
	)
	if err != nil {
		return fmt.Errorf("save shop couriers failed: %w", err)
	}

	return nil
}
