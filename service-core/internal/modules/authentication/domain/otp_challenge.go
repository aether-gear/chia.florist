package domain

import (
	"time"

	"github.com/google/uuid"
)

type OTPType string

const (
	OTPTypeNumeric   OTPType = "numeric"
	OTPTypeMagicLink OTPType = "magic_link"
)

type OTPChannel string

const (
	OTPChannelEmail OTPChannel = "email"
	OTPChannelSMS   OTPChannel = "sms"
)

type OTPPurpose string

const (
	OTPPurposeRegister          OTPPurpose = "register"
	OTPPurposeLogin             OTPPurpose = "login"
	OTPPurposePasswordReset     OTPPurpose = "password_reset"
	OTPPurposeEmailVerification OTPPurpose = "email_verification"
)

type OTPChallenge struct {
	ID uuid.UUID

	UserID *uuid.UUID

	Type    OTPType
	Channel OTPChannel
	Purpose OTPPurpose

	Target string

	CodeHash string

	ExpiresAt  time.Time
	VerifiedAt *time.Time
	ConsumedAt *time.Time

	AttemptCount int

	CreatedAt time.Time
}
