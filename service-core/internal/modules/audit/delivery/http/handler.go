package http

import (
	"net/http"
	"time"

	apperrors "service-core/internal/common/errors"
	apphttp "service-core/internal/common/http"
	"service-core/internal/modules/audit/usecase"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type auditHandler struct {
	findAuditLogs   *usecase.FindAuditLogsUsecase
	getAuditLog     *usecase.GetAuditLogUsecase
	deleteAuditLogs *usecase.DeleteAuditLogsUsecase
}

func NewAuditHandler(
	findAuditLogs *usecase.FindAuditLogsUsecase,
	getAuditLog *usecase.GetAuditLogUsecase,
	deleteAuditLogs *usecase.DeleteAuditLogsUsecase,
) *auditHandler {
	return &auditHandler{
		findAuditLogs:   findAuditLogs,
		getAuditLog:     getAuditLog,
		deleteAuditLogs: deleteAuditLogs,
	}
}

func (h *auditHandler) ListAuditLogs(w http.ResponseWriter, r *http.Request) error {
	page := apphttp.QueryIntDefault(r, "page", 1)
	if page <= 0 {
		page = 1
	}
	limit := apphttp.QueryIntDefault(r, "limit", 10)
	if limit <= 0 {
		limit = 10
	}

	sort := apphttp.Query(r, "sort")
	actionStr := apphttp.Query(r, "action")
	userIDStr := apphttp.Query(r, "user_id")

	var action *string
	if actionStr != "" {
		action = &actionStr
	}

	var actorID *string
	if userIDStr != "" {
		actorID = &userIDStr
	}

	var startDate *time.Time
	if startDateStr := apphttp.Query(r, "start_date"); startDateStr != "" {
		parsed, err := time.Parse(time.RFC3339, startDateStr)
		if err != nil {
			parsedDate, errDate := time.Parse("2006-01-02", startDateStr)
			if errDate != nil {
				return apperrors.NewBadRequest("invalid start_date format, must be RFC3339 or YYYY-MM-DD")
			}
			startDate = &parsedDate
		} else {
			startDate = &parsed
		}
	}

	var endDate *time.Time
	if endDateStr := apphttp.Query(r, "end_date"); endDateStr != "" {
		parsed, err := time.Parse(time.RFC3339, endDateStr)
		if err != nil {
			parsedDate, errDate := time.Parse("2006-01-02", endDateStr)
			if errDate != nil {
				return apperrors.NewBadRequest("invalid end_date format, must be RFC3339 or YYYY-MM-DD")
			}
			// Set time to end of day if only date is provided
			parsedDate = parsedDate.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
			endDate = &parsedDate
		} else {
			endDate = &parsed
		}
	}

	input := usecase.FindAuditLogsInput{
		Page:      page,
		Limit:     limit,
		Action:    action,
		ActorID:   actorID,
		StartDate: startDate,
		EndDate:   endDate,
		Sort:      sort,
	}

	logs, total, err := h.findAuditLogs.Execute(r.Context(), input)
	if err != nil {
		return err
	}

	results := make([]auditLogResponse, 0, len(logs))
	for _, log := range logs {
		results = append(results, auditLogResponse{
			ID:         log.ID,
			Category:   string(log.Category),
			Action:     log.Action,
			Resource:   log.Resource,
			ResourceID: log.ResourceID,
			ActorID:    log.ActorID,
			Outcome:    string(log.Outcome),
			RequestID:  log.RequestID,
			ClientIP:   log.ClientIP,
			Metadata:   log.Metadata,
			CreatedAt:  log.CreatedAt,
		})
	}

	response := map[string]interface{}{
		"audit_logs": results,
		"page":       page,
		"limit":      limit,
		"total":      total,
	}

	apphttp.WriteJSON(w, http.StatusOK, response)
	return nil
}

func (h *auditHandler) GetAuditLog(w http.ResponseWriter, r *http.Request) error {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return apperrors.NewBadRequest("invalid audit log id")
	}

	input := usecase.GetAuditLogInput{
		ID: id,
	}

	log, err := h.getAuditLog.Execute(r.Context(), input)
	if err != nil {
		return err
	}

	response := auditLogResponse{
		ID:         log.ID,
		Category:   string(log.Category),
		Action:     log.Action,
		Resource:   log.Resource,
		ResourceID: log.ResourceID,
		ActorID:    log.ActorID,
		Outcome:    string(log.Outcome),
		RequestID:  log.RequestID,
		ClientIP:   log.ClientIP,
		Metadata:   log.Metadata,
		CreatedAt:  log.CreatedAt,
	}

	apphttp.WriteJSON(w, http.StatusOK, response)
	return nil
}

func (h *auditHandler) DeleteAuditLogs(w http.ResponseWriter, r *http.Request) error {
	var input usecase.DeleteAuditLogsInput
	var msg string

	idStr := chi.URLParam(r, "id")
	if idStr != "" {
		id, err := uuid.Parse(idStr)
		if err != nil {
			return apperrors.NewBadRequest("invalid audit log id")
		}

		input.IDs = []uuid.UUID{id}
		msg = "audit log deleted successfully"
	} else {
		var req deleteAuditLogsRequest
		if err := apphttp.DecodeJSON(r, &req); err != nil {
			return apperrors.NewBadRequest("invalid body request")
		}

		switch {
		case req.All:
			input.DeleteAll = true
			msg = "all audit logs deleted successfully"

		case len(req.IDs) == 0:
			return apperrors.NewBadRequest("no valid audit log IDs provided for deletion")

		default:
			input.IDs = req.IDs
			msg = "audit logs deleted successfully"
		}
	}

	if err := h.deleteAuditLogs.Execute(r.Context(), input); err != nil {
		return err
	}

	apphttp.WriteJSON(w, http.StatusOK, map[string]string{
		"message": msg,
	})
	return nil
}
