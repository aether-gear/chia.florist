package repository

import (
	"context"

	"service-core/internal/modules/location/domain"
)

type LocationRepository interface {
	ListProvinces(
		ctx context.Context,
	) ([]domain.Province, error)
	ListCitiesByProvince(
		ctx context.Context,
		provinceID string,
	) ([]domain.City, error)
	ListDistrictsByCity(
		ctx context.Context,
		cityID string,
	) ([]domain.District, error)
	ListVillagesByDistrict(
		ctx context.Context,
		districtID string,
	) ([]domain.Village, error)
}
