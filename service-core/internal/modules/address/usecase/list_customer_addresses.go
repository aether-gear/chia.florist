package usecase

import (
	"context"
	"fmt"

	"service-core/internal/modules/address/domain"
	"service-core/internal/modules/address/repository"
	transaction "service-core/internal/shared/transaction"

	"github.com/google/uuid"
)

type ListCustomerAddressesUsecase struct {
	customerAddressRepo repository.CustomerAddressRepository
	executor            transaction.Executor
}

func NewListCustomerAddressesUsecase(
	customerAddressRepo repository.CustomerAddressRepository,
	executor transaction.Executor,
) *ListCustomerAddressesUsecase {
	return &ListCustomerAddressesUsecase{
		customerAddressRepo: customerAddressRepo,
		executor:            executor,
	}
}

func (u *ListCustomerAddressesUsecase) ListByCustomerID(
	ctx context.Context,
	customerID uuid.UUID,
) ([]domain.CustomerAddress, error) {
	res, err := u.customerAddressRepo.ListByCustomerID(ctx, u.executor, customerID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve address: %w", err)
	}

	return res, nil
}
