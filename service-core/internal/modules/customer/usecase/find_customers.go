package usecase

import (
	"context"
	"fmt"
	"strings"

	apperrors "service-core/internal/common/errors"
	"service-core/internal/modules/customer/domain"
	"service-core/internal/modules/customer/repository"
	query "service-core/internal/shared/query"
	transaction "service-core/internal/shared/transaction"

	"github.com/google/uuid"
)

type FindCustomersUsecase struct {
	executor     transaction.Executor
	customerRepo repository.CustomerRepository
}

func NewFindCustomersUsecase(
	executor transaction.Executor,
	customerRepo repository.CustomerRepository,
) *FindCustomersUsecase {
	return &FindCustomersUsecase{
		executor:     executor,
		customerRepo: customerRepo,
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
) ([]domain.Customer, int, error) {
	var userSortKeys = map[string]query.SortKey{
		"latest":     repository.CustomerSortLatest,
		"name":       repository.CustomerSortName,
		"username":   repository.CustomerSortUsername,
		"phone":      repository.CustomerSortPhone,
		"modify":     repository.CustomerSortModify,
		"last_login": repository.CustomerSortLastLogin,
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
				By:        repository.CustomerSortLatest,
				Direction: query.SortDesc,
			},
		}
	}

	params := repository.FindCustomerParams{
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

	users, total, err := u.customerRepo.
		FindCustomers(ctx, u.executor, params)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to load users: %w", err)
	}
	if len(users) == 0 {
		return nil, 0, apperrors.NewNotFound("users not available at the moment")
	}

	return users, total, nil
}
