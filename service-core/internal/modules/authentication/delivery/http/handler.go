package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	apperrors "service-core/internal/common/errors"
	apphttp "service-core/internal/common/http"
	appcookie "service-core/internal/common/http/cookie"
	authdomain "service-core/internal/modules/authentication/domain"
	"service-core/internal/modules/authentication/usecase"

	"github.com/google/uuid"
)

type authHandler struct {
	loginCustomer    *usecase.LoginCustomerUsecase
	registerCustomer *usecase.RegisterCustomerUsecase
	verifyAccount    *usecase.VerifyAccountUsecase
	getAccount       *usecase.GetAccountUsecase
}

func NewAuthHandler(
	loginCustomer *usecase.LoginCustomerUsecase,
	registerCustomer *usecase.RegisterCustomerUsecase,
	verifyAccount *usecase.VerifyAccountUsecase,
	getAccount *usecase.GetAccountUsecase,
) *authHandler {
	return &authHandler{
		loginCustomer:    loginCustomer,
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

	acc, err := h.getAccount.Execute(authCtx.UserID)
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

func (h *authHandler) SignInEmail(w http.ResponseWriter, r *http.Request) error {
	var req signInEmailRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
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

	tokens, err := h.loginCustomer.Execute(input)
	if err != nil {
		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     appcookie.AccessTokenCookieName,
		Value:    tokens.AccessToken.Token,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		Expires:  tokens.AccessToken.ExpiresAt,
	})

	response := map[string]string{
		"message": "login success",
	}

	apphttp.WriteJSON(w, http.StatusOK, response)
	return nil
}

func (h *authHandler) SignUpAccount(w http.ResponseWriter, r *http.Request) error {
	var req signUpRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
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

	challengeID, err := h.registerCustomer.Execute(input)
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

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return apperrors.NewBadRequest("invalid body request")
	}

	if req.ChallengeID == "" {
		return apperrors.NewBadRequest("invalid challenge id")
	}
	challengeID, err := uuid.Parse(req.ChallengeID)
	if err != nil {
		return apperrors.NewBadRequest("invalid challenge id")
	}
	if len(strconv.Itoa(req.OTP)) != 6 {
		return apperrors.NewBadRequest("invalid otp")
	}

	input := usecase.VerifyAccountParams{
		UserAgent:   req.UserAgent,
		IPAddress:   req.IPAddress,
		ChallengeID: challengeID,
		OTP:         req.OTP,
	}

	tokens, err := h.verifyAccount.Execute(input)
	if err != nil {
		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     appcookie.AccessTokenCookieName,
		Value:    tokens.AccessToken.Token,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		Expires:  tokens.AccessToken.ExpiresAt,
	})

	response := map[string]string{
		"message": "verify success",
	}

	apphttp.WriteJSON(w, http.StatusCreated, response)
	return nil
}
