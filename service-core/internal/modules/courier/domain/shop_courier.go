package domain

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type CourierVerificationStatus string

const (
	CourierVerificationUnconfigured CourierVerificationStatus = "unconfigured"
	CourierVerificationPending      CourierVerificationStatus = "pending"
	CourierVerificationVerified     CourierVerificationStatus = "verified"
	CourierVerificationRejected     CourierVerificationStatus = "rejected"
)

func (s CourierVerificationStatus) IsValid() bool {
	switch s {
	case CourierVerificationUnconfigured,
		CourierVerificationPending,
		CourierVerificationVerified,
		CourierVerificationRejected:
		return true
	default:
		return false
	}
}

type ShopCourier struct {
	ShopID             uuid.UUID
	Code               string
	BranchName         string
	Name               *string
	LocationAddress    *string
	Active             bool
	VerificationStatus CourierVerificationStatus
	VerifiedAt         *time.Time
	VerifiedBy         *uuid.UUID
	RejectionReason    *string
	CreatedAt          *time.Time
	UpdatedAt          *time.Time
}

func (c ShopCourier) IsOperable() bool {
	return c.Active && c.VerificationStatus == CourierVerificationVerified
}

func (c ShopCourier) Validate() error {
	if c.Active {
		if c.Name == nil || strings.TrimSpace(*c.Name) == "" {
			return ErrCourierNameRequired
		}
		if c.LocationAddress == nil || strings.TrimSpace(*c.LocationAddress) == "" {
			return ErrCourierLocationRequired
		}
	}
	return nil
}
