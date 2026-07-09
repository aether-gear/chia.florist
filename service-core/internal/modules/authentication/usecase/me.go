package usecase

import (
	"context"
	"fmt"

	apperrors "service-core/internal/common/errors"
	"service-core/internal/modules/authentication/domain"
	"service-core/internal/modules/authentication/repository"
	authorDomain "service-core/internal/modules/authorization/domain"
	authorRepo "service-core/internal/modules/authorization/repository"
	userDomain "service-core/internal/modules/user/domain"
	userRepo "service-core/internal/modules/user/repository"
	transaction "service-core/internal/shared/transaction"
)

type MeUsecase struct {
	exec        transaction.Executor
	accountRepo repository.AccountRepository
	userRepo    userRepo.UserRepository
	actorSvc    authorRepo.ActorService
	oauthRepo   repository.OAuthConnectionRepository
}

func NewMeUsecase(
	exec transaction.Executor,
	accountRepo repository.AccountRepository,
	userRepo userRepo.UserRepository,
	actorSvc authorRepo.ActorService,
	oauthRepo repository.OAuthConnectionRepository,
) *MeUsecase {
	return &MeUsecase{
		exec:        exec,
		accountRepo: accountRepo,
		userRepo:    userRepo,
		actorSvc:    actorSvc,
		oauthRepo:   oauthRepo,
	}
}

type MeResult struct {
	Account domain.Account
	Actor   authorDomain.Actor
	User    *userDomain.User
	OAuth   *domain.OAuthConnection
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
		authCtx.StaffID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve actor: %w", err)
	}

	user, err := u.userRepo.
		GetByID(ctx, u.exec, authCtx.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve user: %w", err)
	}

	var oauthConn *domain.OAuthConnection
	if u.oauthRepo != nil {
		oauthConn, err = u.oauthRepo.
			GetByUserID(ctx, u.exec, authCtx.UserID)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve oauth connection: %w", err)
		}
	}

	return &MeResult{
		Account: *account,
		Actor:   *actor,
		User:    user,
		OAuth:   oauthConn,
	}, nil
}
