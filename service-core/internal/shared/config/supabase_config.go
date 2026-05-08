package config

// SupabaseConfig intentionally groups all Supabase-related infrastructure
// configuration, including database and storage concerns.
//
// Although database and object storage serve different capabilities at the
// application level, they share the same infrastructure provider boundary.
// Keeping them together avoids coupling infra credentials to capability-specific
// configs and makes infrastructure changes easier to reason about.
//
// Capability-level behavior (such as storage bucket or signed URL expiry)
// should remain in their respective app-level configs, while provider
// connection details live here.
type SupabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
	DSN      string

	ProjectURL     string
	ServiceRoleKey string
	PublicBaseURL  string

	DBVersion      *string
	StorageVersion *string
}

func LoadSupabaseConfig() SupabaseConfig {
	var (
		host           = GetEnv("DB_HOST", "")
		port           = GetEnv("DB_PORT", "")
		name           = GetEnv("DB_NAME", "")
		user           = GetEnv("DB_USER", "")
		password       = GetEnv("DB_PASSWORD", "")
		sslMode        = GetEnv("DB_SSLMODE", "")
		projectURL     = GetEnv("SUPABASE_PROJECT_URL", "")
		serviceRoleKey = GetEnv("SUPABASE_SERVICE_ROLE_KEY", "")
		publicBaseURL  = GetEnv("SUPABASE_PUBLIC_BASE_URL", "")
		dbVersion      = GetEnv("", "")
		storageVersion = GetEnv("", "")
	)

	dsn := "postgresql://" +
		user + ":" +
		password + "@" +
		host + ":" +
		port + "/" +
		name

	if sslMode != "" {
		dsn += "?sslmode=" + sslMode
	}

	return SupabaseConfig{
		Host:           host,
		Port:           port,
		User:           user,
		Password:       password,
		Name:           name,
		SSLMode:        sslMode,
		DSN:            dsn,
		ProjectURL:     projectURL,
		ServiceRoleKey: serviceRoleKey,
		PublicBaseURL:  publicBaseURL,
		DBVersion:      &dbVersion,
		StorageVersion: &storageVersion,
	}
}
