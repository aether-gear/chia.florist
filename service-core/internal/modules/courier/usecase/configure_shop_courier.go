package usecase

import (
	"context"
	"fmt"

	apperrors "service-core/internal/common/errors"
	"service-core/internal/modules/courier/domain"
	"service-core/internal/modules/courier/repository"
	shopRepo "service-core/internal/modules/shop/repository"

	"github.com/google/uuid"
)

type ConfigureShopCourierUsecase struct {
	courierRepo     repository.CourierRepository
	shopCourierRepo repository.ShopCourierRepository
	shopRepo        shopRepo.ShopRepository
}

func NewConfigureShopCourierUsecase(
	courierRepo repository.CourierRepository,
	shopCourierRepo repository.ShopCourierRepository,
	shopRepo shopRepo.ShopRepository,
) *ConfigureShopCourierUsecase {
	return &ConfigureShopCourierUsecase{
		courierRepo:     courierRepo,
		shopCourierRepo: shopCourierRepo,
		shopRepo:        shopRepo,
	}
}

type ConfigureShopCourierInput struct {
	Code   string
	Active bool
}

func (u *ConfigureShopCourierUsecase) Execute(
	ctx context.Context,
	shopID uuid.UUID,
	inputs []ConfigureShopCourierInput,
) error {
	shop, err := u.shopRepo.GetByID(ctx, shopID)
	if err != nil {
		return fmt.Errorf("failed to retrieve shop: %w", err)
	}
	if shop == nil {
		return apperrors.NewNotFound(domain.ErrShopNotFound.Error())
	}

	codes := make([]string, len(inputs))
	for i, input := range inputs {
		codes[i] = input.Code
	}

	validCodes, err := u.courierRepo.ValidateCouriers(ctx, codes)
	if err != nil {
		return fmt.Errorf("failed to validate couriers: %w", err)
	}
	if len(validCodes) != len(codes) {
		return fmt.Errorf("failed to validate couriers: some couriers are invalid")
	}

	validMap := make(map[string]struct{})
	for _, code := range validCodes {
		validMap[code] = struct{}{}
	}

	var shopCouriers []domain.ShopCourier

	for _, input := range inputs {
		if _, ok := validMap[input.Code]; !ok {
			return fmt.Errorf("invalid courier: %s", input.Code)
		}

		shopCouriers = append(shopCouriers, domain.ShopCourier{
			ShopID: shopID,
			Code:   input.Code,
			Active: input.Active,
		})
	}

	if err := u.shopCourierRepo.SaveShopCouriers(ctx, shopCouriers); err != nil {
		return fmt.Errorf("failed to save shop couriers: %w", err)
	}

	return nil
}
