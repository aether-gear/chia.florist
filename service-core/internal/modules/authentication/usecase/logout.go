package usecase

import (
	"context"
	"fmt"

	"service-core/internal/modules/authentication/domain"
	"service-core/internal/modules/authentication/repository"
	transaction "service-core/internal/shared/transaction"
)

type LogoutUsecase struct {
	executor       transaction.Executor
	transactor     transaction.Transactor
	refreshTknRepo repository.RefreshTokenRepository
	sessionRepo    repository.SessionRepository
}

func NewLogoutUsecase(
	executor transaction.Executor,
	transactor transaction.Transactor,
	refreshTknRepo repository.RefreshTokenRepository,
	sessionRepo repository.SessionRepository,
) *LogoutUsecase {
	return &LogoutUsecase{
		executor:       executor,
		transactor:     transactor,
		refreshTknRepo: refreshTknRepo,
		sessionRepo:    sessionRepo,
	}
}

func (u *LogoutUsecase) Execute(
	ctx context.Context,
	authCtx domain.AuthContext,
) error {
	err := u.transactor.WithinTransaction(
		ctx,
		func(e transaction.Executor) error {
			if err := u.refreshTknRepo.RevokeBySessionID(
				ctx,
				u.executor,
				authCtx.SessionID,
			); err != nil {
				return fmt.Errorf("failed to revoke refresh token %w", err)
			}
			if err := u.sessionRepo.RevokeByID(
				ctx,
				u.executor,
				authCtx.SessionID,
			); err != nil {
				return fmt.Errorf("failed to revoke refresh token %w", err)
			}

			return nil
		},
	)
	if err != nil {
		return err
	}

	return nil
}
