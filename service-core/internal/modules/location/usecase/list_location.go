package usecase

import (
	"context"
	"fmt"

	"service-core/internal/modules/location/domain"
	"service-core/internal/modules/location/repository"
)

type ListLocationUsecase struct {
	locationRepo repository.LocationRepository
}

func NewListLocationUsecase(
	locationRepo repository.LocationRepository,
) *ListLocationUsecase {
	return &ListLocationUsecase{
		locationRepo: locationRepo,
	}
}

func (u *ListLocationUsecase) Province(ctx context.Context) ([]domain.Province, error) {
	res, err := u.locationRepo.ListProvinces(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to load provinces: %w", err)
	}

	return res, nil
}

func (u *ListLocationUsecase) City(ctx context.Context, provinceID string) ([]domain.City, error) {
	res, err := u.locationRepo.ListCitiesByProvince(ctx, provinceID)
	if err != nil {
		return nil, fmt.Errorf("failed to load cities: %w", err)
	}

	return res, nil
}

func (u *ListLocationUsecase) District(ctx context.Context, cityID string) ([]domain.District, error) {
	res, err := u.locationRepo.ListDistrictsByCity(ctx, cityID)
	if err != nil {
		return nil, fmt.Errorf("failed to load districts: %w", err)
	}

	return res, nil
}

func (u *ListLocationUsecase) Village(ctx context.Context, districtID string) ([]domain.Village, error) {
	res, err := u.locationRepo.ListVillagesByDistrict(ctx, districtID)
	if err != nil {
		return nil, fmt.Errorf("failed to load villages: %w", err)
	}

	return res, nil
}
