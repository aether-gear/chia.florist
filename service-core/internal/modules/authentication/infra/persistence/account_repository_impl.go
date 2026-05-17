package persistence

import (
	"context"
	"errors"
	"fmt"
	"time"

	appErr "service-core/internal/common/errors"
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
			email,
			password,
			last_login_at
		FROM accounts
		WHERE email = $1
		LIMIT 1
	`

	var m domain.Account

	err := r.db.QueryRow(ctx, query, email).Scan(
		&m.ID,
		&m.Email,
		&m.Password,
		&m.LastLoginAt,
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
			email,
			password,
			last_login_at
		FROM accounts
		WHERE id = $1
		LIMIT 1
	`

	var m domain.Account

	err := r.db.QueryRow(ctx, query, id).Scan(
		&m.ID,
		&m.Email,
		&m.Password,
		&m.LastLoginAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("query account by id failed: %w", err)
	}

	return &m, nil
}

func (r *accountRepositoryImpl) ActivateByUserID(
	id uuid.UUID,
) error {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
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
		return appErr.NewNotFound(
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
			created_at
		)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.db.Exec(ctx, query,
		acc.ID,
		acc.UserID,
		acc.Email,
		acc.Password,
		acc.Status,
		acc.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("insert account failed: %w", err)
	}

	return nil
}
