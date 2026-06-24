package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	apperrors "service-core/internal/common/errors"
	paymentgateway "service-core/internal/infra/payment-gateway"
	cartRepo "service-core/internal/modules/cart/repository"
	inventoryRepo "service-core/internal/modules/inventory/repository"
	"service-core/internal/modules/order/domain"
	"service-core/internal/modules/order/repository"
	paymentDomain "service-core/internal/modules/payment/domain"
	paymentRepo "service-core/internal/modules/payment/repository"
	markdown "service-core/internal/shared/markdown"
	transaction "service-core/internal/shared/transaction"

	"github.com/google/uuid"
)

type CreateOrderUsecase struct {
	executor               transaction.Executor
	transactor             transaction.Transactor
	orderRepo              repository.OrderRepository
	orderItemRepo          repository.OrderItemRepository
	invoiceRepo            repository.InvoiceRepository
	invoiceItemRepo        repository.InvoiceItemRepository
	paymentRepo            paymentRepo.PaymentRepository
	paymentMethodRepo      paymentRepo.PaymentMethodRepository
	paymentAccRepo         paymentRepo.PaymentAccountRepository
	paymentEventRepo       paymentRepo.PaymentEventRepository
	paymentInstructionRepo paymentRepo.PaymentInstructionRepository
	inventoryRepo          inventoryRepo.InventoryRepository
	cartRepo               cartRepo.CartRepository
	paymentGateway         paymentgateway.Provider
	pricingService         repository.PricingService
}

func NewCreateOrderUsecase(
	executor transaction.Executor,
	transactor transaction.Transactor,
	orderRepo repository.OrderRepository,
	orderItemRepo repository.OrderItemRepository,
	invoiceRepo repository.InvoiceRepository,
	invoiceItemRepo repository.InvoiceItemRepository,
	paymentRepo paymentRepo.PaymentRepository,
	paymentMethodRepo paymentRepo.PaymentMethodRepository,
	paymentAccRepo paymentRepo.PaymentAccountRepository,
	paymentEventRepo paymentRepo.PaymentEventRepository,
	paymentInstructionRepo paymentRepo.PaymentInstructionRepository,
	inventoryRepo inventoryRepo.InventoryRepository,
	cartRepo cartRepo.CartRepository,
	paymentGateway paymentgateway.Provider,
	pricingService repository.PricingService,
) *CreateOrderUsecase {
	return &CreateOrderUsecase{
		executor:               executor,
		transactor:             transactor,
		orderRepo:              orderRepo,
		orderItemRepo:          orderItemRepo,
		invoiceRepo:            invoiceRepo,
		invoiceItemRepo:        invoiceItemRepo,
		paymentRepo:            paymentRepo,
		paymentMethodRepo:      paymentMethodRepo,
		paymentAccRepo:         paymentAccRepo,
		paymentEventRepo:       paymentEventRepo,
		paymentInstructionRepo: paymentInstructionRepo,
		inventoryRepo:          inventoryRepo,
		cartRepo:               cartRepo,
		paymentGateway:         paymentGateway,
		pricingService:         pricingService,
	}
}

type OrderItemInput struct {
	ProductID   uuid.UUID
	ProductName string
	Quantity    int
}

type OrderCourierInput struct {
	Code    string
	Service string
}

type OrderShopInput struct {
	ShopID   uuid.UUID
	ShopName string
	Courier  *OrderCourierInput
	Items    []OrderItemInput
}

type CreateOrderInput struct {
	UserID          uuid.UUID
	AddressID       uuid.UUID
	PaymentMethodID uuid.UUID
	IsManual        bool
	Shops           []OrderShopInput
}

type PaymentAccountResult struct {
	AccountName   string
	AccountNumber *string
	PhoneNumber   string
	QRString      *string
}

type CreateOrderResult struct {
	OrderID        uuid.UUID
	PaymentAccount *PaymentAccountResult
	Instruction    *string
	Total          int64
}

const PAYMENT_PROVIDER = "midtrans"

func (u *CreateOrderUsecase) Execute(
	ctx context.Context,
	input CreateOrderInput,
) (*CreateOrderResult, error) {
	now := time.Now()

	// Build a pricing request from the checkout input so all
	// product, shipping, and payment costs can be calculated
	pricingInput := repository.PricingInput{
		UserID:          input.UserID,
		AddressID:       &input.AddressID,
		PaymentMethodID: &input.PaymentMethodID,
		Shops: make(
			[]repository.PricingShopInput,
			0,
			len(input.Shops),
		),
	}

	for _, shop := range input.Shops {
		var courierCode, courierService *string
		if shop.Courier != nil {
			courierCode = &shop.Courier.Code
			courierService = &shop.Courier.Service
		}

		shopInput := repository.PricingShopInput{
			ShopID:         shop.ShopID,
			CourierCode:    courierCode,
			CourierService: courierService,
			Items:          make([]repository.PricingItemInput, 0, len(shop.Items)),
		}

		for _, item := range shop.Items {
			shopInput.Items = append(shopInput.Items, repository.PricingItemInput{
				ProductID: item.ProductID,
				Quantity:  item.Quantity,
			})
		}

		pricingInput.Shops = append(pricingInput.Shops, shopInput)
	}

	// Calculate the final order pricing, including item subtotals,
	// shipping fees, payment fees, and the grand total
	pricingResult, err :=
		u.pricingService.Calculate(
			ctx,
			u.executor,
			pricingInput,
		)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate order prices: %w", err)
	}

	// Ensure the selected payment method exists and can be used
	// before creating any order-related records
	method, err :=
		u.paymentMethodRepo.GetByID(
			ctx,
			u.executor,
			input.PaymentMethodID,
		)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve payment method: %w", err)
	}

	if method == nil {
		return nil, apperrors.NewNotFound("payment method not found")
	}

	order := domain.Order{
		ID:          uuid.New(),
		Number:      domain.NewOrderNumber(),
		UserID:      input.UserID,
		AddressID:   input.AddressID,
		Status:      domain.OrderStatusPending,
		Subtotal:    pricingResult.Subtotal,
		ShippingFee: pricingResult.TotalShippingFee,
		Total:       pricingResult.GrandTotal,
		CreatedAt:   now,
	}

	invoice := order.NewInvoice()

	var orderItems []domain.OrderItem
	var invoiceItems []domain.InvoiceItem

	// Generate order and invoice items from the pricing result to
	// preserve product, pricing, and shipping details at checkout time
	for _, shopRes := range pricingResult.Shops {
		var courierCode, courierService *string
		if shopRes.SelectedCourier.Code != "" {
			courierCode = &shopRes.SelectedCourier.Code
			courierService = &shopRes.SelectedCourier.Service
		}

		for _, itemRes := range shopRes.Items {
			orderItem := domain.OrderItem{
				ID:             uuid.New(),
				OrderID:        order.ID,
				ShopID:         shopRes.ShopID,
				ShopName:       shopRes.ShopName,
				ProductID:      itemRes.ProductID,
				ProductName:    itemRes.ProductName,
				Quantity:       itemRes.Quantity,
				UnitPrice:      itemRes.UnitPrice,
				Subtotal:       itemRes.Subtotal,
				CourierCode:    courierCode,
				CourierService: courierService,
				ShippingFee:    shopRes.SelectedCourier.Fee,
			}

			invoiceItem := invoice.NewInvoiceItemFromOrderItem(orderItem)

			orderItems = append(orderItems, orderItem)
			invoiceItems = append(invoiceItems, invoiceItem)
		}
	}

	// Delegate payment processing to the configured payment gateway
	// when the selected method is not handled manually.
	if !input.IsManual {
		// gatewayType, err := midtrans.MapPaymentType(method.Name)
		// if err != nil {
		// 	return nil, fmt.Errorf("payment method %q is not supported by the gateway: %w", method.Name, err)
		// }

		// return u.executeGatewayPayment(ctx, now, order, invoice, orderItems, invoiceItems, input, *method, gatewayType)

		return nil, apperrors.NewForbidden("service is not available yet")
	}

	// Select an available payment account for manual payment methods
	// using the current account load distribution strategy
	paymentAccount, err :=
		u.paymentAccRepo.RetrieveLeastLoaded(
			ctx,
			u.executor,
			input.PaymentMethodID,
		)
	if err != nil {
		return nil, fmt.Errorf("failed to acquire payment account: %w", err)
	}

	if paymentAccount == nil {
		return nil, apperrors.NewConflict("no available payment account for the selected method")
	}

	// Create the initial payment record in pending state before
	// any customer payment or provider confirmation is received
	payment := paymentDomain.Payment{
		ID:               uuid.New(),
		OrderID:          order.ID,
		MethodID:         input.PaymentMethodID,
		PaymentAccountID: &paymentAccount.ID,
		Provider:         PAYMENT_PROVIDER,
		Amount:           order.Total,
		Status:           paymentDomain.PaymentStatusPending,
		CreatedAt:        now,
	}

	pendingPayload, err := json.Marshal(map[string]string{"status": "pending"})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal pending payment event payload: %w", err)
	}

	paymentEvent := paymentDomain.PaymentEvent{
		ID:        uuid.New(),
		PaymentID: payment.ID,
		EventName: string(paymentDomain.PaymentEventStatusPending),
		Payload:   pendingPayload,
		CreatedAt: now,
	}

	var instruction *paymentDomain.PaymentInstruction
	err = u.transactor.WithinTransaction(
		ctx,
		func(exec transaction.Executor) error {
			if err :=
				u.orderRepo.Save(
					ctx,
					exec,
					order,
				); err != nil {
				return fmt.Errorf("failed to save order: %w", err)
			}

			if err :=
				u.invoiceRepo.Save(
					ctx,
					exec,
					invoice,
				); err != nil {
				return fmt.Errorf("failed to save invoice: %w", err)
			}

			if err :=
				u.orderItemRepo.SaveBulk(
					ctx,
					exec,
					orderItems,
				); err != nil {
				return fmt.Errorf("failed to save order items: %w", err)
			}

			if err :=
				u.invoiceItemRepo.SaveBulk(
					ctx,
					exec,
					invoiceItems,
				); err != nil {
				return fmt.Errorf("failed to save invoice items: %w", err)
			}

			if err :=
				u.paymentRepo.Save(
					ctx,
					exec,
					payment,
				); err != nil {
				return fmt.Errorf("failed to save payment: %w", err)
			}

			if err :=
				u.paymentEventRepo.Create(
					ctx,
					exec,
					paymentEvent,
				); err != nil {
				return fmt.Errorf("failed to save payment event: %w", err)
			}

			if err :=
				u.paymentAccRepo.IncrementLoad(
					ctx,
					exec,
					paymentAccount.ID,
				); err != nil {
				return fmt.Errorf("failed to increment payment account load: %w", err)
			}

			for _, item := range orderItems {
				if err :=
					u.inventoryRepo.Reserve(
						ctx, exec,
						item.ProductID,
						item.ShopID,
						item.Quantity,
					); err != nil {
					return fmt.Errorf("failed to reserve inventory for product %s: %w", item.ProductID, err)
				}
			}

			cart, err := u.cartRepo.GetWithItemsByUserID(ctx, exec, input.UserID)
			if err != nil {
				return fmt.Errorf("failed to load cart with items: %w", err)
			}
			if cart != nil {
				for _, item := range orderItems {
					cart.RemoveItem(item.ProductID, item.ShopID)
				}

				if err := u.cartRepo.Save(ctx, exec, cart); err != nil {
					return fmt.Errorf("failed to update cart: %w", err)
				}
			}

			ins, err :=
				u.paymentInstructionRepo.GetByPaymentID(
					ctx,
					u.executor,
					input.PaymentMethodID,
				)
			if err != nil {
				return fmt.Errorf("failed to retrieve payment instruction: %w", err)
			}
			instruction = ins

			return nil
		},
	)
	if err != nil {
		return nil, err
	}

	var instructionContent *string

	// Render payment instructions with transaction-specific values
	// such as invoice number, amount, expiration time, and account details
	if instruction != nil {
		content, err := markdown.Render(
			instruction.Content,
			map[string]string{
				"invoice_number": invoice.Number,
				"amount":         strconv.FormatInt(order.Total, 10),
				"expired_at":     now.Add(24 * time.Hour).Format(time.RFC3339),
				"va_number":      paymentAccount.ID.String(),
			},
		)
		if err != nil {
			return nil, fmt.Errorf("failed to format payment instruction: %w", err)
		}

		instructionContent = &content
	}

	return &CreateOrderResult{
		OrderID:     order.ID,
		Instruction: instructionContent,
		PaymentAccount: &PaymentAccountResult{
			AccountName:   paymentAccount.AccountName,
			AccountNumber: paymentAccount.AccountNumber,
			PhoneNumber:   paymentAccount.PhoneNumber,
			QRString:      paymentAccount.QRString,
		},
		Total: order.Total,
	}, nil
}

// executeGatewayPayment handles
// the Midtrans-backed charge path
//
// provider.Charge() is intentionally called OUTSIDE the DB transaction
//
// If the DB write subsequently fails, CancelTransaction is attempted
// as a best-effort cleanup to avoid dangling gateway transactions
func (u *CreateOrderUsecase) executeGatewayPayment(
	ctx context.Context,
	now time.Time,
	order domain.Order,
	invoice domain.Invoice,
	orderItems []domain.OrderItem,
	invoiceItems []domain.InvoiceItem,
	input CreateOrderInput,
	method paymentDomain.PaymentMethod,
	gatewayType string,
) (*CreateOrderResult, error) {
	payment := paymentDomain.Payment{
		ID:        uuid.New(),
		OrderID:   order.ID,
		MethodID:  input.PaymentMethodID,
		Provider:  PAYMENT_PROVIDER,
		Amount:    order.Total,
		Status:    paymentDomain.PaymentStatusPending,
		CreatedAt: now,
	}

	chargeResp, err := u.paymentGateway.
		Charge(ctx, paymentgateway.ChargeRequest{
			PaymentID:   payment.ID,
			OrderID:     order.ID,
			Amount:      order.Total,
			PaymentType: gatewayType,
		})
	if err != nil {
		return nil, fmt.Errorf("payment gateway charge failed: %w", err)
	}

	providerPaymentID := chargeResp.GatewayTransactionID
	providerOrderID := chargeResp.GatewayOrderID
	payment.ProviderPaymentID = &providerPaymentID
	payment.ProviderOrderID = &providerOrderID
	if !chargeResp.ExpiresAt.IsZero() {
		payment.ExpiresAt = &chargeResp.ExpiresAt
	}

	pendingPayload, err := json.Marshal(map[string]string{"status": "pending"})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal pending payment event payload: %w", err)
	}

	paymentEvent := paymentDomain.PaymentEvent{
		ID:        uuid.New(),
		PaymentID: payment.ID,
		EventName: string(paymentDomain.PaymentEventStatusPending),
		Payload:   pendingPayload,
		CreatedAt: now,
	}

	var instruction *paymentDomain.PaymentInstruction
	err = u.transactor.WithinTransaction(
		ctx,
		func(exec transaction.Executor) error {
			if err := u.orderRepo.
				Save(ctx, exec, order); err != nil {
				return fmt.Errorf("failed to save order: %w", err)
			}

			if err := u.invoiceRepo.
				Save(ctx, exec, invoice); err != nil {
				return fmt.Errorf("failed to save invoice: %w", err)
			}

			if err := u.orderItemRepo.
				SaveBulk(ctx, exec, orderItems); err != nil {
				return fmt.Errorf("failed to save order items: %w", err)
			}

			if err := u.invoiceItemRepo.
				SaveBulk(ctx, exec, invoiceItems); err != nil {
				return fmt.Errorf("failed to save invoice items: %w", err)
			}

			if err := u.paymentRepo.
				Save(ctx, exec, payment); err != nil {
				return fmt.Errorf("failed to save payment: %w", err)
			}

			if err := u.paymentEventRepo.
				Create(ctx, exec, paymentEvent); err != nil {
				return fmt.Errorf("failed to save payment event: %w", err)
			}

			for _, item := range orderItems {
				if err := u.inventoryRepo.Reserve(
					ctx, exec,
					item.ProductID,
					item.ShopID,
					item.Quantity,
				); err != nil {
					return fmt.Errorf(
						"failed to reserve inventory for product %s: %w",
						item.ProductID, err,
					)
				}
			}

			cart, err :=
				u.cartRepo.GetWithItemsByUserID(
					ctx,
					exec,
					input.UserID,
				)
			if err != nil {
				return fmt.Errorf("failed to load cart with items: %w", err)
			}

			if cart != nil {
				for _, item := range orderItems {
					cart.RemoveItem(item.ProductID, item.ShopID)
				}

				if err := u.cartRepo.Save(ctx, exec, cart); err != nil {
					return fmt.Errorf("failed to update cart: %w", err)
				}
			}

			ins, err := u.paymentInstructionRepo.
				GetByPaymentID(ctx, u.executor, input.PaymentMethodID)
			if err != nil {
				return fmt.Errorf("failed to retrieve payment instruction: %w", err)
			}
			instruction = ins

			return nil
		},
	)
	if err != nil {
		// Best-effort cancellation:
		// the gateway transaction exists but our
		//
		// DB write failed,
		// attempt to void it so it does not expire unpaid
		if payment.ProviderOrderID != nil {
			_ = u.paymentGateway.
				CancelTransaction(ctx, *payment.ProviderOrderID)
		}
		return nil, err
	}

	var instructionContent *string
	if instruction != nil {
		content, err := markdown.Render(
			instruction.Content,
			map[string]string{
				"invoice_number": invoice.Number,
				"amount":         strconv.FormatInt(order.Total, 10),
				"expired_at":     now.Add(24 * time.Hour).Format(time.RFC3339),
				"va_number":      providerPaymentID,
			},
		)
		if err != nil {
			return nil, fmt.Errorf("failed to format payment instruction: %w", err)
		}

		instructionContent = &content
	}

	return &CreateOrderResult{
		OrderID:        order.ID,
		Instruction:    instructionContent,
		PaymentAccount: nil,
		Total:          order.Total,
	}, nil
}
