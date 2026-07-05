package domain

import (
	"time"
)

// WAFRule is a detection rule matched
// against incoming request payloads.
type WAFRule struct {
	ID          string
	Description string
	Pattern     string
	Tags        []string
	Impact      string
	Enabled     bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
