package apphttp

import (
	"net/http"
	"strconv"

	"github.com/google/uuid"
)

func Query(r *http.Request, key string) string {
	return r.URL.Query().Get(key)
}

func QueryInt(r *http.Request, key string) (int, error) {
	value := r.URL.Query().Get(key)
	if value == "" {
		return 0, nil
	}

	return strconv.Atoi(value)
}

func QueryUUID(r *http.Request, key string) (*uuid.UUID, error) {
	value := r.URL.Query().Get(key)

	parsed, err := uuid.Parse(value)
	if err != nil {
		return nil, err
	}

	return &parsed, nil
}

func QueryBool(r *http.Request, key string) bool {
	value := r.URL.Query().Get(key)

	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return false
	}

	return parsed
}

func QueryIntDefault(r *http.Request, key string, fallback int) int {
	value := r.URL.Query().Get(key)
	if value == "" {
		return fallback
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}

	return parsed
}
