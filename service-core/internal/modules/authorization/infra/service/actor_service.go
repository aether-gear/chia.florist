package service

import (
	"context"
	"fmt"

	authenDomain "service-core/internal/modules/authentication/domain"
	authenRepo "service-core/internal/modules/authentication/repository"
	"service-core/internal/modules/authorization/domain"
	"service-core/internal/modules/authorization/repository"
	merchantRepo "service-core/internal/modules/merchant/repository"
	transaction "service-core/internal/shared/transaction"

	"github.com/google/uuid"
)

type ActorService struct {
	accountRepo  authenRepo.AccountRepository
	merchantRepo merchantRepo.MerchantRepository
}

func NewActorService(
	accountRepo authenRepo.AccountRepository,
	merchantRepo merchantRepo.MerchantRepository,
) repository.ActorService {
	return &ActorService{
		accountRepo:  accountRepo,
		merchantRepo: merchantRepo,
	}
}

func (s *ActorService) Load(
	ctx context.Context,
	exec transaction.Executor,
	userID uuid.UUID,
) (*domain.Actor, error) {
	account, err := s.accountRepo.GetByUserID(ctx, exec, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account: %w", err)
	}

	actor := domain.Actor{
		AccountID: account.ID,
		Type:      account.Type,
	}

	if actor.Type == authenDomain.AccountTypeMerchant {
		merchant, err := s.merchantRepo.GetByAccountID(ctx, exec, account.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve merchant: %w", err)
		}

		actor.MerchantID = &merchant.ID
	}

	return &actor, nil
}
