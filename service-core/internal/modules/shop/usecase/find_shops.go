package usecase

import (
	"context"
	"fmt"
	"strings"

	apperrors "service-core/internal/common/errors"
	"service-core/internal/modules/shop/domain"
	"service-core/internal/modules/shop/repository"
	query "service-core/internal/shared/query"
	transaction "service-core/internal/shared/transaction"
)

type FindShopsUsecase struct {
	executor transaction.Executor
	shopRepo repository.ShopRepository
}

func NewFindShopsUsecase(
	executor transaction.Executor,
	shopRepo repository.ShopRepository,
) *FindShopsUsecase {
	return &FindShopsUsecase{
		executor: executor,
		shopRepo: shopRepo,
	}
}

type FindShopsInput struct {
	Page     int
	Limit    int
	ID       *string
	Name     *string
	IsActive *bool
	Sort     string
}

func (u *FindShopsUsecase) Execute(
	ctx context.Context,
	input FindShopsInput,
) ([]domain.Shop, int, error) {
	var shopSortKeys = map[string]query.SortKey{
		"name":     repository.ShopSortName,
		"active":   repository.ShopSortActive,
		"date":     repository.ShopSortLatest,
		"modified": repository.ShopSortModify,
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

			sortKey, exists := shopSortKeys[key]
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
				By:        repository.ShopSortLatest,
				Direction: query.SortDesc,
			},
		}
	}

	params := repository.FindShopsParams{
		ID:       input.ID,
		Name:     input.Name,
		IsActive: input.IsActive,
		Pagination: query.Pagination{
			Page:  input.Page,
			Limit: input.Limit,
		},
		Sorts: sorts,
	}

	shops, total, err := u.shopRepo.
		FindByParams(ctx, u.executor, params)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to load shops: %w", err)
	}
	if len(shops) == 0 {
		return nil, 0, apperrors.NewNotFound("shops not available at the moment")
	}

	return shops, total, nil
}
