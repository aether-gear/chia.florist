package usecase

import (
	"context"
	"fmt"
	"time"

	apperrors "service-core/internal/common/errors"
	"service-core/internal/modules/authentication/domain"
	"service-core/internal/modules/authentication/repository"
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
	existing, err := u.accountRepo.
		GetByEmail(ctx, u.executor, input.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account: %w", err)
	}

	if existing == nil {
		return nil, apperrors.NewUnauthorized(domain.ErrInvalidCredentials.Error())
	}
	if existing.Status != domain.AccountActive {
		return nil, apperrors.NewForbidden(domain.ErrEmailNotVerified.Error())
	}

	if err := u.pwHasher.Compare(existing.Password, input.Password); err != nil {
		return nil, apperrors.NewUnauthorized(domain.ErrInvalidCredentials.Error())
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
			UserID:    existing.UserID,
			SessionID: session.ID,
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
			Type:      domain.TokenTypeRefresh,
			Duration:  7 * 24 * time.Hour,
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

	err = u.transactor.WithinTransaction(
		ctx,
		func(exec transaction.Executor) error {
			if err := u.sessionRepo.
				Save(ctx, exec, session); err != nil {
				return fmt.Errorf("failed to save session %w", err)
			}
			if err := u.refreshTokenRepo.
				Save(ctx, exec, refreshTknDomain); err != nil {
				return fmt.Errorf("failed to save refresh token %w", err)
			}

			return nil
		},
	)
	if err != nil {
		return nil, err
	}

	result := LoginEmailResult{
		AccessToken:  accessTkn,
		RefreshToken: refreshTkn,
	}

	return &result, nil
}
