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
	authorzDomain "service-core/internal/modules/authorization/domain"
	authorRepo "service-core/internal/modules/authorization/repository"
	staffRepo "service-core/internal/modules/staff/repository"
	transaction "service-core/internal/shared/transaction"

	"github.com/google/uuid"
)

type LoginStaffUsecase struct {
	executor         transaction.Executor
	transactor       transaction.Transactor
	accountRepo      repository.AccountRepository
	pwHasher         repository.PasswordHasher
	tokenHasher      repository.TokenHasher
	tokenSvc         repository.TokenService
	sessionRepo      repository.SessionRepository
	refreshTokenRepo repository.RefreshTokenRepository
	staffRepo        staffRepo.StaffRepository
	membershipRepo   authorRepo.StaffMembershipRepository
	auditLogger      applogger.AuditLogger
}

func NewLoginStaffUsecase(
	executor transaction.Executor,
	transactor transaction.Transactor,
	accountRepo repository.AccountRepository,
	pwHasher repository.PasswordHasher,
	tokenHasher repository.TokenHasher,
	tokenSvc repository.TokenService,
	sessionRepo repository.SessionRepository,
	refreshTokenRepo repository.RefreshTokenRepository,
	staffRepo staffRepo.StaffRepository,
	membershipRepo authorRepo.StaffMembershipRepository,
	auditLogger applogger.AuditLogger,
) *LoginStaffUsecase {
	return &LoginStaffUsecase{
		executor:         executor,
		transactor:       transactor,
		accountRepo:      accountRepo,
		pwHasher:         pwHasher,
		tokenHasher:      tokenHasher,
		tokenSvc:         tokenSvc,
		sessionRepo:      sessionRepo,
		refreshTokenRepo: refreshTokenRepo,
		staffRepo:        staffRepo,
		membershipRepo:   membershipRepo,
		auditLogger:      auditLogger,
	}
}

type LoginStaffParams struct {
	UserAgent *string
	IPAddress *string
	Email     string
	Password  string
}

func (u *LoginStaffUsecase) Execute(
	ctx context.Context,
	input LoginStaffParams,
) (*LoginEmailResult, error) {
	existing, err := u.accountRepo.
		GetByEmail(ctx, u.executor, input.Email)

	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account: %w", err)
	}
	if existing == nil {
		u.auditLogger.Log(ctx, applogger.AuditEvent{
			Category: "user_action",
			Action:   "login_staff",
			Resource: "session",
			Outcome:  applogger.OutcomeFailure,
			Metadata: map[string]any{"email": input.Email, "reason": "account not found"},
		})
		return nil, apperrors.NewUnauthorized(domain.ErrInvalidCredentials.Error())
	}

	if existing.Type != domain.AccountTypeStaff {
		u.auditLogger.Log(ctx, applogger.AuditEvent{
			Category: "user_action",
			Action:   "login_staff",
			Resource: "session",
			Outcome:  applogger.OutcomeFailure,
			Metadata: map[string]any{"email": input.Email, "reason": "invalid account type"},
		})
		return nil, apperrors.NewUnauthorized(domain.ErrInvalidCredentials.Error())
	}
	if existing.Status != domain.AccountActive {
		u.auditLogger.Log(ctx, applogger.AuditEvent{
			Category: "user_action",
			Action:   "login_staff",
			Resource: "session",
			Outcome:  applogger.OutcomeFailure,
			Metadata: map[string]any{"email": input.Email, "reason": "account not active"},
		})
		return nil, apperrors.NewForbidden(domain.ErrEmailNotVerified.Error())
	}

	if err := u.pwHasher.Compare(existing.Password, input.Password); err != nil {
		u.auditLogger.Log(ctx, applogger.AuditEvent{
			Category: "user_action",
			Action:   "login_staff",
			Resource: "session",
			Outcome:  applogger.OutcomeFailure,
			Metadata: map[string]any{"email": input.Email, "reason": "invalid password"},
		})
		return nil, apperrors.NewUnauthorized(domain.ErrInvalidCredentials.Error())
	}

	memberStaff, err := u.membershipRepo.
		GetByAccountID(ctx, u.executor, existing.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve membership: %w", err)
	}
	if memberStaff == nil {
		u.auditLogger.Log(ctx, applogger.AuditEvent{
			Category: "user_action",
			Action:   "login_staff",
			Resource: "session",
			Outcome:  applogger.OutcomeFailure,
			Metadata: map[string]any{"email": input.Email, "reason": "staff membership not found"},
		})
		return nil, apperrors.NewUnauthorized(domain.ErrInvalidCredentials.Error())
	}

	roles, err := u.membershipRepo.
		ListRolesByAccountIDAndStaffID(
			ctx,
			u.executor,
			existing.ID,
			memberStaff.StaffID,
		)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve roles: %w", err)
	}
	if roles == nil {
		u.auditLogger.Log(ctx, applogger.AuditEvent{
			Category: "user_action",
			Action:   "login_staff",
			Resource: "session",
			Outcome:  applogger.OutcomeFailure,
			Metadata: map[string]any{"email": input.Email, "reason": "no roles associated"},
		})
		return nil, apperrors.NewForbidden("no role associated with this account")
	}

	roleCodes := make([]authorzDomain.RoleCode, len(roles))
	for i, r := range roles {
		roleCodes[i] = r.Code
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

	accessTkn, err := u.tokenSvc.
		Generate(repository.GenerateTokenParams{
			UserID:    existing.UserID,
			SessionID: session.ID,
			StaffID:   &memberStaff.StaffID,
			Roles:     roleCodes,
			Type:      domain.TokenTypeAccess,
			Duration:  30 * time.Minute,
		})
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshTkn, err := u.tokenSvc.
		Generate(repository.GenerateTokenParams{
			UserID:    existing.UserID,
			SessionID: session.ID,
			StaffID:   &memberStaff.StaffID,
			Roles:     roleCodes,
			Type:      domain.TokenTypeRefresh,
			Duration:  7 * 24 * time.Hour,
		})
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	refreshTknHashed := u.tokenHasher.Hash(refreshTkn.Token)
	refreshTknDomain := domain.RefreshToken{
		ID:        uuid.New(),
		SessionID: session.ID,
		TokenHash: refreshTknHashed,
		ExpiresAt: now.Add(7 * 24 * time.Hour),
		CreatedAt: now,
	}

	err = u.transactor.WithinTransaction(
		ctx,
		func(exec transaction.Executor) error {
			if err := u.sessionRepo.
				Save(ctx, exec, session); err != nil {
				return fmt.Errorf("failed to save session: %w", err)
			}

			if err := u.refreshTokenRepo.
				Save(ctx, exec, refreshTknDomain); err != nil {
				return fmt.Errorf("failed to save refresh token: %w", err)
			}

			return nil
		},
	)
	if err != nil {
		return nil, err
	}

	u.auditLogger.Log(ctx, applogger.AuditEvent{
		Category:   "user_action",
		Action:     "login_staff",
		Resource:   "session",
		ResourceID: session.ID.String(),
		Outcome:    applogger.OutcomeSuccess,
		Metadata:   map[string]any{"email": input.Email},
	})

	return &LoginEmailResult{
		AccessToken:  accessTkn,
		RefreshToken: refreshTkn,
	}, nil
}
