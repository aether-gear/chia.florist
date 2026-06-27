package main

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	connStr := "postgres://postgres.mqolpawlannysqjokzoq:Chia.Florist@21@aws-1-ap-northeast-2.pooler.supabase.com:6543/postgres?sslmode=disable"
	ctx := context.Background()

	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	// Get all shop IDs
	rows, err := pool.Query(ctx, "SELECT id FROM shops")
	if err != nil {
		log.Fatal(err)
	}

	var shopIDs []string
	for rows.Next() {
		var id string
		rows.Scan(&id)
		shopIDs = append(shopIDs, id)
	}
	rows.Close()

	if len(shopIDs) == 0 {
		fmt.Println("No shops found in database.")
		return
	}

	// Couriers to enable for each shop
	couriers := []string{"jne", "sicepat", "tiki", "pos", "jnt"}

	// Insert couriers for each shop
	for _, shopID := range shopIDs {
		for _, courier := range couriers {
			_, err := pool.Exec(ctx, `
				INSERT INTO shop_couriers (shop_id, code, active)
				VALUES ($1, $2, true)
				ON CONFLICT (shop_id, code) DO UPDATE SET active = true
			`, shopID, courier)
			if err != nil {
				log.Printf("Error inserting courier %s for shop %s: %v\n", courier, shopID, err)
			}
		}
		fmt.Printf("Successfully configured couriers for shop %s\n", shopID)
	}
	fmt.Println("Done inserting shop couriers.")
}
