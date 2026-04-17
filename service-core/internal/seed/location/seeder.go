package location

import (
	"context"
	"log"
	"service-core/internal/shared/loader"
	"strings"

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
	if err := seedLocation(ctx,
		pool,
		"internal/seed/location/source/provinces.csv",
		"provinces",
		"id",
		"name",
	); err != nil {
		return err
	}

	log.Println("seeding cities...")
	if err := seedLocation(ctx,
		pool,
		"internal/seed/location/source/regencies.csv",
		"cities",
		"id",
		"province_id",
		"name",
	); err != nil {
		return err
	}

	log.Println("seeding districts...")
	if err := seedLocation(ctx,
		pool,
		"internal/seed/location/source/districts.csv",
		"districts",
		"id",
		"city_id",
		"name",
	); err != nil {
		return err
	}

	log.Println("seeding villages...")
	if err := seedLocation(ctx,
		pool,
		"internal/seed/location/source/villages.csv",
		"villages",
		"id",
		"district_id",
		"name",
	); err != nil {
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
	`, "location_v1").Scan(&exists)

	return exists, err
}

func markSeeded(ctx context.Context, pool *pgxpool.Pool) error {
	_, err := pool.Exec(ctx, `
		INSERT INTO seed_versions (name, version)
		VALUES ($1, $2)
		ON CONFLICT (name) DO NOTHING
	`, "location_v1", "1.0")

	return err
}

func seedLocation(ctx context.Context, pool *pgxpool.Pool, pathCSV string, table string, cols ...string) error {
	rows, err := loader.LoadCSV(pathCSV)
	if err != nil {
		return err
	}

	totalCols := len(cols)
	if len(rows) <= 1 {
		return nil
	}

	data := rows[1:]
	records := make([][]interface{}, 0, len(data))

	for _, r := range data {
		if len(r) == 1 {
			r = strings.Split(r[0], ";")
		}

		if len(r) < totalCols {
			continue
		}

		record := make([]interface{}, totalCols)
		for i := 0; i < totalCols; i++ {
			record[i] = strings.TrimSpace(r[i])
		}

		records = append(records, record)
	}

	if len(records) == 0 {
		return nil
	}

	_, err = pool.CopyFrom(
		ctx,
		pgx.Identifier{table},
		cols,
		pgx.CopyFromRows(records),
	)
	if err != nil {
		return err
	}

	return nil
}
