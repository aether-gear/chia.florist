package persistence

import (
	"context"
	"errors"
	"fmt"

	"service-core/internal/modules/address/domain"
	"service-core/internal/modules/address/repository"
	transaction "service-core/internal/shared/transaction"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type shopAddressRepositoryImpl struct{}

func NewShopAddressRepositoryImpl() repository.ShopAddressRepository {
	return &shopAddressRepositoryImpl{}
}

func (r *shopAddressRepositoryImpl) GetByID(
	ctx context.Context,
	exec transaction.Executor,
	addressID uuid.UUID,
) (*domain.ShopAddress, error) {
	query := `
		SELECT
			id,
			shop_id,
			label,
			phone,
			province,
			city,
			district,
			village,
			full_address,
			postal_code,
			is_active,
			created_at,
			updated_at,
			deleted_at
		FROM shop_addresses
		WHERE id = $1 AND deleted_at IS NULL
		LIMIT 1
	`

	var a domain.ShopAddress
	err := exec.QueryRow(ctx, query, addressID).Scan(
		&a.ID,
		&a.ShopID,
		&a.Label,
		&a.Phone,
		&a.Detail.ProvinceID,
		&a.Detail.CityID,
		&a.Detail.DistrictID,
		&a.Detail.VillageID,
		&a.Detail.FullAddress,
		&a.Detail.PostalCode,
		&a.IsActive,
		&a.CreatedAt,
		&a.UpdatedAt,
		&a.DeletedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("query address by id failed: %w", err)
	}

	return &a, nil
}

func (r *shopAddressRepositoryImpl) GetDefaultByShopID(
	ctx context.Context,
	exec transaction.Executor,
	shopID uuid.UUID,
) (*domain.ShopAddress, error) {
	query := `
		SELECT
			id,
			shop_id,
			label,
			phone,
			province,
			city,
			district,
			village,
			full_address,
			postal_code,
			is_active,
			created_at,
			updated_at,
			deleted_at
		FROM
			shop_addresses
		WHERE
			shop_id = $1
			AND is_active = true
			AND deleted_at IS NULL
		LIMIT 1
	`

	var a domain.ShopAddress
	err := exec.QueryRow(ctx, query, shopID).Scan(
		&a.ID,
		&a.ShopID,
		&a.Label,
		&a.Phone,
		&a.Detail.ProvinceID,
		&a.Detail.CityID,
		&a.Detail.DistrictID,
		&a.Detail.VillageID,
		&a.Detail.FullAddress,
		&a.Detail.PostalCode,
		&a.IsActive,
		&a.CreatedAt,
		&a.UpdatedAt,
		&a.DeletedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("query address default by shop id failed: %w", err)
	}

	return &a, nil
}

func (r *shopAddressRepositoryImpl) GetDefaultsByShopIDs(
	ctx context.Context,
	exec transaction.Executor,
	shopIDs []uuid.UUID,
) (map[uuid.UUID]domain.ShopAddress, error) {
	addressMap := make(map[uuid.UUID]domain.ShopAddress)
	if len(shopIDs) == 0 {
		return addressMap, nil
	}

	shopIDStrings := make([]string, len(shopIDs))
	for i, id := range shopIDs {
		shopIDStrings[i] = id.String()
	}

	query := `
		SELECT
			id,
			shop_id,
			label,
			phone,
			province,
			city,
			district,
			village,
			full_address,
			postal_code,
			is_active,
			created_at,
			updated_at,
			deleted_at
		FROM
			shop_addresses
		WHERE
			shop_id = ANY($1::uuid[])
			AND is_active = true
			AND deleted_at IS NULL
	`

	rows, err := exec.Query(ctx, query, shopIDStrings)
	if err != nil {
		return nil, fmt.Errorf("query shop addresses by shop ids failed: %w", err)
	}
	defer rows.Close()

	var shopAddresses []domain.ShopAddress
	for rows.Next() {
		var sC domain.ShopAddress
		if err := rows.Scan(
			&sC.ID,
			&sC.ShopID,
			&sC.Label,
			&sC.Phone,
			&sC.Detail.ProvinceID,
			&sC.Detail.CityID,
			&sC.Detail.DistrictID,
			&sC.Detail.VillageID,
			&sC.Detail.FullAddress,
			&sC.Detail.PostalCode,
			&sC.IsActive,
			&sC.CreatedAt,
			&sC.UpdatedAt,
			&sC.DeletedAt,
		); err != nil {
			return nil, fmt.Errorf("mapping shop address model to domain failed: %w", err)
		}

		shopAddresses = append(shopAddresses, sC)
	}

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("query address default by shop id failed: %w", err)
	}

	for _, shopAddr := range shopAddresses {
		addressMap[shopAddr.ShopID] = shopAddr
	}

	return addressMap, nil
}

func (r *shopAddressRepositoryImpl) FindByShopID(
	ctx context.Context,
	exec transaction.Executor,
	shopID uuid.UUID,
) ([]domain.ShopAddress, error) {
	query := `
		SELECT
			id,
			shop_id,
			label,
			phone,
			province,
			city,
			district,
			village,
			full_address,
			postal_code,
			is_active,
			created_at,
			updated_at,
			deleted_at
		FROM shop_addresses
		WHERE shop_id = $1 AND deleted_at IS NULL
	`

	rows, err := exec.Query(ctx, query, shopID)
	if err != nil {
		return nil, fmt.Errorf("query address by shop id failed: %w", err)
	}
	defer rows.Close()

	var addresses []domain.ShopAddress
	for rows.Next() {
		var a domain.ShopAddress

		err := rows.Scan(
			&a.ID,
			&a.ShopID,
			&a.Label,
			&a.Phone,
			&a.Detail.ProvinceID,
			&a.Detail.CityID,
			&a.Detail.DistrictID,
			&a.Detail.VillageID,
			&a.Detail.FullAddress,
			&a.Detail.PostalCode,
			&a.IsActive,
			&a.CreatedAt,
			&a.UpdatedAt,
			&a.DeletedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("mapping address model to domain failed: %w", err)
		}

		addresses = append(addresses, a)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate addresses failed: %w", err)
	}

	return addresses, nil
}

func (r *shopAddressRepositoryImpl) Create(
	ctx context.Context,
	exec transaction.Executor,
	address domain.ShopAddress,
) error {
	query := `
		INSERT INTO shop_addresses (
			id,
			shop_id,
			label,
			phone,
			is_active,
			province,
			city,
			district,
			village,
			full_address,
			postal_code,
			created_at
		) VALUES (
			$1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12
		)
	`

	_, err := exec.Exec(ctx, query,
		address.ID,
		address.ShopID,
		address.Label,
		address.Phone,
		address.IsActive,
		address.Detail.ProvinceID,
		address.Detail.CityID,
		address.Detail.DistrictID,
		address.Detail.VillageID,
		address.Detail.FullAddress,
		address.Detail.PostalCode,
		address.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("insert address failed: %w", err)
	}

	return nil
}
