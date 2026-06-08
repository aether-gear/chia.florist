package usecase

import (
	"context"
	"fmt"
	"time"

	apperrors "service-core/internal/common/errors"
	"service-core/internal/modules/merchant/domain"
	"service-core/internal/modules/merchant/repository"
	transaction "service-core/internal/shared/transaction"

	"github.com/google/uuid"
)

type CreateMerchantUsecase struct {
	merchantRepo repository.MerchantRepository
	executor     transaction.Executor
}

func NewCreateMerchantUsecase(
	merchantRepo repository.MerchantRepository,
	executor transaction.Executor,
) *CreateMerchantUsecase {
	return &CreateMerchantUsecase{
		merchantRepo: merchantRepo,
		executor:     executor,
	}
}

type CreateMerchantInput struct {
	Name        string
	Description *string
	LogoUrl     *string
	BannerUrl   *string
}

func (u *CreateMerchantUsecase) Execute(
	ctx context.Context,
	input CreateMerchantInput,
) error {
	if input.Name == "" {
		return apperrors.NewInvalidInput(domain.ErrInvalidName.Error())
	}

	merchant := domain.Merchant{
		ID:          uuid.New(),
		Name:        input.Name,
		Description: input.Description,
		LogoUrl:     input.LogoUrl,
		BannerUrl:   input.BannerUrl,
		CreatedAt:   time.Now(),
	}

	if err := u.merchantRepo.Create(ctx, u.executor, merchant); err != nil {
		return fmt.Errorf("failed to create merchant: %w", err)
	}

	return nil
}
