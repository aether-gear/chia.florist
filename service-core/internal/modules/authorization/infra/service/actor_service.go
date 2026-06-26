package service

import (
	"context"
	"fmt"

	authenDomain "service-core/internal/modules/authentication/domain"
	authenRepo "service-core/internal/modules/authentication/repository"
	"service-core/internal/modules/authorization/domain"
	"service-core/internal/modules/authorization/repository"
	transaction "service-core/internal/shared/transaction"

	"github.com/google/uuid"
)

type ActorService struct {
	accountRepo    authenRepo.AccountRepository
	membershipRepo repository.StaffMembershipRepository
}

func NewActorService(
	accountRepo authenRepo.AccountRepository,
	membershipRepo repository.StaffMembershipRepository,
) repository.ActorService {
	return &ActorService{
		accountRepo:    accountRepo,
		membershipRepo: membershipRepo,
	}
}

func (s *ActorService) Load(
	ctx context.Context,
	exec transaction.Executor,
	userID uuid.UUID,
	staffID *uuid.UUID,
) (*domain.Actor, error) {
	account, err := s.accountRepo.
		GetByUserID(ctx, exec, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account: %w", err)
	}

	actor := domain.Actor{
		AccountID: account.ID,
		Type:      account.Type,
	}

	if actor.Type == authenDomain.AccountTypeStaff {
		roles, err := s.membershipRepo.
			ListRolesByAccountIDAndStaffID(
				ctx,
				exec,
				account.ID,
				*staffID,
			)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve staff roles: %w", err)
		}

		actor.Roles = roles
	}

	return &actor, nil
}
