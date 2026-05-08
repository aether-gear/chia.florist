package database

import (
	"log"

	"service-core/internal/shared/config"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigrations(cfg config.DatabaseConfig) {
	m, err := migrate.New(
		"file://migrations",
		*cfg.DSN,
	)
	if err != nil {
		log.Fatalf("failed to init migration: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("migration failed: %v", err)
	}

	log.Println("migrations applied successfully")
}
