package usecase

import (
	"context"
	"fmt"

	apperrors "service-core/internal/common/errors"
	applogger "service-core/internal/common/logger"
	"service-core/internal/modules/authentication/domain"
	"service-core/internal/modules/authentication/infra/service"
	"service-core/internal/modules/authentication/repository"
	customerRepo "service-core/internal/modules/customer/repository"
	transaction "service-core/internal/shared/transaction"

	"github.com/google/uuid"
)

type LoginCustomerUsecase struct {
	executor      transaction.Executor
	transactor    transaction.Transactor
	accountRepo   repository.AccountRepository
	pwHasher      repository.PasswordHasher
	customerRepo  customerRepo.CustomerRepository
	sessionIssuer repository.SessionIssuerService
	auditLogger   applogger.AuditLogger
	sysLogger     applogger.Logger
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
	sessionIssuer := service.NewSessionIssuerService(
		transactor,
		tokenSvc,
		tokenHasher,
		sessionRepo,
		refreshTokenRepo,
		accountRepo,
	)

	return &LoginCustomerUsecase{
		executor:      executor,
		transactor:    transactor,
		accountRepo:   accountRepo,
		pwHasher:      pwHasher,
		customerRepo:  customerRepo,
		sessionIssuer: sessionIssuer,
		auditLogger:   auditLogger,
	}
}

func (u *LoginCustomerUsecase) SetSessionIssuer(sessionIssuer repository.SessionIssuerService) {
	u.sessionIssuer = sessionIssuer
}

func (u *LoginCustomerUsecase) SetSysLogger(sysLogger applogger.Logger) {
	u.sysLogger = sysLogger
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
) (result *LoginEmailResult, err error) {
	audit := &applogger.AuditScope{
		Category: "user_action",
		Action:   "login_customer",
		Resource: "session",
		Metadata: map[string]any{"email": input.Email},
	}
	defer applogger.TrackAudit(ctx, u.auditLogger, u.sysLogger, audit, &err)()

	existing, err := u.accountRepo.GetByEmail(ctx, u.executor, input.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account: %w", err)
	}
	if existing == nil {
		audit.SetReason("account not found")
		return nil, apperrors.NewUnauthorized(domain.ErrInvalidCredentials.Error())
	}
	if existing.Type != domain.AccountTypeCustomer {
		audit.SetReason("invalid account type")
		return nil, apperrors.NewUnauthorized(domain.ErrInvalidCredentials.Error())
	}
	if existing.Status != domain.AccountActive {
		audit.SetReason("account not active")
		return nil, apperrors.NewForbidden(domain.ErrEmailNotVerified.Error())
	}

	if err := u.pwHasher.Compare(existing.Password, input.Password); err != nil {
		audit.SetReason("invalid password")
		return nil, apperrors.NewUnauthorized(domain.ErrInvalidCredentials.Error())
	}

	var customerID *uuid.UUID
	cust, err := u.customerRepo.GetByUserID(ctx, u.executor, existing.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve customer: %w", err)
	}
	if cust != nil {
		customerID = &cust.ID
	}

	sessionRes, err := u.sessionIssuer.Issue(ctx, repository.IssueSessionParams{
		UserID:     existing.UserID,
		AccountID:  existing.ID,
		UserAgent:  input.UserAgent,
		IPAddress:  input.IPAddress,
		CustomerID: customerID,
	})
	if err != nil {
		return nil, err
	}

	audit.SetResourceID(sessionRes.SessionID.String())

	return &LoginEmailResult{
		AccessToken:  sessionRes.AccessToken,
		RefreshToken: sessionRes.RefreshToken,
	}, nil
}
