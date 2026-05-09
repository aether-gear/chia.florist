package bootstrap

import (
	"service-core/internal/shared/config"
)

type Config struct {
	App      config.AppConfig
	Shipping config.ShippingConfig
	JWT      config.JWTConfig
	Storage  config.StorageConfig
	Supabase config.SupabaseConfig
	DB       config.DatabaseConfig
}

func LoadConfig() Config {
	supabaseCfg := config.LoadSupabaseConfig()

	return Config{
		App:      config.LoadAppConfig(),
		Shipping: config.LoadShippingConfig(),
		JWT:      config.LoadJWTConfig(),
		Storage:  config.LoadStorageConfig(),
		Supabase: supabaseCfg,
		DB: config.LoadDBConfig(
			supabaseCfg.Host,
			supabaseCfg.Port,
			supabaseCfg.User,
			supabaseCfg.Password,
			supabaseCfg.Name,
			supabaseCfg.SSLMode,
			&supabaseCfg.DSN,
		),
	}
}
