package persistence

import (
	"context"
	"errors"
	"fmt"
	"time"

	"service-core/internal/modules/courier/domain"
	"service-core/internal/modules/courier/repository"
	transaction "service-core/internal/shared/transaction"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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
			c.code,
			c.name AS master_name,
			COALESCE(sc.shop_id, $1) AS shop_id,
			sc.branch_name,
			sc.location_address,
			COALESCE(sc.active, false) AS active,
			COALESCE(sc.verification_status, 'unconfigured') AS verification_status,
			sc.verified_at,
			sc.verified_by,
			sc.rejection_reason,
			sc.created_at,
			sc.updated_at
		FROM couriers c
		LEFT JOIN shop_couriers sc ON sc.code = c.code AND sc.shop_id = $1
		WHERE c.is_active = true
		ORDER BY c.name ASC
	`

	rows, err := exec.Query(ctx, query, shopID)
	if err != nil {
		return nil, fmt.Errorf("query shop courier by shop id failed: %w", err)
	}
	defer rows.Close()

	var shopCouriers []domain.ShopCourier
	for rows.Next() {
		var sC domain.ShopCourier
		var statusStr string

		err := rows.Scan(
			&sC.Code,
			&sC.BranchName,
			&sC.ShopID,
			&sC.Name,
			&sC.LocationAddress,
			&sC.Active,
			&statusStr,
			&sC.VerifiedAt,
			&sC.VerifiedBy,
			&sC.RejectionReason,
			&sC.CreatedAt,
			&sC.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("mapping shop courier model to domain failed: %w", err)
		}

		sC.VerificationStatus = domain.CourierVerificationStatus(statusStr)
		shopCouriers = append(shopCouriers, sC)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate shop couriers failed: %w", err)
	}

	return shopCouriers, nil
}

func (r *shopCourierRepositoryImpl) GetByShopIDAndCode(
	ctx context.Context,
	exec transaction.Executor,
	shopID uuid.UUID,
	code string,
) (*domain.ShopCourier, error) {
	query := `
		SELECT
			c.code,
			c.name AS master_name,
			COALESCE(sc.shop_id, $1) AS shop_id,
			sc.branch_name,
			sc.location_address,
			COALESCE(sc.active, false) AS active,
			COALESCE(sc.verification_status, 'unconfigured') AS verification_status,
			sc.verified_at,
			sc.verified_by,
			sc.rejection_reason,
			sc.created_at,
			sc.updated_at
		FROM couriers c
		LEFT JOIN shop_couriers sc ON sc.code = c.code AND sc.shop_id = $1
		WHERE c.code = $2 AND c.is_active = true
	`

	var sC domain.ShopCourier
	var statusStr string

	err := exec.QueryRow(ctx, query, shopID, code).Scan(
		&sC.Code,
		&sC.BranchName,
		&sC.ShopID,
		&sC.Name,
		&sC.LocationAddress,
		&sC.Active,
		&statusStr,
		&sC.VerifiedAt,
		&sC.VerifiedBy,
		&sC.RejectionReason,
		&sC.CreatedAt,
		&sC.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("query shop courier by code failed: %w", err)
	}

	sC.VerificationStatus = domain.CourierVerificationStatus(statusStr)
	return &sC, nil
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
			active,
			verification_status
		FROM shop_couriers
		WHERE shop_id = ANY($1::uuid[])
		  AND active = true
		  AND verification_status = 'verified'
	`

	rows, err := exec.Query(ctx, query, shopIDStrings)
	if err != nil {
		return nil, fmt.Errorf("query active verified shop couriers failed: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var sC domain.ShopCourier
		var statusStr string
		if err := rows.Scan(
			&sC.ShopID,
			&sC.Code,
			&sC.Active,
			&statusStr,
		); err != nil {
			return nil, fmt.Errorf("mapping shop courier model to domain failed: %w", err)
		}
		sC.VerificationStatus = domain.CourierVerificationStatus(statusStr)
		courierShopMap[sC.ShopID] = append(courierShopMap[sC.ShopID], sC)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate shop couriers failed: %w", err)
	}

	return courierShopMap, nil
}

func (r *shopCourierRepositoryImpl) SaveShopCourier(
	ctx context.Context,
	exec transaction.Executor,
	shopCourier domain.ShopCourier,
) error {
	if err := shopCourier.Validate(); err != nil {
		return fmt.Errorf("invalid shop courier entity: %w", err)
	}

	query := `
		INSERT INTO shop_couriers (
			shop_id,
			code,
			branch_name,
			location_address,
			active,
			verification_status,
			verified_at,
			verified_by,
			rejection_reason,
			created_at,
			updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, NOW(), NOW()
		)
		ON CONFLICT (shop_id, code)
		DO UPDATE SET
			branch_name = EXCLUDED.branch_name,
			location_address = EXCLUDED.location_address,
			active = EXCLUDED.active,
			verification_status = EXCLUDED.verification_status,
			verified_at = EXCLUDED.verified_at,
			verified_by = EXCLUDED.verified_by,
			rejection_reason = EXCLUDED.rejection_reason,
			updated_at = NOW()
	`

	_, err := exec.Exec(ctx, query,
		shopCourier.ShopID,
		shopCourier.Code,
		shopCourier.Name,
		shopCourier.LocationAddress,
		shopCourier.Active,
		string(shopCourier.VerificationStatus),
		shopCourier.VerifiedAt,
		shopCourier.VerifiedBy,
		shopCourier.RejectionReason,
	)
	if err != nil {
		return fmt.Errorf("save shop couriers failed: %w", err)
	}

	return nil
}

func (r *shopCourierRepositoryImpl) VerifyShopCourier(
	ctx context.Context,
	exec transaction.Executor,
	shopID uuid.UUID,
	code string,
	status domain.CourierVerificationStatus,
	active bool,
	verifiedBy uuid.UUID,
	rejectionReason *string,
) error {
	now := time.Now()
	query := `
		UPDATE shop_couriers
		SET
			verification_status = $3,
			active = $4,
			verified_at = $5,
			verified_by = $6,
			rejection_reason = $7,
			updated_at = $5
		WHERE shop_id = $1 AND code = $2
	`

	tag, err := exec.Exec(ctx, query,
		shopID,
		code,
		string(status),
		active,
		now,
		verifiedBy,
		rejectionReason,
	)
	if err != nil {
		return fmt.Errorf("verify shop courier failed: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("shop courier record not found to verify")
	}

	return nil
}
