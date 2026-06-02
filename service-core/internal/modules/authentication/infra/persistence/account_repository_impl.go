package persistence

import (
	"context"
	"errors"
	"fmt"
	"time"

	apperrors "service-core/internal/common/errors"
	database "service-core/internal/infra/db"
	"service-core/internal/modules/authentication/domain"
	"service-core/internal/modules/authentication/repository"

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

func (r *accountRepositoryImpl) GetByEmail(email string) (*domain.Account, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

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

func (r *accountRepositoryImpl) GetByID(id uuid.UUID) (*domain.Account, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

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

func (r *accountRepositoryImpl) GetByUserID(id uuid.UUID) (*domain.Account, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

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

func (r *accountRepositoryImpl) ActivateByUserID(id uuid.UUID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		UPDATE accounts
		SET
			status = 'active'
		WHERE user_id = $1
	`

	result, err := r.db.Exec(ctx, query, id)
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

func (r *accountRepositoryImpl) Create(acc domain.Account) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

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

	_, err := r.db.Exec(ctx, query,
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
