package usecase

import (
	"context"
	"fmt"

	apperrors "service-core/internal/common/errors"
	"service-core/internal/modules/security_policy/domain"
	"service-core/internal/modules/security_policy/repository"
	transaction "service-core/internal/shared/transaction"
)

type UpdateIPActionUsecase struct {
	executor     transaction.Executor
	securityRepo repository.SecurityPolicyRepository
}

func NewUpdateIPActionUsecase(
	executor transaction.Executor,
	securityRepo repository.SecurityPolicyRepository,
) *UpdateIPActionUsecase {
	return &UpdateIPActionUsecase{
		executor:     executor,
		securityRepo: securityRepo,
	}
}

type UpdateIPActionInput struct {
	IP     string
	Action string // only "ban", "whitelist", "ignore", "reset" are allowed
	Reason string
}

func (u *UpdateIPActionUsecase) Execute(ctx context.Context, input UpdateIPActionInput) error {
	switch input.Action {
	case "ban":
		record := domain.IPRecord{
			IP:     input.IP,
			Status: domain.IPStatusBanned,
			Reason: input.Reason,
		}
		if err := u.securityRepo.
			UpsertIPRecord(ctx, u.executor, record); err != nil {
			return fmt.Errorf("update ip action (ban): %w", err)
		}

	case "whitelist":
		record := domain.IPRecord{
			IP:     input.IP,
			Status: domain.IPStatusWhitelisted,
		}
		if err := u.securityRepo.
			UpsertIPRecord(ctx, u.executor, record); err != nil {
			return fmt.Errorf("update ip action (whitelist): %w", err)
		}

	case "ignore":
		record := domain.IPRecord{
			IP:     input.IP,
			Status: domain.IPStatusIgnored,
		}
		if err := u.securityRepo.
			UpsertIPRecord(ctx, u.executor, record); err != nil {
			return fmt.Errorf("update ip action (ignore): %w", err)
		}

	case "reset":
		if err := u.securityRepo.
			DeleteIPRecord(ctx, u.executor, input.IP); err != nil {
			return fmt.Errorf("update ip action (reset): %w", err)
		}

	default:
		return apperrors.NewBadRequest(
			fmt.Sprintf("unknown ip action %q: must be one of ban, whitelist, ignore, reset", input.Action),
		)
	}

	return nil
}
