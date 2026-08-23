package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// OrderItemCustomDesign stores the full design snapshot and extracted
// metadata produced by the custom board designer.
type OrderItemCustomDesign struct {
	ID              uuid.UUID
	OrderItemID     uuid.UUID
	Version         string
	PhysicalSizeID  string
	PreviewURL      *string
	HeaderTextUpper *string
	BodyTextUpper   *string
	HeaderTextLower *string
	BodyTextLower   *string
	DesignSnapshot  json.RawMessage
	CreatedAt       time.Time
}
