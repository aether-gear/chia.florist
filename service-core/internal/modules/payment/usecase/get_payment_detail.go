package usecase

import (
	"context"
	"fmt"
	"strconv"
	"time"

	apperrors "service-core/internal/common/errors"
	orderRepo "service-core/internal/modules/order/repository"
	paymentDomain "service-core/internal/modules/payment/domain"
	paymentRepo "service-core/internal/modules/payment/repository"
	markdown "service-core/internal/shared/markdown"
	transaction "service-core/internal/shared/transaction"

	"github.com/google/uuid"
)

type GetPaymentDetailUsecase struct {
	executor               transaction.Executor
	orderRepo              orderRepo.OrderRepository
	invoiceRepo            orderRepo.InvoiceRepository
	paymentRepo            paymentRepo.PaymentRepository
	paymentMethodRepo      paymentRepo.PaymentMethodRepository
	paymentInstructionRepo paymentRepo.PaymentInstructionRepository
	paymentChannelDataRepo paymentRepo.PaymentChannelDataRepository
}

func NewGetPaymentDetailUsecase(
	executor transaction.Executor,
	orderRepo orderRepo.OrderRepository,
	invoiceRepo orderRepo.InvoiceRepository,
	paymentRepo paymentRepo.PaymentRepository,
	paymentMethodRepo paymentRepo.PaymentMethodRepository,
	paymentInstructionRepo paymentRepo.PaymentInstructionRepository,
	paymentChannelDataRepo paymentRepo.PaymentChannelDataRepository,
) *GetPaymentDetailUsecase {
	return &GetPaymentDetailUsecase{
		executor:               executor,
		orderRepo:              orderRepo,
		invoiceRepo:            invoiceRepo,
		paymentRepo:            paymentRepo,
		paymentMethodRepo:      paymentMethodRepo,
		paymentInstructionRepo: paymentInstructionRepo,
		paymentChannelDataRepo: paymentChannelDataRepo,
	}
}

type GetPaymentDetailInput struct {
	OrderID uuid.UUID

	// Optional customer ownership enforcement
	CustomerID *uuid.UUID
}

type GetPaymentDetailResult struct {
	Payment        paymentDomain.Payment
	ChannelData    *paymentDomain.PaymentChannelData
	Instruction    *string
}

func (u *GetPaymentDetailUsecase) Execute(
	ctx context.Context,
	input GetPaymentDetailInput,
) (*GetPaymentDetailResult, error) {
	order, err := u.orderRepo.
		GetByID(ctx, u.executor,
			input.OrderID,
		)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve order: %w", err)
	}

	if order == nil {
		return nil, apperrors.NewNotFound("order not found")
	}

	if input.CustomerID != nil &&
		order.CustomerID != *input.CustomerID {
		return nil, apperrors.NewNotFound("order not found")
	}

	payment, err := u.paymentRepo.
		GetByOrderID(ctx, u.executor,
			order.ID,
		)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve payment: %w", err)
	}

	if payment == nil {
		return nil, apperrors.NewNotFound("payment not found")
	}

	var (
		channelData    *paymentDomain.PaymentChannelData
		vaNumber       string
	)

	cd, err := u.paymentChannelDataRepo.
		GetByPaymentID(ctx, u.executor,
			payment.ID,
		)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve payment channel data: %w", err)
	}

	channelData = cd
	if cd != nil && cd.ActionURL != nil {
		vaNumber = *cd.ActionURL
	}

	var renderedInstruction *string
	instruction, err := u.paymentInstructionRepo.
		GetByPaymentMethodID(ctx, u.executor,
			payment.MethodID,
		)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve payment instruction: %w", err)
	}

	if instruction != nil {
		invoice, err := u.invoiceRepo.
			GetByOrderID(ctx, u.executor,
				order.ID,
			)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve invoice: %w", err)
		}

		invoiceNumber := order.Number
		if invoice != nil {
			invoiceNumber = invoice.Number
		}

		expiredAtStr := ""
		if payment.ExpiresAt != nil {
			expiredAtStr = payment.ExpiresAt.Format(time.RFC3339)
		} else {
			expiredAtStr = payment.CreatedAt.Add(24 * time.Hour).Format(time.RFC3339)
		}

		vars := map[string]string{
			"invoice_number": invoiceNumber,
			"amount":         strconv.FormatInt(order.Total, 10),
			"expired_at":     expiredAtStr,
			"va_number":      vaNumber,
		}

		content, err := markdown.Render(instruction.Content, vars)
		if err != nil {
			return nil, fmt.Errorf("failed to format payment instruction: %w", err)
		}
		renderedInstruction = &content
	}

	return &GetPaymentDetailResult{
		Payment:        *payment,
		ChannelData:    channelData,
		Instruction:    renderedInstruction,
	}, nil
}
