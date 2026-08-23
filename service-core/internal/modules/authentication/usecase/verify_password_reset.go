package usecase

import (
	"context"
	"fmt"

	appclock "service-core/internal/common/clock"
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
	sysLogger     applogger.Logger
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

func (u *VerifyPasswordResetUsecase) SetSysLogger(sysLogger applogger.Logger) {
	u.sysLogger = sysLogger
}

type VerifyPasswordResetParams struct {
	ChallengeID uuid.UUID
	OTP         string
}

func (u *VerifyPasswordResetUsecase) Execute(
	ctx context.Context,
	input VerifyPasswordResetParams,
) (verifiedID *uuid.UUID, err error) {
	now := appclock.Now()

	audit := &applogger.AuditScope{
		Category: "user_action",
		Action:   "verify_password_reset",
		Resource: "account",
		Metadata: map[string]any{"challenge_id": input.ChallengeID.String()},
	}
	defer applogger.TrackAudit(ctx, u.auditLogger, u.sysLogger, audit, &err)()

	challenge, err := u.challengeRepo.
		GetByID(ctx, u.executor, input.ChallengeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get challenge: %w", err)
	}

	if challenge == nil {
		audit.SetReason("challenge not found")
		return nil, apperrors.NewNotFound(domain.ErrNotFoundChallenge.Error())
	}

	if challenge.Purpose != domain.OTPPurposePasswordReset {
		audit.SetReason("wrong challenge purpose")
		return nil, apperrors.NewNotFound(domain.ErrNotFoundChallenge.Error())
	}

	if challenge.ConsumedAt != nil {
		audit.SetReason("challenge already consumed")
		return nil, apperrors.NewConflict(domain.ErrConsumedChallenge.Error())
	}

	if challenge.VerifiedAt != nil {
		audit.SetReason("challenge already verified")
		return nil, apperrors.NewConflict(domain.ErrVerifiedChallenge.Error())
	}

	if challenge.ExpiresAt.Before(now) {
		audit.SetReason("challenge expired")
		return nil, apperrors.NewConflict(domain.ErrExpiredChallenge.Error())
	}

	if challenge.AttemptCount >= 5 {
		audit.SetReason("maximum attempts reached")
		return nil, apperrors.NewConflict(domain.ErrMaxAttemptReached.Error())
	}

	if err := u.pwHasher.Compare(challenge.CodeHash, input.OTP); err != nil {
		challenge.AttemptCount++

		if err := u.challengeRepo.
			Save(ctx, u.executor, *challenge); err != nil {
			return nil, fmt.Errorf("failed to update challenge attempts: %w", err)
		}

		audit.SetReason("invalid otp")
		audit.SetMeta("attempt_count", challenge.AttemptCount)
		return nil, apperrors.NewUnauthorized(domain.ErrInvalidOTP.Error())
	}

	challenge.VerifiedAt = &now
	if err := u.challengeRepo.
		Save(ctx, u.executor, *challenge); err != nil {
		return nil, fmt.Errorf("failed to update challenge verification state: %w", err)
	}

	return &challenge.ID, nil
}
