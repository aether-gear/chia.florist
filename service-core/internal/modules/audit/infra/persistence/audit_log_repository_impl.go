package persistence

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"service-core/internal/modules/audit/domain"
	"service-core/internal/modules/audit/repository"

	query "service-core/internal/shared/query"
	transaction "service-core/internal/shared/transaction"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type auditLogRepositoryImpl struct{}

func NewAuditLogRepository() repository.AuditLogRepository {
	return &auditLogRepositoryImpl{}
}

func (r *auditLogRepositoryImpl) Save(
	ctx context.Context,
	exec transaction.Executor,
	log domain.AuditLog,
) error {
	metadata, err := json.Marshal(log.Metadata)
	if err != nil {
		return fmt.Errorf("audit log: failed to marshal metadata: %w", err)
	}

	if log.ID == uuid.Nil {
		log.ID = uuid.New()
	}
	if log.CreatedAt.IsZero() {
		log.CreatedAt = time.Now().UTC()
	}

	queryStr := `
		INSERT INTO audit_logs (
			id,
			category,
			action,
			resource,
			resource_id,
			actor_id,
			outcome,
			request_id,
			client_ip,
			metadata,
			created_at
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
	`

	_, err = exec.Exec(ctx, queryStr,
		log.ID,
		log.Category,
		log.Action,
		log.Resource,
		nullableString(log.ResourceID),
		nullableString(log.ActorID),
		log.Outcome,
		nullableString(log.RequestID),
		nullableString(log.ClientIP),
		json.RawMessage(metadata),
		log.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("audit log: failed to save record: %w", err)
	}

	return nil
}

func (r *auditLogRepositoryImpl) GetByID(
	ctx context.Context,
	exec transaction.Executor,
	id uuid.UUID,
) (*domain.AuditLog, error) {
	queryStr := `
		SELECT
			id,
			category,
			action,
			resource,
			resource_id,
			actor_id,
			outcome,
			request_id,
			client_ip,
			metadata,
			created_at
		FROM
			audit_logs
		WHERE
			id = $1
		LIMIT 1
	`

	var (
		item          domain.AuditLog
		resourceID    *string
		actorID       *string
		requestID     *string
		clientIP      *string
		metadataBytes []byte
	)

	err := exec.QueryRow(ctx, queryStr, id).Scan(
		&item.ID,
		&item.Category,
		&item.Action,
		&item.Resource,
		&resourceID,
		&actorID,
		&item.Outcome,
		&requestID,
		&clientIP,
		&metadataBytes,
		&item.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("query audit log by id failed: %w", err)
	}

	if resourceID != nil {
		item.ResourceID = *resourceID
	}
	if actorID != nil {
		item.ActorID = *actorID
	}
	if requestID != nil {
		item.RequestID = *requestID
	}
	if clientIP != nil {
		item.ClientIP = *clientIP
	}
	if len(metadataBytes) > 0 {
		if err := json.Unmarshal(metadataBytes, &item.Metadata); err != nil {
			return nil, fmt.Errorf("audit log: failed to unmarshal metadata: %w", err)
		}
	}

	return &item, nil
}

func (r *auditLogRepositoryImpl) Find(
	ctx context.Context,
	exec transaction.Executor,
	params repository.FindAuditLogsParams,
) ([]domain.AuditLog, int, error) {
	baseQuery := `
		FROM audit_logs a
	`

	selectQuery := `
		SELECT
			a.id,
			a.category,
			a.action,
			a.resource,
			a.resource_id,
			a.actor_id,
			a.outcome,
			a.request_id,
			a.client_ip,
			a.metadata,
			a.created_at
	`

	whereClause := ""
	var (
		conditions []string
		args       []any
		argPos     = 1
	)

	if params.Action != nil && *params.Action != "" {
		conditions = append(conditions, fmt.Sprintf("a.action = $%d", argPos))
		args = append(args, *params.Action)
		argPos++
	}

	if params.ActorID != nil && *params.ActorID != "" {
		conditions = append(conditions, fmt.Sprintf("a.actor_id = $%d", argPos))
		args = append(args, *params.ActorID)
		argPos++
	}

	if params.StartDate != nil {
		conditions = append(conditions, fmt.Sprintf("a.created_at >= $%d", argPos))
		args = append(args, *params.StartDate)
		argPos++
	}

	if params.EndDate != nil {
		conditions = append(conditions, fmt.Sprintf("a.created_at <= $%d", argPos))
		args = append(args, *params.EndDate)
		argPos++
	}

	if len(conditions) > 0 {
		whereClause = " WHERE " + strings.Join(conditions, " AND ")
	}

	// Count matching audit logs
	countQuery := `
		SELECT COUNT(a.id)
	` + baseQuery + whereClause

	countArgs := append([]any{}, args...)
	var total int
	err := exec.QueryRow(ctx, countQuery, countArgs...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("query count audit logs failed: %w", err)
	}

	// Build sorting expressions
	var sortKeys = map[query.SortKey]string{
		repository.AuditLogSortDate:   "a.created_at",
		repository.AuditLogSortAction: "a.action",
	}

	var sortClauses []string
	for _, sort := range params.Sorts {
		colName, exists := sortKeys[sort.By]
		if !exists {
			continue
		}

		direction := "DESC"
		if sort.Direction == query.SortAsc {
			direction = "ASC"
		}

		sortClauses = append(
			sortClauses,
			fmt.Sprintf("%s %s", colName, direction),
		)
	}

	orderBy := "ORDER BY a.created_at DESC"
	if len(sortClauses) > 0 {
		orderBy = "ORDER BY " + strings.Join(sortClauses, ", ")
	}

	// Apply pagination
	limitPos := argPos
	offsetPos := argPos + 1

	limit := params.Pagination.Limit
	if limit <= 0 {
		limit = 10
	}

	page := params.Pagination.Page
	if page <= 0 {
		page = 1
	}

	offset := (page - 1) * limit
	args = append(args, limit, offset)

	queryStr := selectQuery +
		baseQuery +
		whereClause +
		" " +
		orderBy +
		fmt.Sprintf(
			" LIMIT $%d OFFSET $%d",
			limitPos,
			offsetPos,
		)

	rows, err := exec.Query(ctx, queryStr, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("query audit logs failed: %w", err)
	}
	defer rows.Close()

	var results []domain.AuditLog
	for rows.Next() {
		var (
			item          domain.AuditLog
			resourceID    *string
			actorID       *string
			requestID     *string
			clientIP      *string
			metadataBytes []byte
		)

		err := rows.Scan(
			&item.ID,
			&item.Category,
			&item.Action,
			&item.Resource,
			&resourceID,
			&actorID,
			&item.Outcome,
			&requestID,
			&clientIP,
			&metadataBytes,
			&item.CreatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("scan audit log failed: %w", err)
		}

		if resourceID != nil {
			item.ResourceID = *resourceID
		}
		if actorID != nil {
			item.ActorID = *actorID
		}
		if requestID != nil {
			item.RequestID = *requestID
		}
		if clientIP != nil {
			item.ClientIP = *clientIP
		}
		if len(metadataBytes) > 0 {
			if err := json.Unmarshal(metadataBytes, &item.Metadata); err != nil {
				return nil, 0, fmt.Errorf("audit log: failed to unmarshal metadata: %w", err)
			}
		}

		results = append(results, item)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("iterate audit logs failed: %w", err)
	}

	return results, total, nil
}

// nullableString converts an empty string to nil so it maps to SQL NULL
// rather than an empty string in nullable text columns.
func nullableString(s string) interface{} {
	if s == "" {
		return nil
	}
	return s
}

func (r *auditLogRepositoryImpl) DeleteMultiple(
	ctx context.Context,
	exec transaction.Executor,
	ids []uuid.UUID,
) error {
	if len(ids) == 0 {
		return nil
	}
	idStrings := make([]string, len(ids))
	for i, id := range ids {
		idStrings[i] = id.String()
	}
	queryStr := `DELETE FROM audit_logs WHERE id = ANY($1::uuid[])`
	_, err := exec.Exec(ctx, queryStr, idStrings)
	if err != nil {
		return fmt.Errorf("audit log: failed to delete multiple logs: %w", err)
	}
	return nil
}

func (r *auditLogRepositoryImpl) DeleteAll(
	ctx context.Context,
	exec transaction.Executor,
) error {
	queryStr := `DELETE FROM audit_logs`
	_, err := exec.Exec(ctx, queryStr)
	if err != nil {
		return fmt.Errorf("audit log: failed to delete all logs: %w", err)
	}
	return nil
}
