package database

import (
	"fmt"
	"log"

	"service-core/internal/shared/config"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigration(cfg config.DatabaseConfig) error {
	log.Printf("database: running migration")

	dsn := *cfg.DSN

	m, err := migrate.New(
		"file://migrations",
		dsn,
	)
	if err != nil {
		return fmt.Errorf("failed to init migration: %w", err)
	}

	err = m.Up()
	if err != nil {
		if err == migrate.ErrNoChange {
			log.Printf("database: already up to date")
			return nil
		}

		return fmt.Errorf("migration failed: %w", err)
	}

	log.Println("database: migration applied")
	return nil
}
