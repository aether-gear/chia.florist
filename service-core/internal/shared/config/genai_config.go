package config

import (
	"strconv"
)

type GenAIConfig struct {
	BaseURL            string
	APIKey             string
	Model              string
	TimeoutMS          int
	MaxRequestsPerHour int
	Enabled            bool
}

func LoadGenAIConfig() GenAIConfig {
	timeoutMS, err := strconv.Atoi(GetEnv("GENAI_TIMEOUT_MS", "15000"))
	if err != nil || timeoutMS <= 0 {
		timeoutMS = 15000
	}

	maxReqs, err := strconv.Atoi(GetEnv("GENAI_MAX_REQUESTS_PER_HOUR", "10"))
	if err != nil || maxReqs <= 0 {
		maxReqs = 10
	}

	enabledStr := GetEnv("GENAI_ENABLED", "true")
	enabled := enabledStr == "true" || enabledStr == "1"

	return GenAIConfig{
		BaseURL:            GetEnv("GENAI_BASE_URL", "https://api.banana.dev/v1"),
		APIKey:             GetEnv("GENAI_API_KEY", ""),
		Model:              GetEnv("GENAI_MODEL", "flower-board-v3"),
		TimeoutMS:          timeoutMS,
		MaxRequestsPerHour: maxReqs,
		Enabled:            enabled,
	}
}
