package usecase

import (
	"context"
	"fmt"
	"time"

	appclock "service-core/internal/common/clock"
	apperrors "service-core/internal/common/errors"
	applogger "service-core/internal/common/logger"
	"service-core/internal/modules/authentication/domain"
	"service-core/internal/modules/authentication/repository"
	customerRepo "service-core/internal/modules/customer/repository"
	transaction "service-core/internal/shared/transaction"

	"github.com/google/uuid"
)

type LoginCustomerUsecase struct {
	executor         transaction.Executor
	transactor       transaction.Transactor
	accountRepo      repository.AccountRepository
	pwHasher         repository.PasswordHasher
	tokenHasher      repository.TokenHasher
	tokenSvc         repository.TokenService
	sessionRepo      repository.SessionRepository
	refreshTokenRepo repository.RefreshTokenRepository
	customerRepo     customerRepo.CustomerRepository
	auditLogger      applogger.AuditLogger
}

func NewLoginCustomerUsecase(
	executor transaction.Executor,
	transactor transaction.Transactor,
	accountRepo repository.AccountRepository,
	pwHasher repository.PasswordHasher,
	tokenHasher repository.TokenHasher,
	tokenSvc repository.TokenService,
	sessionRepo repository.SessionRepository,
	refreshTokenRepo repository.RefreshTokenRepository,
	customerRepo customerRepo.CustomerRepository,
	auditLogger applogger.AuditLogger,
) *LoginCustomerUsecase {
	return &LoginCustomerUsecase{
		executor:         executor,
		transactor:       transactor,
		accountRepo:      accountRepo,
		pwHasher:         pwHasher,
		tokenHasher:      tokenHasher,
		tokenSvc:         tokenSvc,
		sessionRepo:      sessionRepo,
		refreshTokenRepo: refreshTokenRepo,
		customerRepo:     customerRepo,
		auditLogger:      auditLogger,
	}
}

type LoginCustomerParams struct {
	UserAgent *string
	IPAddress *string
	Email     string
	Password  string
}

type LoginEmailResult struct {
	AccessToken, RefreshToken repository.GeneratedToken
}

func (u *LoginCustomerUsecase) Execute(
	ctx context.Context,
	input LoginCustomerParams,
) (*LoginEmailResult, error) {
	existing, err := u.accountRepo.GetByEmail(ctx, u.executor, input.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account: %w", err)
	}
	if existing == nil {
		u.auditLogger.Log(ctx, applogger.AuditEvent{
			Category: "user_action",
			Action:   "login_customer",
			Resource: "session",
			Outcome:  applogger.OutcomeFailure,
			Metadata: map[string]any{"email": input.Email, "reason": "account not found"},
		})
		return nil, apperrors.NewUnauthorized(domain.ErrInvalidCredentials.Error())
	}

	if existing.Type != domain.AccountTypeCustomer {
		u.auditLogger.Log(ctx, applogger.AuditEvent{
			Category: "user_action",
			Action:   "login_customer",
			Resource: "session",
			Outcome:  applogger.OutcomeFailure,
			Metadata: map[string]any{"email": input.Email, "reason": "invalid account type"},
		})
		return nil, apperrors.NewUnauthorized(domain.ErrInvalidCredentials.Error())
	}

	if existing.Status != domain.AccountActive {
		u.auditLogger.Log(ctx, applogger.AuditEvent{
			Category: "user_action",
			Action:   "login_customer",
			Resource: "session",
			Outcome:  applogger.OutcomeFailure,
			Metadata: map[string]any{"email": input.Email, "reason": "account not active"},
		})
		return nil, apperrors.NewForbidden(domain.ErrEmailNotVerified.Error())
	}

	if err := u.pwHasher.Compare(existing.Password, input.Password); err != nil {
		u.auditLogger.Log(ctx, applogger.AuditEvent{
			Category: "user_action",
			Action:   "login_customer",
			Resource: "session",
			Outcome:  applogger.OutcomeFailure,
			Metadata: map[string]any{"email": input.Email, "reason": "invalid password"},
		})
		return nil, apperrors.NewUnauthorized(domain.ErrInvalidCredentials.Error())
	}

	now := appclock.Now()
	session := domain.Session{
		ID:        uuid.New(),
		UserID:    existing.UserID,
		UserAgent: input.UserAgent,
		IPAddress: input.IPAddress,
		ExpiresAt: now.Add(7 * 24 * time.Hour),
		CreatedAt: now,
	}

	var customerID *uuid.UUID
	cust, err := u.customerRepo.GetByUserID(ctx, u.executor, existing.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve customer: %w", err)
	}
	if cust != nil {
		customerID = &cust.ID
	}

	accessTkn, err := u.tokenSvc.Generate(repository.GenerateTokenParams{
		UserID:     existing.UserID,
		SessionID:  session.ID,
		CustomerID: customerID,
		Type:       domain.TokenTypeAccess,
		Duration:   30 * time.Minute,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshTkn, err := u.tokenSvc.Generate(repository.GenerateTokenParams{
		UserID:     existing.UserID,
		SessionID:  session.ID,
		CustomerID: customerID,
		Type:       domain.TokenTypeRefresh,
		Duration:   7 * 24 * time.Hour,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshTknHashed := u.tokenHasher.Hash(refreshTkn.Token)
	refreshTknDomain := domain.RefreshToken{
		ID:        uuid.New(),
		SessionID: session.ID,
		TokenHash: refreshTknHashed,
		ExpiresAt: now.Add(7 * 24 * time.Hour),
		CreatedAt: now,
	}

	err = u.transactor.WithinTransaction(ctx, func(exec transaction.Executor) error {
		if err := u.sessionRepo.Save(ctx, exec, session); err != nil {
			return fmt.Errorf("failed to save session %w", err)
		}

		if err := u.refreshTokenRepo.Save(ctx, exec, refreshTknDomain); err != nil {
			return fmt.Errorf("failed to save refresh token %w", err)
		}

		if err := u.accountRepo.UpdateLastLoginAt(ctx, exec,
			existing.ID,
			now,
		); err != nil {
			return fmt.Errorf("failed to update last login at: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	result := LoginEmailResult{
		AccessToken:  accessTkn,
		RefreshToken: refreshTkn,
	}

	u.auditLogger.Log(ctx, applogger.AuditEvent{
		Category:   "user_action",
		Action:     "login_customer",
		Resource:   "session",
		ResourceID: session.ID.String(),
		Outcome:    applogger.OutcomeSuccess,
		Metadata:   map[string]any{"email": input.Email},
	})

	return &result, nil
}
