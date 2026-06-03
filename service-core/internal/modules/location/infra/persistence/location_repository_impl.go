package persistence

import (
	"context"
	"fmt"

	database "service-core/internal/infra/db"
	"service-core/internal/modules/location/domain"
	"service-core/internal/modules/location/repository"

	"github.com/jackc/pgx/v5/pgxpool"
)

type locationRepositoryImpl struct {
	db *pgxpool.Pool
}

func NewLocationRepositoryImpl(conn *database.Connection) repository.LocationRepository {
	return &locationRepositoryImpl{
		db: conn.Pool,
	}
}

func (r *locationRepositoryImpl) ListProvinces(ctx context.Context) ([]domain.Province, error) {
	query := `
		SELECT
			id, name
		FROM provinces
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query provinces failed: %w", err)
	}
	defer rows.Close()

	var districts []domain.Province
	for rows.Next() {
		var v domain.Province

		err := rows.Scan(
			&v.ID,
			&v.Name,
		)
		if err != nil {
			return nil, fmt.Errorf("mapping province model to domain failed: %w", err)
		}

		districts = append(districts, v)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate provinces failed: %w", err)
	}

	return districts, nil
}

func (r *locationRepositoryImpl) ListCitiesByProvince(
	ctx context.Context, provinceID string,
) ([]domain.City, error) {
	query := `
		SELECT 
			id,
			province_id,
			name
		FROM cities
		WHERE province_id = $1
	`

	rows, err := r.db.Query(ctx, query, provinceID)
	if err != nil {
		return nil, fmt.Errorf("query cities by province id failed: %w", err)
	}
	defer rows.Close()

	var districts []domain.City
	for rows.Next() {
		var v domain.City

		err := rows.Scan(
			&v.ID,
			&v.ProvinceID,
			&v.Name,
		)
		if err != nil {
			return nil, fmt.Errorf("mapping city model to domain failed: %w", err)
		}

		districts = append(districts, v)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate cities failed: %w", err)
	}

	return districts, nil
}

func (r *locationRepositoryImpl) ListDistrictsByCity(
	ctx context.Context, cityID string,
) ([]domain.District, error) {
	query := `
		SELECT 
			id,
			city_id,
			name
		FROM districts
		WHERE city_id = $1
	`

	rows, err := r.db.Query(ctx, query, cityID)
	if err != nil {
		return nil, fmt.Errorf("query districts by city id failed: %w", err)
	}
	defer rows.Close()

	var districts []domain.District
	for rows.Next() {
		var v domain.District

		err := rows.Scan(
			&v.ID,
			&v.CityID,
			&v.Name,
		)
		if err != nil {
			return nil, fmt.Errorf("mapping district model to domain failed: %w", err)
		}

		districts = append(districts, v)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate districts failed: %w", err)
	}

	return districts, nil
}

func (r *locationRepositoryImpl) ListVillagesByDistrict(
	ctx context.Context, districtID string,
) ([]domain.Village, error) {
	query := `
		SELECT 
			id,
			district_id,
			name
		FROM villages
		WHERE district_id = $1
	`

	rows, err := r.db.Query(ctx, query, districtID)
	if err != nil {
		return nil, fmt.Errorf("query villages by district id failed: %w", err)
	}
	defer rows.Close()

	var villages []domain.Village
	for rows.Next() {
		var v domain.Village

		err := rows.Scan(
			&v.ID,
			&v.DistrictID,
			&v.Name,
		)
		if err != nil {
			return nil, fmt.Errorf("mapping village model to domain failed: %w", err)
		}

		villages = append(villages, v)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate villages failed: %w", err)
	}

	return villages, nil
}
