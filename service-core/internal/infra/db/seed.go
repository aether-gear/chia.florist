package database

import (
	"context"
	"log"
	"service-core/internal/seed/courier"
	"service-core/internal/seed/location"
)

func RunSeed(conn *Connection) {
	log.Println("starting seed process...")

	ctx := context.Background()

	if err := location.SeedAll(ctx, conn.Pool); err != nil {
		log.Fatal(err)
	}

	if err := courier.SeedAll(ctx, conn.Pool); err != nil {
		log.Fatal(err)
	}
}
