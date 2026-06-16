package persistence

import (
	"context"
	"errors"
	"fmt"

	"service-core/internal/modules/authorization/domain"
	"service-core/internal/modules/authorization/repository"
	transaction "service-core/internal/shared/transaction"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type roleRepositoryImpl struct{}

func NewRoleRepositoryImpl() repository.RoleRepository {
	return &roleRepositoryImpl{}
}

func (r *roleRepositoryImpl) GetRolesByAccountAndMerchant(
	ctx context.Context,
	exec transaction.Executor,
	accountID uuid.UUID,
	merchantID uuid.UUID,
) ([]domain.Role, error) {
	query := `
		SELECT
			ro.id,
			ro.code,
			ro.name
		FROM
			merchant_memberships mm
		JOIN roles ro ON ro.id = mm.role_id
		WHERE
			mm.account_id = $1 AND mm.merchant_id = $2
	`

	rows, err := exec.Query(ctx, query,
		accountID,
		merchantID,
	)
	if err != nil {
		return nil, fmt.Errorf("query roles by merchant membership failed: %w", err)
	}
	defer rows.Close()

	var roles []domain.Role
	for rows.Next() {
		var role domain.Role
		if err := rows.Scan(&role.ID, &role.Code, &role.Name); err != nil {
			return nil, fmt.Errorf("scan role failed: %w", err)
		}
		roles = append(roles, role)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate roles failed: %w", err)
	}

	return roles, nil
}

func (r *roleRepositoryImpl) GetByCode(
	ctx context.Context,
	exec transaction.Executor,
	code domain.RoleCode,
) (*domain.Role, error) {
	query := `
		SELECT id, code, name
		FROM roles
		WHERE code = $1
		LIMIT 1
	`

	var role domain.Role
	err := exec.QueryRow(ctx, query, code).Scan(
		&role.ID,
		&role.Code,
		&role.Name,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("query role by code failed: %w", err)
	}

	return &role, nil
}
