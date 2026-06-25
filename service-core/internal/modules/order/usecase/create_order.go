package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	apperrors "service-core/internal/common/errors"
	paymentgateway "service-core/internal/infra/payment-gateway"
	authenRepo "service-core/internal/modules/authentication/repository"
	cartRepo "service-core/internal/modules/cart/repository"
	inventoryRepo "service-core/internal/modules/inventory/repository"
	"service-core/internal/modules/order/domain"
	"service-core/internal/modules/order/repository"
	paymentDomain "service-core/internal/modules/payment/domain"
	paymentRepo "service-core/internal/modules/payment/repository"
	userRepo "service-core/internal/modules/user/repository"
	markdown "service-core/internal/shared/markdown"
	transaction "service-core/internal/shared/transaction"

	"github.com/google/uuid"
)

type CreateOrderUsecase struct {
	executor               transaction.Executor
	transactor             transaction.Transactor
	accountRepo            authenRepo.AccountRepository
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
	userRepo               userRepo.UserRepository
	paymentGateway         paymentgateway.Provider
	pricingService         repository.PricingService
}

func NewCreateOrderUsecase(
	executor transaction.Executor,
	transactor transaction.Transactor,
	accountRepo authenRepo.AccountRepository,
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
	userRepo userRepo.UserRepository,
	paymentGateway paymentgateway.Provider,
	pricingService repository.PricingService,
) *CreateOrderUsecase {
	return &CreateOrderUsecase{
		executor:               executor,
		transactor:             transactor,
		accountRepo:            accountRepo,
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
		userRepo:               userRepo,
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

	// Ensure the selected payment method exists and can be used
	// before creating any order-related records
	method, err := u.paymentMethodRepo.
		GetByID(ctx, u.executor, input.PaymentMethodID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve payment method: %w", err)
	}

	if method == nil {
		return nil, apperrors.NewNotFound("payment method not found")
	}

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
			shopInput.Items = append(
				shopInput.Items,
				repository.PricingItemInput{
					ProductID: item.ProductID,
					Quantity:  item.Quantity,
				},
			)
		}

		pricingInput.Shops = append(pricingInput.Shops, shopInput)
	}

	// Calculate the final order pricing, including item subtotals,
	// shipping fees, payment fees, and the grand total
	pricingResult, err := u.pricingService.
		Calculate(ctx, u.executor, pricingInput)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate order prices: %w", err)
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

	payment := paymentDomain.Payment{
		ID:        uuid.New(),
		OrderID:   order.ID,
		MethodID:  input.PaymentMethodID,
		Amount:    order.Total,
		Status:    paymentDomain.PaymentStatusPending,
		CreatedAt: now,
	}

	user, err := u.userRepo.
		GetByID(ctx, u.executor, input.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve user: %w", err)
	}

	account, err := u.accountRepo.
		GetByUserID(ctx, u.executor, user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account: %w", err)
	}

	// Process payments through the configured gateway for
	// methods that require external payment handling.
	var (
		paymentAccount       *paymentDomain.PaymentAccount
		paymentAccountResult *PaymentAccountResult
		chargeResp           *paymentgateway.ChargeResponse
	)

	if !input.IsManual {
		provider := PAYMENT_PROVIDER
		payment.Provider = provider

		var chargeItems []paymentgateway.ChargeItem
		for _, item := range orderItems {
			chargeItems = append(chargeItems, paymentgateway.ChargeItem{
				ID:       item.ProductID.String(),
				Name:     item.ProductName,
				Quantity: item.Quantity,
				Price:    item.UnitPrice,
			})
		}

		if pricingResult.TotalShippingFee > 0 {
			chargeItems = append(chargeItems, paymentgateway.ChargeItem{
				ID:       "shipping_fee",
				Name:     "Shipping Fee",
				Quantity: 1,
				Price:    pricingResult.TotalShippingFee,
			})
		}

		var itemSum int64
		for _, ci := range chargeItems {
			itemSum += ci.Price * int64(ci.Quantity)
		}

		if diff := order.Total - itemSum; diff != 0 {
			chargeItems = append(chargeItems, paymentgateway.ChargeItem{
				ID:       "adjustment",
				Name:     "Fees & Adjustments",
				Quantity: 1,
				Price:    diff,
			})
		}

		var err error
		chargeResp, err = u.paymentGateway.
			Charge(
				ctx,
				paymentgateway.ChargeRequest{
					PaymentID:     payment.ID,
					OrderID:       order.ID,
					Amount:        order.Total,
					PaymentType:   method.Name,
					ExpiresAt:     time.Now().Add(time.Hour * 24),
					CustomerEmail: account.Email,
					CustomerName:  user.Name,
					CustomerPhone: *user.Phone,
					Items:         chargeItems,
				},
			)
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

		if len(chargeResp.Instructions) > 0 {
			inst := chargeResp.Instructions[0]
			accountResult := &PaymentAccountResult{
				AccountName: inst.Label,
			}
			switch inst.Type {
			case "qris", "ewallet":
				accountResult.QRString = &inst.Value
			case "bank_transfer":
				accountResult.AccountNumber = &inst.Value
			}
			paymentAccountResult = accountResult
		}

	} else {
		// Assign a payment account for manual payment methods
		// using the current load-balancing strategy
		var err error
		paymentAccount, err = u.paymentAccRepo.
			RetrieveLeastLoaded(ctx, u.executor, input.PaymentMethodID)
		if err != nil {
			return nil, fmt.Errorf("failed to acquire payment account: %w", err)
		}

		if paymentAccount == nil {
			return nil, apperrors.NewConflict("no available payment account for the selected method")
		}

		provider := "manual"
		payment.Provider = provider
		payment.PaymentAccountID = &paymentAccount.ID

		paymentAccountResult = &PaymentAccountResult{
			AccountName:   paymentAccount.AccountName,
			AccountNumber: paymentAccount.AccountNumber,
			PhoneNumber:   paymentAccount.PhoneNumber,
			QRString:      paymentAccount.QRString,
		}
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

			if input.IsManual {
				if err := u.paymentAccRepo.
					IncrementLoad(ctx, exec, *payment.PaymentAccountID); err != nil {
					return fmt.Errorf("failed to increment payment account load: %w", err)
				}
			}

			for _, item := range orderItems {
				if err := u.inventoryRepo.
					Reserve(ctx, exec, item.ProductID, item.ShopID, item.Quantity); err != nil {
					return fmt.Errorf("failed to reserve inventory for product %s: %w", item.ProductID, err)
				}
			}

			cart, err := u.cartRepo.
				GetWithItemsByUserID(ctx, exec, input.UserID)
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
				GetByPaymentMethodID(ctx, u.executor, input.PaymentMethodID)
			if err != nil {
				return fmt.Errorf("failed to retrieve payment instruction: %w", err)
			}
			instruction = ins

			return nil
		},
	)
	if err != nil {
		if !input.IsManual &&
			payment.ProviderOrderID != nil {
			// Best-effort cancellation:
			// the gateway transaction exists but our
			//
			// DB write failed,
			// attempt to void it so it does not expire unpaid
			_ = u.paymentGateway.
				CancelTransaction(ctx, *payment.ProviderOrderID)
		}

		return nil, err
	}

	var instructionContent *string

	// Render payment instructions with transaction-specific values
	// such as invoice number, amount, expiration time, and account details
	if instruction != nil {
		var vaNumber string
		if input.IsManual {
			if paymentAccount != nil {
				if paymentAccount.AccountNumber != nil {
					vaNumber = *paymentAccount.AccountNumber
				} else if paymentAccount.PhoneNumber != "" {
					vaNumber = paymentAccount.PhoneNumber
				} else if paymentAccount.QRString != nil {
					vaNumber = *paymentAccount.QRString
				}
			}
		} else {
			if chargeResp != nil && len(chargeResp.Instructions) > 0 {
				vaNumber = chargeResp.Instructions[0].Value
			}
		}

		content, err := markdown.Render(
			instruction.Content,
			map[string]string{
				"invoice_number": invoice.Number,
				"amount":         strconv.FormatInt(order.Total, 10),
				"expired_at":     now.Add(24 * time.Hour).Format(time.RFC3339),
				"va_number":      vaNumber,
			},
		)
		if err != nil {
			return nil, fmt.Errorf("failed to format payment instruction: %w", err)
		}

		instructionContent = &content
	}

	result := CreateOrderResult{
		OrderID:        order.ID,
		PaymentAccount: paymentAccountResult,
		Instruction:    instructionContent,
		Total:          order.Total,
	}

	return &result, nil
}
