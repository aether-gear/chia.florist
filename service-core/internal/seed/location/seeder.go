package location

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func SeedAll(ctx context.Context, pool *pgxpool.Pool) error {
	seeded, err := isAlreadySeeded(ctx, pool)
	if err != nil {
		return err
	}

	if seeded {
		log.Println("location already seeded, skipping...")
		return nil
	}

	log.Println("seeding provinces...")
	if err := seedProvinces(ctx, pool); err != nil {
		return err
	}

	log.Println("seeding cities...")
	if err := seedCities(ctx, pool); err != nil {
		return err
	}

	log.Println("seeding districts...")
	if err := seedDistricts(ctx, pool); err != nil {
		return err
	}

	log.Println("seeding villages...")
	if err := seedDistricts(ctx, pool); err != nil {
		return err
	}

	return markSeeded(ctx, pool)
}

func isAlreadySeeded(ctx context.Context, pool *pgxpool.Pool) (bool, error) {
	var exists bool

	err := pool.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT 1 FROM seed_versions WHERE name = $1
		)
	`, "location_v2").Scan(&exists)

	return exists, err
}

func markSeeded(ctx context.Context, pool *pgxpool.Pool) error {
	_, err := pool.Exec(ctx, `
		INSERT INTO seed_versions (name, version)
		VALUES ($1, $2)
		ON CONFLICT (name) DO NOTHING
	`, "location_v2", "1.0")

	return err
}

func seedProvinces(ctx context.Context, pool *pgxpool.Pool) error {
	rows, err := LoadCSV("internal/seed/location/source/provinces.csv")
	if err != nil {
		return err
	}

	tx, err := pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	for _, r := range rows[1:] {
		_, err := tx.Exec(ctx, `
			INSERT INTO provinces (id, name)
			VALUES ($1, $2)
			ON CONFLICT (id) DO NOTHING
		`, r[0], r[1])

		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func seedCities(ctx context.Context, pool *pgxpool.Pool) error {
	rows, err := LoadCSV("internal/seed/location/source/regencies.csv")
	if err != nil {
		return err
	}

	tx, err := pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	for _, r := range rows[1:] {
		_, err := tx.Exec(ctx, `
			INSERT INTO cities (id, province_id, name)
			VALUES ($1, $2, $3)
			ON CONFLICT (id) DO NOTHING
		`, r[0], r[1], r[2])

		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func seedDistricts(ctx context.Context, pool *pgxpool.Pool) error {
	rows, err := LoadCSV("internal/seed/location/source/districts.csv")
	if err != nil {
		return err
	}

	tx, err := pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	for _, r := range rows[1:] {
		_, err := tx.Exec(ctx, `
			INSERT INTO districts (id, city_id, name)
			VALUES ($1, $2, $3)
			ON CONFLICT (id) DO NOTHING
		`, r[0], r[1], r[2])

		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func seedVillages(ctx context.Context, pool *pgxpool.Pool) error {
	rows, err := LoadCSV("internal/seed/location/source/villages.csv")
	if err != nil {
		return err
	}

	tx, err := pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	for _, r := range rows[1:] {
		_, err := tx.Exec(ctx, `
			INSERT INTO villages (id, district_id, name)
			VALUES ($1, $2, $3)
			ON CONFLICT (id) DO NOTHING
		`, r[0], r[1], r[2])

		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}
