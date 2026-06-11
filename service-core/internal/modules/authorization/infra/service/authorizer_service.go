package service

import (
	"context"
	"net/http"

	apperrors "service-core/internal/common/errors"
	apphttp "service-core/internal/common/http"
	appmiddleware "service-core/internal/common/middleware"
	authendomain "service-core/internal/modules/authentication/domain"
	"service-core/internal/modules/authorization/domain"
	"service-core/internal/modules/authorization/repository"
	"service-core/internal/shared/transaction"
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

func (s *authorizer) RequireAccountType(allowedTypes ...authendomain.AccountType) appmiddleware.Middleware {
	return appmiddleware.RequireAccountType(
		func(ctx context.Context) (authendomain.AccountType, bool) {
			actor, ok := GetActor(ctx)
			if !ok {
				return "", false
			}
			return actor.Type, true
		},
		allowedTypes...,
	)
}

func (s *authorizer) RequireMerchantRole(allowedRoles ...string) appmiddleware.Middleware {
	return func(next apphttp.AppHandler) apphttp.AppHandler {
		return func(w http.ResponseWriter, r *http.Request) error {
			actor, ok := GetActor(r.Context())
			if !ok {
				return apperrors.NewUnauthorized("authentication required")
			}

			if actor.Type != authendomain.AccountTypeMerchant {
				return apperrors.NewForbidden(domain.ErrMerchantRequired.Error())
			}

			allowed := false
			for _, actorRole := range actor.Roles {
				for _, requiredRole := range allowedRoles {
					if actorRole.Code == requiredRole {
						allowed = true
						break
					}
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

			actor, err := s.actorSvc.Load(r.Context(), exec,
				authCtx.UserID,
				authCtx.MerchantID,
			)
			if err != nil {
				return err
			}

			ctx := context.WithValue(r.Context(), actorContextKey{}, actor)

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
