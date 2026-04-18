package persistence

import (
	"context"
	"fmt"
	database "service-core/internal/infra/db"
	"service-core/internal/modules/address/domain"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type addressRepositoryImpl struct {
	db *pgxpool.Pool
}

func NewAddressRepositoryImpl(conn *database.Connection) *addressRepositoryImpl {
	return &addressRepositoryImpl{
		db: conn.Pool,
	}
}

func (r *addressRepositoryImpl) GetByUserID(userID uuid.UUID) ([]domain.Address, error) {
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
		FROM addresses
		WHERE user_id = $1 AND deleted_at IS NULL
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
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
			&a.ProvinceID,
			&a.CityID,
			&a.DistrictID,
			&a.VillageID,
			&a.FullAddress,
			&a.PostalCode,
			&a.IsDefault,
			&a.CreatedAt,
			&a.UpdatedAt,
			&a.DeletedAt,
		)
		if err != nil {
			return nil, err
		}

		addresses = append(addresses, a)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return addresses, nil
}

func (r *addressRepositoryImpl) Save(address domain.Address) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if address.DeletedAt != nil {
		query := `
			UPDATE addresses
			SET
				deleted_at = NOW(),
				updated_at = NOW()
			WHERE id = $1 AND deleted_at IS NULL
		`
		res, err := tx.Exec(ctx, query, address.ID)
		if err != nil {
			return err
		}

		if res.RowsAffected() == 0 {
			return fmt.Errorf("address not found or already deleted")
		}
	} else {
		query := `
			INSERT INTO addresses (
				id, user_id, recipient_name, phone, is_default,
				province, city, district, village,
				full_address, postal_code, created_at, updated_at
			) VALUES (
				$1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13
			)
		`

		_, err := tx.Exec(ctx, query,
			address.ID,
			address.UserID,
			address.ReceiverName,
			address.Phone,
			address.IsDefault,
			address.ProvinceID,
			address.CityID,
			address.DistrictID,
			address.VillageID,
			address.FullAddress,
			address.PostalCode,
			address.CreatedAt,
			address.UpdatedAt,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}
