package service

import (
	"context"
	"fmt"
	"net/http"
	"time"

	appclock "service-core/internal/common/clock"
	apperrors "service-core/internal/common/errors"
	apphttp "service-core/internal/common/http"
	appcookie "service-core/internal/common/http/cookie"
	applogger "service-core/internal/common/logger"
	commonmiddleware "service-core/internal/common/middleware"
	"service-core/internal/modules/authentication/domain"
	"service-core/internal/modules/authentication/repository"
	authorzDomain "service-core/internal/modules/authorization/domain"
	transaction "service-core/internal/shared/transaction"
)

type jwtAuthenticator struct {
	tokenSvc         repository.TokenService
	sessionRepo      repository.SessionRepository
	tokenHasher      repository.TokenHasher
	refreshTokenRepo repository.RefreshTokenRepository
}

func NewJWTAuthenticator(
	tokenSvc repository.TokenService,
	sessionRepo repository.SessionRepository,
	tokenHasher repository.TokenHasher,
	refreshTokenRepo repository.RefreshTokenRepository,
) repository.Authenticator {
	return &jwtAuthenticator{
		tokenSvc:         tokenSvc,
		sessionRepo:      sessionRepo,
		tokenHasher:      tokenHasher,
		refreshTokenRepo: refreshTokenRepo,
	}
}

func (aM *jwtAuthenticator) RequireAuth(
	exec transaction.Executor,
	tran transaction.Transactor,
	cookie appcookie.CookieName,
) commonmiddleware.Middleware {
	return func(next apphttp.AppHandler) apphttp.AppHandler {
		return func(w http.ResponseWriter, r *http.Request) error {
			var authCtx *domain.AuthContext

			token, err := appcookie.Extract(r, cookie)
			if err == nil {
				authCtx, err = aM.
					authenticate(r.Context(), exec, token)
			}

			if err != nil {
				var refreshErr error
				authCtx, refreshErr = aM.
					trySilentRefresh(
						r.Context(),
						exec,
						tran,
						w,
						r,
						cookie,
					)
				if refreshErr != nil {
					return apperrors.NewUnauthorized(domain.ErrAuthenticationRequired.Error())
				}
			}

			if !isValidCookieForAuth(cookie, authCtx) {
				return apperrors.NewUnauthorized(domain.ErrAuthenticationRequired.Error())
			}

			ctx := domain.WithAuthContext(
				r.Context(),
				authCtx,
			)
			ctx = applogger.WithActorID(ctx, authCtx.UserID.String())

			return next(w, r.WithContext(ctx))
		}
	}
}

func (aM *jwtAuthenticator) RequireAnyAuth(
	exec transaction.Executor,
	tran transaction.Transactor,
	cookies ...appcookie.CookieName,
) commonmiddleware.Middleware {
	return func(next apphttp.AppHandler) apphttp.AppHandler {
		return func(w http.ResponseWriter, r *http.Request) error {
			for _, cookie := range cookies {
				var authCtx *domain.AuthContext

				token, err := appcookie.Extract(r, cookie)
				if err == nil {
					authCtx, err = aM.
						authenticate(r.Context(), exec, token)
				}

				if err != nil {
					var refreshErr error
					authCtx, refreshErr = aM.
						trySilentRefresh(
							r.Context(),
							exec,
							tran,
							w,
							r,
							cookie,
						)
					if refreshErr != nil {
						continue
					}
				}

				if !isValidCookieForAuth(cookie, authCtx) {
					continue
				}

				ctx := domain.WithAuthContext(
					r.Context(),
					authCtx,
				)
				ctx = applogger.WithActorID(ctx, authCtx.UserID.String())

				return next(w, r.WithContext(ctx))
			}

			return apperrors.NewUnauthorized(domain.ErrAuthenticationRequired.Error())
		}
	}
}

func (aM *jwtAuthenticator) RequireMultiAuth(
	exec transaction.Executor,
	tran transaction.Transactor,
	cookies ...appcookie.CookieName,
) commonmiddleware.Middleware {
	return func(next apphttp.AppHandler) apphttp.AppHandler {
		return func(w http.ResponseWriter, r *http.Request) error {
			var collected []*domain.AuthContext

			for _, cookie := range cookies {
				var authCtx *domain.AuthContext

				token, err := appcookie.Extract(r, cookie)
				if err == nil {
					authCtx, err = aM.
						authenticate(r.Context(), exec, token)
				}

				if err != nil {
					var refreshErr error
					authCtx, refreshErr = aM.
						trySilentRefresh(
							r.Context(),
							exec,
							tran,
							w,
							r,
							cookie,
						)
					if refreshErr != nil {
						continue
					}
				}

				if !isValidCookieForAuth(cookie, authCtx) {
					continue
				}

				collected = append(collected, authCtx)
			}

			if len(collected) == 0 {
				return apperrors.NewUnauthorized(domain.ErrAuthenticationRequired.Error())
			}

			ctx := domain.WithMultiAuthContext(r.Context(), collected)
			ctx = domain.WithAuthContext(ctx, collected[0])
			ctx = applogger.WithActorID(ctx, collected[0].UserID.String())

			return next(w, r.WithContext(ctx))
		}
	}
}

func (aM *jwtAuthenticator) OptionalAuth(
	exec transaction.Executor,
	tran transaction.Transactor,
	cookies ...appcookie.CookieName,
) commonmiddleware.Middleware {
	return func(next apphttp.AppHandler) apphttp.AppHandler {
		return func(w http.ResponseWriter, r *http.Request) error {
			var collected []*domain.AuthContext

			for _, cookie := range cookies {
				var authCtx *domain.AuthContext

				token, err := appcookie.Extract(r, cookie)
				if err == nil {
					authCtx, err = aM.authenticate(r.Context(), exec, token)
				}

				if err != nil {
					var refreshErr error
					authCtx, refreshErr = aM.trySilentRefresh(
						r.Context(),
						exec,
						tran,
						w,
						r,
						cookie,
					)
					if refreshErr != nil {
						continue
					}
				}

				if !isValidCookieForAuth(cookie, authCtx) {
					continue
				}

				collected = append(collected, authCtx)
			}

			if len(collected) > 0 {
				ctx := domain.WithMultiAuthContext(r.Context(), collected)
				ctx = domain.WithAuthContext(ctx, collected[0])
				ctx = applogger.WithActorID(ctx, collected[0].UserID.String())
				return next(w, r.WithContext(ctx))
			}

			return next(w, r)
		}
	}
}

func isValidCookieForAuth(
	cookie appcookie.CookieName,
	authCtx *domain.AuthContext,
) bool {
	if authCtx == nil {
		return false
	}
	switch cookie {
	case appcookie.CookieCustomer:
		return authCtx.CustomerID != nil
	case appcookie.CookieStaff:
		return authCtx.StaffID != nil
	default:
		return true
	}
}

func (aM *jwtAuthenticator) trySilentRefresh(
	ctx context.Context,
	exec transaction.Executor,
	tran transaction.Transactor,
	w http.ResponseWriter,
	r *http.Request,
	cookieName appcookie.CookieName,
) (*domain.AuthContext, error) {
	var refreshCookie appcookie.CookieName
	switch cookieName {
	case appcookie.CookieCustomer:
		refreshCookie = appcookie.CookieCustomerRefresh
	case appcookie.CookieStaff:
		refreshCookie = appcookie.CookieStaffRefresh
	default:
		return nil, fmt.Errorf("unknown cookie type: %s", cookieName)
	}

	refreshTokenStr, err := appcookie.
		Extract(r, refreshCookie)
	if err != nil {
		return nil, fmt.Errorf("refresh token not found: %w", err)
	}

	claims, err := aM.tokenSvc.
		Validate(refreshTokenStr)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	if claims.Type != domain.TokenTypeRefresh {
		return nil, fmt.Errorf("invalid token type for refresh: %s", claims.Type)
	}

	session, err := aM.sessionRepo.
		GetByID(ctx, exec, claims.SessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	if session == nil ||
		session.RevokedAt != nil ||
		session.ExpiresAt.Before(appclock.Now()) {

		return nil, fmt.Errorf("session is invalid, revoked or expired")
	}

	dbRefreshToken, err := aM.refreshTokenRepo.
		GetBySessionID(
			ctx,
			exec,
			claims.SessionID,
		)
	if err != nil {
		return nil, fmt.Errorf("failed to get refresh token: %w", err)
	}

	if dbRefreshToken == nil ||
		dbRefreshToken.RevokedAt != nil ||
		dbRefreshToken.ExpiresAt.Before(appclock.Now()) {

		return nil, fmt.Errorf("refresh token in db is invalid, revoked or expired")
	}

	if !aM.tokenHasher.Compare(
		dbRefreshToken.TokenHash,
		refreshTokenStr,
	) {
		return nil, fmt.Errorf("refresh token hash mismatch")
	}

	roleCodes := make([]authorzDomain.RoleCode, len(claims.Roles))
	for i, rCode := range claims.Roles {
		roleCodes[i] = authorzDomain.RoleCode(rCode)
	}

	newAccessTkn, err := aM.tokenSvc.
		Generate(repository.GenerateTokenParams{
			UserID:     claims.UserID,
			SessionID:  claims.SessionID,
			StaffID:    claims.StaffID,
			CustomerID: claims.CustomerID,
			Roles:      roleCodes,
			Type:       domain.TokenTypeAccess,
			Duration:   30 * time.Minute,
		})
	if err != nil {
		return nil, fmt.Errorf("failed to generate new access token: %w", err)
	}

	newRefreshTkn, err := aM.tokenSvc.
		Generate(repository.GenerateTokenParams{
			UserID:     claims.UserID,
			SessionID:  claims.SessionID,
			StaffID:    claims.StaffID,
			CustomerID: claims.CustomerID,
			Roles:      roleCodes,
			Type:       domain.TokenTypeRefresh,
			Duration:   7 * 24 * time.Hour,
		})
	if err != nil {
		return nil, fmt.Errorf("failed to generate new refresh token: %w", err)
	}

	now := appclock.Now()
	session.ExpiresAt = now.Add(7 * 24 * time.Hour)
	session.LastActivityAt = &now

	newRefreshTknHashed := aM.tokenHasher.Hash(newRefreshTkn.Token)
	dbRefreshToken.TokenHash = newRefreshTknHashed
	dbRefreshToken.ExpiresAt = now.Add(7 * 24 * time.Hour)

	if err := tran.WithinTransaction(ctx, func(e transaction.Executor) error {
		if err := aM.sessionRepo.
			Save(ctx, e, *session); err != nil {
			return fmt.Errorf("failed to save updated session: %w", err)
		}
		if err := aM.refreshTokenRepo.
			Save(ctx, e, *dbRefreshToken); err != nil {
			return fmt.Errorf("failed to save updated refresh token: %w", err)
		}
		return nil
	}); err != nil {
		return nil, err
	}

	appcookie.Bind(w, cookieName, newAccessTkn.Token, newAccessTkn.ExpiresAt)
	appcookie.Bind(w, refreshCookie, newRefreshTkn.Token, newRefreshTkn.ExpiresAt)

	result := domain.AuthContext{
		UserID:          claims.UserID,
		SessionID:       claims.SessionID,
		TokenType:       domain.TokenTypeAccess,
		IsAuthenticated: true,
		StaffID:         claims.StaffID,
		CustomerID:      claims.CustomerID,
		Roles:           claims.Roles,
	}

	return &result, nil
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

	session, err := aM.sessionRepo.
		GetByID(
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
		session.ExpiresAt.Before(appclock.Now()) {

		return nil, apperrors.NewUnauthorized(domain.ErrInvalidSession.Error())
	}

	return &domain.AuthContext{
		UserID:          claims.UserID,
		SessionID:       claims.SessionID,
		TokenType:       claims.Type,
		IsAuthenticated: true,
		StaffID:         claims.StaffID,
		CustomerID:      claims.CustomerID,
		Roles:           claims.Roles,
	}, nil
}
