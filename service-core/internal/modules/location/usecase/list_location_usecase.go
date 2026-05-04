package usecase

import (
	"fmt"

	"service-core/internal/modules/location/domain"
	"service-core/internal/modules/location/repository"
)

type ListLocationUsecase struct {
	locationRepo repository.LocationRepository
}

func NewListLocationUsecase(
	lR repository.LocationRepository,
) *ListLocationUsecase {
	return &ListLocationUsecase{
		locationRepo: lR,
	}
}

func (u *ListLocationUsecase) Province() ([]domain.Province, error) {
	res, err := u.locationRepo.ListProvinces()
	if err != nil {
		return nil, fmt.Errorf("failed to load provinces: %w", err)
	}

	return res, nil
}

func (u *ListLocationUsecase) City(provinceID string) ([]domain.City, error) {
	res, err := u.locationRepo.ListCitiesByProvince(provinceID)
	if err != nil {
		return nil, fmt.Errorf("failed to load cities: %w", err)
	}

	return res, nil
}

func (u *ListLocationUsecase) District(cityID string) ([]domain.District, error) {
	res, err := u.locationRepo.ListDistrictsByCity(cityID)
	if err != nil {
		return nil, fmt.Errorf("failed to load districts: %w", err)
	}

	return res, nil
}

func (u *ListLocationUsecase) Village(districtID string) ([]domain.Village, error) {
	res, err := u.locationRepo.ListVillagesByDistrict(districtID)
	if err != nil {
		return nil, fmt.Errorf("failed to load villages: %w", err)
	}

	return res, nil
}
