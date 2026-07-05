package http

import (
	"time"

	"github.com/google/uuid"
)

type auditLogResponse struct {
	ID         uuid.UUID      `json:"id"`
	Category   string         `json:"category"`
	Action     string         `json:"action"`
	Resource   string         `json:"resource"`
	ResourceID string         `json:"resource_id,omitempty"`
	ActorID    string         `json:"actor_id,omitempty"`
	Outcome    string         `json:"outcome"`
	RequestID  string         `json:"request_id,omitempty"`
	ClientIP   string         `json:"client_ip,omitempty"`
	Metadata   map[string]any `json:"metadata,omitempty"`
	CreatedAt  time.Time      `json:"created_at"`
}
