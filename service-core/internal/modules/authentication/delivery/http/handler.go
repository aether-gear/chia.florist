package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"service-core/internal/modules/authentication/usecase"

	"service-core/internal/common/errors"
	apphttp "service-core/internal/common/http"

	"github.com/google/uuid"
)

type authHandler struct {
	signInEmail *usecase.LoginEmailUsecase
	signUp      *usecase.RegisterUsecase
	verify      *usecase.VerifyAccountUsecase
	getAccount  *usecase.GetAccountUsecase
}

func NewAuthHandler(
	signInEmail *usecase.LoginEmailUsecase,
	signUp *usecase.RegisterUsecase,
	verify *usecase.VerifyAccountUsecase,
	getAccount *usecase.GetAccountUsecase,
) *authHandler {
	return &authHandler{
		signInEmail: signInEmail,
		signUp:      signUp,
		verify:      verify,
		getAccount:  getAccount,
	}
}

func (h *authHandler) GetByID(w http.ResponseWriter, r *http.Request) error {
	UserId, ok := r.Context().Value("user_id").(string)
	if !ok {
		return errors.ErrUnauthorized
	}

	userID, err := uuid.Parse(UserId)
	if err != nil {
		return errors.ErrBadRequest
	}

	acc, err := h.getAccount.Execute(userID)
	if err != nil {
		return err
	}

	if acc == nil {
		return errors.ErrNotFound
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
	var req SignInEmailParams

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return errors.ErrBadRequest
	}

	if req.Email == "" || req.Password == "" {
		return errors.ErrBadRequest
	}

	tokens, err := h.signInEmail.Execute(usecase.LoginEmailParams{
		UserAgent: req.UserAgent,
		IPAddress: req.IPAddress,
		Email:     req.Email,
		Password:  req.Password,
	})
	if err != nil {
		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session",
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
	var req SignUpParams

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return errors.ErrBadRequest
	}

	if req.Email == "" || req.Password == "" || req.Username == "" {
		return errors.ErrBadRequest
	}

	challengeID, err := h.signUp.Execute(usecase.SignUpParams{
		Email:    req.Email,
		Password: req.Password,
		Name:     req.Name,
		Username: req.Username,
		Phone:    req.Phone,
	})
	if err != nil {
		return err
	}

	response := SignUpResponse{
		Message:     "verification code sent",
		ChallengeID: *challengeID,
	}

	apphttp.WriteJSON(w, http.StatusCreated, response)
	return nil
}

func (h *authHandler) VerifyAccount(w http.ResponseWriter, r *http.Request) error {
	var req VerifyParams

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return errors.ErrBadRequest
	}

	if req.ChallengeID == "" {
		return errors.ErrBadRequest
	}
	challengeID, err := uuid.Parse(req.ChallengeID)
	if err != nil {
		return errors.ErrBadRequest
	}
	if len(strconv.Itoa(req.OTP)) != 6 {
		return errors.ErrBadRequest
	}

	tokens, err := h.verify.Execute(usecase.VerifyAccountParams{
		UserAgent:   req.UserAgent,
		IPAddress:   req.IPAddress,
		ChallengeID: challengeID,
		OTP:         req.OTP,
	})
	if err != nil {
		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "chast",
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

	apphttp.WriteJSON(w, http.StatusCreated, response)
	return nil
}
