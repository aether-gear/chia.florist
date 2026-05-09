package database

import (
	"fmt"

	"service-core/internal/shared/config"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigration(cfg config.DatabaseConfig) error {
	m, err := migrate.New(
		"file://migrations",
		*cfg.DSN,
	)
	if err != nil {
		return fmt.Errorf("failed to init migration: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migration failed: %w", err)
	}

	return nil
}
