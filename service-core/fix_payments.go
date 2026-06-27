package main

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	connStr := "postgres://postgres.mqolpawlannysqjokzoq:Chia.Florist@21@aws-1-ap-northeast-2.pooler.supabase.com:6543/postgres?sslmode=disable&statement_cache_capacity=0"
	ctx := context.Background()

	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	// Insert BCA Bank Transfer
	var bcaID string
	err = pool.QueryRow(ctx, `
		INSERT INTO payment_methods (name, type, is_active, description, fee_type, fee_amount)
		VALUES ('BCA Bank Transfer', 'bank_transfer', true, 'Transfer to BCA account', 'flat', 0)
		ON CONFLICT (name) DO UPDATE SET is_active = true
		RETURNING id
	`).Scan(&bcaID)
	if err != nil {
		log.Fatal("Failed to insert BCA:", err)
	}

	// Insert BCA Account
	_, err = pool.Exec(ctx, `
		INSERT INTO payment_accounts (method_id, account_name, account_number, is_active)
		VALUES ($1, 'Chia Florist', '1234567890', true)
		ON CONFLICT (method_id, account_number) DO NOTHING
	`, bcaID)
	if err != nil {
		log.Fatal("Failed to insert BCA account:", err)
	}

	// Insert GoPay E-Wallet
	var gopayID string
	err = pool.QueryRow(ctx, `
		INSERT INTO payment_methods (name, type, is_active, description, fee_type, fee_amount)
		VALUES ('GoPay', 'ewallet', true, 'Pay using GoPay', 'flat', 0)
		ON CONFLICT (name) DO UPDATE SET is_active = true
		RETURNING id
	`).Scan(&gopayID)
	if err != nil {
		log.Fatal("Failed to insert GoPay:", err)
	}

	// Insert GoPay Account
	_, err = pool.Exec(ctx, `
		INSERT INTO payment_accounts (method_id, account_name, phone_number, is_active)
		VALUES ($1, 'Chia Florist', '081234567891', true)
	`, gopayID)
	if err != nil {
		log.Printf("Failed to insert GoPay account: %v", err)
	}

	// Insert QRIS
	var qrisID string
	err = pool.QueryRow(ctx, `
		INSERT INTO payment_methods (name, type, is_active, description, fee_type, fee_amount)
		VALUES ('QRIS', 'qr_code', true, 'Scan QRIS', 'flat', 0)
		ON CONFLICT (name) DO UPDATE SET is_active = true
		RETURNING id
	`).Scan(&qrisID)
	if err != nil {
		log.Fatal("Failed to insert QRIS:", err)
	}

	// Insert QRIS Account
	_, err = pool.Exec(ctx, `
		INSERT INTO payment_accounts (method_id, account_name, qr_string, is_active)
		VALUES ($1, 'Chia Florist QRIS', 'QRIS_STRING_HERE', true)
	`, qrisID)
	if err != nil {
		log.Printf("Failed to insert QRIS account: %v", err)
	}

	fmt.Println("Successfully seeded payment methods and accounts!")
}
