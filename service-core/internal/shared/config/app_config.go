package config

import "strings"

type AppConfig struct {
	Env                string
	CORSAllowedOrigins []string
}

func LoadAppConfig() AppConfig {
	return AppConfig{
		Env:                GetEnv("APP_ENV"),
		CORSAllowedOrigins: splitCSVEnv("APP_CORS_ALLOWED_ORIGINS"),
	}
}

func splitCSVEnv(key string) []string {
	raw := GetEnv(key, "")
	if raw == "" {
		return nil
	}

	parts := strings.Split(raw, ",")

	values := make([]string, 0, len(parts))
	for _, part := range parts {
		value := strings.TrimSpace(part)
		if value == "" {
			continue
		}

		values = append(values, value)
	}

	return values
}
