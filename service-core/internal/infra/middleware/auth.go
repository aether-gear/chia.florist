package middleware

import (
	"fmt"
	"net/http"
	"time"

	apperrors "service-core/internal/common/errors"
	apphttp "service-core/internal/common/http"
	appcookie "service-core/internal/common/http/cookie"
	commonmiddleware "service-core/internal/common/middleware"
	authendomain "service-core/internal/modules/authentication/domain"
	authenrepo "service-core/internal/modules/authentication/repository"
)

func RequireAuth(
	tokenSvc authenrepo.TokenService,
	sessionRepo authenrepo.SessionRepository,
) commonmiddleware.Middleware {
	return func(next apphttp.AppHandler) apphttp.AppHandler {
		return func(w http.ResponseWriter, r *http.Request) error {
			token, err := appcookie.CookieValue(r, appcookie.AccessTokenCookieName)
			if err != nil {
				return apperrors.NewUnauthorized(ErrAuthenticationRequired.Error())
			}

			claims, err := tokenSvc.Validate(token)
			if err != nil {
				return apperrors.NewUnauthorized(ErrInvalidToken.Error())
			}
			if claims.Type != authendomain.TokenTypeAccess {
				return apperrors.NewUnauthorized(ErrInvalidToken.Error())
			}

			session, err := sessionRepo.GetByID(claims.SessionID)
			if err != nil {
				return fmt.Errorf("failed to load session: %w", err)
			}
			if session == nil ||
				session.UserID != claims.UserID ||
				session.RevokedAt != nil ||
				session.ExpiresAt.Before(time.Now()) {
				return apperrors.NewUnauthorized(ErrInvalidSession.Error())
			}

			authCtx := authendomain.AuthContext{
				UserID:          claims.UserID,
				SessionID:       claims.SessionID,
				TokenType:       claims.Type,
				IsAuthenticated: true,
			}

			r = r.WithContext(authendomain.WithAuthContext(r.Context(), &authCtx))

			return next(w, r)
		}
	}
}
