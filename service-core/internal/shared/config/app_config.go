package config

type AppConfig struct {
	Env string
}

func LoadAppConfig() AppConfig {
	return AppConfig{
		Env: GetEnv("APP_ENV"),
	}
}
