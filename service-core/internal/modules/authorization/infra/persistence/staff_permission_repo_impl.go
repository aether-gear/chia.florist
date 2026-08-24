package persistence

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"service-core/internal/modules/authorization/domain"
	"service-core/internal/modules/authorization/repository"
	transaction "service-core/internal/shared/transaction"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type staffPermissionRepositoryImpl struct{}

func NewStaffPermissionRepositoryImpl() repository.StaffPermissionRepository {
	return &staffPermissionRepositoryImpl{}
}

func (r *staffPermissionRepositoryImpl) Save(
	ctx context.Context,
	exec transaction.Executor,
	p domain.StaffPermission,
) error {
	if p.Permissions == nil {
		p.Permissions = []string{}
	}
	permsBytes, err := json.Marshal(p.Permissions)
	if err != nil {
		return fmt.Errorf("failed to marshal permissions: %w", err)
	}

	if p.Rules == nil {
		p.Rules = map[string]any{}
	}
	rulesBytes, err := json.Marshal(p.Rules)
	if err != nil {
		return fmt.Errorf("failed to marshal rules: %w", err)
	}

	query := `
		INSERT INTO staff_permissions (
			id,
			staff_id,
			shop_id,
			permissions,
			rules,
			created_at,
			updated_at
		) VALUES (
			COALESCE(NULLIF($1, '00000000-0000-0000-0000-000000000000'::uuid), gen_random_uuid()),
			$2,
			$3,
			$4::jsonb,
			$5::jsonb,
			NOW(),
			NOW()
		)
		ON CONFLICT (staff_id, shop_id) DO UPDATE SET
			permissions = EXCLUDED.permissions,
			rules = EXCLUDED.rules,
			updated_at = NOW()
	`

	_, err = exec.Exec(ctx, query,
		p.ID,
		p.StaffID,
		p.ShopID,
		string(permsBytes),
		string(rulesBytes),
	)
	if err != nil {
		return fmt.Errorf("save staff shop permission failed: %w", err)
	}

	return nil
}

func (r *staffPermissionRepositoryImpl) GetByStaffIDAndShopID(
	ctx context.Context,
	exec transaction.Executor,
	staffID uuid.UUID,
	shopID uuid.UUID,
) (*domain.StaffPermission, error) {
	query := `
		SELECT
			ssp.id,
			ssp.staff_id,
			ssp.shop_id,
			s.name as shop_name,
			ssp.permissions,
			ssp.rules,
			ssp.created_at,
			ssp.updated_at
		FROM
			staff_permissions ssp
		JOIN
			shops s ON s.id = ssp.shop_id
		WHERE
			ssp.staff_id = $1 AND ssp.shop_id = $2
		LIMIT 1
	`

	var (
		p          domain.StaffPermission
		permsBytes []byte
		rulesBytes []byte
	)

	err := exec.QueryRow(ctx, query, staffID, shopID).Scan(
		&p.ID,
		&p.StaffID,
		&p.ShopID,
		&p.ShopName,
		&permsBytes,
		&rulesBytes,
		&p.CreatedAt,
		&p.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get staff shop permission failed: %w", err)
	}

	if len(permsBytes) > 0 {
		_ = json.Unmarshal(permsBytes, &p.Permissions)
	}
	if len(rulesBytes) > 0 {
		_ = json.Unmarshal(rulesBytes, &p.Rules)
	}

	return &p, nil
}

func (r *staffPermissionRepositoryImpl) ListByStaffID(
	ctx context.Context,
	exec transaction.Executor,
	staffID uuid.UUID,
) ([]domain.StaffPermission, error) {
	query := `
		SELECT
			ssp.id,
			ssp.staff_id,
			ssp.shop_id,
			s.name as shop_name,
			ssp.permissions,
			ssp.rules,
			ssp.created_at,
			ssp.updated_at
		FROM
			staff_permissions ssp
		JOIN
			shops s ON s.id = ssp.shop_id
		WHERE
			ssp.staff_id = $1
		ORDER BY
			s.name ASC
	`

	rows, err := exec.Query(ctx, query, staffID)
	if err != nil {
		return nil, fmt.Errorf("list staff shop permissions failed: %w", err)
	}
	defer rows.Close()

	var result []domain.StaffPermission
	for rows.Next() {
		var (
			p          domain.StaffPermission
			permsBytes []byte
			rulesBytes []byte
		)

		if err := rows.Scan(
			&p.ID,
			&p.StaffID,
			&p.ShopID,
			&p.ShopName,
			&permsBytes,
			&rulesBytes,
			&p.CreatedAt,
			&p.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan staff shop permission failed: %w", err)
		}

		if len(permsBytes) > 0 {
			_ = json.Unmarshal(permsBytes, &p.Permissions)
		}
		if len(rulesBytes) > 0 {
			_ = json.Unmarshal(rulesBytes, &p.Rules)
		}

		result = append(result, p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error in list staff shop permissions: %w", err)
	}

	return result, nil
}

func (r *staffPermissionRepositoryImpl) Delete(
	ctx context.Context,
	exec transaction.Executor,
	staffID uuid.UUID,
	shopID uuid.UUID,
) error {
	query := `
		DELETE FROM staff_permissions
		WHERE staff_id = $1 AND shop_id = $2
	`

	_, err := exec.Exec(ctx, query, staffID, shopID)
	if err != nil {
		return fmt.Errorf("delete staff shop permission failed: %w", err)
	}

	return nil
}
