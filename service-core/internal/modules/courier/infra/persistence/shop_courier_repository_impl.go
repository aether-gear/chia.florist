package persistence

import (
	"context"
	"fmt"

	"service-core/internal/modules/courier/domain"
	"service-core/internal/modules/courier/repository"
	transaction "service-core/internal/shared/transaction"

	"github.com/google/uuid"
)

type shopCourierRepositoryImpl struct{}

func NewShopCourierRepositoryImpl() repository.ShopCourierRepository {
	return &shopCourierRepositoryImpl{}
}

func (r *shopCourierRepositoryImpl) ListByShopID(
	ctx context.Context,
	exec transaction.Executor,
	shopID uuid.UUID,
) ([]domain.ShopCourier, error) {
	query := `
		SELECT
			shop_id,
			code,
			active
		FROM shop_couriers
		WHERE shop_id = $1
	`

	rows, err := exec.Query(ctx, query, shopID)
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

func (r *shopCourierRepositoryImpl) ListsByShopIDs(
	ctx context.Context,
	exec transaction.Executor,
	shopIDs []uuid.UUID,
) (map[uuid.UUID][]domain.ShopCourier, error) {
	courierShopMap := make(map[uuid.UUID][]domain.ShopCourier)
	if len(shopIDs) == 0 {
		return courierShopMap, nil
	}

	shopIDStrings := make([]string, len(shopIDs))
	for i, id := range shopIDs {
		shopIDStrings[i] = id.String()
	}

	query := `
		SELECT
			shop_id,
			code,
			active
		FROM shop_couriers
		WHERE shop_id = ANY($1::uuid[])
	`

	rows, err := exec.Query(ctx, query, shopIDStrings)
	if err != nil {
		return nil, fmt.Errorf("query shop courier by shop id failed: %w", err)
	}
	defer rows.Close()

	var shopCouriers []domain.ShopCourier
	for rows.Next() {
		var sC domain.ShopCourier
		if err := rows.Scan(
			&sC.ShopID,
			&sC.Code,
			&sC.Active,
		); err != nil {
			return nil, fmt.Errorf("mapping shop courier model to domain failed: %w", err)
		}

		shopCouriers = append(shopCouriers, sC)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate shop couriers failed: %w", err)
	}

	for _, courierShop := range shopCouriers {
		courierShopMap[courierShop.ShopID] = append(courierShopMap[courierShop.ShopID], courierShop)
	}

	return courierShopMap, nil
}

func (r *shopCourierRepositoryImpl) SaveShopCourier(
	ctx context.Context,
	exec transaction.Executor,
	shopCourier domain.ShopCourier,
) error {
	query := `
		INSERT INTO shop_couriers (
			shop_id,
			code,
			active
		) VALUES (
		 	$1,$2,$3
		)
		ON CONFLICT (shop_id, code)
		DO UPDATE SET
			active = EXCLUDED.active
	`

	_, err := exec.Exec(ctx, query,
		shopCourier.ShopID,
		shopCourier.Code,
		shopCourier.Active,
	)
	if err != nil {
		return fmt.Errorf("save shop couriers failed: %w", err)
	}

	return nil
}
