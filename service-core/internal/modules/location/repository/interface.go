package repository

import (
	"context"

	"service-core/internal/modules/location/domain"
	transaction "service-core/internal/shared/transaction"
)

type LocationRepository interface {
	ListProvinces(
		ctx context.Context,
		exec transaction.Executor,
	) ([]domain.Province, error)
	ListCitiesByProvince(
		ctx context.Context,
		exec transaction.Executor,
		provinceID string,
	) ([]domain.City, error)
	ListDistrictsByCity(
		ctx context.Context,
		exec transaction.Executor,
		cityID string,
	) ([]domain.District, error)
	ListVillagesByDistrict(
		ctx context.Context,
		exec transaction.Executor,
		districtID string,
	) ([]domain.Village, error)
}
