package database

import (
	"context"
	"fmt"
	"service-core/internal/seed/courier"
	"service-core/internal/seed/location"
)

func RunSeed(conn *Connection) error {
	ctx := context.Background()

	if err := location.SeedAll(ctx, conn.Pool); err != nil {
		return fmt.Errorf("seed locations: %w", err)
	}

	if err := courier.SeedAll(ctx, conn.Pool); err != nil {
		return fmt.Errorf("seed couriers: %w", err)
	}

	return nil
}
