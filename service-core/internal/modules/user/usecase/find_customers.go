package usecase

import (
	"context"
	"fmt"
	"strings"

	apperrors "service-core/internal/common/errors"
	"service-core/internal/modules/user/domain"
	"service-core/internal/modules/user/repository"
	query "service-core/internal/shared/query"
	transaction "service-core/internal/shared/transaction"

	"github.com/google/uuid"
)

type FindCustomersUsecase struct {
	executor transaction.Executor
	userRepo repository.UserRepository
}

func NewFindCustomersUsecase(
	executor transaction.Executor,
	userRepo repository.UserRepository,
) *FindCustomersUsecase {
	return &FindCustomersUsecase{
		executor: executor,
		userRepo: userRepo,
	}
}

type FindCustomersInput struct {
	Page     int
	Limit    int
	ID       *uuid.UUID
	Name     *string
	Username *string
	Email    *string
	Sort     string
}

func (u *FindCustomersUsecase) Execute(
	ctx context.Context,
	input FindCustomersInput,
) ([]domain.User, int, error) {
	var userSortKeys = map[string]query.SortKey{
		"latest":     repository.UserSortLatest,
		"name":       repository.UserSortName,
		"username":   repository.UserSortUsername,
		"phone":      repository.UserSortPhone,
		"modify":     repository.UserSortModify,
		"last_login": repository.UserSortLastLogin,
	}

	var sorts query.Sorts
	if input.Sort != "" {
		parts := strings.SplitSeq(input.Sort, ",")
		for part := range parts {
			part = strings.TrimSpace(part)
			if part == "" {
				continue
			}

			subparts := strings.Split(part, ":")
			key := strings.TrimSpace(subparts[0])

			var dir query.SortDirection = query.SortDesc
			if len(subparts) > 1 {
				d := strings.ToLower(strings.TrimSpace(subparts[1]))
				if d == "asc" {
					dir = query.SortAsc
				}
			}

			sortKey, exists := userSortKeys[key]
			if exists {
				sorts = append(sorts, query.Sort{
					By:        sortKey,
					Direction: dir,
				})
			}
		}
	}

	if len(sorts) == 0 {
		sorts = query.Sorts{
			{
				By:        repository.UserSortLatest,
				Direction: query.SortDesc,
			},
		}
	}

	params := repository.FindUserParams{
		ID:       input.ID,
		Name:     input.Name,
		Username: input.Username,
		Email:    input.Email,
		Pagination: query.Pagination{
			Page:  input.Page,
			Limit: input.Limit,
		},
		Sorts: sorts,
	}

	users, total, err := u.userRepo.
		FindCustomers(ctx, u.executor, params)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to load users: %w", err)
	}
	if len(users) == 0 {
		return nil, 0, apperrors.NewNotFound("users not available at the moment")
	}

	return users, total, nil
}
