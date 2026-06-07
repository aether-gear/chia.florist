package service

import (
	"fmt"
	"net/http"
	"time"

	apperrors "service-core/internal/common/errors"
	apphttp "service-core/internal/common/http"
	appcookie "service-core/internal/common/http/cookie"
	commonmiddleware "service-core/internal/common/middleware"
	"service-core/internal/modules/authentication/domain"
	"service-core/internal/modules/authentication/repository"
	transaction "service-core/internal/shared/transaction"
)

type jwtAuthenticator struct {
	tokenSvc    repository.TokenService
	sessionRepo repository.SessionRepository
}

func NewJWTAuthenticator(
	tokenSvc repository.TokenService,
	sessionRepo repository.SessionRepository,
) repository.Authenticator {
	return &jwtAuthenticator{
		tokenSvc:    tokenSvc,
		sessionRepo: sessionRepo,
	}
}

func (aM *jwtAuthenticator) RequireAuth(
	exec transaction.Executor,
) commonmiddleware.Middleware {
	return func(next apphttp.AppHandler) apphttp.AppHandler {
		return func(w http.ResponseWriter, r *http.Request) error {
			token, err := appcookie.CookieValue(r, appcookie.AccessTokenCookieName)
			if err != nil {
				return apperrors.NewUnauthorized(domain.ErrAuthenticationRequired.Error())
			}

			claims, err := aM.tokenSvc.Validate(token)
			if err != nil {
				return apperrors.NewUnauthorized(domain.ErrInvalidToken.Error())
			}
			if claims.Type != domain.TokenTypeAccess {
				return apperrors.NewUnauthorized(domain.ErrInvalidToken.Error())
			}

			session, err := aM.sessionRepo.GetByID(r.Context(), exec, claims.SessionID)
			if err != nil {
				return fmt.Errorf("failed to load session: %w", err)
			}
			if session == nil ||
				session.UserID != claims.UserID ||
				session.RevokedAt != nil ||
				session.ExpiresAt.Before(time.Now()) {
				return apperrors.NewUnauthorized(domain.ErrInvalidSession.Error())
			}

			authCtx := domain.AuthContext{
				UserID:          claims.UserID,
				SessionID:       claims.SessionID,
				TokenType:       claims.Type,
				IsAuthenticated: true,
			}

			r = r.WithContext(domain.WithAuthContext(r.Context(), &authCtx))

			return next(w, r)
		}
	}
}
