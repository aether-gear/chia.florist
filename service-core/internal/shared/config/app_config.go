package config

import (
	"strconv"
	"strings"
)

type AppConfig struct {
	Codename           string
	Host               string
	Port               string
	Env                string
	CORSAllowedOrigins []string
	ServerTimeout      int // in seconds
	ShutdownTimeout    int // in seconds
}

func LoadAppConfig() AppConfig {
	timeout, err := strconv.Atoi(GetEnv("SERVER_TIMEOUT", "30"))
	if err != nil {
		timeout = 30
	}

	shutdownTimeout, err := strconv.Atoi(GetEnv("SHUTDOWN_TIMEOUT", "15"))
	if err != nil {
		shutdownTimeout = 15
	}

	return AppConfig{
		Codename:           GetEnv("MY_APP_SAYA", ""),
		Host:               GetEnv("MY_APP_SAYA_HOST", "0.0.0.0"),
		Port:               GetEnv("MY_APP_SAYA_PORT", "8000"),
		Env:                GetEnv("APP_ENV"),
		CORSAllowedOrigins: splitCSVEnv("APP_CORS_ALLOWED_ORIGINS"),
		ServerTimeout:      timeout,
		ShutdownTimeout:    shutdownTimeout,
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
