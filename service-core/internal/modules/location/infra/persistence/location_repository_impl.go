package persistence

import (
	"context"
	database "service-core/internal/infra/db"
	"service-core/internal/modules/location/domain"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type locationRepositoryImpl struct {
	db *pgxpool.Pool
}

func NewLocationRepositoryImpl(conn *database.Connection) *locationRepositoryImpl {
	return &locationRepositoryImpl{
		db: conn.Pool,
	}
}

func (r *locationRepositoryImpl) ListProvinces() ([]domain.Province, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT id, name
		FROM provinces
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
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
			return nil, err
		}

		districts = append(districts, v)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return districts, nil
}

func (r *locationRepositoryImpl) ListCitiesByProvince(provinceID string) ([]domain.City, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT id, province_id, name
		FROM cities
		WHERE province_id = $1
	`

	rows, err := r.db.Query(ctx, query, provinceID)
	if err != nil {
		return nil, err
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
			return nil, err
		}

		districts = append(districts, v)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return districts, nil
}

func (r *locationRepositoryImpl) ListDistrictsByCity(cityID string) ([]domain.District, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT id, city_id, name
		FROM districts
		WHERE city_id = $1
	`

	rows, err := r.db.Query(ctx, query, cityID)
	if err != nil {
		return nil, err
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
			return nil, err
		}

		districts = append(districts, v)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return districts, nil
}

func (r *locationRepositoryImpl) ListVillagesByDistrict(districtID string) ([]domain.Village, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT id, district_id, name
		FROM villages
		WHERE district_id = $1
	`

	rows, err := r.db.Query(ctx, query, districtID)
	if err != nil {
		return nil, err
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
			return nil, err
		}

		villages = append(villages, v)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return villages, nil
}
