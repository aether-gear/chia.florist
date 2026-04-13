package persistence

import (
	"context"
	"database/sql"
	"errors"
	"service-core/internal/features/auth/domain"
	"service-core/internal/features/auth/repository"
	database "service-core/internal/infra/db"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type authRepositoryImpl struct {
	db *pgxpool.Pool
}

func NewAuthRepository(conn *database.Connection) *authRepositoryImpl {
	return &authRepositoryImpl{db: conn.Pool}
}

func (r *authRepositoryImpl) GetByEmail(email string) (*domain.Account, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT id, email, password, last_login_at
		FROM users
		WHERE email = $1
		LIMIT 1
	`

	var m AccountModel

	err := r.db.QueryRow(ctx, query, email).Scan(
		&m.ID,
		&m.Email,
		&m.Password,
		&m.LastLoginAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	d, err := m.ToDomain()
	if err != nil {
		return nil, err
	}

	return d, nil
}

func (r *authRepositoryImpl) GetByID(id uuid.UUID) (*domain.Account, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT id, email, password,
			last_login_at
		FROM users
		WHERE id = $1
		LIMIT 1
	`

	var m AccountModel

	err := r.db.QueryRow(ctx, query, id).Scan(
		&m.ID,
		&m.Email,
		&m.Password,
		&m.LastLoginAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	d, err := m.ToDomain()
	if err != nil {
		return nil, err
	}

	return d, nil
}

func (r *authRepositoryImpl) Create(acc repository.CreateAccountProps) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		INSERT INTO users (
			id,
			name,
			username,
			email,
			password,
			phone,
			created_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := r.db.Exec(ctx, query,
		acc.ID,
		acc.Name,
		acc.Username,
		acc.Email,
		acc.PasswordHash,
		acc.Phone,
		acc.CreatedAt,
	)

	return err
}
