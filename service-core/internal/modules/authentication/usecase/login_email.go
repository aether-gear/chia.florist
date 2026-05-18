package usecase

import (
	"fmt"
	"time"

	appErr "service-core/internal/common/errors"
	authenDomain "service-core/internal/modules/authentication/domain"
	authenRepo "service-core/internal/modules/authentication/repository"

	"github.com/google/uuid"
)

type LoginEmailUsecase struct {
	accountRepo      authenRepo.AccountRepository
	pwHasher         authenRepo.PasswordHasher
	tokenHasher      authenRepo.TokenHasher
	tokenSvc         authenRepo.TokenService
	sessionRepo      authenRepo.SessionRepository
	refreshTokenRepo authenRepo.RefreshTokenRepository
}

func NewLoginEmailUsecase(
	accountRepo authenRepo.AccountRepository,
	pwHasher authenRepo.PasswordHasher,
	tokenHasher authenRepo.TokenHasher,
	tokenSvc authenRepo.TokenService,
	sessionRepo authenRepo.SessionRepository,
	refreshTokenRepo authenRepo.RefreshTokenRepository,
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
	AccessToken, RefreshToken authenRepo.GeneratedToken
}

func (u *LoginEmailUsecase) Execute(input LoginEmailParams) (*LoginEmailResult, error) {
	existing, err := u.accountRepo.GetByEmail(input.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account: %w", err)
	}
	if existing == nil {
		return nil, appErr.NewUnauthorized(authenDomain.ErrInvalidCredentials.Error())
	}

	if err := u.pwHasher.Compare(existing.Password, input.Password); err != nil {
		return nil, appErr.NewUnauthorized(authenDomain.ErrInvalidCredentials.Error())
	}

	sessionID := uuid.New()
	accessTkn, err := u.tokenSvc.Generate(authenRepo.GenerateTokenParams{
		UserID:    existing.UserID,
		SessionID: sessionID,
		Type:      authenDomain.TokenTypeAccess,
		Duration:  30 * time.Minute,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshTkn, err := u.tokenSvc.Generate(authenRepo.GenerateTokenParams{
		UserID:    existing.UserID,
		SessionID: sessionID,
		Type:      authenDomain.TokenTypeRefresh,
		Duration:  7 * 24 * time.Hour,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	now := time.Now()
	session := authenDomain.Session{
		ID:        sessionID,
		UserID:    existing.UserID,
		UserAgent: input.UserAgent,
		IPAddress: input.IPAddress,
		ExpiresAt: now.Add(7 * 24 * time.Hour),
		CreatedAt: now,
	}

	refreshTknHashed := u.tokenHasher.Hash(refreshTkn.Token)
	refreshToken := authenDomain.RefreshToken{
		ID:        uuid.New(),
		SessionID: session.ID,
		TokenHash: refreshTknHashed,
		ExpiresAt: now.Add(7 * 24 * time.Hour),
		CreatedAt: now,
	}

	if err := u.sessionRepo.Save(session); err != nil {
		return nil, fmt.Errorf("failed to save session %w", err)
	}
	if err := u.refreshTokenRepo.Save(refreshToken); err != nil {
		return nil, fmt.Errorf("failed to save refresh token %w", err)
	}

	result := LoginEmailResult{
		AccessToken:  accessTkn,
		RefreshToken: refreshTkn,
	}

	return &result, nil
}
