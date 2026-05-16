package persistence

import (
	"context"
	"fmt"
	"time"

	database "service-core/internal/infra/db"
	authDomain "service-core/internal/modules/authentication/domain"
	authRepo "service-core/internal/modules/authentication/repository"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ChallengeRepositoryImpl struct {
	db *pgxpool.Pool
}

func NewChallengeRepository(conn *database.Connection) authRepo.VerificationChallengeRepository {
	return &ChallengeRepositoryImpl{
		db: conn.Pool,
	}
}

func (r *ChallengeRepositoryImpl) Create(challenge authDomain.VerificationChallenge) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		INSERT INTO verification_challenges (
			id,
			user_id,
			type,
			channel,
			purpose,
			target,
			code_hash,
			expires_at,
			attempt_count,
			created_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	_, err := r.db.Exec(ctx, query,
		challenge.ID,
		challenge.UserID,
		challenge.Type,
		challenge.Channel,
		challenge.Purpose,
		challenge.Target,
		challenge.CodeHash,
		challenge.ExpiresAt,
		challenge.AttemptCount,
		challenge.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("query to create verification challenge: %w", err)
	}

	return nil
}
