package usecase

import (
	"context"
	"fmt"
	"strings"
	"time"

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
	shop, err := u.shopRepo.GetByID(ctx, u.executor, shopID)
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

	validCodes, err := u.courierRepo.ValidateCouriers(ctx, u.executor, codes)
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

		existing, err := u.shopCourierRepo.GetByShopIDAndCode(ctx, u.executor,
			shopID,
			input.Code,
		)
		if err != nil {
			return fmt.Errorf("failed to retrieve courier: %w", err)
		}

		var (
			existingName    *string
			existingAddress *string
			existingStatus  = domain.CourierVerificationUnconfigured
		)

		if existing != nil {
			existingName = existing.Name
			existingAddress = existing.LocationAddress
			existingStatus = existing.VerificationStatus
		}

		sc := domain.ShopCourier{
			ShopID:             shopID,
			Code:               input.Code,
			Name:               existingName,
			LocationAddress:    existingAddress,
			Active:             input.Active,
			VerificationStatus: existingStatus,
		}
		if err := sc.Validate(); err != nil {
			return apperrors.NewBadRequest(fmt.Sprintf("%s: %s", input.Code, err.Error()))
		}

		shopCouriers = append(shopCouriers, sc)
	}

	err = u.transactor.WithinTransaction(ctx, func(exec transaction.Executor) error {
		for _, shopCourier := range shopCouriers {
			if err := u.shopCourierRepo.SaveShopCourier(ctx, exec, shopCourier); err != nil {
				return fmt.Errorf("failed to save shop couriers: %w", err)
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

type UpdateSingleShopCourierInput struct {
	ShopID          uuid.UUID
	Code            string
	Name            *string
	LocationAddress *string
	Active          bool
	IsAdmin         bool
	AdminStaffID    *uuid.UUID
}

func (u *ConfigureShopCourierUsecase) UpdateSingle(
	ctx context.Context,
	input UpdateSingleShopCourierInput,
) (*domain.ShopCourier, error) {
	shop, err := u.shopRepo.GetByID(ctx, u.executor, input.ShopID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve shop: %w", err)
	}
	if shop == nil {
		return nil, apperrors.NewNotFound(domain.ErrShopNotFound.Error())
	}

	validCodes, err := u.courierRepo.ValidateCouriers(ctx, u.executor, []string{input.Code})
	if err != nil || len(validCodes) == 0 {
		return nil, apperrors.NewBadRequest(domain.ErrCourierNotFound.Error())
	}

	var trimmedName *string
	if input.Name != nil {
		clean := strings.TrimSpace(*input.Name)
		trimmedName = &clean
	}

	var trimmedAddress *string
	if input.LocationAddress != nil {
		clean := strings.TrimSpace(*input.LocationAddress)
		trimmedAddress = &clean
	}

	if input.Active {
		if trimmedName == nil || *trimmedName == "" {
			return nil, apperrors.NewBadRequest(domain.ErrCourierNameRequired.Error())
		}
		if trimmedAddress == nil || *trimmedAddress == "" {
			return nil, apperrors.NewBadRequest(domain.ErrCourierLocationRequired.Error())
		}
	}

	existing, err := u.shopCourierRepo.GetByShopIDAndCode(ctx, u.executor, input.ShopID, input.Code)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing courier: %w", err)
	}

	if !input.Active {
		if (trimmedName == nil || *trimmedName == "") && existing != nil {
			trimmedName = existing.Name
		}
		if (trimmedAddress == nil || *trimmedAddress == "") && existing != nil {
			trimmedAddress = existing.LocationAddress
		}
	}

	var (
		verifiedAt      *time.Time
		verifiedBy      *uuid.UUID
		rejectionReason *string
	)

	finalActive := false
	verificationStatus := domain.CourierVerificationUnconfigured
	if input.Active {
		if input.IsAdmin {
			finalActive = true
			verificationStatus = domain.CourierVerificationVerified

			now := time.Now()
			verifiedAt = &now
			verifiedBy = input.AdminStaffID
		} else {
			verificationStatus = domain.CourierVerificationPending
		}
	}

	shopCourier := domain.ShopCourier{
		ShopID:             input.ShopID,
		Code:               input.Code,
		Name:               trimmedName,
		LocationAddress:    trimmedAddress,
		Active:             finalActive,
		VerificationStatus: verificationStatus,
		VerifiedAt:         verifiedAt,
		VerifiedBy:         verifiedBy,
		RejectionReason:    rejectionReason,
	}
	if err := shopCourier.Validate(); err != nil {
		return nil, apperrors.NewBadRequest(err.Error())
	}

	err = u.transactor.WithinTransaction(ctx, func(exec transaction.Executor) error {
		return u.shopCourierRepo.SaveShopCourier(ctx, exec, shopCourier)
	})
	if err != nil {
		return nil, fmt.Errorf("failed to save shop courier: %w", err)
	}

	updated, err := u.shopCourierRepo.GetByShopIDAndCode(ctx, u.executor,
		input.ShopID,
		input.Code,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve updated shop courier: %w", err)
	}

	return updated, nil
}
