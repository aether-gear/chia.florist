package config

import "time"

type StorageConfig struct {
	SignedURLExpiry time.Duration
}

func LoadStorageConfig() StorageConfig {
	expiry, err := time.ParseDuration(
		GetEnv("STORAGE_SIGNED_URL_EXPIRY", ""),
	)
	if err != nil {
		expiry = 15 * time.Minute
	}

	return StorageConfig{
		SignedURLExpiry: expiry,
	}
}
