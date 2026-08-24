package usecase

import (
	"context"
	"fmt"

	appclock "service-core/internal/common/clock"
	apperrors "service-core/internal/common/errors"
	applogger "service-core/internal/common/logger"
	"service-core/internal/modules/authentication/domain"
	"service-core/internal/modules/authentication/infra/service"
	"service-core/internal/modules/authentication/repository"
	authorzDomain "service-core/internal/modules/authorization/domain"
	authorzRepo "service-core/internal/modules/authorization/repository"
	customerRepo "service-core/internal/modules/customer/repository"
	userRepo "service-core/internal/modules/user/repository"
	transaction "service-core/internal/shared/transaction"

	"github.com/google/uuid"
)

type VerifyAccountUsecase struct {
	executor       transaction.Executor
	transactor     transaction.Transactor
	accountRepo    repository.AccountRepository
	pwHasher       repository.PasswordHasher
	userRepo       userRepo.UserRepository
	customerRepo   customerRepo.CustomerRepository
	membershipRepo authorzRepo.StaffMembershipRepository
	challengeRepo  repository.VerificationChallengeRepository
	sessionIssuer  repository.SessionIssuerService
	auditLogger    applogger.AuditLogger
	sysLogger      applogger.Logger
}

func NewVerifyAccountUsecase(
	executor transaction.Executor,
	transactor transaction.Transactor,
	accountRepo repository.AccountRepository,
	pwHasher repository.PasswordHasher,
	tokenHasher repository.TokenHasher,
	userRepo userRepo.UserRepository,
	customerRepo customerRepo.CustomerRepository,
	membershipRepo authorzRepo.StaffMembershipRepository,
	challengeRepo repository.VerificationChallengeRepository,
	tokenSvc repository.TokenService,
	sessionRepo repository.SessionRepository,
	refreshTokenRepo repository.RefreshTokenRepository,
	auditLogger applogger.AuditLogger,
) *VerifyAccountUsecase {
	sessionIssuer := service.NewSessionIssuerService(
		transactor,
		tokenSvc,
		tokenHasher,
		sessionRepo,
		refreshTokenRepo,
		accountRepo,
	)

	return &VerifyAccountUsecase{
		executor:       executor,
		transactor:     transactor,
		accountRepo:    accountRepo,
		pwHasher:       pwHasher,
		userRepo:       userRepo,
		customerRepo:   customerRepo,
		membershipRepo: membershipRepo,
		challengeRepo:  challengeRepo,
		sessionIssuer:  sessionIssuer,
		auditLogger:    auditLogger,
	}
}

func (u *VerifyAccountUsecase) SetSessionIssuer(sessionIssuer repository.SessionIssuerService) {
	u.sessionIssuer = sessionIssuer
}

func (u *VerifyAccountUsecase) SetSysLogger(sysLogger applogger.Logger) {
	u.sysLogger = sysLogger
}

type VerifyAccountParams struct {
	UserAgent   *string
	IPAddress   *string
	ChallengeID uuid.UUID
	OTP         string
}

type VerifyAccountResult struct {
	AccessToken, RefreshToken repository.GeneratedToken
}

func (u *VerifyAccountUsecase) Execute(
	ctx context.Context,
	input VerifyAccountParams,
) (result *VerifyAccountResult, err error) {
	now := appclock.Now()

	audit := &applogger.AuditScope{
		Category: "user_action",
		Action:   "verify_account",
		Resource: "account",
		Metadata: map[string]any{"challenge_id": input.ChallengeID.String()},
	}
	defer applogger.TrackAudit(ctx, u.auditLogger, u.sysLogger, audit, &err)()

	challenge, err := u.challengeRepo.GetByID(ctx, u.executor, input.ChallengeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get challenge: %w", err)
	}
	if challenge == nil {
		audit.SetReason("challenge not found")
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

	if err := u.pwHasher.Compare(
		challenge.CodeHash,
		input.OTP,
	); err != nil {
		challenge.AttemptCount++

		if err := u.challengeRepo.Save(ctx, u.executor, *challenge); err != nil {
			return nil, fmt.Errorf("failed to update challenge attempts: %w", err)
		}

		audit.SetReason("invalid otp")
		audit.SetMeta("attempt_count", challenge.AttemptCount)
		return nil, apperrors.NewUnauthorized(domain.ErrInvalidOTP.Error())
	}

	challenge.VerifiedAt = &now
	challenge.ConsumedAt = &now

	var (
		accountID  uuid.UUID
		staffID    *uuid.UUID
		customerID *uuid.UUID
		roleCodes  []authorzDomain.RoleCode
	)

	account, err := u.accountRepo.GetByUserID(ctx, u.executor, *challenge.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get account: %w", err)
	}
	if account != nil {
		accountID = account.ID
		audit.SetResourceID(account.ID.String())

		switch account.Type {
		case domain.AccountTypeCustomer:
			cust, err := u.customerRepo.GetByUserID(ctx, u.executor, account.UserID)
			if err != nil {
				return nil, fmt.Errorf("failed to get customer profile: %w", err)
			}
			if cust != nil {
				customerID = &cust.ID
			}

		case domain.AccountTypeStaff:
			memberStaff, err := u.membershipRepo.GetByAccountID(ctx, u.executor, account.ID)
			if err != nil {
				return nil, fmt.Errorf("failed to get staff membership: %w", err)
			}
			if memberStaff != nil {
				staffID = &memberStaff.StaffID
				roles, err := u.membershipRepo.ListRolesByAccountIDAndStaffID(ctx, u.executor,
					account.ID,
					memberStaff.StaffID,
				)
				if err != nil {
					return nil, fmt.Errorf("failed to list staff roles: %w", err)
				}
				roleCodes = make([]authorzDomain.RoleCode, len(roles))
				for i, r := range roles {
					roleCodes[i] = r.Code
				}
			}
		}
	}

	err = u.transactor.WithinTransaction(ctx, func(exec transaction.Executor) error {
		if err := u.challengeRepo.Save(ctx, exec, *challenge); err != nil {
			return fmt.Errorf("failed to consume challenge: %w", err)
		}

		if err := u.accountRepo.ActivateByUserID(ctx, exec, *challenge.UserID); err != nil {
			return fmt.Errorf("failed to activate account: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	sessionRes, err := u.sessionIssuer.Issue(ctx, repository.IssueSessionParams{
		UserID:     *challenge.UserID,
		AccountID:  accountID,
		UserAgent:  input.UserAgent,
		IPAddress:  input.IPAddress,
		StaffID:    staffID,
		CustomerID: customerID,
		Roles:      roleCodes,
	})
	if err != nil {
		return nil, err
	}

	return &VerifyAccountResult{
		AccessToken:  sessionRes.AccessToken,
		RefreshToken: sessionRes.RefreshToken,
	}, nil
}
