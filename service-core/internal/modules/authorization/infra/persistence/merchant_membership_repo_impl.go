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

type staffMembershipRepositoryImpl struct{}

func NewStaffMembershipRepositoryImpl() repository.StaffMembershipRepository {
	return &staffMembershipRepositoryImpl{}
}

func (r *staffMembershipRepositoryImpl) GetByAccountID(
	ctx context.Context,
	exec transaction.Executor,
	accountID uuid.UUID,
) (*domain.StaffMembership, error) {
	query := `
		SELECT
			id,
			staff_id,
			account_id,
			role_id,
			created_by,
			created_at
		FROM
			staff_memberships
		WHERE
			account_id = $1
		LIMIT 1
	`

	var m domain.StaffMembership
	err := exec.QueryRow(ctx, query, accountID).Scan(
		&m.ID,
		&m.StaffID,
		&m.AccountID,
		&m.RoleID,
		&m.CreatedBy,
		&m.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("query staff membership failed: %w", err)
	}

	return &m, nil
}

func (r *staffMembershipRepositoryImpl) GetRolesByAccountID(
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
			staff_memberships sm
		JOIN roles ro ON ro.id = sm.role_id
		WHERE
			sm.account_id = $1
	`

	rows, err := exec.Query(ctx, query,
		accountID,
	)
	if err != nil {
		return nil, fmt.Errorf("query roles by staff membership failed: %w", err)
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

func (r *staffMembershipRepositoryImpl) GetByAccountIDAndStaffID(
	ctx context.Context,
	exec transaction.Executor,
	accountID uuid.UUID,
	staffID uuid.UUID,
) (*domain.StaffMembership, error) {
	query := `
		SELECT
			id,
			staff_id,
			account_id,
			role_id,
			created_by,
			created_at
		FROM
			staff_memberships
		WHERE
			account_id = $1 AND staff_id = $2
		LIMIT 1
	`

	var m domain.StaffMembership
	err := exec.QueryRow(ctx, query, accountID, staffID).Scan(
		&m.ID,
		&m.StaffID,
		&m.AccountID,
		&m.RoleID,
		&m.CreatedBy,
		&m.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("query staff membership failed: %w", err)
	}

	return &m, nil
}

func (r *staffMembershipRepositoryImpl) ListRolesByAccountIDAndStaffID(
	ctx context.Context,
	exec transaction.Executor,
	accountID uuid.UUID,
	staffID uuid.UUID,
) ([]domain.Role, error) {
	query := `
		SELECT
			ro.id,
			ro.code,
			ro.name
		FROM
			staff_memberships sm
		JOIN roles ro ON ro.id = sm.role_id
		WHERE
			sm.account_id = $1 AND sm.staff_id = $2
	`

	rows, err := exec.Query(ctx, query,
		accountID,
		staffID,
	)
	if err != nil {
		return nil, fmt.Errorf("query roles by staff membership failed: %w", err)
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

func (r *staffMembershipRepositoryImpl) Save(
	ctx context.Context,
	exec transaction.Executor,
	m domain.StaffMembership,
) error {
	query := `
		INSERT INTO staff_memberships (
			id,
			staff_id,
			account_id,
			role_id,
			created_by,
			created_at
		)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := exec.Exec(ctx, query,
		m.ID,
		m.StaffID,
		m.AccountID,
		m.RoleID,
		m.CreatedBy,
		m.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("insert staff membership failed: %w", err)
	}

	return nil
}
