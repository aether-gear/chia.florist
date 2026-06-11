package usecase

import (
	"context"
	"fmt"

	apperrors "service-core/internal/common/errors"
	"service-core/internal/modules/authentication/domain"
	"service-core/internal/modules/authentication/repository"
	authorDomain "service-core/internal/modules/authorization/domain"
	authorRepo "service-core/internal/modules/authorization/repository"
	transaction "service-core/internal/shared/transaction"
)

type MeUsecase struct {
	exec        transaction.Executor
	accountRepo repository.AccountRepository
	actorSvc    authorRepo.ActorService
}

func NewMeUsecase(
	exec transaction.Executor,
	accountRepo repository.AccountRepository,
	actorSvc authorRepo.ActorService,
) *MeUsecase {
	return &MeUsecase{
		exec:        exec,
		accountRepo: accountRepo,
		actorSvc:    actorSvc,
	}
}

type MeResult struct {
	Account domain.Account
	Actor   authorDomain.Actor
}

func (u *MeUsecase) Execute(
	ctx context.Context,
	authCtx domain.AuthContext,
) (*MeResult, error) {
	account, err := u.accountRepo.GetByUserID(
		ctx,
		u.exec,
		authCtx.UserID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account: %w", err)
	}

	if account == nil {
		return nil, apperrors.NewNotFound("account not found")
	}
	if account.Status != domain.AccountActive {
		return nil, apperrors.NewForbidden(domain.ErrEmailNotVerified.Error())
	}

	actor, err := u.actorSvc.Load(
		ctx,
		u.exec,
		authCtx.UserID,
		authCtx.MerchantID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve actor: %w", err)
	}

	return &MeResult{
		*account,
		*actor,
	}, nil
}
