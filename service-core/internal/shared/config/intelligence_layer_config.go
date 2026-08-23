package config

import (
	"strconv"
)

type IntelligenceLayerConfig struct {
	BaseURL   string
	TimeoutMS int
	Enabled   bool
}

func LoadIntelligenceLayerConfig() IntelligenceLayerConfig {
	timeoutMS, err := strconv.Atoi(GetEnv("INTELLIGENCE_LAYER_TIMEOUT_MS", "500"))
	if err != nil || timeoutMS <= 0 {
		timeoutMS = 500
	}

	enabledStr := GetEnv("INTELLIGENCE_LAYER_ENABLED", "true")
	enabled := enabledStr == "true" || enabledStr == "1"

	return IntelligenceLayerConfig{
		BaseURL:   GetEnv("INTELLIGENCE_LAYER_BASE_URL", "http://localhost:8000/api/v1"),
		TimeoutMS: timeoutMS,
		Enabled:   enabled,
	}
}
