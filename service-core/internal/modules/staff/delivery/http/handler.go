package http

import (
	"net/http"

	apperrors "service-core/internal/common/errors"
	apphttp "service-core/internal/common/http"
	authenDomain "service-core/internal/modules/authentication/domain"
	authzSvc "service-core/internal/modules/authorization/infra/service"
	"service-core/internal/modules/staff/usecase"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type staffHandler struct {
	addStaffAccount *usecase.AddStaffAccountUsecase
	createStaff     *usecase.CreateStaffUsecase
	findStaff       *usecase.FindStaffUsecase
}

func NewStaffHandler(
	addStaffAccount *usecase.AddStaffAccountUsecase,
	createStaff *usecase.CreateStaffUsecase,
	findStaff *usecase.FindStaffUsecase,
) *staffHandler {
	return &staffHandler{
		addStaffAccount: addStaffAccount,
		createStaff:     createStaff,
		findStaff:       findStaff,
	}
}

func (h *staffHandler) AddStaffAccount(w http.ResponseWriter, r *http.Request) error {
	staffIDStr := chi.URLParam(r, "staffID")
	staffID, err := uuid.Parse(staffIDStr)
	if err != nil {
		return apperrors.NewBadRequest("invalid staff id")
	}

	authCtx, ok := authenDomain.GetAuthContext(r.Context())
	if !ok {
		return apperrors.NewUnauthorized("authentication required")
	}
	_ = authCtx

	actor, ok := authzSvc.GetActor(r.Context())
	if !ok {
		return apperrors.NewUnauthorized("authentication required")
	}

	var req addStaffAccountRequest
	if err := apphttp.DecodeJSON(r, &req); err != nil {
		return apperrors.NewBadRequest("invalid body request")
	}

	if req.Email == "" {
		return apperrors.NewBadRequest("email is required")
	}
	if req.Name == "" {
		return apperrors.NewBadRequest("name is required")
	}
	if req.Username == "" {
		return apperrors.NewBadRequest("username is required")
	}

	input := usecase.AddStaffAccountParams{
		ActorAccountID: actor.AccountID,
		ActorStaffID:   *actor.StaffID,
		StaffID:        staffID,
		Email:          req.Email,
		Name:           req.Name,
		Password:       req.Password,
		Username:       req.Username,
		Phone:          req.Phone,
	}

	err = h.addStaffAccount.Execute(r.Context(), input)
	if err != nil {
		return err
	}

	response := map[string]string{
		"message": "verify success",
	}

	apphttp.WriteJSON(w, http.StatusCreated, response)
	return nil
}

func (h *staffHandler) CreateStaff(w http.ResponseWriter, r *http.Request) error {
	response := map[string]string{
		"message": "staff successfully created",
	}

	apphttp.WriteJSON(w, http.StatusOK, response)
	return nil
}

func (h *staffHandler) FindStaff(w http.ResponseWriter, r *http.Request) error {
	page := apphttp.QueryIntDefault(r, "page", 1)
	if page <= 0 {
		page = 1
	}
	limit := apphttp.QueryIntDefault(r, "limit", 10)
	if limit <= 0 {
		limit = 10
	}

	idStr := apphttp.Query(r, "id")
	sort := apphttp.Query(r, "sort")

	input := usecase.FindStaffInput{
		Page:  page,
		Limit: limit,
		Sort:  sort,
	}
	if idStr != "" {
		id, err := uuid.Parse(idStr)
		if err != nil {
			return apperrors.NewBadRequest("invalid staff id")
		}
		input.ID = &id
	}

	staff, total, err := h.findStaff.Execute(r.Context(), input)
	if err != nil {
		return err
	}

	results := make([]staffResponse, 0, len(staff))
	for _, m := range staff {
		results = append(results, staffResponse{
			ID:        m.ID,
			UserID:    m.UserID,
			Name:      m.Name,
			Username:  m.Username,
			Phone:     m.Phone,
			AvatarURL: m.AvatarURL,
			CreatedAt: m.CreatedAt,
		})
	}

	response := map[string]interface{}{
		"staff": results,
		"page":  page,
		"limit": limit,
		"total": total,
	}

	apphttp.WriteJSON(w, http.StatusOK, response)
	return nil
}
