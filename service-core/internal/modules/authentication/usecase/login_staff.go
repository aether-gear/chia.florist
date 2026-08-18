package usecase

import (
	"context"
	"fmt"

	apperrors "service-core/internal/common/errors"
	applogger "service-core/internal/common/logger"
	"service-core/internal/modules/authentication/domain"
	"service-core/internal/modules/authentication/infra/service"
	"service-core/internal/modules/authentication/repository"
	authorzDomain "service-core/internal/modules/authorization/domain"
	authorRepo "service-core/internal/modules/authorization/repository"
	staffRepo "service-core/internal/modules/staff/repository"
	transaction "service-core/internal/shared/transaction"
)

type LoginStaffUsecase struct {
	executor       transaction.Executor
	transactor     transaction.Transactor
	accountRepo    repository.AccountRepository
	pwHasher       repository.PasswordHasher
	staffRepo      staffRepo.StaffRepository
	membershipRepo authorRepo.StaffMembershipRepository
	sessionIssuer  repository.SessionIssuerService
	auditLogger    applogger.AuditLogger
	sysLogger      applogger.Logger
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
	sessionIssuer := service.NewSessionIssuerService(
		transactor,
		tokenSvc,
		tokenHasher,
		sessionRepo,
		refreshTokenRepo,
		accountRepo,
	)

	return &LoginStaffUsecase{
		executor:       executor,
		transactor:     transactor,
		accountRepo:    accountRepo,
		pwHasher:       pwHasher,
		staffRepo:      staffRepo,
		membershipRepo: membershipRepo,
		sessionIssuer:  sessionIssuer,
		auditLogger:    auditLogger,
	}
}

func (u *LoginStaffUsecase) SetSessionIssuer(sessionIssuer repository.SessionIssuerService) {
	u.sessionIssuer = sessionIssuer
}

func (u *LoginStaffUsecase) SetSysLogger(sysLogger applogger.Logger) {
	u.sysLogger = sysLogger
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
) (result *LoginEmailResult, err error) {
	audit := &applogger.AuditScope{
		Category: "user_action",
		Action:   "login_staff",
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
	if existing.Type != domain.AccountTypeStaff {
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

	memberStaff, err := u.membershipRepo.GetByAccountID(ctx, u.executor, existing.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve membership: %w", err)
	}
	if memberStaff == nil {
		audit.SetReason("staff membership not found")
		return nil, apperrors.NewUnauthorized(domain.ErrInvalidCredentials.Error())
	}

	roles, err := u.membershipRepo.ListRolesByAccountIDAndStaffID(ctx, u.executor,
		existing.ID,
		memberStaff.StaffID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve roles: %w", err)
	}
	if roles == nil {
		audit.SetReason("no roles associated")
		return nil, apperrors.NewForbidden("no role associated with this account")
	}

	roleCodes := make([]authorzDomain.RoleCode, len(roles))
	for i, r := range roles {
		roleCodes[i] = r.Code
	}

	sessionRes, err := u.sessionIssuer.Issue(ctx, repository.IssueSessionParams{
		UserID:    existing.UserID,
		AccountID: existing.ID,
		UserAgent: input.UserAgent,
		IPAddress: input.IPAddress,
		StaffID:   &memberStaff.StaffID,
		Roles:     roleCodes,
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
