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

type ResetPasswordUsecase struct {
	executor      transaction.Executor
	transactor    transaction.Transactor
	accountRepo   repository.AccountRepository
	sessionRepo   repository.SessionRepository
	challengeRepo repository.VerificationChallengeRepository
	pwHasher      repository.PasswordHasher
	auditLogger   applogger.AuditLogger
	sysLogger     applogger.Logger
}

func NewResetPasswordUsecase(
	executor transaction.Executor,
	transactor transaction.Transactor,
	accountRepo repository.AccountRepository,
	sessionRepo repository.SessionRepository,
	challengeRepo repository.VerificationChallengeRepository,
	pwHasher repository.PasswordHasher,
	auditLogger applogger.AuditLogger,
) *ResetPasswordUsecase {
	return &ResetPasswordUsecase{
		executor:      executor,
		transactor:    transactor,
		accountRepo:   accountRepo,
		sessionRepo:   sessionRepo,
		challengeRepo: challengeRepo,
		pwHasher:      pwHasher,
		auditLogger:   auditLogger,
	}
}

func (u *ResetPasswordUsecase) SetSysLogger(sysLogger applogger.Logger) {
	u.sysLogger = sysLogger
}

type ResetPasswordParams struct {
	ChallengeID uuid.UUID
	NewPassword string
}

func (u *ResetPasswordUsecase) Execute(
	ctx context.Context,
	params ResetPasswordParams,
) (err error) {
	now := appclock.Now()

	audit := &applogger.AuditScope{
		Category: "user_action",
		Action:   "reset_password",
		Resource: "account",
		Metadata: map[string]any{"challenge_id": params.ChallengeID.String()},
	}
	defer applogger.TrackAudit(ctx, u.auditLogger, u.sysLogger, audit, &err)()

	challenge, err := u.challengeRepo.GetByID(ctx, u.executor, params.ChallengeID)
	if err != nil {
		return fmt.Errorf("failed to get challenge: %w", err)
	}
	if challenge == nil {
		audit.SetReason("challenge not found")
		return apperrors.NewNotFound(domain.ErrNotFoundChallenge.Error())
	}
	if challenge.Purpose != domain.OTPPurposePasswordReset {
		audit.SetReason("wrong challenge purpose")
		return apperrors.NewNotFound(domain.ErrNotFoundChallenge.Error())
	}
	if challenge.ConsumedAt != nil {
		audit.SetReason("challenge already consumed")
		return apperrors.NewConflict(domain.ErrConsumedChallenge.Error())
	}
	if challenge.VerifiedAt == nil {
		audit.SetReason("challenge not verified yet")
		return apperrors.NewConflict("challenge is not verified")
	}
	if challenge.ExpiresAt.Before(now) {
		audit.SetReason("challenge expired")
		return apperrors.NewConflict(domain.ErrExpiredChallenge.Error())
	}
	if challenge.UserID == nil {
		audit.SetReason("challenge missing user_id")
		return apperrors.NewInternal(fmt.Errorf("challenge has no user_id bound"))
	}

	audit.SetResourceID(challenge.UserID.String())

	hashedPassword, err := u.pwHasher.Hash(params.NewPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	challenge.ConsumedAt = &now

	err = u.transactor.WithinTransaction(ctx, func(exec transaction.Executor) error {
		if err := u.accountRepo.UpdatePasswordByUserID(ctx, exec,
			*challenge.UserID,
			hashedPassword,
		); err != nil {
			return fmt.Errorf("failed to update password: %w", err)
		}

		if err := u.sessionRepo.RevokeAllByUserID(ctx, exec,
			*challenge.UserID,
		); err != nil {
			return fmt.Errorf("failed to revoke sessions: %w", err)
		}

		if err := u.challengeRepo.Save(ctx, exec, *challenge); err != nil {
			return fmt.Errorf("failed to consume challenge: %w", err)
		}

		return nil
	})

	return err
}
