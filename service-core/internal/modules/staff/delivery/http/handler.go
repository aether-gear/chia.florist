package http

import (
	"net/http"

	apperrors "service-core/internal/common/errors"
	apphttp "service-core/internal/common/http"
	authzDomain "service-core/internal/modules/authorization/domain"
	authzSvc "service-core/internal/modules/authorization/infra/service"
	"service-core/internal/modules/staff/usecase"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type staffHandler struct {
	addStaffAccount       *usecase.AddStaffAccountUsecase
	createStaff           *usecase.CreateStaffUsecase
	findStaff             *usecase.FindStaffUsecase
	listStaffAccounts     *usecase.ListStaffAccountsUsecase
	updateStaff           *usecase.UpdateStaffUsecase
	deleteStaff           *usecase.DeleteStaffUsecase
	removeStaffAccount    *usecase.RemoveStaffAccountUsecase
	listStaffPermissions  *usecase.ListStaffPermissionsUsecase
	saveStaffPermission   *usecase.SaveStaffPermissionUsecase
	deleteStaffPermission *usecase.DeleteStaffPermissionUsecase
}

func NewStaffHandler(
	addStaffAccount *usecase.AddStaffAccountUsecase,
	createStaff *usecase.CreateStaffUsecase,
	findStaff *usecase.FindStaffUsecase,
	listStaffAccounts *usecase.ListStaffAccountsUsecase,
	updateStaff *usecase.UpdateStaffUsecase,
	deleteStaff *usecase.DeleteStaffUsecase,
	removeStaffAccount *usecase.RemoveStaffAccountUsecase,
	listStaffPermissions *usecase.ListStaffPermissionsUsecase,
	saveStaffPermission *usecase.SaveStaffPermissionUsecase,
	deleteStaffPermission *usecase.DeleteStaffPermissionUsecase,
) *staffHandler {
	return &staffHandler{
		addStaffAccount:       addStaffAccount,
		createStaff:           createStaff,
		findStaff:             findStaff,
		listStaffAccounts:     listStaffAccounts,
		updateStaff:           updateStaff,
		deleteStaff:           deleteStaff,
		removeStaffAccount:    removeStaffAccount,
		listStaffPermissions:  listStaffPermissions,
		saveStaffPermission:   saveStaffPermission,
		deleteStaffPermission: deleteStaffPermission,
	}
}

func (h *staffHandler) AddStaffAccount(w http.ResponseWriter, r *http.Request) error {
	staffIDStr := chi.URLParam(r, "staffID")
	staffID, err := uuid.Parse(staffIDStr)
	if err != nil {
		return apperrors.NewBadRequest("invalid staff id")
	}

	actor, ok := authzSvc.GetActor(r.Context())
	if !ok {
		return apperrors.NewUnauthorized("authentication required")
	}
	if actor.StaffID == nil {
		return apperrors.NewForbidden(authzDomain.ErrInsufficientRole.Error())
	}

	var req addStaffAccountRequest
	if err := apphttp.DecodeJSON(r, &req); err != nil {
		return apperrors.NewBadRequest("invalid body request")
	}

	if req.Email == "" {
		return apperrors.NewBadRequest("email is required")
	}
	if req.Password == "" {
		return apperrors.NewBadRequest("password is required")
	}

	input := usecase.AddStaffAccountParams{
		ActorAccountID: actor.AccountID,
		ActorStaffID:   *actor.StaffID,
		StaffID:        staffID,
		Email:          req.Email,
		Password:       req.Password,
	}

	err = h.addStaffAccount.Execute(r.Context(), input)
	if err != nil {
		return err
	}

	response := map[string]string{
		"message": "staff account successfully created",
	}

	apphttp.WriteJSON(w, http.StatusCreated, response)
	return nil
}

func (h *staffHandler) CreateStaff(w http.ResponseWriter, r *http.Request) error {
	var req createStaffRequest
	if err := apphttp.DecodeJSON(r, &req); err != nil {
		return apperrors.NewBadRequest("invalid body request")
	}

	if req.Name == "" {
		return apperrors.NewBadRequest("name is required")
	}
	if req.Username == "" {
		return apperrors.NewBadRequest("username is required")
	}

	input := usecase.CreateStaffInput{
		Name:        req.Name,
		Username:    req.Username,
		Description: req.Description,
		LogoUrl:     req.LogoUrl,
		BannerUrl:   req.BannerUrl,
	}

	err := h.createStaff.Execute(r.Context(), input)
	if err != nil {
		return err
	}

	response := map[string]string{
		"message": "staff successfully created",
	}

	apphttp.WriteJSON(w, http.StatusCreated, response)
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

func (h *staffHandler) ListStaffAccounts(w http.ResponseWriter, r *http.Request) error {
	staffIDStr := chi.URLParam(r, "staffID")
	staffID, err := uuid.Parse(staffIDStr)
	if err != nil {
		return apperrors.NewBadRequest("invalid staff id")
	}

	actor, ok := authzSvc.GetActor(r.Context())
	if !ok {
		return apperrors.NewUnauthorized("authentication required")
	}
	if actor.StaffID == nil {
		return apperrors.NewForbidden(authzDomain.ErrInsufficientRole.Error())
	}

	input := usecase.ListStaffAccountsParams{
		ActorAccountID: actor.AccountID,
		ActorStaffID:   *actor.StaffID,
		StaffID:        staffID,
	}

	accounts, err := h.listStaffAccounts.Execute(r.Context(), input)
	if err != nil {
		return err
	}

	results := make([]staffAccountResponse, 0, len(accounts))
	for _, m := range accounts {
		results = append(results, staffAccountResponse{
			AccountID: m.AccountID,
			UserID:    m.UserID,
			Email:     m.Email,
			Name:      m.Name,
			Username:  m.Username,
			Phone:     m.Phone,
			AvatarURL: m.AvatarURL,
			Role: staffAccountRoleResponse{
				ID:   m.Role.ID,
				Code: string(m.Role.Code),
				Name: m.Role.Name,
			},
			LastLoginAt: m.LastLoginAt,
			CreatedAt:   m.CreatedAt,
		})
	}

	response := listStaffAccountsResponse{
		StaffID:  staffID,
		Total:    len(results),
		Accounts: results,
	}

	apphttp.WriteJSON(w, http.StatusOK, response)
	return nil
}

func (h *staffHandler) UpdateStaff(w http.ResponseWriter, r *http.Request) error {
	staffIDStr := chi.URLParam(r, "staffID")
	staffID, err := uuid.Parse(staffIDStr)
	if err != nil {
		return apperrors.NewBadRequest("invalid staff id")
	}

	actor, ok := authzSvc.GetActor(r.Context())
	if !ok {
		return apperrors.NewUnauthorized("authentication required")
	}
	if actor.StaffID == nil {
		return apperrors.NewForbidden(authzDomain.ErrInsufficientRole.Error())
	}

	var req updateStaffRequest
	if err := apphttp.DecodeJSON(r, &req); err != nil {
		return apperrors.NewBadRequest("invalid body request")
	}

	if req.Name == "" {
		return apperrors.NewBadRequest("name is required")
	}

	input := usecase.UpdateStaffInput{
		ActorAccountID: actor.AccountID,
		ActorStaffID:   *actor.StaffID,
		StaffID:        staffID,
		Name:           req.Name,
		Description:    req.Description,
		LogoUrl:        req.LogoUrl,
		BannerUrl:      req.BannerUrl,
	}

	err = h.updateStaff.Execute(r.Context(), input)
	if err != nil {
		return err
	}

	response := map[string]string{
		"message": "staff successfully updated",
	}

	apphttp.WriteJSON(w, http.StatusOK, response)
	return nil
}

func (h *staffHandler) DeleteStaff(w http.ResponseWriter, r *http.Request) error {
	staffIDStr := chi.URLParam(r, "staffID")
	staffID, err := uuid.Parse(staffIDStr)
	if err != nil {
		return apperrors.NewBadRequest("invalid staff id")
	}

	actor, ok := authzSvc.GetActor(r.Context())
	if !ok {
		return apperrors.NewUnauthorized("authentication required")
	}
	if actor.StaffID == nil {
		return apperrors.NewForbidden(authzDomain.ErrInsufficientRole.Error())
	}

	input := usecase.DeleteStaffInput{
		ActorAccountID: actor.AccountID,
		ActorStaffID:   *actor.StaffID,
		StaffID:        staffID,
	}

	err = h.deleteStaff.Execute(r.Context(), input)
	if err != nil {
		return err
	}

	response := map[string]string{
		"message": "staff successfully deleted",
	}

	apphttp.WriteJSON(w, http.StatusOK, response)
	return nil
}

func (h *staffHandler) RemoveStaffAccount(w http.ResponseWriter, r *http.Request) error {
	staffIDStr := chi.URLParam(r, "staffID")
	staffID, err := uuid.Parse(staffIDStr)
	if err != nil {
		return apperrors.NewBadRequest("invalid staff id")
	}

	accountIDStr := chi.URLParam(r, "accountID")
	accountID, err := uuid.Parse(accountIDStr)
	if err != nil {
		return apperrors.NewBadRequest("invalid account id")
	}

	actor, ok := authzSvc.GetActor(r.Context())
	if !ok {
		return apperrors.NewUnauthorized("authentication required")
	}
	if actor.StaffID == nil {
		return apperrors.NewForbidden(authzDomain.ErrInsufficientRole.Error())
	}

	input := usecase.RemoveStaffAccountInput{
		ActorAccountID: actor.AccountID,
		ActorStaffID:   *actor.StaffID,
		StaffID:        staffID,
		AccountID:      accountID,
	}

	err = h.removeStaffAccount.Execute(r.Context(), input)
	if err != nil {
		return err
	}

	response := map[string]string{
		"message": "staff account successfully removed",
	}

	apphttp.WriteJSON(w, http.StatusOK, response)
	return nil
}

func (h *staffHandler) ListStaffPermissions(w http.ResponseWriter, r *http.Request) error {
	staffIDStr := chi.URLParam(r, "staffID")
	staffID, err := uuid.Parse(staffIDStr)
	if err != nil {
		return apperrors.NewBadRequest("invalid staff id")
	}

	perms, err := h.listStaffPermissions.Execute(r.Context(), staffID)
	if err != nil {
		return err
	}

	res := listStaffPermissionsResponse{
		StaffID:     staffID,
		Total:       len(perms),
		Permissions: make([]staffPermissionResponse, 0, len(perms)),
	}

	for _, p := range perms {
		res.Permissions = append(res.Permissions, staffPermissionResponse{
			ID:          p.ID,
			StaffID:     p.StaffID,
			ShopID:      p.ShopID,
			ShopName:    p.ShopName,
			Permissions: p.Permissions,
			Rules:       p.Rules,
			CreatedAt:   p.CreatedAt,
			UpdatedAt:   p.UpdatedAt,
		})
	}

	apphttp.WriteJSON(w, http.StatusOK, res)
	return nil
}

func (h *staffHandler) SaveStaffPermission(w http.ResponseWriter, r *http.Request) error {
	staffIDStr := chi.URLParam(r, "staffID")
	staffID, err := uuid.Parse(staffIDStr)
	if err != nil {
		return apperrors.NewBadRequest("invalid staff id")
	}

	var req saveStaffPermissionRequest
	if err := apphttp.DecodeJSON(r, &req); err != nil {
		return apperrors.NewBadRequest("invalid request body")
	}

	if req.ShopID == uuid.Nil {
		return apperrors.NewBadRequest("shop_id is required")
	}

	input := usecase.SaveStaffPermissionParams{
		StaffID:     staffID,
		ShopID:      req.ShopID,
		Permissions: req.Permissions,
		Rules:       req.Rules,
	}

	if err := h.saveStaffPermission.Execute(r.Context(), input); err != nil {
		return err
	}

	apphttp.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "staff shop permission successfully updated",
	})
	return nil
}

func (h *staffHandler) DeleteStaffPermission(w http.ResponseWriter, r *http.Request) error {
	staffIDStr := chi.URLParam(r, "staffID")
	staffID, err := uuid.Parse(staffIDStr)
	if err != nil {
		return apperrors.NewBadRequest("invalid staff id")
	}

	shopIDStr := chi.URLParam(r, "shopID")
	shopID, err := uuid.Parse(shopIDStr)
	if err != nil {
		return apperrors.NewBadRequest("invalid shop id")
	}

	if err := h.deleteStaffPermission.Execute(r.Context(), staffID, shopID); err != nil {
		return err
	}

	apphttp.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "staff shop permission successfully deleted",
	})
	return nil
}
