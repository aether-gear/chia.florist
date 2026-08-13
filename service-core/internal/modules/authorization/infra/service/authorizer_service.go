package service

import (
	"context"
	"net/http"
	"slices"

	apperrors "service-core/internal/common/errors"
	apphttp "service-core/internal/common/http"
	appmiddleware "service-core/internal/common/middleware"
	authendomain "service-core/internal/modules/authentication/domain"
	"service-core/internal/modules/authorization/domain"
	"service-core/internal/modules/authorization/repository"
	transaction "service-core/internal/shared/transaction"
)

type actorContextKey struct{}

type authorizer struct {
	actorSvc repository.ActorService
}

func NewAuthorizer(
	actorSvc repository.ActorService,
) repository.Authorizer {
	return &authorizer{
		actorSvc: actorSvc,
	}
}

func (s *authorizer) RequireAccountType(
	allowedTypes ...authendomain.AccountType,
) appmiddleware.Middleware {
	return func(next apphttp.AppHandler) apphttp.AppHandler {
		return func(w http.ResponseWriter, r *http.Request) error {
			if multiCtxs, ok := authendomain.GetMultiAuthContext(r.Context()); ok && len(multiCtxs) > 1 {
				headerHint := apphttp.GetAccountTypeHeader(r)

				var chosenCandidate *authendomain.AuthContext
				if headerHint != "" {
					for _, candidate := range multiCtxs {
						candidateType := authendomain.AccountTypeCustomer
						if candidate.StaffID != nil {
							candidateType = authendomain.AccountTypeStaff
						}
						if string(candidateType) == headerHint && slices.Contains(allowedTypes, candidateType) {
							chosenCandidate = candidate
							break
						}
					}
				}

				if chosenCandidate == nil {
					for _, candidate := range multiCtxs {
						candidateType := authendomain.AccountTypeCustomer
						if candidate.StaffID != nil {
							candidateType = authendomain.AccountTypeStaff
						}
						if slices.Contains(allowedTypes, candidateType) {
							chosenCandidate = candidate
							break
						}
					}
				}

				if chosenCandidate != nil {
					r = r.WithContext(authendomain.WithAuthContext(r.Context(), chosenCandidate))
				}
			}

			actor, ok := GetActor(r.Context())
			if !ok {
				authCtx, authOk := authendomain.GetAuthContext(r.Context())
				if !authOk {
					return apperrors.NewUnauthorized("authentication required")
				}

				accType := authendomain.AccountTypeCustomer
				if authCtx.StaffID != nil {
					accType = authendomain.AccountTypeStaff
				}
				if slices.Contains(allowedTypes, accType) {
					return next(w, r)
				}

				return apperrors.NewForbidden("insufficient account type")
			}

			if slices.Contains(allowedTypes, actor.Type) {
				return next(w, r)
			}

			return apperrors.NewForbidden("insufficient account type")
		}
	}
}

func (s *authorizer) RequireStaffRole(
	allowedRoles ...domain.RoleCode,
) appmiddleware.Middleware {
	return func(next apphttp.AppHandler) apphttp.AppHandler {
		return func(w http.ResponseWriter, r *http.Request) error {
			actor, ok := GetActor(r.Context())
			if !ok {
				return apperrors.NewUnauthorized("authentication required")
			}
			if actor.Type != authendomain.AccountTypeStaff {
				return apperrors.NewForbidden(domain.ErrStaffRequired.Error())
			}

			allowed := false
			for _, actorRole := range actor.Roles {
				if slices.Contains(allowedRoles, actorRole.Code) {
					allowed = true
				}
				if allowed {
					break
				}
			}

			if !allowed {
				return apperrors.NewForbidden(domain.ErrInsufficientRole.Error())
			}

			return next(w, r)
		}
	}
}

func (s *authorizer) LoadActor(
	exec transaction.Executor,
) appmiddleware.Middleware {
	return func(next apphttp.AppHandler) apphttp.AppHandler {
		return func(w http.ResponseWriter, r *http.Request) error {
			authCtx, ok := authendomain.GetAuthContext(r.Context())
			if !ok {
				return apperrors.NewUnauthorized("authentication required")
			}

			actor, err := s.actorSvc.
				Load(
					r.Context(),
					exec,
					authCtx.UserID,
					authCtx.StaffID,
				)
			if err != nil {
				return err
			}

			ctx := context.WithValue(
				r.Context(),
				actorContextKey{},
				actor,
			)

			return next(w, r.WithContext(ctx))
		}
	}
}

func WithActor(ctx context.Context, actor *domain.Actor) context.Context {
	return context.WithValue(ctx, actorContextKey{}, actor)
}

func GetActor(ctx context.Context) (*domain.Actor, bool) {
	actor, ok := ctx.Value(actorContextKey{}).(*domain.Actor)
	return actor, ok
}
