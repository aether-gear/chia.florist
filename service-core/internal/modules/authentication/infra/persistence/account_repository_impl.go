package persistence

import (
	"context"
	"errors"
	"fmt"

	apperrors "service-core/internal/common/errors"
	database "service-core/internal/infra/db"
	"service-core/internal/modules/authentication/domain"
	"service-core/internal/modules/authentication/repository"
	transaction "service-core/internal/shared/transaction"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type accountRepositoryImpl struct {
	db *pgxpool.Pool
}

func NewAccountRepository(conn *database.Connection) repository.AccountRepository {
	return &accountRepositoryImpl{
		db: conn.Pool,
	}
}

func (r *accountRepositoryImpl) GetByEmail(
	ctx context.Context, email string,
) (*domain.Account, error) {
	query := `
		SELECT
			id,
			user_id,
			email,
			password,
			status,
			type,
			last_login_at,
			created_at,
			updated_at
		FROM
			accounts
		WHERE
			email = $1
		LIMIT 1
	`

	var m domain.Account
	err := r.db.QueryRow(ctx, query, email).Scan(
		&m.ID,
		&m.UserID,
		&m.Email,
		&m.Password,
		&m.Status,
		&m.Type,
		&m.LastLoginAt,
		&m.CreatedAt,
		&m.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("query account by email failed: %w", err)
	}

	return &m, nil
}

func (r *accountRepositoryImpl) GetByID(
	ctx context.Context, id uuid.UUID,
) (*domain.Account, error) {
	query := `
		SELECT
			id,
			user_id,
			email,
			password,
			status,
			type,
			last_login_at,
			created_at,
			updated_at
		FROM
			accounts
		WHERE
			id = $1
		LIMIT 1
	`

	var m domain.Account
	err := r.db.QueryRow(ctx, query, id).Scan(
		&m.ID,
		&m.UserID,
		&m.Email,
		&m.Password,
		&m.Status,
		&m.Type,
		&m.LastLoginAt,
		&m.CreatedAt,
		&m.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("query account by id failed: %w", err)
	}

	return &m, nil
}

func (r *accountRepositoryImpl) GetByUserID(
	ctx context.Context, id uuid.UUID,
) (*domain.Account, error) {
	query := `
		SELECT
			id,
			user_id,
			email,
			password,
			status,
			type,
			last_login_at,
			created_at,
			updated_at
		FROM
			accounts
		WHERE
			user_id = $1
		LIMIT 1
	`

	var m domain.Account
	err := r.db.QueryRow(ctx, query, id).Scan(
		&m.ID,
		&m.UserID,
		&m.Email,
		&m.Password,
		&m.Status,
		&m.Type,
		&m.LastLoginAt,
		&m.CreatedAt,
		&m.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("query account by user id failed: %w", err)
	}

	return &m, nil
}

func (r *accountRepositoryImpl) ActivateByUserID(
	ctx context.Context,
	exec transaction.Executor,
	id uuid.UUID,
) error {
	query := `
		UPDATE accounts
		SET
			status = 'active'
		WHERE user_id = $1
	`

	result, err := exec.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf(
			"query to activate account: %w",
			err,
		)
	}

	if result.RowsAffected() == 0 {
		return apperrors.NewNotFound(
			domain.ErrNotFoundAccount.Error(),
		)
	}

	return nil
}

func (r *accountRepositoryImpl) Create(
	ctx context.Context,
	exec transaction.Executor,
	acc domain.Account,
) error {
	query := `
		INSERT INTO accounts (
			id,
			user_id,
			email,
			password,
			status,
			type,
			created_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := exec.Exec(ctx, query,
		acc.ID,
		acc.UserID,
		acc.Email,
		acc.Password,
		acc.Status,
		acc.Type,
		acc.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("insert account failed: %w", err)
	}

	return nil
}
