package service

import (
	"context"
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
	cookie string,
) commonmiddleware.Middleware {
	return func(next apphttp.AppHandler) apphttp.AppHandler {
		return func(w http.ResponseWriter, r *http.Request) error {
			token, err := appcookie.CookieValue(r, cookie)
			if err != nil {
				return apperrors.
					NewUnauthorized(domain.ErrAuthenticationRequired.Error())
			}

			authCtx, err := aM.authenticate(
				r.Context(),
				exec,
				token,
			)
			if err != nil {
				return err
			}

			r = r.WithContext(
				domain.WithAuthContext(
					r.Context(),
					authCtx,
				),
			)

			return next(w, r)
		}
	}
}

func (aM *jwtAuthenticator) RequireAnyAuth(
	exec transaction.Executor,
	cookies ...string,
) commonmiddleware.Middleware {
	return func(next apphttp.AppHandler) apphttp.AppHandler {
		return func(w http.ResponseWriter, r *http.Request) error {
			for _, cookie := range cookies {
				token, err := appcookie.CookieValue(r, cookie)
				if err != nil {
					continue
				}

				authCtx, err := aM.authenticate(
					r.Context(),
					exec,
					token,
				)
				if err != nil {
					continue
				}

				r = r.WithContext(
					domain.WithAuthContext(
						r.Context(),
						authCtx,
					),
				)

				return next(w, r)
			}

			return apperrors.
				NewUnauthorized(domain.ErrAuthenticationRequired.Error())
		}
	}
}

func (aM *jwtAuthenticator) authenticate(
	ctx context.Context,
	exec transaction.Executor,
	token string,
) (*domain.AuthContext, error) {
	claims, err := aM.tokenSvc.Validate(token)
	if err != nil {
		return nil, apperrors.NewUnauthorized(domain.ErrInvalidToken.Error())
	}

	if claims.Type != domain.TokenTypeAccess {
		return nil, apperrors.NewUnauthorized(domain.ErrInvalidToken.Error())
	}

	session, err := aM.sessionRepo.GetByID(
		ctx,
		exec,
		claims.SessionID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load session: %w", err)
	}

	if session == nil ||
		session.UserID != claims.UserID ||
		session.RevokedAt != nil ||
		session.ExpiresAt.Before(time.Now()) {
		return nil, apperrors.NewUnauthorized(domain.ErrInvalidSession.Error())
	}

	return &domain.AuthContext{
		UserID:          claims.UserID,
		SessionID:       claims.SessionID,
		TokenType:       claims.Type,
		IsAuthenticated: true,
		MerchantID:      claims.MerchantID,
		Roles:           claims.Roles,
	}, nil
}
