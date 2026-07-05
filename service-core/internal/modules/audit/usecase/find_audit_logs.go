package usecase

import (
	"context"
	"fmt"
	"strings"
	"time"

	"service-core/internal/modules/audit/domain"
	"service-core/internal/modules/audit/repository"
	query "service-core/internal/shared/query"
	transaction "service-core/internal/shared/transaction"
)

type FindAuditLogsUsecase struct {
	executor  transaction.Executor
	auditRepo repository.AuditLogRepository
}

func NewFindAuditLogsUsecase(
	executor transaction.Executor,
	auditRepo repository.AuditLogRepository,
) *FindAuditLogsUsecase {
	return &FindAuditLogsUsecase{
		executor:  executor,
		auditRepo: auditRepo,
	}
}

type FindAuditLogsInput struct {
	Page      int
	Limit     int
	Action    *string
	ActorID   *string
	StartDate *time.Time
	EndDate   *time.Time
	Sort      string
}

func (u *FindAuditLogsUsecase) Execute(
	ctx context.Context,
	input FindAuditLogsInput,
) ([]domain.AuditLog, int, error) {
	var sortKeys = map[string]query.SortKey{
		"date":   repository.AuditLogSortDate,
		"action": repository.AuditLogSortAction,
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

			sortKey, exists := sortKeys[key]
			if exists {
				sorts = append(sorts, query.Sort{
					By:        query.SortKey(sortKey),
					Direction: dir,
				})
			}
		}
	}

	if len(sorts) == 0 {
		sorts = query.Sorts{
			{
				By:        query.SortKey(repository.AuditLogSortDate),
				Direction: query.SortDesc,
			},
		}
	}

	params := repository.FindAuditLogsParams{
		Action:    input.Action,
		ActorID:   input.ActorID,
		StartDate: input.StartDate,
		EndDate:   input.EndDate,
		Pagination: query.Pagination{
			Page:  input.Page,
			Limit: input.Limit,
		},
		Sorts: sorts,
	}

	logs, total, err := u.auditRepo.
		Find(ctx, u.executor, params)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to load audit logs: %w", err)
	}

	return logs, total, nil
}
