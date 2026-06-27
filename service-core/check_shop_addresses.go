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

	rows, err := pool.Query(ctx, "SELECT id, shop_id, district, is_active FROM shop_addresses")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	fmt.Println("=== SHOP ADDRESSES ===")
	count := 0
	for rows.Next() {
		var id, shopId, district string
		var isDefault bool
		rows.Scan(&id, &shopId, &district, &isDefault)
		fmt.Printf("- AddressID: %s, ShopID: %s, District: '%s', Active: %v\n", id, shopId, district, isDefault)
		count++
	}
	
	if count == 0 {
		fmt.Println("No shop addresses found in database!")
	}
}
