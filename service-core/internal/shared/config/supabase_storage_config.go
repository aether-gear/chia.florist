package config

import (
	"time"

	"service-core/internal/infra/storage"
)

type StorageConfig struct {
	Provider        storage.ProviderName
	LocalRootDir    string
	LocalPublicBase string
	Supabase        SupabaseStorageConfig
}

type SupabaseStorageConfig struct {
	ProjectURL      string
	ServiceRoleKey  string
	BucketName      string
	PublicBaseURL   string
	SignedURLExpiry time.Duration
}

func LoadStorageConfig() StorageConfig {
	return StorageConfig{
		Provider:        storage.ProviderName(GetEnv("STORAGE_PROVIDER", string(storage.ProviderLocal))),
		LocalRootDir:    GetEnv("STORAGE_LOCAL_ROOT", "storage"),
		LocalPublicBase: GetEnv("STORAGE_LOCAL_PUBLIC_BASE", ""),
		Supabase:        LoadSupabaseStorageConfig(),
	}
}

func LoadSupabaseStorageConfig() SupabaseStorageConfig {
	return SupabaseStorageConfig{
		ProjectURL:      GetEnv("SUPABASE_STORAGE_PROJECT_URL", ""),
		ServiceRoleKey:  GetEnv("SUPABASE_STORAGE_SERVICE_ROLE_KEY", ""),
		BucketName:      GetEnv("SUPABASE_STORAGE_BUCKET", ""),
		PublicBaseURL:   GetEnv("SUPABASE_STORAGE_PUBLIC_BASE_URL", ""),
		SignedURLExpiry: loadStorageDuration("SUPABASE_STORAGE_SIGNED_URL_EXPIRY", 15*time.Minute),
	}
}

func loadStorageDuration(key string, fallback time.Duration) time.Duration {
	raw := GetEnv(key, fallback.String())

	duration, err := time.ParseDuration(raw)
	if err != nil {
		return fallback
	}

	return duration
}
