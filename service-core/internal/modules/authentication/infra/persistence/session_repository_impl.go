package persistence

import (
	"context"
	"errors"
	"fmt"

	errorCommon "service-core/internal/common/errors"
	database "service-core/internal/infra/db"
	"service-core/internal/modules/authentication/domain"
	"service-core/internal/modules/authentication/repository"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type sessionRepositoryImpl struct {
	db *pgxpool.Pool
}

func NewSessionRepositoryImpl(conn *database.Connection) repository.SessionRepository {
	return &sessionRepositoryImpl{
		db: conn.Pool,
	}
}

func (r *sessionRepositoryImpl) GetByID(
	ctx context.Context, id uuid.UUID,
) (*domain.Session, error) {
	query := `
		SELECT
			id,
			user_id,
			user_agent,
			ip_address,
			expires_at,
			revoked_at,
			created_at,
			last_activity_at
		FROM
			sessions
		WHERE
			id = $1
		LIMIT 1
	`

	var s domain.Session
	err := r.db.QueryRow(ctx, query, id).Scan(
		&s.ID,
		&s.UserID,
		&s.UserAgent,
		&s.IPAddress,
		&s.ExpiresAt,
		&s.RevokedAt,
		&s.CreatedAt,
		&s.LastActivityAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("query session by id failed: %w", err)
	}

	return &s, nil
}

func (r *sessionRepositoryImpl) RevokeByID(
	ctx context.Context, id uuid.UUID,
) error {
	query := `
		UPDATE
			sessions
		SET
			revoked_at = now()
		WHERE
			id = $1
			AND revoked_at IS NULL
	`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("query to revoke session: %w", err)
	}

	if result.RowsAffected() == 0 {
		return errorCommon.ErrNotFound
	}

	return nil
}

func (r *sessionRepositoryImpl) UpdateLastActivityByID(
	ctx context.Context, id uuid.UUID,
) error {
	query := `
		UPDATE sessions
		SET
			last_activity_at = now()
		WHERE
			id = $1
	`

	_, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("query to update last activity failed: %w", err)
	}

	return nil
}

func (r *sessionRepositoryImpl) Save(
	ctx context.Context, session domain.Session,
) error {
	query := `
		INSERT INTO sessions (
			id,
			user_id,
			user_agent,
			ip_address,
			expires_at,
			revoked_at,
			created_at,
			last_activity_at
		)
		VALUES
			($1,$2,$3,$4,$5,$6,$7,$8)
		ON CONFLICT
			(id)
		DO UPDATE SET
			user_id = EXCLUDED.user_id,
			user_agent = EXCLUDED.user_agent,
			ip_address = EXCLUDED.ip_address,
			expires_at = EXCLUDED.expires_at,
			revoked_at = EXCLUDED.revoked_at,
			created_at = EXCLUDED.created_at,
			last_activity_at = EXCLUDED.last_activity_at
	`

	_, err := r.db.Exec(ctx, query,
		session.ID,
		session.UserID,
		session.UserAgent,
		session.IPAddress,
		session.ExpiresAt,
		session.RevokedAt,
		session.CreatedAt,
		session.LastActivityAt,
	)
	if err != nil {
		return fmt.Errorf("query to save session failed: %w", err)
	}

	return nil
}
