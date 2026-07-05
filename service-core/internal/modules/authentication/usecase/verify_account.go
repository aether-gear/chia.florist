package usecase

import (
	"context"
	"fmt"
	"time"

	apperrors "service-core/internal/common/errors"
	applogger "service-core/internal/common/logger"
	"service-core/internal/modules/authentication/domain"
	"service-core/internal/modules/authentication/repository"
	authorzDomain "service-core/internal/modules/authorization/domain"
	authorzRepo "service-core/internal/modules/authorization/repository"
	customerRepo "service-core/internal/modules/customer/repository"
	userRepo "service-core/internal/modules/user/repository"
	transaction "service-core/internal/shared/transaction"

	"github.com/google/uuid"
)

type VerifyAccountUsecase struct {
	executor         transaction.Executor
	transactor       transaction.Transactor
	accountRepo      repository.AccountRepository
	pwHasher         repository.PasswordHasher
	tokenHasher      repository.TokenHasher
	userRepo         userRepo.UserRepository
	customerRepo     customerRepo.CustomerRepository
	membershipRepo   authorzRepo.StaffMembershipRepository
	challengeRepo    repository.VerificationChallengeRepository
	tokenSvc         repository.TokenService
	sessionRepo      repository.SessionRepository
	refreshTokenRepo repository.RefreshTokenRepository
	auditLogger      applogger.AuditLogger
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
	return &VerifyAccountUsecase{
		executor:         executor,
		transactor:       transactor,
		accountRepo:      accountRepo,
		pwHasher:         pwHasher,
		tokenHasher:      tokenHasher,
		userRepo:         userRepo,
		customerRepo:     customerRepo,
		membershipRepo:   membershipRepo,
		challengeRepo:    challengeRepo,
		tokenSvc:         tokenSvc,
		sessionRepo:      sessionRepo,
		refreshTokenRepo: refreshTokenRepo,
		auditLogger:      auditLogger,
	}
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
) (*VerifyAccountResult, error) {
	now := time.Now()

	challenge, err := u.challengeRepo.
		GetByID(ctx, u.executor, input.ChallengeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get challenge: %w", err)
	}
	if challenge == nil {
		u.auditLogger.Log(ctx, applogger.AuditEvent{
			Category: "user_action",
			Action:   "verify_account",
			Resource: "account",
			Outcome:  applogger.OutcomeFailure,
			Metadata: map[string]any{"challenge_id": input.ChallengeID.String(), "reason": "challenge not found"},
		})
		return nil, apperrors.NewNotFound(domain.ErrNotFoundChallenge.Error())
	}

	if challenge.ConsumedAt != nil {
		u.auditLogger.Log(ctx, applogger.AuditEvent{
			Category: "user_action",
			Action:   "verify_account",
			Resource: "account",
			Outcome:  applogger.OutcomeFailure,
			Metadata: map[string]any{"challenge_id": input.ChallengeID.String(), "reason": "challenge already consumed"},
		})
		return nil, apperrors.NewConflict(domain.ErrConsumedChallenge.Error())
	}
	if challenge.VerifiedAt != nil {
		u.auditLogger.Log(ctx, applogger.AuditEvent{
			Category: "user_action",
			Action:   "verify_account",
			Resource: "account",
			Outcome:  applogger.OutcomeFailure,
			Metadata: map[string]any{"challenge_id": input.ChallengeID.String(), "reason": "challenge already verified"},
		})
		return nil, apperrors.NewConflict(domain.ErrVerifiedChallenge.Error())
	}
	if challenge.ExpiresAt.Before(now) {
		u.auditLogger.Log(ctx, applogger.AuditEvent{
			Category: "user_action",
			Action:   "verify_account",
			Resource: "account",
			Outcome:  applogger.OutcomeFailure,
			Metadata: map[string]any{"challenge_id": input.ChallengeID.String(), "reason": "challenge expired"},
		})
		return nil, apperrors.NewConflict(domain.ErrExpiredChallenge.Error())
	}
	if challenge.AttemptCount >= 5 {
		u.auditLogger.Log(ctx, applogger.AuditEvent{
			Category: "user_action",
			Action:   "verify_account",
			Resource: "account",
			Outcome:  applogger.OutcomeFailure,
			Metadata: map[string]any{"challenge_id": input.ChallengeID.String(), "reason": "maximum attempts reached"},
		})
		return nil, apperrors.NewConflict(domain.ErrMaxAttemptReached.Error())
	}

	if err := u.pwHasher.
		Compare(
			challenge.CodeHash,
			input.OTP,
		); err != nil {
		challenge.AttemptCount++

		if err := u.challengeRepo.
			Save(
				ctx,
				u.executor,
				*challenge,
			); err != nil {
			return nil, fmt.Errorf(
				"failed to update challenge attempts: %w",
				err,
			)
		}

		u.auditLogger.Log(ctx, applogger.AuditEvent{
			Category: "user_action",
			Action:   "verify_account",
			Resource: "account",
			Outcome:  applogger.OutcomeFailure,
			Metadata: map[string]any{"challenge_id": input.ChallengeID.String(), "attempt_count": challenge.AttemptCount, "reason": "invalid otp"},
		})

		return nil, apperrors.NewUnauthorized(domain.ErrInvalidOTP.Error())
	}

	challenge.VerifiedAt = &now
	challenge.ConsumedAt = &now

	var (
		staffID    *uuid.UUID
		customerID *uuid.UUID
		roleCodes  []authorzDomain.RoleCode
	)

	account, err := u.accountRepo.
		GetByUserID(ctx, u.executor, *challenge.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get account: %w", err)
	}
	if account != nil {
		switch account.Type {
		case domain.AccountTypeCustomer:
			cust, err := u.customerRepo.
				GetByUserID(ctx, u.executor, account.UserID)
			if err != nil {
				return nil, fmt.Errorf("failed to get customer profile: %w", err)
			}
			if cust != nil {
				customerID = &cust.ID
			}
		case domain.AccountTypeStaff:
			memberStaff, err := u.membershipRepo.
				GetByAccountID(ctx, u.executor, account.ID)
			if err != nil {
				return nil, fmt.Errorf("failed to get staff membership: %w", err)
			}
			if memberStaff != nil {
				staffID = &memberStaff.StaffID
				roles, err := u.membershipRepo.
					ListRolesByAccountIDAndStaffID(
						ctx,
						u.executor,
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

	sessionID := uuid.New()
	accessTkn, err := u.tokenSvc.
		Generate(repository.GenerateTokenParams{
			UserID:     *challenge.UserID,
			SessionID:  sessionID,
			StaffID:    staffID,
			CustomerID: customerID,
			Roles:      roleCodes,
			Type:       domain.TokenTypeAccess,
			Duration:   30 * time.Minute,
		})
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshTkn, err := u.tokenSvc.
		Generate(repository.GenerateTokenParams{
			UserID:     *challenge.UserID,
			SessionID:  sessionID,
			StaffID:    staffID,
			CustomerID: customerID,
			Roles:      roleCodes,
			Type:       domain.TokenTypeRefresh,
			Duration:   7 * 24 * time.Hour,
		})
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	session := domain.Session{
		ID:        sessionID,
		UserID:    *challenge.UserID,
		UserAgent: input.UserAgent,
		IPAddress: input.IPAddress,
		ExpiresAt: now.Add(7 * 24 * time.Hour),
		CreatedAt: now,
	}

	refreshTknHashed := u.tokenHasher.Hash(refreshTkn.Token)
	refreshToken := domain.RefreshToken{
		ID:        uuid.New(),
		SessionID: session.ID,
		TokenHash: refreshTknHashed,
		ExpiresAt: now.Add(7 * 24 * time.Hour),
		CreatedAt: now,
	}

	err = u.transactor.WithinTransaction(
		ctx,
		func(exec transaction.Executor) error {
			if err := u.challengeRepo.
				Save(ctx, exec, *challenge); err != nil {
				return fmt.Errorf("failed to consume challenge: %w", err)
			}

			if err := u.accountRepo.
				ActivateByUserID(ctx, exec, *challenge.UserID); err != nil {
				return fmt.Errorf("failed to activate account: %w", err)
			}

			if err := u.sessionRepo.
				Save(ctx, exec, session); err != nil {
				return fmt.Errorf("failed to save session %w", err)
			}

			if err := u.refreshTokenRepo.
				Save(ctx, exec, refreshToken); err != nil {
				return fmt.Errorf("failed to save refresh token %w", err)
			}

			return nil
		},
	)
	if err != nil {
		return nil, err
	}

	result := VerifyAccountResult{
		AccessToken:  accessTkn,
		RefreshToken: refreshTkn,
	}

	if account != nil {
		u.auditLogger.Log(ctx, applogger.AuditEvent{
			Category:   "user_action",
			Action:     "verify_account",
			Resource:   "account",
			ResourceID: account.ID.String(),
			Outcome:    applogger.OutcomeSuccess,
			Metadata:   map[string]any{"challenge_id": input.ChallengeID.String()},
		})
	}

	return &result, nil
}
