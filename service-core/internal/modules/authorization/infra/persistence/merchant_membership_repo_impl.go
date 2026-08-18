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

func (r *staffMembershipRepositoryImpl) ListAccountsByStaffID(
	ctx context.Context,
	exec transaction.Executor,
	staffID uuid.UUID,
) ([]domain.StaffAccountMember, error) {
	query := `
		SELECT
			a.id AS account_id,
			u.id AS user_id,
			a.email,
			u.name,
			u.username,
			u.phone,
			u.avatar_url,
			r.id AS role_id,
			r.code AS role_code,
			r.name AS role_name,
			a.last_login_at,
			sm.created_at
		FROM staff_memberships sm
		JOIN accounts a ON a.id = sm.account_id AND a.deleted_at IS NULL
		JOIN users u ON u.id = a.user_id AND u.deleted_at IS NULL
		JOIN roles r ON r.id = sm.role_id
		WHERE sm.staff_id = $1
		ORDER BY sm.created_at ASC
	`

	rows, err := exec.Query(ctx, query, staffID)
	if err != nil {
		return nil, fmt.Errorf("query staff accounts by staff id failed: %w", err)
	}
	defer rows.Close()

	var members []domain.StaffAccountMember
	for rows.Next() {
		var m domain.StaffAccountMember
		if err := rows.Scan(
			&m.AccountID,
			&m.UserID,
			&m.Email,
			&m.Name,
			&m.Username,
			&m.Phone,
			&m.AvatarURL,
			&m.Role.ID,
			&m.Role.Code,
			&m.Role.Name,
			&m.LastLoginAt,
			&m.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan staff account member failed: %w", err)
		}
		members = append(members, m)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate staff account members failed: %w", err)
	}

	return members, nil
}

func (r *staffMembershipRepositoryImpl) DeleteByAccountIDAndStaffID(
	ctx context.Context,
	exec transaction.Executor,
	accountID uuid.UUID,
	staffID uuid.UUID,
) error {
	query := `
		DELETE FROM staff_memberships
		WHERE account_id = $1 AND staff_id = $2
	`

	_, err := exec.Exec(ctx, query, accountID, staffID)
	if err != nil {
		return fmt.Errorf("delete staff membership failed: %w", err)
	}

	return nil
}

func (r *staffMembershipRepositoryImpl) DeleteByStaffID(
	ctx context.Context,
	exec transaction.Executor,
	staffID uuid.UUID,
) error {
	query := `
		DELETE FROM staff_memberships
		WHERE staff_id = $1
	`

	_, err := exec.Exec(ctx, query, staffID)
	if err != nil {
		return fmt.Errorf("delete staff memberships by staff id failed: %w", err)
	}

	return nil
}
