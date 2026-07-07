package usecase

import (
	"context"
	"fmt"
	"time"

	apperrors "service-core/internal/common/errors"
	applogger "service-core/internal/common/logger"
	"service-core/internal/modules/authentication/domain"
	"service-core/internal/modules/authentication/repository"
	mailer "service-core/internal/shared/mailer"
	otp "service-core/internal/shared/otp"
	transaction "service-core/internal/shared/transaction"

	"github.com/google/uuid"
)

type RequestPasswordResetUsecase struct {
	executor      transaction.Executor
	transactor    transaction.Transactor
	accountRepo   repository.AccountRepository
	challengeRepo repository.VerificationChallengeRepository
	pwHasher      repository.PasswordHasher
	otpGen        otp.Generator
	mailer        mailer.Sender
	auditLogger   applogger.AuditLogger
}

func NewRequestPasswordResetUsecase(
	executor transaction.Executor,
	transactor transaction.Transactor,
	accountRepo repository.AccountRepository,
	challengeRepo repository.VerificationChallengeRepository,
	pwHasher repository.PasswordHasher,
	otpGen otp.Generator,
	mailer mailer.Sender,
	auditLogger applogger.AuditLogger,
) *RequestPasswordResetUsecase {
	return &RequestPasswordResetUsecase{
		executor:      executor,
		transactor:    transactor,
		accountRepo:   accountRepo,
		challengeRepo: challengeRepo,
		pwHasher:      pwHasher,
		otpGen:        otpGen,
		mailer:        mailer,
		auditLogger:   auditLogger,
	}
}

type RequestPasswordResetParams struct {
	Email       string
	AccountType domain.AccountType
}

func (u *RequestPasswordResetUsecase) Execute(
	ctx context.Context,
	params RequestPasswordResetParams,
) (*uuid.UUID, error) {
	now := time.Now()

	account, err := u.accountRepo.
		GetByEmail(ctx, u.executor, params.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to get account: %w", err)
	}

	if account == nil {
		u.auditLogger.Log(ctx, applogger.AuditEvent{
			Category: "user_action",
			Action:   "request_password_reset",
			Resource: "account",
			Outcome:  applogger.OutcomeFailure,
			Metadata: map[string]any{
				"email":  params.Email,
				"reason": "account not found (silent)",
			},
		})

		return nil, nil
	}

	if account.Type != params.AccountType {
		u.auditLogger.Log(ctx, applogger.AuditEvent{
			Category: "user_action",
			Action:   "request_password_reset",
			Resource: "account",
			Outcome:  applogger.OutcomeFailure,
			Metadata: map[string]any{
				"email":  params.Email,
				"reason": "account type mismatch",
			},
		})

		return nil, apperrors.NewForbidden("account type mismatch")
	}

	if account.Status != domain.AccountActive {
		u.auditLogger.Log(ctx, applogger.AuditEvent{
			Category: "user_action",
			Action:   "request_password_reset",
			Resource: "account",
			Outcome:  applogger.OutcomeFailure,
			Metadata: map[string]any{
				"email":  params.Email,
				"reason": "account not active",
			},
		})
		return nil, apperrors.NewForbidden(domain.ErrEmailNotVerified.Error())
	}

	otpCode, err := u.otpGen.Generate()
	if err != nil {
		return nil, fmt.Errorf("failed to generate otp: %w", err)
	}

	otpHash, err := u.pwHasher.Hash(otpCode)
	if err != nil {
		return nil, fmt.Errorf("failed to hash otp: %w", err)
	}

	challenge := domain.VerificationChallenge{
		ID:           uuid.New(),
		UserID:       &account.UserID,
		Type:         domain.OTPTypeNumeric,
		Channel:      domain.OTPChannelEmail,
		Purpose:      domain.OTPPurposePasswordReset,
		Target:       params.Email,
		CodeHash:     otpHash,
		ExpiresAt:    now.Add(15 * time.Minute),
		AttemptCount: 0,
		CreatedAt:    now,
	}

	if err := u.challengeRepo.
		Save(ctx, u.executor, challenge); err != nil {
		return nil, fmt.Errorf("failed to save challenge: %w", err)
	}

	mail := mailer.SendInput{
		To:      params.Email,
		Subject: "Reset your password",
		Text:    fmt.Sprintf("Your password reset code is %s", otpCode),
	}
	if err := u.mailer.Send(mail); err != nil {
		return nil, fmt.Errorf("failed to send otp email: %w", err)
	}

	u.auditLogger.Log(ctx, applogger.AuditEvent{
		Category:   "user_action",
		Action:     "request_password_reset",
		Resource:   "account",
		ResourceID: account.ID.String(),
		Outcome:    applogger.OutcomeSuccess,
		Metadata: map[string]any{
			"email": params.Email,
		},
	})

	return &challenge.ID, nil
}
