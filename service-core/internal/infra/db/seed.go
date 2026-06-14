package database

import (
	"context"
	"fmt"
	"log"

	"service-core/seeds"
)

func RunSeed(conn *Connection) error {
	ctx := context.Background()

	if err := seeds.SeedLocations(ctx, conn.Pool); err != nil {
		return fmt.Errorf("seed locations: %w", err)
	}
	log.Printf("database: locations seeded")

	if err := seeds.SeedCouriers(ctx, conn.Pool); err != nil {
		return fmt.Errorf("seed couriers: %w", err)
	}
	log.Printf("database: couriers seeded")

	return nil
}
