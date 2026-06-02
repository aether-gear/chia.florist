package middleware

import (
	"context"
	"net/http"

	apperrors "service-core/internal/common/errors"
	apphttp "service-core/internal/common/http"
	authendomain "service-core/internal/modules/authentication/domain"
)

func RequireAccountType(
	extractAccountType func(ctx context.Context) (authendomain.AccountType, bool),
	allowedTypes ...authendomain.AccountType,
) Middleware {
	return func(next apphttp.AppHandler) apphttp.AppHandler {
		return func(w http.ResponseWriter, r *http.Request) error {
			accountType, ok := extractAccountType(r.Context())
			if !ok {
				return apperrors.NewUnauthorized("actor unauthorized")
			}

			isAllowed := false
			for _, allowed := range allowedTypes {
				if accountType == allowed {
					isAllowed = true
					break
				}
			}

			if !isAllowed {
				return apperrors.NewForbidden("insufficient account type")
			}

			return next(w, r)
		}
	}
}
