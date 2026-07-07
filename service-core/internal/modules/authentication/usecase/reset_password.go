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

type ResetPasswordUsecase struct {
	executor      transaction.Executor
	transactor    transaction.Transactor
	accountRepo   repository.AccountRepository
	sessionRepo   repository.SessionRepository
	challengeRepo repository.VerificationChallengeRepository
	pwHasher      repository.PasswordHasher
	auditLogger   applogger.AuditLogger
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

type ResetPasswordParams struct {
	ChallengeID uuid.UUID
	NewPassword string
}

func (u *ResetPasswordUsecase) Execute(
	ctx context.Context,
	params ResetPasswordParams,
) error {
	now := time.Now()

	challenge, err := u.challengeRepo.
		GetByID(ctx, u.executor, params.ChallengeID)
	if err != nil {
		return fmt.Errorf("failed to get challenge: %w", err)
	}

	if challenge == nil {
		u.auditLogger.Log(ctx, applogger.AuditEvent{
			Category: "user_action",
			Action:   "reset_password",
			Resource: "account",
			Outcome:  applogger.OutcomeFailure,
			Metadata: map[string]any{
				"challenge_id": params.ChallengeID.String(),
				"reason":       "challenge not found",
			},
		})

		return apperrors.NewNotFound(domain.ErrNotFoundChallenge.Error())
	}

	if challenge.Purpose != domain.OTPPurposePasswordReset {
		u.auditLogger.Log(ctx, applogger.AuditEvent{
			Category: "user_action",
			Action:   "reset_password",
			Resource: "account",
			Outcome:  applogger.OutcomeFailure,
			Metadata: map[string]any{
				"challenge_id": params.ChallengeID.String(),
				"reason":       "wrong challenge purpose",
			},
		})

		return apperrors.NewNotFound(domain.ErrNotFoundChallenge.Error())
	}

	if challenge.ConsumedAt != nil {
		u.auditLogger.Log(ctx, applogger.AuditEvent{
			Category: "user_action",
			Action:   "reset_password",
			Resource: "account",
			Outcome:  applogger.OutcomeFailure,
			Metadata: map[string]any{
				"challenge_id": params.ChallengeID.String(),
				"reason":       "challenge already consumed",
			},
		})

		return apperrors.NewConflict(domain.ErrConsumedChallenge.Error())
	}

	if challenge.VerifiedAt == nil {
		u.auditLogger.Log(ctx, applogger.AuditEvent{
			Category: "user_action",
			Action:   "reset_password",
			Resource: "account",
			Outcome:  applogger.OutcomeFailure,
			Metadata: map[string]any{
				"challenge_id": params.ChallengeID.String(),
				"reason":       "otp not verified",
			},
		})

		return apperrors.NewForbidden("otp must be verified before resetting password")
	}

	if challenge.ExpiresAt.Before(now) {
		u.auditLogger.Log(ctx, applogger.AuditEvent{
			Category: "user_action",
			Action:   "reset_password",
			Resource: "account",
			Outcome:  applogger.OutcomeFailure,
			Metadata: map[string]any{
				"challenge_id": params.ChallengeID.String(),
				"reason":       "challenge expired",
			},
		})

		return apperrors.NewConflict(domain.ErrExpiredChallenge.Error())
	}

	hashedPassword, err := u.pwHasher.Hash(params.NewPassword)
	if err != nil {
		return fmt.Errorf("failed to hash new password: %w", err)
	}

	userID := *challenge.UserID
	challenge.ConsumedAt = &now

	err = u.transactor.WithinTransaction(
		ctx,
		func(exec transaction.Executor) error {
			if err := u.accountRepo.
				UpdatePasswordByUserID(
					ctx,
					exec,
					userID,
					hashedPassword,
				); err != nil {
				return fmt.Errorf("failed to update password: %w", err)
			}

			// Revoke all active sessions so the old credentials
			// are immediately invalidated on every device.
			if err := u.sessionRepo.
				RevokeAllByUserID(ctx, exec, userID); err != nil {
				return fmt.Errorf("failed to revoke sessions: %w", err)
			}

			if err := u.challengeRepo.
				Save(ctx, exec, *challenge); err != nil {
				return fmt.Errorf("failed to consume challenge: %w", err)
			}

			return nil
		},
	)
	if err != nil {
		return err
	}

	u.auditLogger.Log(ctx, applogger.AuditEvent{
		Category: "user_action",
		Action:   "reset_password",
		Resource: "account",
		Outcome:  applogger.OutcomeSuccess,
		Metadata: map[string]any{
			"challenge_id": params.ChallengeID.String(),
			"user_id":      userID.String(),
		},
	})

	return nil
}
