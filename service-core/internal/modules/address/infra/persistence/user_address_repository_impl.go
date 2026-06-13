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

type userAddressRepositoryImpl struct{}

func NewUserAddressRepositoryImpl() repository.UserAddressRepository {
	return &userAddressRepositoryImpl{}
}

func (r *userAddressRepositoryImpl) GetByUserID(
	ctx context.Context,
	exec transaction.Executor,
	userID uuid.UUID,
) ([]domain.Address, error) {
	query := `
		SELECT
			id,
			user_id,
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
		FROM user_addresses
		WHERE user_id = $1 AND deleted_at IS NULL
	`

	rows, err := exec.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("query user_addresses by user id failed: %w", err)
	}
	defer rows.Close()

	var addresses []domain.Address
	for rows.Next() {
		var a domain.Address

		err := rows.Scan(
			&a.ID,
			&a.UserID,
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

func (r *userAddressRepositoryImpl) CountByUserID(
	ctx context.Context,
	exec transaction.Executor,
	userID uuid.UUID,
) (*int, error) {
	query := `
		SELECT COUNT(*)
		FROM user_addresses
		WHERE user_id = $1
	`

	var count int
	err := exec.QueryRow(ctx, query, userID).Scan(&count)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("query count address failed: %w", err)
	}

	return &count, nil
}

func (r *userAddressRepositoryImpl) Create(
	ctx context.Context,
	exec transaction.Executor,
	address domain.Address,
) error {
	query := `
		INSERT INTO user_addresses (
			id,
			user_id,
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
	`

	_, err := exec.Exec(ctx, query,
		address.ID,
		address.UserID,
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
