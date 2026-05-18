package domain

import "errors"

var (
	ErrInvalidCredentials   = errors.New("invalid credentials")
	ErrNotFoundAccount      = errors.New("account not found")
	ErrAccountAlreadyExists = errors.New("account already exists")
)

var (
	ErrNotFoundChallenge = errors.New("challenge not found")
	ErrExpiredChallenge  = errors.New("challenge expired")
	ErrConsumedChallenge = errors.New("challenge already consumed")
	ErrVerifiedChallenge = errors.New("challenge already verified")
	ErrInvalidOTP        = errors.New("invalid otp")
	ErrMaxAttemptReached = errors.New("max otp attempts reached")
)

var (
	ErrNotFoundSession = errors.New("not found session")
)
