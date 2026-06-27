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

type customerAddressRepositoryImpl struct{}

func NewCustomerAddressRepositoryImpl() repository.CustomerAddressRepository {
	return &customerAddressRepositoryImpl{}
}

func (r *customerAddressRepositoryImpl) GetByID(
	ctx context.Context,
	exec transaction.Executor,
	addressID uuid.UUID,
) (*domain.CustomerAddress, error) {
	query := `
		SELECT
			id,
			customer_id,
			recipient_name,
			phone,
			province,
			city,
			district,
			village,
			full_address,
			postal_code,
			is_default,
			created_at,
			updated_at
		FROM customer_addresses
		WHERE id = $1 AND deleted_at IS NULL
	`

	var add domain.CustomerAddress
	err := exec.QueryRow(ctx, query, addressID).Scan(
		&add.ID,
		&add.CustomerID,
		&add.ReceiverName,
		&add.Phone,
		&add.Detail.ProvinceID,
		&add.Detail.CityID,
		&add.Detail.DistrictID,
		&add.Detail.VillageID,
		&add.Detail.FullAddress,
		&add.Detail.PostalCode,
		&add.IsDefault,
		&add.CreatedAt,
		&add.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("query get address by id failed: %w", err)
	}

	return &add, nil
}

func (r *customerAddressRepositoryImpl) GetDefaultByCustomerID(
	ctx context.Context,
	exec transaction.Executor,
	customerID uuid.UUID,
) (*domain.CustomerAddress, error) {
	query := `
		SELECT
			id,
			customer_id,
			recipient_name,
			phone,
			province,
			city,
			district,
			village,
			full_address,
			postal_code,
			is_default,
			created_at,
			updated_at
		FROM
			customer_addresses
		WHERE
			customer_id = $1
			AND is_default = true
			AND deleted_at IS NULL
		LIMIT 1
	`

	var add domain.CustomerAddress
	err := exec.QueryRow(ctx, query, customerID).Scan(
		&add.ID,
		&add.CustomerID,
		&add.ReceiverName,
		&add.Phone,
		&add.Detail.ProvinceID,
		&add.Detail.CityID,
		&add.Detail.DistrictID,
		&add.Detail.VillageID,
		&add.Detail.FullAddress,
		&add.Detail.PostalCode,
		&add.IsDefault,
		&add.CreatedAt,
		&add.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("query get default address by customer id failed: %w", err)
	}

	return &add, nil
}

func (r *customerAddressRepositoryImpl) ListByCustomerID(
	ctx context.Context,
	exec transaction.Executor,
	customerID uuid.UUID,
) ([]domain.CustomerAddress, error) {
	query := `
		SELECT
			id,
			customer_id,
			recipient_name,
			phone,
			province,
			city,
			district,
			village,
			full_address,
			postal_code,
			is_default,
			created_at,
			updated_at,
			deleted_at
		FROM customer_addresses
		WHERE customer_id = $1 AND deleted_at IS NULL
	`

	rows, err := exec.Query(ctx, query, customerID)
	if err != nil {
		return nil, fmt.Errorf("query customer_addresses by customer id failed: %w", err)
	}
	defer rows.Close()

	var addresses []domain.CustomerAddress
	for rows.Next() {
		var a domain.CustomerAddress

		err := rows.Scan(
			&a.ID,
			&a.CustomerID,
			&a.ReceiverName,
			&a.Phone,
			&a.Detail.ProvinceID,
			&a.Detail.CityID,
			&a.Detail.DistrictID,
			&a.Detail.VillageID,
			&a.Detail.FullAddress,
			&a.Detail.PostalCode,
			&a.IsDefault,
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

func (r *customerAddressRepositoryImpl) CountByCustomerID(
	ctx context.Context,
	exec transaction.Executor,
	customerID uuid.UUID,
) (*int, error) {
	query := `
		SELECT COUNT(*)
		FROM
			customer_addresses
		WHERE
			customer_id = $1
			AND deleted_at IS NULL
	`

	var count int
	err := exec.QueryRow(ctx, query, customerID).Scan(&count)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("query count address failed: %w", err)
	}

	return &count, nil
}

func (r *customerAddressRepositoryImpl) UnsetDefaultByCustomerID(
	ctx context.Context,
	exec transaction.Executor,
	customerID uuid.UUID,
) error {
	query := `
		UPDATE
			customer_addresses
		SET
			is_default = false,
			updated_at = NOW()
		WHERE
			customer_id = $1
			AND is_default = true
			AND deleted_at IS NULL
	`

	_, err := exec.Exec(ctx, query, customerID)
	if err != nil {
		return fmt.Errorf("query unset default address failed: %w", err)
	}

	return nil
}

func (r *customerAddressRepositoryImpl) Save(
	ctx context.Context,
	exec transaction.Executor,
	address domain.CustomerAddress,
) error {
	query := `
		INSERT INTO customer_addresses (
			id,
			customer_id,
			recipient_name,
			phone,
			is_default,
			province,
			city,
			district,
			village,
			full_address,
			postal_code,
			created_at,
			updated_at
		) VALUES (
			$1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13
		)
		ON CONFLICT (id)
		DO UPDATE SET
			recipient_name = EXCLUDED.recipient_name,
			phone = EXCLUDED.phone,
			is_default = EXCLUDED.is_default,
			province = EXCLUDED.province,
			city = EXCLUDED.city,
			district = EXCLUDED.district,
			village = EXCLUDED.village,
			full_address = EXCLUDED.full_address,
			postal_code = EXCLUDED.postal_code,
			updated_at = EXCLUDED.updated_at
	`

	_, err := exec.Exec(ctx, query,
		address.ID,
		address.CustomerID,
		address.ReceiverName,
		address.Phone,
		address.IsDefault,
		address.Detail.ProvinceID,
		address.Detail.CityID,
		address.Detail.DistrictID,
		address.Detail.VillageID,
		address.Detail.FullAddress,
		address.Detail.PostalCode,
		address.CreatedAt,
		address.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("insert address failed: %w", err)
	}

	return nil
}

func (r *customerAddressRepositoryImpl) Delete(
	ctx context.Context,
	exec transaction.Executor,
	addressID uuid.UUID,
) error {
	query := `
		UPDATE
			customer_addresses
		SET
			is_default = false,
			deleted_at = now()
		WHERE
			id = $1
			AND deleted_at IS NULL
	`

	_, err := exec.Exec(ctx, query,
		addressID,
	)
	if err != nil {
		return fmt.Errorf("delete address failed: %w", err)
	}

	return nil
}
