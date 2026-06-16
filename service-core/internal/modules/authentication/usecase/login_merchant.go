package usecase

import (
	"context"
	"fmt"
	"time"

	apperrors "service-core/internal/common/errors"
	"service-core/internal/modules/authentication/domain"
	"service-core/internal/modules/authentication/repository"
	authorzDomain "service-core/internal/modules/authorization/domain"
	authorRepo "service-core/internal/modules/authorization/repository"
	merchantRepo "service-core/internal/modules/merchant/repository"
	transaction "service-core/internal/shared/transaction"

	"github.com/google/uuid"
)

type LoginMerchantUsecase struct {
	executor         transaction.Executor
	transactor       transaction.Transactor
	accountRepo      repository.AccountRepository
	pwHasher         repository.PasswordHasher
	tokenHasher      repository.TokenHasher
	tokenSvc         repository.TokenService
	sessionRepo      repository.SessionRepository
	refreshTokenRepo repository.RefreshTokenRepository
	merchantRepo     merchantRepo.MerchantRepository
	membershipRepo   authorRepo.MerchantMembershipRepository
}

func NewLoginMerchantUsecase(
	executor transaction.Executor,
	transactor transaction.Transactor,
	accountRepo repository.AccountRepository,
	pwHasher repository.PasswordHasher,
	tokenHasher repository.TokenHasher,
	tokenSvc repository.TokenService,
	sessionRepo repository.SessionRepository,
	refreshTokenRepo repository.RefreshTokenRepository,
	merchantRepo merchantRepo.MerchantRepository,
	membershipRepo authorRepo.MerchantMembershipRepository,
) *LoginMerchantUsecase {
	return &LoginMerchantUsecase{
		executor:         executor,
		transactor:       transactor,
		accountRepo:      accountRepo,
		pwHasher:         pwHasher,
		tokenHasher:      tokenHasher,
		tokenSvc:         tokenSvc,
		sessionRepo:      sessionRepo,
		refreshTokenRepo: refreshTokenRepo,
		merchantRepo:     merchantRepo,
		membershipRepo:   membershipRepo,
	}
}

type LoginMerchantParams struct {
	UserAgent *string
	IPAddress *string
	Email     string
	Password  string
}

func (u *LoginMerchantUsecase) Execute(
	ctx context.Context,
	input LoginMerchantParams,
) (*LoginEmailResult, error) {
	existing, err := u.accountRepo.
		GetByEmail(ctx, u.executor, input.Email)

	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account: %w", err)
	}
	if existing == nil {
		return nil, apperrors.NewUnauthorized(domain.ErrInvalidCredentials.Error())
	}

	if existing.Type != domain.AccountTypeMerchant {
		return nil, apperrors.NewUnauthorized(domain.ErrInvalidCredentials.Error())
	}
	if existing.Status != domain.AccountActive {
		return nil, apperrors.NewForbidden(domain.ErrEmailNotVerified.Error())
	}

	if err := u.pwHasher.Compare(existing.Password, input.Password); err != nil {
		return nil, apperrors.NewUnauthorized(domain.ErrInvalidCredentials.Error())
	}

	memberMerchant, err := u.membershipRepo.
		GetByAccountID(ctx, u.executor, existing.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve membership: %w", err)
	}
	if memberMerchant == nil {
		return nil, apperrors.NewUnauthorized(domain.ErrInvalidCredentials.Error())
	}

	roles, err := u.membershipRepo.
		ListRolesByAccountIDAndMerchantID(
			ctx,
			u.executor,
			existing.ID,
			memberMerchant.MerchantID,
		)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve roles: %w", err)
	}
	if roles == nil {
		return nil, apperrors.NewForbidden("no role associated with this account")
	}

	roleCodes := make([]authorzDomain.RoleCode, len(roles))
	for i, r := range roles {
		roleCodes[i] = r.Code
	}

	now := time.Now()
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
			UserID:     existing.UserID,
			SessionID:  session.ID,
			MerchantID: &memberMerchant.MerchantID,
			Roles:      roleCodes,
			Type:       domain.TokenTypeAccess,
			Duration:   30 * time.Minute,
		})
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshTkn, err := u.tokenSvc.
		Generate(repository.GenerateTokenParams{
			UserID:     existing.UserID,
			SessionID:  session.ID,
			MerchantID: &memberMerchant.ID,
			Roles:      roleCodes,
			Type:       domain.TokenTypeRefresh,
			Duration:   7 * 24 * time.Hour,
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

	return &LoginEmailResult{
		AccessToken:  accessTkn,
		RefreshToken: refreshTkn,
	}, nil
}
