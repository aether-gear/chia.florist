package usecase

import (
	"context"
	"fmt"
	"strings"

	apperrors "service-core/internal/common/errors"
	"service-core/internal/modules/merchant/domain"
	"service-core/internal/modules/merchant/repository"
	query "service-core/internal/shared/query"
	transaction "service-core/internal/shared/transaction"

	"github.com/google/uuid"
)

type FindMerchantsUsecase struct {
	executor     transaction.Executor
	merchantRepo repository.MerchantRepository
}

func NewFindMerchantsUsecase(
	executor transaction.Executor,
	merchantRepo repository.MerchantRepository,
) *FindMerchantsUsecase {
	return &FindMerchantsUsecase{
		executor:     executor,
		merchantRepo: merchantRepo,
	}
}

type FindMerchantsInput struct {
	Page  int
	Limit int
	ID    *uuid.UUID
	Name  *string
	Sort  string
}

func (u *FindMerchantsUsecase) Execute(
	ctx context.Context,
	input FindMerchantsInput,
) ([]domain.Merchant, int, error) {
	var merchantSortKeys = map[string]query.SortKey{
		"latest": repository.MerchantSortLatest,
		"name":   repository.MerchantSortName,
		"modify": repository.MerchantSortModify,
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

			sortKey, exists := merchantSortKeys[key]
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
				By:        repository.MerchantSortLatest,
				Direction: query.SortDesc,
			},
		}
	}

	params := repository.FindMerchantParams{
		ID:   input.ID,
		Name: input.Name,
		Pagination: query.Pagination{
			Page:  input.Page,
			Limit: input.Limit,
		},
		Sorts: sorts,
	}

	merchants, total, err := u.merchantRepo.FindMerchants(ctx, u.executor, params)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to load merchants: %w", err)
	}
	if len(merchants) == 0 {
		return nil, 0, apperrors.NewNotFound("merchants not available at the moment")
	}

	return merchants, total, nil
}
