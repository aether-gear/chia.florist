package usecase

import (
	"context"
	"fmt"

	apperrors "service-core/internal/common/errors"
	"service-core/internal/modules/security_policy/repository"
	transaction "service-core/internal/shared/transaction"
)

type UpdateFilterUsecase struct {
	executor     transaction.Executor
	securityRepo repository.SecurityPolicyRepository
}

func NewUpdateFilterUsecase(
	executor transaction.Executor,
	securityRepo repository.SecurityPolicyRepository,
) *UpdateFilterUsecase {
	return &UpdateFilterUsecase{
		executor:     executor,
		securityRepo: securityRepo,
	}
}

type UpdateFilterInput struct {
	Type   string // must be "keyword" or "url".
	Action string // must be "add" or "remove".
	Value  string
}

func (u *UpdateFilterUsecase) Execute(
	ctx context.Context,
	input UpdateFilterInput,
) error {
	if input.Type != "keyword" && input.Type != "url" {
		return apperrors.NewBadRequest(
			fmt.Sprintf("unknown filter type %q: must be keyword or url", input.Type),
		)
	}

	switch input.Action {
	case "add":
		if err := u.securityRepo.
			UpsertFilterEntry(
				ctx,
				u.executor,
				input.Type,
				input.Value,
			); err != nil {
			return fmt.Errorf("update filter (add): %w", err)
		}

	case "remove":
		if err := u.securityRepo.
			DeleteFilterEntry(
				ctx,
				u.executor,
				input.Type,
				input.Value,
			); err != nil {
			return fmt.Errorf("update filter (remove): %w", err)
		}

	default:
		return apperrors.NewBadRequest(
			fmt.Sprintf("unknown filter action %q: must be add or remove", input.Action),
		)
	}

	return nil
}
