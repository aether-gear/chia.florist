package service

import (
	"context"
	"fmt"

	apperrors "service-core/internal/common/errors"
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
	staffPermRepo  repository.StaffPermissionRepository
}

func NewActorService(
	accountRepo authenRepo.AccountRepository,
	membershipRepo repository.StaffMembershipRepository,
	staffPermRepo repository.StaffPermissionRepository,
) repository.ActorService {
	return &ActorService{
		accountRepo:    accountRepo,
		membershipRepo: membershipRepo,
		staffPermRepo:  staffPermRepo,
	}
}

func (s *ActorService) Load(
	ctx context.Context,
	exec transaction.Executor,
	userID uuid.UUID,
	staffID *uuid.UUID,
) (*domain.Actor, error) {
	account, err := s.accountRepo.GetByUserID(ctx, exec, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account: %w", err)
	}

	actor := domain.Actor{
		AccountID:   account.ID,
		Type:        account.Type,
		StaffID:     staffID,
		Permissions: make(map[uuid.UUID][]string),
		Rules:       make(map[uuid.UUID]map[string]any),
	}

	if actor.Type == authenDomain.AccountTypeStaff {
		if staffID == nil {
			return nil, apperrors.NewUnauthorized("staff account is not associated with a staff profile")
		}

		roles, err := s.membershipRepo.ListRolesByAccountIDAndStaffID(ctx, exec,
			account.ID,
			*staffID,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve staff roles: %w", err)
		}

		actor.Roles = roles

		if s.staffPermRepo != nil {
			shopPerms, err := s.staffPermRepo.ListByStaffID(ctx, exec, *staffID)
			if err == nil {
				for _, sp := range shopPerms {
					actor.Permissions[sp.ShopID] = sp.Permissions
					if sp.Rules != nil {
						actor.Rules[sp.ShopID] = sp.Rules
					}
				}
			}
		}
	}

	return &actor, nil
}
