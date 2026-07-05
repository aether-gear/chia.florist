package bootstrap

import (
	"service-core/internal/shared/config"
)

type Config struct {
	App         config.AppConfig
	Shipping    config.ShippingConfig
	JWT         config.JWTConfig
	Storage     config.StorageConfig
	Supabase    config.SupabaseConfig
	Postgres    config.PostgresConfig
	DB          config.DatabaseConfig
	SMTP        config.SMTPConfig
	MidTrans    config.MidTransConfig
	RajaOngkir  config.RajaOngkirConfig
	Komerce     config.KomerceConfig
	GoogleOAuth config.GoogleOAuthConfig
	WAF         config.WAFConfig
}

func LoadConfig() Config {
	supabaseCfg := config.LoadSupabaseConfig()
	postgresCfg := config.LoadPostgresConfig()

	dbConf := config.LoadDBConfig(
		supabaseCfg.Host,
		supabaseCfg.Port,
		supabaseCfg.User,
		supabaseCfg.Password,
		supabaseCfg.Name,
		supabaseCfg.SSLMode,
		&supabaseCfg.DSN,
	)
	// dbConf := config.LoadDBConfig(
	// 	postgresCfg.Host,
	// 	postgresCfg.Port,
	// 	postgresCfg.User,
	// 	postgresCfg.Password,
	// 	postgresCfg.Name,
	// 	postgresCfg.SSLMode,
	// 	&postgresCfg.DSN,
	// )

	return Config{
		App:         config.LoadAppConfig(),
		Shipping:    config.LoadShippingConfig(),
		JWT:         config.LoadJWTConfig(),
		Storage:     config.LoadStorageConfig(),
		Supabase:    supabaseCfg,
		Postgres:    postgresCfg,
		DB:          dbConf,
		SMTP:        config.LoadSMTPConfig(),
		MidTrans:    config.LoadMidTransConfig(),
		RajaOngkir:  config.LoadRajaOngkirConfig(),
		Komerce:     config.LoadKomerceConfig(),
		GoogleOAuth: config.LoadGoogleOAuthConfig(),
		WAF:         config.LoadWAFConfig(),
	}
}
