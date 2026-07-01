package http

import (
	"net/http"

	apperrors "service-core/internal/common/errors"
	apphttp "service-core/internal/common/http"
	appcookie "service-core/internal/common/http/cookie"
	authdomain "service-core/internal/modules/authentication/domain"
	"service-core/internal/modules/authentication/usecase"

	"github.com/google/uuid"
)

type authHandler struct {
	me               *usecase.MeUsecase
	logout           *usecase.LogoutUsecase
	loginCustomer    *usecase.LoginCustomerUsecase
	loginStaff       *usecase.LoginStaffUsecase
	registerCustomer *usecase.RegisterCustomerUsecase
	verifyAccount    *usecase.VerifyAccountUsecase
	getAccount       *usecase.GetAccountUsecase
}

func NewAuthHandler(
	me *usecase.MeUsecase,
	logout *usecase.LogoutUsecase,
	loginCustomer *usecase.LoginCustomerUsecase,
	loginStaff *usecase.LoginStaffUsecase,
	registerCustomer *usecase.RegisterCustomerUsecase,
	verifyAccount *usecase.VerifyAccountUsecase,
	getAccount *usecase.GetAccountUsecase,
) *authHandler {
	return &authHandler{
		me:               me,
		logout:           logout,
		loginCustomer:    loginCustomer,
		loginStaff:       loginStaff,
		registerCustomer: registerCustomer,
		verifyAccount:    verifyAccount,
		getAccount:       getAccount,
	}
}

func (h *authHandler) GetByID(w http.ResponseWriter, r *http.Request) error {
	authCtx, ok := authdomain.GetAuthContext(r.Context())
	if !ok {
		return apperrors.NewUnauthorized("authentication required")
	}

	acc, err := h.getAccount.Execute(r.Context(), authCtx.UserID)
	if err != nil {
		return err
	}
	if acc == nil {
		return apperrors.NewNotFound("account not found")
	}

	response := map[string]interface{}{
		"id":            acc.ID,
		"email":         acc.Email,
		"last_login_at": acc.LastLoginAt,
	}

	apphttp.WriteJSON(w, http.StatusOK, response)
	return nil
}

func (h *authHandler) Me(w http.ResponseWriter, r *http.Request) error {
	authCtx, ok := authdomain.GetAuthContext(r.Context())
	if !ok || !authCtx.IsAuthenticated {
		return apperrors.NewUnauthorized("authentication required")
	}

	me, err := h.me.Execute(r.Context(), *authCtx)
	if err != nil {
		return err
	}

	roles := make([]roleResponse, 0, len(me.Actor.Roles))
	for _, role := range me.Actor.Roles {
		roles = append(roles, roleResponse{
			Code: string(role.Code),
			Name: role.Name,
		})
	}

	permissions := make([]permissionResponse, 0, len(me.Actor.Permissions))
	for _, role := range me.Actor.Roles {
		permissions = append(permissions, permissionResponse{
			Code: string(role.Code),
		})
	}

	response := meResponse{
		AccountID:       me.Account.ID,
		AccountType:     string(me.Account.Type),
		IsAuthenticated: true,
		StaffID:         me.Actor.StaffID,
		Roles:           roles,
		Permissions:     permissions,
	}

	apphttp.WriteJSON(w, http.StatusOK, response)
	return nil
}

func (h *authHandler) SignInEmail(w http.ResponseWriter, r *http.Request) error {
	var req signInEmailRequest

	if err := apphttp.DecodeJSON(r, &req); err != nil {
		return apperrors.NewBadRequest("invalid body request")
	}

	if req.Email == "" {
		return apperrors.NewBadRequest("invalid email")
	}
	if req.Password == "" {
		return apperrors.NewBadRequest("invalid password")
	}

	input := usecase.LoginCustomerParams{
		UserAgent: req.UserAgent,
		IPAddress: req.IPAddress,
		Email:     req.Email,
		Password:  req.Password,
	}

	tokens, err := h.loginCustomer.Execute(r.Context(), input)
	if err != nil {
		return err
	}

	appcookie.Bind(
		w,
		appcookie.CookieAccess,
		tokens.AccessToken.Token,
		tokens.AccessToken.ExpiresAt,
	)

	response := map[string]string{
		"message": "login success",
	}

	apphttp.WriteJSON(w, http.StatusOK, response)
	return nil
}

func (h *authHandler) SignUpAccount(w http.ResponseWriter, r *http.Request) error {
	var req signUpRequest

	if err := apphttp.DecodeJSON(r, &req); err != nil {
		return apperrors.NewBadRequest("invalid body request")
	}

	if req.Email == "" {
		return apperrors.NewBadRequest("invalid email")
	}
	if req.Password == "" {
		return apperrors.NewBadRequest("invalid password")
	}
	if req.Username == "" {
		return apperrors.NewBadRequest("invalid user name")
	}

	input := usecase.RegisterCustomerParams{
		Email:    req.Email,
		Password: req.Password,
		Name:     req.Name,
		Username: req.Username,
		Phone:    req.Phone,
	}

	challengeID, err := h.registerCustomer.Execute(r.Context(), input)
	if err != nil {
		return err
	}

	response := signUpResponse{
		Message:     "verification code sent",
		ChallengeID: *challengeID,
	}

	apphttp.WriteJSON(w, http.StatusCreated, response)
	return nil
}

func (h *authHandler) VerifyAccount(w http.ResponseWriter, r *http.Request) error {
	var req verifyAccountRequest

	if err := apphttp.DecodeJSON(r, &req); err != nil {
		return apperrors.NewBadRequest("invalid body request")
	}

	if req.ChallengeID == "" {
		return apperrors.NewBadRequest("invalid challenge id")
	}
	challengeID, err := uuid.Parse(req.ChallengeID)
	if err != nil {
		return apperrors.NewBadRequest("invalid challenge id")
	}
	if len(req.OTP) != 6 {
		return apperrors.NewBadRequest("invalid otp")
	}

	input := usecase.VerifyAccountParams{
		UserAgent:   req.UserAgent,
		IPAddress:   req.IPAddress,
		ChallengeID: challengeID,
		OTP:         req.OTP,
	}

	tokens, err := h.verifyAccount.Execute(r.Context(), input)
	if err != nil {
		return err
	}

	appcookie.Bind(
		w,
		appcookie.CookieAccess,
		tokens.AccessToken.Token,
		tokens.AccessToken.ExpiresAt,
	)

	response := map[string]string{
		"message": "verify success",
	}

	apphttp.WriteJSON(w, http.StatusCreated, response)
	return nil
}

func (h *authHandler) SignInStaffEmail(w http.ResponseWriter, r *http.Request) error {
	var req signInEmailRequest

	if err := apphttp.DecodeJSON(r, &req); err != nil {
		return apperrors.NewBadRequest("invalid body request")
	}

	if req.Email == "" {
		return apperrors.NewBadRequest("invalid email")
	}
	if req.Password == "" {
		return apperrors.NewBadRequest("invalid password")
	}

	input := usecase.LoginStaffParams{
		UserAgent: req.UserAgent,
		IPAddress: req.IPAddress,
		Email:     req.Email,
		Password:  req.Password,
	}

	tokens, err := h.loginStaff.Execute(r.Context(), input)
	if err != nil {
		return err
	}

	appcookie.Bind(
		w,
		appcookie.CookieStaff,
		tokens.AccessToken.Token,
		tokens.AccessToken.ExpiresAt,
	)

	response := map[string]string{
		"message": "login success",
	}

	apphttp.WriteJSON(w, http.StatusOK, response)
	return nil
}

func (h *authHandler) Logout(w http.ResponseWriter, r *http.Request) error {
	authCtx, ok := authdomain.GetAuthContext(r.Context())
	if !ok || !authCtx.IsAuthenticated {
		return apperrors.NewUnauthorized("authentication required")
	}

	err := h.logout.Execute(r.Context(), *authCtx)
	if err != nil {
		return err
	}

	if authCtx.CustomerID != nil {
		appcookie.Clear(w, appcookie.CookieAccess)
	} else if authCtx.StaffID != nil {
		appcookie.Clear(w, appcookie.CookieStaff)
	}

	response := map[string]string{
		"message": "logout success",
	}

	apphttp.WriteJSON(w, http.StatusOK, response)
	return nil
}
