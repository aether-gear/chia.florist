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
	accountRepo   authRepo.AccountRepository
	hasher        authDomain.PasswordHasher
	userRepo      userRepo.UserRepository
	challengeRepo authRepo.VerificationChallengeRepository
	tokenSvc      authDomain.TokenService
}

func NewVerifyAccountUsecase(
	accountRepo authRepo.AccountRepository,
	hasher authDomain.PasswordHasher,
	userRepo userRepo.UserRepository,
	challengeRepo authRepo.VerificationChallengeRepository,
	tokenSvc authDomain.TokenService,
) *VerifyAccountUsecase {
	return &VerifyAccountUsecase{
		accountRepo:   accountRepo,
		hasher:        hasher,
		userRepo:      userRepo,
		challengeRepo: challengeRepo,
		tokenSvc:      tokenSvc,
	}
}

type VerifyAccountParams struct {
	ChallengeID uuid.UUID
	OTP         int
}

func (u *VerifyAccountUsecase) Execute(
	input VerifyAccountParams,
) (*string, time.Time, error) {
	challenge, err := u.challengeRepo.GetByID(input.ChallengeID)
	if err != nil {
		return nil, time.Time{}, fmt.Errorf(
			"failed to get challenge: %w",
			err,
		)
	}
	if challenge == nil {
		return nil, time.Time{}, appErr.NewNotFound(
			authDomain.ErrNotFoundChallenge.Error(),
		)
	}

	now := time.Now()
	if challenge.ConsumedAt != nil {
		return nil, time.Time{}, appErr.NewConflict(
			authDomain.ErrConsumedChallenge.Error(),
		)
	}
	if challenge.VerifiedAt != nil {
		return nil, time.Time{}, appErr.NewConflict(
			authDomain.ErrVerifiedChallenge.Error(),
		)
	}
	if challenge.ExpiresAt.Before(now) {
		return nil, time.Time{}, appErr.NewConflict(
			authDomain.ErrExpiredChallenge.Error(),
		)
	}
	if challenge.AttemptCount >= 5 {
		return nil, time.Time{}, appErr.NewConflict(
			authDomain.ErrMaxAttemptReached.Error(),
		)
	}

	if err := u.hasher.Compare(challenge.CodeHash, strconv.Itoa(input.OTP)); err != nil {
		challenge.AttemptCount++

		if err := u.challengeRepo.Save(*challenge); err != nil {
			return nil, time.Time{}, fmt.Errorf(
				"failed to update challenge attempts: %w",
				err,
			)
		}

		return nil, time.Time{}, appErr.NewUnauthorized(
			authDomain.ErrInvalidOTP.Error(),
		)
	}

	challenge.VerifiedAt = &now
	challenge.ConsumedAt = &now
	if err := u.challengeRepo.Save(*challenge); err != nil {
		return nil, time.Time{}, fmt.Errorf(
			"failed to consume challenge: %w",
			err,
		)
	}

	if err := u.accountRepo.ActivateByUserID(*challenge.UserID); err != nil {
		return nil, time.Time{}, fmt.Errorf(
			"failed to activate account: %w",
			err,
		)
	}

	token, expiry, err := u.tokenSvc.Generate(*challenge.UserID)
	if err != nil {
		return nil, time.Time{}, fmt.Errorf("failed to generate access token: %w", err)
	}

	return &token, expiry, nil
}
