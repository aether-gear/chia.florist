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

	// Insert default address for each shop
	for _, shopID := range shopIDs {
		_, err := pool.Exec(ctx, `
			INSERT INTO shop_addresses (
				shop_id, label, phone, is_active,
				province, city, district, village, full_address, postal_code
			)
			VALUES (
				$1, 'Utama', '081234567890', true,
				'11', '149', '344', 'Desa XYZ', 'Jl. Kebon Jeruk No. 1', '12345'
			)
			ON CONFLICT DO NOTHING
		`, shopID)
		if err != nil {
			log.Printf("Error inserting address for shop %s: %v\n", shopID, err)
		} else {
			fmt.Printf("Successfully inserted address for shop %s\n", shopID)
		}
	}
	fmt.Println("Done inserting shop addresses.")
}
