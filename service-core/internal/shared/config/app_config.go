package config

import "strings"

type AppConfig struct {
	Host               string
	Port               string
	Env                string
	CORSAllowedOrigins []string
}

func LoadAppConfig() AppConfig {
	return AppConfig{
		Host:               GetEnv("MY_APP_SAYA_HOST", "0.0.0.0"),
		Port:               GetEnv("MY_APP_SAYA_PORT", "8000"),
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
