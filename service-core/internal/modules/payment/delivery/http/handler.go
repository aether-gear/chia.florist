package http

import (
	"fmt"
	"net/http"
	"strconv"

	apperrors "service-core/internal/common/errors"
	apphttp "service-core/internal/common/http"
	"service-core/internal/modules/payment/usecase"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type PaymentHandler struct {
	createPaymentAccount  *usecase.CreatePaymentAccountUsecase
	listPaymentAccount    *usecase.ListPaymentAccountUsecase
	createPaymentMethod   *usecase.CreatePaymentMethodUsecase
	listPaymentMethod     *usecase.ListPaymentMethodUsecase
	processPaymentWebhook *usecase.ProcessPaymentWebhookUsecase
	processManualPayment  *usecase.ProcessManualPaymentUsecase
}

func NewPaymentHandler(
	createPaymentAccount *usecase.CreatePaymentAccountUsecase,
	listPaymentAccount *usecase.ListPaymentAccountUsecase,
	createPaymentMethod *usecase.CreatePaymentMethodUsecase,
	listPaymentMethod *usecase.ListPaymentMethodUsecase,
	processPaymentWebhook *usecase.ProcessPaymentWebhookUsecase,
	processManualPayment *usecase.ProcessManualPaymentUsecase,
) *PaymentHandler {
	return &PaymentHandler{
		createPaymentAccount:  createPaymentAccount,
		listPaymentAccount:    listPaymentAccount,
		createPaymentMethod:   createPaymentMethod,
		listPaymentMethod:     listPaymentMethod,
		processPaymentWebhook: processPaymentWebhook,
		processManualPayment:  processManualPayment,
	}
}

func (h *PaymentHandler) CreatePaymentAccount(w http.ResponseWriter, r *http.Request) error {
	var req createPaymentAccountRequest

	if err := apphttp.DecodeJSON(r, &req); err != nil {
		return apperrors.NewBadRequest("invalid body request")
	}

	methodID, err := uuid.Parse(req.MethodID)
	if err != nil {
		return apperrors.NewBadRequest("invalid method id")
	}
	isActive, err := strconv.ParseBool(req.IsActive)
	if err != nil {
		return apperrors.NewBadRequest("invalid active status")
	}
	if req.AccountName == "" {
		return apperrors.NewBadRequest("invalid account name")
	}

	input := usecase.CreatePaymentAccountInput{
		MethodID:      methodID,
		AccountName:   req.AccountName,
		AccountNumber: req.AccountNumber,
		PhoneNumber:   req.PhoneNumber,
		QRString:      req.QRString,
		IsActive:      isActive,
	}

	err = h.createPaymentAccount.Execute(r.Context(), input)
	if err != nil {
		return err
	}

	response := map[string]string{
		"message": "payment account successfully created",
	}

	apphttp.WriteJSON(w, http.StatusOK, response)
	return nil
}

func (h *PaymentHandler) ListPaymentAccount(w http.ResponseWriter, r *http.Request) error {
	payAccs, err := h.listPaymentAccount.ListAll(r.Context())
	if err != nil {
		return err
	}

	paymentAccs := make([]paymentAccountResponse, 0, len(payAccs))
	for _, p := range payAccs {
		pA := paymentAccountResponse{
			ID:            p.ID,
			MethodID:      p.MethodID,
			AccountName:   p.AccountName,
			AccountNumber: p.AccountNumber,
			PhoneNumber:   p.PhoneNumber,
			QRString:      p.QRString,
		}

		paymentAccs = append(paymentAccs, pA)
	}

	response := map[string]interface{}{
		"accounts": paymentAccs,
	}

	apphttp.WriteJSON(w, http.StatusOK, response)
	return nil
}

func (h *PaymentHandler) CreatePaymentMethod(w http.ResponseWriter, r *http.Request) error {
	var req createPaymentMethodRequest

	if err := apphttp.DecodeJSON(r, &req); err != nil {
		return apperrors.NewBadRequest("invalid body request")
	}
	isActive, err := strconv.ParseBool(req.IsActive)
	if err != nil {
		return apperrors.NewBadRequest("invalid active status")
	}
	if req.Name == "" {
		return apperrors.NewBadRequest("invalid name")
	}
	if req.Type == "" {
		return apperrors.NewBadRequest("invalid type")
	}
	if req.FeeFixed == nil || *req.FeeFixed == "" {
		return apperrors.NewBadRequest("invalid fee amount")
	}
	if req.FeePercentage == nil || *req.FeePercentage == "" {
		return apperrors.NewBadRequest("invalid fee percentage")
	}

	var feeFixed int64
	var feePercentage float64

	if req.FeeFixed == nil || *req.FeeFixed == "" {
		feeFixed = 0
	} else {
		val, err := strconv.ParseInt(*req.FeeFixed, 10, 64)
		if err != nil {
			return apperrors.NewBadRequest("invalid fee amount")
		}
		feeFixed = val
	}

	if req.FeePercentage == nil || *req.FeePercentage == "" {
		feePercentage = 0
	} else {
		val, err := strconv.ParseFloat(*req.FeePercentage, 64)
		if err != nil {
			return apperrors.NewBadRequest("invalid fee percentage")
		}
		feePercentage = val
	}

	input := usecase.CreatePaymentMethodInput{
		Name:          req.Name,
		Type:          req.Type,
		IsActive:      isActive,
		Description:   req.Description,
		FeeType:       req.FeeType,
		FeeFixed:      feeFixed,
		FeePercentage: feePercentage,
	}

	err = h.createPaymentMethod.Execute(r.Context(), input)
	if err != nil {
		return err
	}

	response := map[string]string{
		"message": "payment method successfully created",
	}

	apphttp.WriteJSON(w, http.StatusOK, response)
	return nil
}

func (h *PaymentHandler) ListPaymentMethod(w http.ResponseWriter, r *http.Request) error {
	payMethods, err := h.listPaymentMethod.ListAll(r.Context())
	if err != nil {
		return err
	}

	paymentMthds := make([]paymentMethodResponse, 0, len(payMethods))
	for _, p := range payMethods {
		pM := paymentMethodResponse{
			ID:            p.ID,
			Name:          p.Name,
			Type:          string(p.Type),
			IsActive:      p.IsActive,
			Description:   p.Description,
			FeeType:       string(p.FeeType),
			FeeFixed:      p.FeeFixed,
			FeePercentage: p.FeePercentage,
		}

		paymentMthds = append(paymentMthds, pM)
	}

	response := map[string]interface{}{
		"methods": paymentMthds,
	}

	apphttp.WriteJSON(w, http.StatusOK, response)
	return nil
}

func (h *PaymentHandler) HandleMidtransWebhook(w http.ResponseWriter, r *http.Request) error {
	var payload map[string]any
	if err := apphttp.DecodeJSON(r, &payload); err != nil {
		return apperrors.NewBadRequest("invalid payload")
	}

	input := usecase.ProcessPaymentWebhookInput{
		Payload: payload,
	}

	if err := h.processPaymentWebhook.
		Execute(r.Context(), input); err != nil {
		return err
	}

	response := map[string]string{
		"message": "webhook processed successfully",
	}

	apphttp.WriteJSON(w, http.StatusOK, response)
	return nil
}

func (h *PaymentHandler) ProcessManualPayment(w http.ResponseWriter, r *http.Request) error {
	idStr := chi.URLParam(r, "id")
	paymentID, err := uuid.Parse(idStr)
	if err != nil {
		return apperrors.NewBadRequest("invalid payment ID")
	}

	var req manualPaymentActionRequest
	if err := apphttp.DecodeJSON(r, &req); err != nil {
		return apperrors.NewBadRequest("invalid body request")
	}

	input := usecase.ProcessManualPaymentInput{
		PaymentID: paymentID,
		Action:    req.Action,
	}

	if err := h.processManualPayment.
		Execute(r.Context(), input); err != nil {
		return err
	}

	response := map[string]string{
		"message": fmt.Sprintf("manual payment successfully updated with action: %s", req.Action),
	}

	apphttp.WriteJSON(w, http.StatusOK, response)
	return nil
}
