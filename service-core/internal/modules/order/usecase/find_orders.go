package usecase

import (
	"context"
	"fmt"
	"strings"

	"service-core/internal/modules/order/domain"
	"service-core/internal/modules/order/repository"
	query "service-core/internal/shared/query"
	transaction "service-core/internal/shared/transaction"

	"github.com/google/uuid"
)

type FindOrdersUsecase struct {
	executor      transaction.Executor
	orderRepo     repository.OrderRepository
	orderItemRepo repository.OrderItemRepository
}

func NewFindOrdersUsecase(
	executor transaction.Executor,
	orderRepo repository.OrderRepository,
	orderItemRepo repository.OrderItemRepository,
) *FindOrdersUsecase {
	return &FindOrdersUsecase{
		executor:      executor,
		orderRepo:     orderRepo,
		orderItemRepo: orderItemRepo,
	}
}

type FindOrdersInput struct {
	Page       int
	Limit      int
	ID         *uuid.UUID
	Number     *string
	CustomerID *uuid.UUID
	Status     *string
	Sort       string
}

type OrderSearchResult struct {
	Order domain.Order
	Items []domain.OrderItem
}

func (u *FindOrdersUsecase) Execute(
	ctx context.Context,
	input FindOrdersInput,
) ([]OrderSearchResult, int, error) {
	var orderSortKeys = map[string]query.SortKey{
		"latest":   repository.OrderSortLatest,
		"date":     repository.OrderSortLatest,
		"number":   repository.OrderSortNumber,
		"total":    repository.OrderSortTotal,
		"status":   repository.OrderSortStatus,
		"modified": repository.OrderSortModify,
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

			sortKey, exists := orderSortKeys[key]
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
				By:        repository.OrderSortLatest,
				Direction: query.SortDesc,
			},
		}
	}

	params := repository.FindOrderParams{
		ID:         input.ID,
		Number:     input.Number,
		CustomerID: input.CustomerID,
		Status:     input.Status,
		Pagination: query.Pagination{
			Page:  input.Page,
			Limit: input.Limit,
		},
		Sorts: sorts,
	}

	orders, total, err := u.orderRepo.FindOrders(ctx, u.executor, params)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list orders: %w", err)
	}

	if len(orders) == 0 {
		return []OrderSearchResult{}, total, nil
	}

	orderIDs := make([]uuid.UUID, len(orders))
	for i, o := range orders {
		orderIDs[i] = o.ID
	}

	orderItems, err := u.orderItemRepo.ListByOrderIDs(ctx, u.executor, orderIDs)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list order items: %w", err)
	}

	itemsMap := make(map[uuid.UUID][]domain.OrderItem)
	for _, item := range orderItems {
		itemsMap[item.OrderID] = append(itemsMap[item.OrderID], item)
	}

	results := make([]OrderSearchResult, len(orders))
	for i, o := range orders {
		items := itemsMap[o.ID]
		if items == nil {
			items = []domain.OrderItem{}
		}
		results[i] = OrderSearchResult{
			Order: o,
			Items: items,
		}
	}

	return results, total, nil
}
