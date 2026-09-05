package usecase

import (
	"context"
	"fmt"
	"strings"

	apperrors "service-core/internal/common/errors"
	"service-core/internal/modules/courier/domain"
	"service-core/internal/modules/courier/repository"
	shopRepo "service-core/internal/modules/shop/repository"
	transaction "service-core/internal/shared/transaction"

	"github.com/google/uuid"
)

type VerifyShopCourierUsecase struct {
	executor        transaction.Executor
	transactor      transaction.Transactor
	courierRepo     repository.CourierRepository
	shopCourierRepo repository.ShopCourierRepository
	shopRepo        shopRepo.ShopRepository
}

func NewVerifyShopCourierUsecase(
	executor transaction.Executor,
	transactor transaction.Transactor,
	courierRepo repository.CourierRepository,
	shopCourierRepo repository.ShopCourierRepository,
	shopRepo shopRepo.ShopRepository,
) *VerifyShopCourierUsecase {
	return &VerifyShopCourierUsecase{
		executor:        executor,
		transactor:      transactor,
		courierRepo:     courierRepo,
		shopCourierRepo: shopCourierRepo,
		shopRepo:        shopRepo,
	}
}

type VerifyShopCourierInput struct {
	ShopID          uuid.UUID
	Code            string
	Action          string // "verify" or "reject"
	RejectionReason *string
	AdminStaffID    uuid.UUID
}

func (u *VerifyShopCourierUsecase) Execute(
	ctx context.Context,
	input VerifyShopCourierInput,
) (*domain.ShopCourier, error) {
	shop, err := u.shopRepo.GetByID(ctx, u.executor, input.ShopID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve shop: %w", err)
	}
	if shop == nil {
		return nil, apperrors.NewNotFound(domain.ErrShopNotFound.Error())
	}

	normalizedAction := strings.ToLower(strings.TrimSpace(input.Action))
	if normalizedAction != "verify" && normalizedAction != "reject" {
		return nil, apperrors.NewBadRequest(domain.ErrInvalidVerificationAction.Error())
	}

	existing, err := u.shopCourierRepo.GetByShopIDAndCode(ctx, u.executor, input.ShopID, input.Code)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve courier record: %w", err)
	}
	if existing == nil {
		return nil, apperrors.NewNotFound(domain.ErrCourierNotFound.Error())
	}

	if normalizedAction == "verify" {
		target := *existing
		target.Active = true
		if err := target.Validate(); err != nil {
			return nil, apperrors.NewBadRequest(err.Error())
		}

		err = u.transactor.WithinTransaction(ctx, func(exec transaction.Executor) error {
			return u.shopCourierRepo.VerifyShopCourier(ctx, exec,
				input.ShopID,
				input.Code,
				domain.CourierVerificationVerified,
				true,
				input.AdminStaffID,
				nil,
			)
		})
		if err != nil {
			return nil, fmt.Errorf("failed to verify shop courier: %w", err)
		}
	} else {
		// Reject action
		var rejectionReason *string
		if input.RejectionReason != nil {
			reason := strings.TrimSpace(*input.RejectionReason)
			if reason != "" {
				rejectionReason = &reason
			}
		}

		err = u.transactor.WithinTransaction(ctx, func(exec transaction.Executor) error {
			return u.shopCourierRepo.VerifyShopCourier(ctx, exec,
				input.ShopID,
				input.Code,
				domain.CourierVerificationRejected,
				false,
				input.AdminStaffID,
				rejectionReason,
			)
		})
		if err != nil {
			return nil, fmt.Errorf("failed to reject shop courier: %w", err)
		}
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
