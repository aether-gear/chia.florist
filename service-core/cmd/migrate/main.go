package main

import (
	"log"

	database "service-core/internal/infra/db"
	"service-core/internal/shared/config"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, fallback to system env")
	}

	supabaseCfg := config.LoadSupabaseConfig()
	cfg := config.LoadDBConfig(
		supabaseCfg.Host,
		supabaseCfg.Port,
		supabaseCfg.User,
		supabaseCfg.Password,
		supabaseCfg.Name,
		supabaseCfg.SSLMode,
		&supabaseCfg.DSN,
	)

	database.RunMigrations(cfg)

	log.Println("migration completed")
}
