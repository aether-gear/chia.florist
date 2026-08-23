package usecase

import (
	"context"
	"fmt"
	"time"

	apperrors "service-core/internal/common/errors"
	applogger "service-core/internal/common/logger"
	"service-core/internal/modules/authentication/domain"
	"service-core/internal/modules/authentication/infra/service"
	"service-core/internal/modules/authentication/repository"
	mailer "service-core/internal/shared/mailer"
	otp "service-core/internal/shared/otp"
	transaction "service-core/internal/shared/transaction"

	"github.com/google/uuid"
)

type RequestPasswordResetUsecase struct {
	executor     transaction.Executor
	transactor   transaction.Transactor
	accountRepo  repository.AccountRepository
	challengeSvc repository.ChallengeService
	auditLogger  applogger.AuditLogger
	sysLogger    applogger.Logger
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
	challengeSvc := service.NewChallengeService(
		transactor,
		challengeRepo,
		pwHasher,
		otpGen,
		mailer,
	)

	return &RequestPasswordResetUsecase{
		executor:     executor,
		transactor:   transactor,
		accountRepo:  accountRepo,
		challengeSvc: challengeSvc,
		auditLogger:  auditLogger,
	}
}

func (u *RequestPasswordResetUsecase) SetChallengeService(challengeSvc repository.ChallengeService) {
	u.challengeSvc = challengeSvc
}

func (u *RequestPasswordResetUsecase) SetSysLogger(sysLogger applogger.Logger) {
	u.sysLogger = sysLogger
}

type RequestPasswordResetParams struct {
	Email       string
	AccountType domain.AccountType
}

func (u *RequestPasswordResetUsecase) Execute(
	ctx context.Context,
	params RequestPasswordResetParams,
) (challengeID *uuid.UUID, err error) {
	audit := &applogger.AuditScope{
		Category: "user_action",
		Action:   "request_password_reset",
		Resource: "account",
		Metadata: map[string]any{"email": params.Email},
	}
	defer applogger.TrackAudit(ctx, u.auditLogger, u.sysLogger, audit, &err)()

	account, err := u.accountRepo.GetByEmail(ctx, u.executor, params.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to get account: %w", err)
	}

	if account == nil {
		// Silent success for non-existent accounts to prevent email enumeration
		audit.SetReason("account not found (silent)")
		return nil, nil
	}

	audit.SetResourceID(account.ID.String())

	if account.Type != params.AccountType {
		audit.SetReason("account type mismatch")
		return nil, apperrors.NewForbidden("account type mismatch")
	}
	if account.Status != domain.AccountActive {
		audit.SetReason("account not active")
		return nil, apperrors.NewForbidden(domain.ErrEmailNotVerified.Error())
	}

	chID, err := u.challengeSvc.CreateAndSend(ctx, repository.CreateChallengeParams{
		UserID:   &account.UserID,
		Email:    params.Email,
		Purpose:  domain.OTPPurposePasswordReset,
		Duration: 15 * time.Minute,
	})
	if err != nil {
		return nil, err
	}

	return chID, nil
}
