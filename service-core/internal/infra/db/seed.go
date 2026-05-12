package database

import (
	"context"
	"fmt"
	"log"
	"service-core/internal/seed/courier"
	"service-core/internal/seed/location"
)

func RunSeed(conn *Connection) error {
	ctx := context.Background()

	if err := location.SeedAll(ctx, conn.Pool); err != nil {
		return fmt.Errorf("seed locations: %w", err)
	}
	log.Printf("database: locations seeded")

	if err := courier.SeedAll(ctx, conn.Pool); err != nil {
		return fmt.Errorf("seed couriers: %w", err)
	}
	log.Printf("database: couriers seeded")

	return nil
}
