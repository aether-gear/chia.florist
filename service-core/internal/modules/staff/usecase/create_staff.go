package usecase

import (
	"context"
	"fmt"
	"time"

	"service-core/internal/modules/staff/domain"
	"service-core/internal/modules/staff/repository"
	transaction "service-core/internal/shared/transaction"

	"github.com/google/uuid"
)

type CreateStaffUsecase struct {
	staffRepo repository.StaffRepository
	executor  transaction.Executor
}

func NewCreateStaffUsecase(
	staffRepo repository.StaffRepository,
	executor transaction.Executor,
) *CreateStaffUsecase {
	return &CreateStaffUsecase{
		staffRepo: staffRepo,
		executor:  executor,
	}
}

type CreateStaffInput struct {
	UserID      string
	Description *string
	LogoUrl     *string
	BannerUrl   *string
}

func (u *CreateStaffUsecase) Execute(
	ctx context.Context,
	input CreateStaffInput,
) error {
	staff := domain.Staff{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
	}

	if err := u.staffRepo.Create(ctx, u.executor, staff); err != nil {
		return fmt.Errorf("failed to create staff: %w", err)
	}

	return nil
}
