package usecase

import (
	"context"
	"fmt"

	shipping "service-core/internal/infra/shipping"
	"service-core/internal/modules/location/domain"
	transaction "service-core/internal/shared/transaction"
)

type ListLocationUsecase struct {
	shipping shipping.Provider
	executor transaction.Executor
}

func NewListLocationUsecase(
	shipping shipping.Provider,
	executor transaction.Executor,
) *ListLocationUsecase {
	return &ListLocationUsecase{
		shipping: shipping,
		executor: executor,
	}
}

func (u *ListLocationUsecase) Province(ctx context.Context) ([]domain.Province, error) {
	res, err := u.shipping.ListProvinces(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to load provinces: %w", err)
	}

	return res, nil
}

func (u *ListLocationUsecase) City(ctx context.Context, provinceID string) ([]domain.City, error) {
	res, err := u.shipping.ListCities(ctx, provinceID)
	if err != nil {
		return nil, fmt.Errorf("failed to load cities: %w", err)
	}

	return res, nil
}

func (u *ListLocationUsecase) District(ctx context.Context, cityID string) ([]domain.District, error) {
	res, err := u.shipping.ListDistricts(ctx, cityID)
	if err != nil {
		return nil, fmt.Errorf("failed to load districts: %w", err)
	}

	return res, nil
}

func (u *ListLocationUsecase) Village(ctx context.Context, districtID string) ([]domain.Village, error) {
	res, err := u.shipping.ListVillages(ctx, districtID)
	if err != nil {
		return nil, fmt.Errorf("failed to load villages: %w", err)
	}

	return res, nil
}
