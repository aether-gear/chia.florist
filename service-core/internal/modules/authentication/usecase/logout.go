package usecase

import (
	"context"
	"fmt"

	applogger "service-core/internal/common/logger"
	"service-core/internal/modules/authentication/domain"
	"service-core/internal/modules/authentication/repository"
	transaction "service-core/internal/shared/transaction"
)

type LogoutUsecase struct {
	transactor     transaction.Transactor
	refreshTknRepo repository.RefreshTokenRepository
	sessionRepo    repository.SessionRepository
	auditLogger    applogger.AuditLogger
}

func NewLogoutUsecase(
	transactor transaction.Transactor,
	refreshTknRepo repository.RefreshTokenRepository,
	sessionRepo repository.SessionRepository,
	auditLogger applogger.AuditLogger,
) *LogoutUsecase {
	return &LogoutUsecase{
		transactor:     transactor,
		refreshTknRepo: refreshTknRepo,
		sessionRepo:    sessionRepo,
		auditLogger:    auditLogger,
	}
}

func (u *LogoutUsecase) Execute(
	ctx context.Context,
	authCtx domain.AuthContext,
) (err error) {
	audit := &applogger.AuditScope{
		Category:   "user_action",
		Action:     "logout",
		Resource:   "session",
		ResourceID: authCtx.SessionID.String(),
		Metadata:   map[string]any{"user_id": authCtx.UserID.String()},
	}
	defer applogger.TrackAudit(ctx, u.auditLogger, nil, audit, &err)()

	return u.transactor.WithinTransaction(
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
}
