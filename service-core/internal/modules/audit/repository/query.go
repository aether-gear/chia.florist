package repository

import (
	"time"

	query "service-core/internal/shared/query"
)

const (
	AuditLogSortDate   query.SortKey = "date"
	AuditLogSortAction query.SortKey = "action"
)

type FindAuditLogsParams struct {
	Action     *string
	ActorID    *string
	StartDate  *time.Time
	EndDate    *time.Time
	Pagination query.Pagination
	Sorts      query.Sorts
}
