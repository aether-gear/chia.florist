package usecase

import (
	"context"
	"fmt"

	apperrors "service-core/internal/common/errors"
	"service-core/internal/modules/courier/domain"
	"service-core/internal/modules/courier/repository"
	shopRepo "service-core/internal/modules/shop/repository"
	transaction "service-core/internal/shared/transaction"

	"github.com/google/uuid"
)

type ConfigureShopCourierUsecase struct {
	executor        transaction.Executor
	transactor      transaction.Transactor
	courierRepo     repository.CourierRepository
	shopCourierRepo repository.ShopCourierRepository
	shopRepo        shopRepo.ShopRepository
}

func NewConfigureShopCourierUsecase(
	executor transaction.Executor,
	transactor transaction.Transactor,
	courierRepo repository.CourierRepository,
	shopCourierRepo repository.ShopCourierRepository,
	shopRepo shopRepo.ShopRepository,
) *ConfigureShopCourierUsecase {
	return &ConfigureShopCourierUsecase{
		executor:        executor,
		transactor:      transactor,
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
	shop, err := u.shopRepo.
		GetByID(ctx, u.executor, shopID)
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

	validCodes, err := u.courierRepo.
		ValidateCouriers(ctx, u.executor, codes)
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

	err = u.transactor.WithinTransaction(
		ctx,
		func(exec transaction.Executor) error {
			for _, shopCourier := range shopCouriers {
				if err := u.shopCourierRepo.
					SaveShopCourier(
						ctx,
						exec,
						shopCourier,
					); err != nil {
					return fmt.Errorf("failed to save shop couriers: %w", err)
				}
			}

			return nil
		},
	)
	if err != nil {
		return err
	}

	return nil
}
