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

type merchantMembershipRepositoryImpl struct{}

func NewMerchantMembershipRepositoryImpl() repository.MerchantMembershipRepository {
	return &merchantMembershipRepositoryImpl{}
}

func (r *merchantMembershipRepositoryImpl) GetByAccountID(
	ctx context.Context,
	exec transaction.Executor,
	accountID uuid.UUID,
) (*domain.MerchantMembership, error) {
	query := `
		SELECT
			id,
			merchant_id,
			account_id,
			role_id,
			created_by,
			created_at
		FROM
			merchant_memberships
		WHERE
			account_id = $1
		LIMIT 1
	`

	var m domain.MerchantMembership
	err := exec.QueryRow(ctx, query, accountID).Scan(
		&m.ID,
		&m.MerchantID,
		&m.AccountID,
		&m.RoleID,
		&m.CreatedBy,
		&m.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("query merchant membership failed: %w", err)
	}

	return &m, nil
}

func (r *merchantMembershipRepositoryImpl) GetRolesByAccountID(
	ctx context.Context,
	exec transaction.Executor,
	accountID uuid.UUID,
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
			mm.account_id = $1
	`

	rows, err := exec.Query(ctx, query,
		accountID,
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

func (r *merchantMembershipRepositoryImpl) GetByAccountIDAndMerchantID(
	ctx context.Context,
	exec transaction.Executor,
	accountID uuid.UUID,
	merchantID uuid.UUID,
) (*domain.MerchantMembership, error) {
	query := `
		SELECT
			id,
			merchant_id,
			account_id,
			role_id,
			created_by,
			created_at
		FROM
			merchant_memberships
		WHERE
			account_id = $1 AND merchant_id = $2
		LIMIT 1
	`

	var m domain.MerchantMembership
	err := exec.QueryRow(ctx, query, accountID, merchantID).Scan(
		&m.ID,
		&m.MerchantID,
		&m.AccountID,
		&m.RoleID,
		&m.CreatedBy,
		&m.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("query merchant membership failed: %w", err)
	}

	return &m, nil
}

func (r *merchantMembershipRepositoryImpl) ListRolesByAccountIDAndMerchantID(
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

func (r *merchantMembershipRepositoryImpl) Save(
	ctx context.Context,
	exec transaction.Executor,
	m domain.MerchantMembership,
) error {
	query := `
		INSERT INTO merchant_memberships (
			id,
			merchant_id,
			account_id,
			role_id,
			created_by,
			created_at
		)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := exec.Exec(ctx, query,
		m.ID,
		m.MerchantID,
		m.AccountID,
		m.RoleID,
		m.CreatedBy,
		m.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("insert merchant membership failed: %w", err)
	}

	return nil
}
