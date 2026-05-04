package persistence

import (
	"context"
	"fmt"
	"time"

	database "service-core/internal/infra/db"
	"service-core/internal/modules/address/domain"
	"service-core/internal/modules/address/repository"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type userAddressRepositoryImpl struct {
	db *pgxpool.Pool
}

func NewUserAddressRepositoryImpl(conn *database.Connection) repository.UserAddressRepository {
	return &userAddressRepositoryImpl{
		db: conn.Pool,
	}
}

func (r *userAddressRepositoryImpl) GetByUserID(userID uuid.UUID) ([]domain.Address, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

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

	rows, err := r.db.Query(ctx, query, userID)
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

func (r *userAddressRepositoryImpl) Create(address domain.Address) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		INSERT INTO user_addresses (
			id, user_id, recipient_name, phone, is_default,
			province, city, district, village,
			full_address, postal_code, created_at, updated_at
		) VALUES (
			$1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13
		)
	`

	_, err := r.db.Exec(ctx, query,
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
