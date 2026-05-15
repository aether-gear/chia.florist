package persistence

import (
	"context"
	"errors"
	"fmt"
	"time"

	database "service-core/internal/infra/db"
	"service-core/internal/modules/authentication/domain"
	"service-core/internal/modules/authentication/repository"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type authRepositoryImpl struct {
	db *pgxpool.Pool
}

func NewAuthRepository(conn *database.Connection) repository.AuthRepository {
	return &authRepositoryImpl{
		db: conn.Pool,
	}
}

func (r *authRepositoryImpl) GetByEmail(email string) (*domain.Account, error) {
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

func (r *authRepositoryImpl) GetByID(id uuid.UUID) (*domain.Account, error) {
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

func (r *authRepositoryImpl) Create(acc repository.CreateAccountProps) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		INSERT INTO accounts (
			id,
			user_id,
			email,
			password,
			created_at
		)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.db.Exec(ctx, query,
		acc.ID,
		acc.UserID,
		acc.Email,
		acc.Password,
		acc.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("insert account failed: %w", err)
	}

	return nil
}
