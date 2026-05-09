package config

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
	DSN      *string
}

func LoadDBConfig(
	Host string,
	Port string,
	User string,
	Password string,
	Name string,
	SSLMode string,
	DSN *string,
) DatabaseConfig {
	return DatabaseConfig{
		Host:     Host,
		Port:     Port,
		User:     User,
		Password: Password,
		Name:     Name,
		SSLMode:  SSLMode,
		DSN:      DSN,
	}
}
