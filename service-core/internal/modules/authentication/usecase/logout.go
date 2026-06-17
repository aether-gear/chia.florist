package usecase

import (
	"context"
	"fmt"

	"service-core/internal/modules/authentication/domain"
	"service-core/internal/modules/authentication/repository"
	transaction "service-core/internal/shared/transaction"
)

type LogoutUsecase struct {
	transactor     transaction.Transactor
	refreshTknRepo repository.RefreshTokenRepository
	sessionRepo    repository.SessionRepository
}

func NewLogoutUsecase(
	transactor transaction.Transactor,
	refreshTknRepo repository.RefreshTokenRepository,
	sessionRepo repository.SessionRepository,
) *LogoutUsecase {
	return &LogoutUsecase{
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
		func(exec transaction.Executor) error {
			if err := u.refreshTknRepo.
				RevokeBySessionID(
					ctx,
					exec,
					authCtx.SessionID,
				); err != nil {
				return fmt.Errorf("failed to revoke refresh token %w", err)
			}

			if err := u.sessionRepo.
				RevokeByID(
					ctx,
					exec,
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
