package usecase

import (
	"fmt"
	"time"

	appErr "service-core/internal/common/errors"
	"service-core/internal/modules/authentication/domain"
	"service-core/internal/modules/authentication/repository"

	"github.com/google/uuid"
)

type LoginEmailUsecase struct {
	accountRepo      repository.AccountRepository
	pwHasher         repository.PasswordHasher
	tokenHasher      repository.TokenHasher
	tokenSvc         repository.TokenService
	sessionRepo      repository.SessionRepository
	refreshTokenRepo repository.RefreshTokenRepository
}

func NewLoginEmailUsecase(
	accountRepo repository.AccountRepository,
	pwHasher repository.PasswordHasher,
	tokenHasher repository.TokenHasher,
	tokenSvc repository.TokenService,
	sessionRepo repository.SessionRepository,
	refreshTokenRepo repository.RefreshTokenRepository,
) *LoginEmailUsecase {
	return &LoginEmailUsecase{
		accountRepo:      accountRepo,
		pwHasher:         pwHasher,
		tokenHasher:      tokenHasher,
		tokenSvc:         tokenSvc,
		sessionRepo:      sessionRepo,
		refreshTokenRepo: refreshTokenRepo,
	}
}

type LoginEmailParams struct {
	UserAgent *string
	IPAddress *string
	Email     string
	Password  string
}

type LoginEmailResult struct {
	AccessToken, RefreshToken repository.GeneratedToken
}

func (u *LoginEmailUsecase) Execute(input LoginEmailParams) (*LoginEmailResult, error) {
	existing, err := u.accountRepo.GetByEmail(input.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account: %w", err)
	}

	if existing == nil {
		return nil, appErr.NewUnauthorized(domain.ErrInvalidCredentials.Error())
	}
	if existing.Status != domain.AccountActive {
		return nil, appErr.NewForbidden(domain.ErrEmailNotVerified.Error())
	}

	if err := u.pwHasher.Compare(existing.Password, input.Password); err != nil {
		return nil, appErr.NewUnauthorized(domain.ErrInvalidCredentials.Error())
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

	accessTkn, err := u.tokenSvc.Generate(repository.GenerateTokenParams{
		UserID:    existing.UserID,
		SessionID: session.ID,
		Type:      domain.TokenTypeAccess,
		Duration:  30 * time.Minute,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshTkn, err := u.tokenSvc.Generate(repository.GenerateTokenParams{
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

	if err := u.sessionRepo.Save(session); err != nil {
		return nil, fmt.Errorf("failed to save session %w", err)
	}
	if err := u.refreshTokenRepo.Save(refreshTknDomain); err != nil {
		return nil, fmt.Errorf("failed to save refresh token %w", err)
	}

	result := LoginEmailResult{
		AccessToken:  accessTkn,
		RefreshToken: refreshTkn,
	}

	return &result, nil
}
