package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"service-core/internal/common/errors"
	apphttp "service-core/internal/common/http"
	"service-core/internal/modules/payment/usecase"

	"github.com/google/uuid"
)

type PaymentHandler struct {
	createPaymentAccount *usecase.CreatePaymentAccount
	listPaymentAccount   *usecase.ListPaymentAccount
	createPaymentMethod  *usecase.CreatePaymentMethod
	listPaymentMethod    *usecase.ListPaymentMethod
}

func NewPaymentHandler(
	cPA *usecase.CreatePaymentAccount,
	lPA *usecase.ListPaymentAccount,
	cPM *usecase.CreatePaymentMethod,
	lPM *usecase.ListPaymentMethod,
) *PaymentHandler {
	return &PaymentHandler{
		createPaymentAccount: cPA,
		listPaymentAccount:   lPA,
		createPaymentMethod:  cPM,
		listPaymentMethod:    lPM,
	}
}

func (h *PaymentHandler) CreatePaymentAccount(w http.ResponseWriter, r *http.Request) error {
	var req CreatePaymentAccountRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return errors.ErrBadRequest
	}

	methodID, err := uuid.Parse(req.MethodID)
	if err != nil {
		return errors.ErrBadRequest
	}

	isActive, err := strconv.ParseBool(req.IsActive)
	if err != nil {
		return errors.ErrBadRequest
	}

	if req.AccountName == "" {
		return errors.ErrBadRequest
	}

	input := usecase.CreatePaymentAccountInput{
		MethodID:      methodID,
		AccountName:   req.AccountName,
		AccountNumber: req.AccountNumber,
		PhoneNumber:   req.PhoneNumber,
		QRString:      req.QRString,
		IsActive:      isActive,
	}

	fmt.Println("kdlneaopdeadpaedpiay")
	err = h.createPaymentAccount.Execute(input)
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
	payAccs, err := h.listPaymentAccount.ListAll()
	if err != nil {
		return err
	}

	result := make([]PaymentAccountResponse, 0, len(payAccs))
	for _, p := range payAccs {
		res := PaymentAccountResponse{
			ID:            p.ID,
			MethodID:      p.MethodID,
			AccountName:   p.AccountName,
			AccountNumber: p.AccountNumber,
			PhoneNumber:   p.PhoneNumber,
			QRString:      p.QRString,
		}

		result = append(result, res)
	}

	response := map[string]interface{}{
		"accounts": result,
	}

	apphttp.WriteJSON(w, http.StatusOK, response)
	return nil
}

func (h *PaymentHandler) CreatePaymentMethod(w http.ResponseWriter, r *http.Request) error {

	var req CreatePaymentMethodRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return errors.ErrBadRequest
	}

	isActive, err := strconv.ParseBool(req.IsActive)
	if err != nil {
		return errors.ErrBadRequest
	}

	if req.Name == "" || req.Type == "" {
		return errors.ErrBadRequest
	}

	if (req.FeeFixed == nil || *req.FeeFixed == "") &&
		(req.FeePercentage == nil || *req.FeePercentage == "") {
		return errors.ErrBadRequest
	}

	var feeFixed int64
	var feePercentage float64

	if req.FeeFixed == nil || *req.FeeFixed == "" {
		feeFixed = 0
	} else {
		val, err := strconv.ParseInt(*req.FeeFixed, 10, 64)
		if err != nil {
			return errors.ErrBadRequest
		}
		feeFixed = val
	}

	if req.FeePercentage == nil || *req.FeePercentage == "" {
		feePercentage = 0
	} else {
		val, err := strconv.ParseFloat(*req.FeePercentage, 64)
		if err != nil {
			return errors.ErrBadRequest
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

	fmt.Println("ajodjeaiponoai1")

	err = h.createPaymentMethod.Execute(input)
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
	payMethods, err := h.listPaymentMethod.ListAll()
	if err != nil {
		return err
	}

	result := make([]PaymentMethodResponse, 0, len(payMethods))
	for _, p := range payMethods {
		res := PaymentMethodResponse{
			ID:            p.ID,
			Name:          p.Name,
			Type:          string(p.Type),
			IsActive:      p.IsActive,
			Description:   p.Description,
			FeeType:       string(p.FeeType),
			FeeFixed:      p.FeeFixed,
			FeePercentage: p.FeePercentage,
		}

		result = append(result, res)
	}

	response := map[string]interface{}{
		"methods": result,
	}

	apphttp.WriteJSON(w, http.StatusOK, response)
	return nil
}
