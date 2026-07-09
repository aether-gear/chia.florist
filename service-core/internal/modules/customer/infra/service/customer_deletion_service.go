package service

import (
	"context"
	"fmt"

	addressRepo "service-core/internal/modules/address/repository"
	cartRepo "service-core/internal/modules/cart/repository"
	customerRepo "service-core/internal/modules/customer/repository"
	transaction "service-core/internal/shared/transaction"

	"github.com/google/uuid"
)

type customerDeletionServiceImpl struct {
	customerRepo customerRepo.CustomerRepository
	addressRepo  addressRepo.CustomerAddressRepository
	cartRepo     cartRepo.CartRepository
}

func NewCustomerDeletionService(
	customerRepo customerRepo.CustomerRepository,
	addressRepo  addressRepo.CustomerAddressRepository,
	cartRepo     cartRepo.CartRepository,
) customerRepo.CustomerDeletionService {
	return &customerDeletionServiceImpl{
		customerRepo: customerRepo,
		addressRepo:  addressRepo,
		cartRepo:     cartRepo,
	}
}

func (s *customerDeletionServiceImpl) DeleteCustomerRecord(
	ctx context.Context,
	exec transaction.Executor,
	customerID uuid.UUID,
) error {
	// 1. Soft delete customer record
	if err := s.customerRepo.Delete(ctx, exec, customerID); err != nil {
		return fmt.Errorf("failed to soft delete customer: %w", err)
	}

	// 2. Soft delete customer addresses
	if err := s.addressRepo.DeleteByCustomerID(ctx, exec, customerID); err != nil {
		return fmt.Errorf("failed to soft delete customer addresses: %w", err)
	}

	// 3. Soft delete customer cart items
	if err := s.cartRepo.DeleteByCustomerID(ctx, exec, customerID); err != nil {
		return fmt.Errorf("failed to soft delete cart items: %w", err)
	}

	return nil
}
