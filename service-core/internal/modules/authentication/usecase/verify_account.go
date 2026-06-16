package usecase

import (
	"context"
	"fmt"
	"strconv"
	"time"

	apperrors "service-core/internal/common/errors"
	"service-core/internal/modules/authentication/domain"
	"service-core/internal/modules/authentication/repository"
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
	challengeRepo    repository.VerificationChallengeRepository
	tokenSvc         repository.TokenService
	sessionRepo      repository.SessionRepository
	refreshTokenRepo repository.RefreshTokenRepository
}

func NewVerifyAccountUsecase(
	executor transaction.Executor,
	transactor transaction.Transactor,
	accountRepo repository.AccountRepository,
	pwHasher repository.PasswordHasher,
	tokenHasher repository.TokenHasher,
	userRepo userRepo.UserRepository,
	challengeRepo repository.VerificationChallengeRepository,
	tokenSvc repository.TokenService,
	sessionRepo repository.SessionRepository,
	refreshTokenRepo repository.RefreshTokenRepository,
) *VerifyAccountUsecase {
	return &VerifyAccountUsecase{
		executor:         executor,
		transactor:       transactor,
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
		return nil, apperrors.NewNotFound(domain.ErrNotFoundChallenge.Error())
	}

	if challenge.ConsumedAt != nil {
		return nil, apperrors.NewConflict(domain.ErrConsumedChallenge.Error())
	}
	if challenge.VerifiedAt != nil {
		return nil, apperrors.NewConflict(domain.ErrVerifiedChallenge.Error())
	}
	if challenge.ExpiresAt.Before(now) {
		return nil, apperrors.NewConflict(domain.ErrExpiredChallenge.Error())
	}
	if challenge.AttemptCount >= 5 {
		return nil, apperrors.NewConflict(domain.ErrMaxAttemptReached.Error())
	}

	if err := u.pwHasher.
		Compare(
			challenge.CodeHash,
			strconv.Itoa(input.OTP),
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

		return nil, apperrors.NewUnauthorized(domain.ErrInvalidOTP.Error())
	}

	challenge.VerifiedAt = &now
	challenge.ConsumedAt = &now

	sessionID := uuid.New()
	accessTkn, err := u.tokenSvc.
		Generate(repository.GenerateTokenParams{
			UserID:    *challenge.UserID,
			SessionID: sessionID,
			Type:      domain.TokenTypeAccess,
			Duration:  30 * time.Minute,
		})
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshTkn, err := u.tokenSvc.
		Generate(repository.GenerateTokenParams{
			UserID:    *challenge.UserID,
			SessionID: sessionID,
			Type:      domain.TokenTypeRefresh,
			Duration:  7 * 24 * time.Hour,
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

	return &result, nil
}
