package persistence

import (
	"context"
	"fmt"

	"service-core/internal/modules/authentication/domain"
	"service-core/internal/modules/authentication/repository"
	transaction "service-core/internal/shared/transaction"

	"github.com/google/uuid"
)

type refreshTokenRepositoryImpl struct{}

func NewRefreshTokenRepositoryImpl() repository.RefreshTokenRepository {
	return &refreshTokenRepositoryImpl{}
}

func (r *refreshTokenRepositoryImpl) RevokeBySessionID(
	ctx context.Context,
	exec transaction.Executor,
	sessionID uuid.UUID,
) error {
	query := `
			UPDATE
				refresh_tokens
			SET
				revoked_at = now()
			WHERE
				session_id = $1
				AND revoked_at IS NULL
		`

	_, err := exec.Exec(ctx, query, sessionID)
	if err != nil {
		return fmt.Errorf("query to revoke refresh token: %w", err)
	}

	return nil
}

func (r *refreshTokenRepositoryImpl) Save(
	ctx context.Context,
	exec transaction.Executor,
	refreshToken domain.RefreshToken,
) error {
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

	_, err := exec.Exec(ctx, query,
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
