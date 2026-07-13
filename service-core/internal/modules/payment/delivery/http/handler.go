package http

import (
	"fmt"
	"net/http"
	"strconv"

	apperrors "service-core/internal/common/errors"
	apphttp "service-core/internal/common/http"
	authenDomain "service-core/internal/modules/authentication/domain"
	"service-core/internal/modules/payment/usecase"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type PaymentHandler struct {
	createPaymentAccount   *usecase.CreatePaymentAccountUsecase
	listPaymentAccount     *usecase.ListPaymentAccountUsecase
	savePaymentMethod      *usecase.SavePaymentMethodUsecase
	listPaymentMethod      *usecase.ListPaymentMethodUsecase
	processPaymentWebhook  *usecase.ProcessPaymentWebhookUsecase
	processManualPayment   *usecase.ProcessManualPaymentUsecase
	savePaymentInstruction *usecase.SavePaymentInstructionUsecase
	getPaymentDetail       *usecase.GetPaymentDetailUsecase
	checkPaymentStatus     *usecase.CheckPaymentStatusUsecase
}

func NewPaymentHandler(
	createPaymentAccount *usecase.CreatePaymentAccountUsecase,
	listPaymentAccount *usecase.ListPaymentAccountUsecase,
	savePaymentMethod *usecase.SavePaymentMethodUsecase,
	listPaymentMethod *usecase.ListPaymentMethodUsecase,
	processPaymentWebhook *usecase.ProcessPaymentWebhookUsecase,
	processManualPayment *usecase.ProcessManualPaymentUsecase,
	savePaymentInstruction *usecase.SavePaymentInstructionUsecase,
	getPaymentDetail *usecase.GetPaymentDetailUsecase,
	checkPaymentStatus *usecase.CheckPaymentStatusUsecase,
) *PaymentHandler {
	return &PaymentHandler{
		createPaymentAccount:   createPaymentAccount,
		listPaymentAccount:     listPaymentAccount,
		savePaymentMethod:      savePaymentMethod,
		listPaymentMethod:      listPaymentMethod,
		processPaymentWebhook:  processPaymentWebhook,
		processManualPayment:   processManualPayment,
		savePaymentInstruction: savePaymentInstruction,
		getPaymentDetail:       getPaymentDetail,
		checkPaymentStatus:     checkPaymentStatus,
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

func (h *PaymentHandler) SavePaymentMethod(w http.ResponseWriter, r *http.Request) error {
	var req savePaymentMethodRequest

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
	if req.Code == "" {
		return apperrors.NewBadRequest("invalid code")
	}
	if req.Provider == "" {
		return apperrors.NewBadRequest("invalid provider")
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

	var methodID *uuid.UUID
	if req.ID != nil && *req.ID != "" {
		id, err := uuid.Parse(*req.ID)
		if err != nil {
			return apperrors.NewBadRequest("invalid id format")
		}
		methodID = &id
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
		ID:            methodID,
		Name:          req.Name,
		Code:          req.Code,
		Provider:      req.Provider,
		Type:          req.Type,
		IsActive:      isActive,
		Description:   req.Description,
		FeeType:       req.FeeType,
		FeeFixed:      feeFixed,
		FeePercentage: feePercentage,
	}

	err = h.savePaymentMethod.Execute(r.Context(), input)
	if err != nil {
		return err
	}

	msg := "payment method successfully created"
	if methodID != nil {
		msg = "payment method successfully updated"
	}

	response := map[string]string{
		"message": msg,
	}

	apphttp.WriteJSON(w, http.StatusOK, response)
	return nil
}

func (h *PaymentHandler) ListPaymentMethod(w http.ResponseWriter, r *http.Request) error {
	sortParam := r.URL.Query().Get("sort")
	input := usecase.ListPaymentMethodInput{
		Sort: sortParam,
	}

	payMethods, err := h.listPaymentMethod.ListAll(r.Context(), input)
	if err != nil {
		return err
	}

	paymentMthds := make([]paymentMethodResponse, 0, len(payMethods))
	for _, p := range payMethods {
		pM := paymentMethodResponse{
			ID:            p.ID,
			Name:          p.Name,
			Code:          string(p.Code),
			Provider:      p.Provider,
			Type:          string(p.Type),
			IsActive:      p.IsActive,
			Description:   p.Description,
			FeeType:       string(p.FeeType),
			FeeFixed:      p.FeeFixed,
			FeePercentage: p.FeePercentage,
		}

		if p.Instruction != nil {
			pM.Instruction = &paymentInstructionResponse{
				ID:        p.Instruction.ID,
				Content:   p.Instruction.Content,
				CreatedAt: p.Instruction.CreatedAt,
			}
		}

		paymentMthds = append(paymentMthds, pM)
	}

	response := map[string][]paymentMethodResponse{
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

func (h *PaymentHandler) SavePaymentInstruction(w http.ResponseWriter, r *http.Request) error {
	methodIDStr := chi.URLParam(r, "methodID")
	methodID, err := uuid.Parse(methodIDStr)
	if err != nil {
		return apperrors.NewBadRequest("invalid payment method ID")
	}

	var req savePaymentInstructionRequest
	if err := apphttp.DecodeJSON(r, &req); err != nil {
		return apperrors.NewBadRequest("invalid body request")
	}

	if req.Content == "" {
		return apperrors.NewBadRequest("content cannot be empty")
	}

	input := usecase.SavePaymentInstructionInput{
		PaymentMethodID: methodID,
		Content:         req.Content,
	}

	err = h.savePaymentInstruction.Execute(r.Context(), input)
	if err != nil {
		return err
	}

	response := map[string]string{
		"message": "payment instruction successfully saved",
	}

	apphttp.WriteJSON(w, http.StatusOK, response)
	return nil
}

func (h *PaymentHandler) GetMyOrderPayment(w http.ResponseWriter, r *http.Request) error {
	authCtx, ok := authenDomain.GetAuthContext(r.Context())
	if !ok || !authCtx.IsAuthenticated {
		return apperrors.NewUnauthorized("authentication required")
	}
	if authCtx.CustomerID == nil {
		return apperrors.NewForbidden("customer account required")
	}

	customerID := *authCtx.CustomerID

	orderIDStr := chi.URLParam(r, "orderID")
	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		return apperrors.NewBadRequest("invalid order id")
	}

	result, err := h.getPaymentDetail.Execute(r.Context(), usecase.GetPaymentDetailInput{
		OrderID:    orderID,
		CustomerID: &customerID,
	})
	if err != nil {
		return err
	}

	resp := getPaymentDetailResponse{
		PaymentID: result.Payment.ID.String(),
		Status:    string(result.Payment.Status),
		Amount:    result.Payment.Amount,
		ExpiresAt: result.Payment.ExpiresAt,
	}

	if result.ChannelData != nil {
		channelTypeStr := string(result.ChannelData.ChannelType)
		resp.ChannelType = &channelTypeStr
		resp.DisplayName = &result.ChannelData.DisplayName
		resp.ActionURL = result.ChannelData.ActionURL
	}

	if result.PaymentAccount != nil {
		resp.AccountName = &result.PaymentAccount.AccountName
		resp.AccountNumber = result.PaymentAccount.AccountNumber
		resp.PhoneNumber = result.PaymentAccount.PhoneNumber
		resp.QRString = result.PaymentAccount.QRString
	}

	resp.Instruction = result.Instruction

	apphttp.WriteJSON(w, http.StatusOK, resp)
	return nil
}

// CheckMyOrderPaymentStatus is the customer-triggered payment sync endpoint.
//
// POST /users/me/orders/{orderID}/payment/check
//
// When a customer's payment appears stuck as 'pending' after they have paid,
// calling this endpoint immediately queries Midtrans for the current status
// and resolves the payment — without waiting for the background reconciler.
func (h *PaymentHandler) CheckMyOrderPaymentStatus(w http.ResponseWriter, r *http.Request) error {
	authCtx, ok := authenDomain.GetAuthContext(r.Context())
	if !ok || !authCtx.IsAuthenticated {
		return apperrors.NewUnauthorized("authentication required")
	}
	if authCtx.CustomerID == nil {
		return apperrors.NewForbidden("customer account required")
	}

	orderIDStr := chi.URLParam(r, "orderID")
	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		return apperrors.NewBadRequest("invalid order id")
	}

	result, err := h.checkPaymentStatus.Execute(
		r.Context(),
		usecase.CheckPaymentStatusInput{
			OrderID:    orderID,
			CustomerID: *authCtx.CustomerID,
		},
	)
	if err != nil {
		return err
	}

	apphttp.WriteJSON(w, http.StatusOK, checkPaymentStatusResponse{
		Status: string(result.Status),
		Synced: result.Synced,
	})
	return nil
}
