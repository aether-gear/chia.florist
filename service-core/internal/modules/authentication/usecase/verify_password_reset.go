package usecase

import (
	"context"
	"fmt"
	"time"

	apperrors "service-core/internal/common/errors"
	applogger "service-core/internal/common/logger"
	"service-core/internal/modules/authentication/domain"
	"service-core/internal/modules/authentication/repository"
	transaction "service-core/internal/shared/transaction"

	"github.com/google/uuid"
)

type VerifyPasswordResetUsecase struct {
	executor      transaction.Executor
	challengeRepo repository.VerificationChallengeRepository
	pwHasher      repository.PasswordHasher
	auditLogger   applogger.AuditLogger
}

func NewVerifyPasswordResetUsecase(
	executor transaction.Executor,
	challengeRepo repository.VerificationChallengeRepository,
	pwHasher repository.PasswordHasher,
	auditLogger applogger.AuditLogger,
) *VerifyPasswordResetUsecase {
	return &VerifyPasswordResetUsecase{
		executor:      executor,
		challengeRepo: challengeRepo,
		pwHasher:      pwHasher,
		auditLogger:   auditLogger,
	}
}

type VerifyPasswordResetParams struct {
	ChallengeID uuid.UUID
	OTP         string
}

func (u *VerifyPasswordResetUsecase) Execute(
	ctx context.Context,
	input VerifyPasswordResetParams,
) (*uuid.UUID, error) {
	now := time.Now()

	challenge, err := u.challengeRepo.
		GetByID(ctx, u.executor, input.ChallengeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get challenge: %w", err)
	}

	if challenge == nil {
		u.auditLogger.Log(ctx, applogger.AuditEvent{
			Category: "user_action",
			Action:   "verify_password_reset",
			Resource: "account",
			Outcome:  applogger.OutcomeFailure,
			Metadata: map[string]any{
				"challenge_id": input.ChallengeID.String(),
				"reason":       "challenge not found",
			},
		})

		return nil, apperrors.NewNotFound(domain.ErrNotFoundChallenge.Error())
	}

	if challenge.Purpose != domain.OTPPurposePasswordReset {
		u.auditLogger.Log(ctx, applogger.AuditEvent{
			Category: "user_action",
			Action:   "verify_password_reset",
			Resource: "account",
			Outcome:  applogger.OutcomeFailure,
			Metadata: map[string]any{
				"challenge_id": input.ChallengeID.String(),
				"reason":       "wrong challenge purpose",
			},
		})

		return nil, apperrors.NewNotFound(domain.ErrNotFoundChallenge.Error())
	}

	if challenge.ConsumedAt != nil {
		u.auditLogger.Log(ctx, applogger.AuditEvent{
			Category: "user_action",
			Action:   "verify_password_reset",
			Resource: "account",
			Outcome:  applogger.OutcomeFailure,
			Metadata: map[string]any{
				"challenge_id": input.ChallengeID.String(),
				"reason":       "challenge already consumed",
			},
		})

		return nil, apperrors.NewConflict(domain.ErrConsumedChallenge.Error())
	}

	if challenge.VerifiedAt != nil {
		u.auditLogger.Log(ctx, applogger.AuditEvent{
			Category: "user_action",
			Action:   "verify_password_reset",
			Resource: "account",
			Outcome:  applogger.OutcomeFailure,
			Metadata: map[string]any{
				"challenge_id": input.ChallengeID.String(),
				"reason":       "challenge already verified",
			},
		})

		return nil, apperrors.NewConflict(domain.ErrVerifiedChallenge.Error())
	}

	if challenge.ExpiresAt.Before(now) {
		u.auditLogger.Log(ctx, applogger.AuditEvent{
			Category: "user_action",
			Action:   "verify_password_reset",
			Resource: "account",
			Outcome:  applogger.OutcomeFailure,
			Metadata: map[string]any{
				"challenge_id": input.ChallengeID.String(),
				"reason":       "challenge expired",
			},
		})

		return nil, apperrors.NewConflict(domain.ErrExpiredChallenge.Error())
	}

	if challenge.AttemptCount >= 5 {
		u.auditLogger.Log(ctx, applogger.AuditEvent{
			Category: "user_action",
			Action:   "verify_password_reset",
			Resource: "account",
			Outcome:  applogger.OutcomeFailure,
			Metadata: map[string]any{
				"challenge_id": input.ChallengeID.String(),
				"reason":       "maximum attempts reached",
			},
		})

		return nil, apperrors.NewConflict(domain.ErrMaxAttemptReached.Error())
	}

	if err := u.pwHasher.Compare(challenge.CodeHash, input.OTP); err != nil {
		challenge.AttemptCount++

		if err := u.challengeRepo.Save(ctx, u.executor, *challenge); err != nil {
			return nil, fmt.Errorf("failed to update challenge attempts: %w", err)
		}

		u.auditLogger.Log(ctx, applogger.AuditEvent{
			Category: "user_action",
			Action:   "verify_password_reset",
			Resource: "account",
			Outcome:  applogger.OutcomeFailure,
			Metadata: map[string]any{
				"challenge_id":  input.ChallengeID.String(),
				"attempt_count": challenge.AttemptCount,
				"reason":        "invalid otp",
			},
		})

		return nil, apperrors.NewUnauthorized(domain.ErrInvalidOTP.Error())
	}

	// Mark as verified only — NOT consumed yet.
	//
	// ResetPasswordUsecase will consume it
	// after the new password is set.
	challenge.VerifiedAt = &now

	if err := u.challengeRepo.
		Save(ctx, u.executor, *challenge); err != nil {
		return nil, fmt.Errorf("failed to mark challenge as verified: %w", err)
	}

	u.auditLogger.Log(ctx, applogger.AuditEvent{
		Category: "user_action",
		Action:   "verify_password_reset",
		Resource: "account",
		Outcome:  applogger.OutcomeSuccess,
		Metadata: map[string]any{
			"challenge_id": input.ChallengeID.String(),
		},
	})

	return &challenge.ID, nil
}
