package usecase

import (
	"fmt"
	"strconv"
	"time"

	appErr "service-core/internal/common/errors"
	authDomain "service-core/internal/modules/authentication/domain"
	authRepo "service-core/internal/modules/authentication/repository"
	userRepo "service-core/internal/modules/user/repository"

	"github.com/google/uuid"
)

type VerifyAccountUsecase struct {
	accountRepo      authRepo.AccountRepository
	pwHasher         authRepo.PasswordHasher
	tokenHasher      authRepo.TokenHasher
	userRepo         userRepo.UserRepository
	challengeRepo    authRepo.VerificationChallengeRepository
	tokenSvc         authRepo.TokenService
	sessionRepo      authRepo.SessionRepository
	refreshTokenRepo authRepo.RefreshTokenRepository
}

func NewVerifyAccountUsecase(
	accountRepo authRepo.AccountRepository,
	pwHasher authRepo.PasswordHasher,
	tokenHasher authRepo.TokenHasher,
	userRepo userRepo.UserRepository,
	challengeRepo authRepo.VerificationChallengeRepository,
	tokenSvc authRepo.TokenService,
	sessionRepo authRepo.SessionRepository,
	refreshTokenRepo authRepo.RefreshTokenRepository,
) *VerifyAccountUsecase {
	return &VerifyAccountUsecase{
		accountRepo:      accountRepo,
		pwHasher:         pwHasher,
		tokenHasher:      tokenHasher,
		userRepo:         userRepo,
		challengeRepo:    challengeRepo,
		tokenSvc:         tokenSvc,
		sessionRepo:      sessionRepo,
		refreshTokenRepo: refreshTokenRepo,
	}
}

type VerifyAccountParams struct {
	UserAgent   *string
	IPAddress   *string
	ChallengeID uuid.UUID
	OTP         int
}

type VerifyAccountResult struct {
	AccessToken, RefreshToken authRepo.GeneratedToken
}

func (u *VerifyAccountUsecase) Execute(
	input VerifyAccountParams,
) (*VerifyAccountResult, error) {
	now := time.Now()

	challenge, err := u.challengeRepo.GetByID(input.ChallengeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get challenge: %w", err)
	}
	if challenge == nil {
		return nil, appErr.NewNotFound(authDomain.ErrNotFoundChallenge.Error())
	}

	if challenge.ConsumedAt != nil {
		return nil, appErr.NewConflict(authDomain.ErrConsumedChallenge.Error())
	}
	if challenge.VerifiedAt != nil {
		return nil, appErr.NewConflict(authDomain.ErrVerifiedChallenge.Error())
	}
	if challenge.ExpiresAt.Before(now) {
		return nil, appErr.NewConflict(authDomain.ErrExpiredChallenge.Error())
	}
	if challenge.AttemptCount >= 5 {
		return nil, appErr.NewConflict(authDomain.ErrMaxAttemptReached.Error())
	}

	if err := u.pwHasher.Compare(challenge.CodeHash, strconv.Itoa(input.OTP)); err != nil {
		challenge.AttemptCount++

		if err := u.challengeRepo.Save(*challenge); err != nil {
			return nil, fmt.Errorf(
				"failed to update challenge attempts: %w",
				err,
			)
		}

		return nil, appErr.NewUnauthorized(
			authDomain.ErrInvalidOTP.Error(),
		)
	}

	challenge.VerifiedAt = &now
	challenge.ConsumedAt = &now
	if err := u.challengeRepo.Save(*challenge); err != nil {
		return nil, fmt.Errorf("failed to consume challenge: %w", err)
	}

	if err := u.accountRepo.ActivateByUserID(*challenge.UserID); err != nil {
		return nil, fmt.Errorf("failed to activate account: %w", err)
	}

	sessionID := uuid.New()
	accessTkn, err := u.tokenSvc.Generate(authRepo.GenerateTokenParams{
		UserID:    *challenge.UserID,
		SessionID: sessionID,
		Type:      authDomain.TokenTypeAccess,
		Duration:  30 * time.Minute,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshTkn, err := u.tokenSvc.Generate(authRepo.GenerateTokenParams{
		UserID:    *challenge.UserID,
		SessionID: sessionID,
		Type:      authDomain.TokenTypeRefresh,
		Duration:  7 * 24 * time.Hour,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	session := authDomain.Session{
		ID:        sessionID,
		UserID:    *challenge.UserID,
		UserAgent: input.UserAgent,
		IPAddress: input.IPAddress,
		ExpiresAt: now.Add(7 * 24 * time.Hour),
		CreatedAt: now,
	}

	refreshTknHashed := u.tokenHasher.Hash(refreshTkn.Token)
	refreshToken := authDomain.RefreshToken{
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

	result := VerifyAccountResult{
		AccessToken:  accessTkn,
		RefreshToken: refreshTkn,
	}

	return &result, nil
}
