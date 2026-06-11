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
	accountRepo    authenRepo.AccountRepository
	merchantRepo   merchantRepo.MerchantRepository
	membershipRepo repository.MerchantMembershipRepository
}

func NewActorService(
	accountRepo authenRepo.AccountRepository,
	merchantRepo merchantRepo.MerchantRepository,
	membershipRepo repository.MerchantMembershipRepository,
) repository.ActorService {
	return &ActorService{
		accountRepo:    accountRepo,
		merchantRepo:   merchantRepo,
		membershipRepo: membershipRepo,
	}
}

func (s *ActorService) Load(
	ctx context.Context,
	exec transaction.Executor,
	userID uuid.UUID,
	merchantID *uuid.UUID,
) (*domain.Actor, error) {
	account, err := s.accountRepo.
		GetByUserID(ctx, exec,
			userID,
		)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account: %w", err)
	}

	actor := domain.Actor{
		AccountID: account.ID,
		Type:      account.Type,
	}

	if actor.Type == authenDomain.AccountTypeMerchant {
		roles, err := s.membershipRepo.
			ListRolesByAccountIDAndMerchantID(ctx, exec,
				account.ID,
				*merchantID,
			)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve merchant roles: %w", err)
		}

		actor.Roles = roles
	}

	return &actor, nil
}
