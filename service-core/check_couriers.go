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

	rows, err := pool.Query(ctx, "SELECT id, name FROM shops")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	fmt.Println("=== SHOPS ===")
	for rows.Next() {
		var id, name string
		rows.Scan(&id, &name)
		fmt.Printf("- %s: %s\n", id, name)
	}

	cRows, err := pool.Query(ctx, "SELECT shop_id, code, active FROM shop_couriers")
	if err != nil {
		log.Fatal(err)
	}
	defer cRows.Close()

	fmt.Println("\n=== SHOP COURIERS ===")
	count := 0
	for cRows.Next() {
		var shopId, code string
		var active bool
		cRows.Scan(&shopId, &code, &active)
		fmt.Printf("- %s: %s (active: %v)\n", shopId, code, active)
		count++
	}
	if count == 0 {
		fmt.Println("No shop couriers found!")
	}
}
