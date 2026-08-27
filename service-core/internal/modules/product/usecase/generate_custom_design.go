package usecase

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	apperrors "service-core/internal/common/errors"
	"service-core/internal/infra/genai"

	"github.com/google/uuid"
)

type GenerateCustomDesignInput struct {
	CustomerID       uuid.UUID
	Prompt           string
	Occasion         string
	PreferredPalette string
	Recipient        string
	Sender           string
	PhysicalSizeID   string
}

type GenerateCustomDesignUsecase struct {
	aiProv             genai.Provider
	maxRequestsPerHour int
	mu                 sync.Mutex
	userRequestHistory map[uuid.UUID][]time.Time
}

func NewGenerateCustomDesignUsecase(
	aiProv genai.Provider,
	maxRequestsPerHour int,
) *GenerateCustomDesignUsecase {
	if maxRequestsPerHour <= 0 {
		maxRequestsPerHour = 10
	}

	return &GenerateCustomDesignUsecase{
		aiProv:             aiProv,
		maxRequestsPerHour: maxRequestsPerHour,
		userRequestHistory: make(map[uuid.UUID][]time.Time),
	}
}

func (u *GenerateCustomDesignUsecase) Execute(
	ctx context.Context,
	input GenerateCustomDesignInput,
) (*genai.CustomDesignPayloadV3, error) {
	if input.CustomerID == uuid.Nil {
		return nil, apperrors.NewUnauthorized("customer authentication required")
	}

	trimmedPrompt := strings.TrimSpace(input.Prompt)
	if len(trimmedPrompt) < 3 {
		return nil, apperrors.NewInvalidInput("prompt must be at least 3 characters long")
	}
	if len(trimmedPrompt) > 1000 {
		return nil, apperrors.NewInvalidInput("prompt exceeds maximum length of 1000 characters")
	}

	// Check customer sliding-window rate limit (per hour)
	if err := u.checkAndRecordRateLimit(input.CustomerID); err != nil {
		return nil, err
	}

	designCtx := genai.DesignContext{
		Occasion:         strings.TrimSpace(input.Occasion),
		PreferredPalette: strings.TrimSpace(input.PreferredPalette),
		Recipient:        strings.TrimSpace(input.Recipient),
		Sender:           strings.TrimSpace(input.Sender),
		PhysicalSizeID:   strings.TrimSpace(input.PhysicalSizeID),
		Locale:           "id-ID",
	}

	result, err := u.aiProv.GenerateBoardDesign(ctx, trimmedPrompt, designCtx)
	if err != nil {
		return nil, fmt.Errorf("failed to generate custom board design: %w", err)
	}

	return result, nil
}

func (u *GenerateCustomDesignUsecase) checkAndRecordRateLimit(customerID uuid.UUID) error {
	u.mu.Lock()
	defer u.mu.Unlock()

	now := time.Now()
	oneHourAgo := now.Add(-1 * time.Hour)

	// Clean up timestamps older than 1 hour
	history := u.userRequestHistory[customerID]
	var validHistory []time.Time
	for _, t := range history {
		if t.After(oneHourAgo) {
			validHistory = append(validHistory, t)
		}
	}

	if len(validHistory) >= u.maxRequestsPerHour {
		u.userRequestHistory[customerID] = validHistory
		return apperrors.NewTooManyRequests(fmt.Sprintf("hourly AI generation limit reached (%d requests/hour). Please try again later.", u.maxRequestsPerHour))
	}

	validHistory = append(validHistory, now)
	u.userRequestHistory[customerID] = validHistory
	return nil
}
