package persistence

import (
	"context"
	"fmt"
	"time"

	database "service-core/internal/infra/db"
	"service-core/internal/modules/authentication/domain"
	"service-core/internal/modules/authentication/repository"

	"github.com/jackc/pgx/v5/pgxpool"
)

type refreshTokenRepositoryImpl struct {
	db *pgxpool.Pool
}

func NewRefreshTokenRepositoryImpl(conn *database.Connection) repository.RefreshTokenRepository {
	return &refreshTokenRepositoryImpl{
		db: conn.Pool,
	}
}

func (r *refreshTokenRepositoryImpl) Save(refreshToken domain.RefreshToken) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		INSERT INTO refresh_tokens (
			id,
			session_id,
			token_hash,
			expires_at,
			revoked_at,
			created_at
		)
		VALUES
			($1,$2,$3,$4,$5,$6)
		ON CONFLICT
			(id)
		DO UPDATE SET
			session_id = EXCLUDED.session_id,
			token_hash = EXCLUDED.token_hash,
			expires_at = EXCLUDED.expires_at,
			revoked_at = EXCLUDED.revoked_at,
			created_at = EXCLUDED.created_at
	`

	_, err := r.db.Exec(ctx, query,
		refreshToken.ID,
		refreshToken.SessionID,
		refreshToken.TokenHash,
		refreshToken.ExpiresAt,
		refreshToken.RevokedAt,
		refreshToken.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("query to save refresh token: %w", err)
	}

	return nil
}
