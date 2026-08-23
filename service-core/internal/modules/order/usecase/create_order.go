package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	appclock "service-core/internal/common/clock"
	apperrors "service-core/internal/common/errors"
	paymentgateway "service-core/internal/infra/payment-gateway"
	authenRepo "service-core/internal/modules/authentication/repository"
	cartDomain "service-core/internal/modules/cart/domain"
	cartRepo "service-core/internal/modules/cart/repository"
	inventoryRepo "service-core/internal/modules/inventory/repository"
	"service-core/internal/modules/order/domain"
	"service-core/internal/modules/order/repository"
	paymentDomain "service-core/internal/modules/payment/domain"
	paymentRepo "service-core/internal/modules/payment/repository"
	productDomain "service-core/internal/modules/product/domain"
	userRepo "service-core/internal/modules/user/repository"
	markdown "service-core/internal/shared/markdown"
	transaction "service-core/internal/shared/transaction"

	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
)

const (
	PAYMENT_PROVIDER   = "midtrans"
	PAYMENT_EXPIRATION = time.Hour * 24
)

type OrderItemInput struct {
	ProductID    *uuid.UUID
	CartItemID   *uuid.UUID
	IsCustom     bool
	CustomDesign json.RawMessage
	ProductName  string
	Quantity     int
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
	CustomerID      uuid.UUID
	AddressID       uuid.UUID
	PaymentMethodID uuid.UUID
	Shops           []OrderShopInput
}

type PaymentAccountResult struct {
	AccountName   string
	AccountNumber *string
	PhoneNumber   *string
	QRString      *string
}

type CreateOrderResult struct {
	OrderID        uuid.UUID
	PaymentAccount *PaymentAccountResult
	ChannelData    *paymentDomain.PaymentChannelData
	Instruction    *string
	Total          int64
}

type CreateOrderUsecase struct {
	executor               transaction.Executor
	transactor             transaction.Transactor
	accountRepo            authenRepo.AccountRepository
	orderRepo              repository.OrderRepository
	orderItemRepo          repository.OrderItemRepository
	customDesignRepo       repository.OrderItemCustomDesignRepository
	invoiceRepo            repository.InvoiceRepository
	invoiceItemRepo        repository.InvoiceItemRepository
	paymentRepo            paymentRepo.PaymentRepository
	paymentMethodRepo      paymentRepo.PaymentMethodRepository
	paymentEventRepo       paymentRepo.PaymentEventRepository
	paymentInstructionRepo paymentRepo.PaymentInstructionRepository
	paymentChannelDataRepo paymentRepo.PaymentChannelDataRepository
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
	customDesignRepo repository.OrderItemCustomDesignRepository,
	invoiceRepo repository.InvoiceRepository,
	invoiceItemRepo repository.InvoiceItemRepository,
	paymentRepo paymentRepo.PaymentRepository,
	paymentMethodRepo paymentRepo.PaymentMethodRepository,
	paymentEventRepo paymentRepo.PaymentEventRepository,
	paymentInstructionRepo paymentRepo.PaymentInstructionRepository,
	paymentChannelDataRepo paymentRepo.PaymentChannelDataRepository,
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
		customDesignRepo:       customDesignRepo,
		invoiceRepo:            invoiceRepo,
		invoiceItemRepo:        invoiceItemRepo,
		paymentRepo:            paymentRepo,
		paymentMethodRepo:      paymentMethodRepo,
		paymentEventRepo:       paymentEventRepo,
		paymentInstructionRepo: paymentInstructionRepo,
		paymentChannelDataRepo: paymentChannelDataRepo,
		inventoryRepo:          inventoryRepo,
		cartRepo:               cartRepo,
		userRepo:               userRepo,
		paymentGateway:         paymentGateway,
		pricingService:         pricingService,
	}
}

func (u *CreateOrderUsecase) Execute(
	ctx context.Context,
	input CreateOrderInput,
) (*CreateOrderResult, error) {
	now := appclock.Now()

	var (
		method        *paymentDomain.PaymentMethod
		pricingResult *repository.PricingResult
		customerName  string
		customerEmail string
		customerPhone string
	)

	// Concurrently fetch & validate payment method/pricing
	// and customer details
	g, gCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		var err error
		method, pricingResult, err = u.validateAndCalculatePricing(gCtx, input)
		return err
	})

	g.Go(func() error {
		var err error
		customerName, customerEmail, customerPhone, err = u.fetchCustomerDetails(gCtx, input.UserID)
		return err
	})

	if err := g.Wait(); err != nil {
		return nil, err
	}

	order := domain.Order{
		ID:          uuid.New(),
		Number:      domain.NewOrderNumber(),
		CustomerID:  input.CustomerID,
		AddressID:   input.AddressID,
		Status:      domain.OrderStatusPending,
		Subtotal:    pricingResult.Subtotal,
		ShippingFee: pricingResult.TotalShippingFee,
		Total:       pricingResult.GrandTotal,
		CreatedAt:   now,
	}
	if err := order.Validate(); err != nil {
		return nil, fmt.Errorf("invalid order domain state: %w", err)
	}

	invoice := order.NewInvoice()

	var (
		orderItems    []domain.OrderItem
		invoiceItems  []domain.InvoiceItem
		customDesigns []domain.OrderItemCustomDesign
		chargeItems   []paymentgateway.ChargeItem

		expiresAt = now.Add(PAYMENT_EXPIRATION)
	)

	// Generate order and invoice items from the pricing result to
	// preserve product, pricing, and shipping details at checkout time
	for _, shopRes := range pricingResult.Shops {
		var courierCode, courierService *string
		if shopRes.SelectedCourier.Code != "" {
			courierCode = &shopRes.SelectedCourier.Code
			courierService = &shopRes.SelectedCourier.Service
		}

		var inputShop *OrderShopInput
		for i := range input.Shops {
			if input.Shops[i].ShopID == shopRes.ShopID {
				inputShop = &input.Shops[i]
				break
			}
		}

		for itemIdx, itemRes := range shopRes.Items {
			var variantType cartDomain.ProductVariantType = cartDomain.ProductVariantTypeStandard
			if itemRes.IsCustom || itemRes.ProductID == nil {
				variantType = cartDomain.ProductVariantTypeCustom
			}

			orderItem := domain.OrderItem{
				ID:                 uuid.New(),
				OrderID:            order.ID,
				ProductVariantType: variantType,
				ShopID:             shopRes.ShopID,
				ShopName:           shopRes.ShopName,
				ProductID:          itemRes.ProductID,
				ProductName:        itemRes.ProductName,
				Quantity:           itemRes.Quantity,
				UnitPrice:          itemRes.UnitPrice,
				Subtotal:           itemRes.Subtotal,
				CourierCode:        courierCode,
				CourierService:     courierService,
				ShippingFee:        shopRes.SelectedCourier.Fee,
			}

			invoiceItem := invoice.NewInvoiceItemFromOrderItem(orderItem)

			orderItems = append(orderItems, orderItem)
			invoiceItems = append(invoiceItems, invoiceItem)

			if variantType == cartDomain.ProductVariantTypeCustom {
				var rawDesign json.RawMessage
				if inputShop != nil {
					if itemRes.CartItemID != nil {
						for _, inItm := range inputShop.Items {
							if inItm.CartItemID != nil && *inItm.CartItemID == *itemRes.CartItemID && len(inItm.CustomDesign) > 0 {
								rawDesign = inItm.CustomDesign
								break
							}
						}
					}
					if len(rawDesign) == 0 && itemIdx < len(inputShop.Items) && len(inputShop.Items[itemIdx].CustomDesign) > 0 {
						rawDesign = inputShop.Items[itemIdx].CustomDesign
					}
				}

				if len(rawDesign) > 0 {
					version := "3.0.0"
					physicalSizeID := productDomain.DEFAULT_PHYSICAL_SIZE_ID
					var previewURL *string
					var hUpper, bUpper, hLower, bLower *string

					if parsed, err := productDomain.ParseCustomDesignPayload(rawDesign); err == nil {
						if parsed.Metadata.Version != "" {
							version = parsed.Metadata.Version
						}
						size, prev, hu, bu, hl, bl := productDomain.ExtractDesignSummary(*parsed)
						physicalSizeID = size
						previewURL = prev
						hUpper, bUpper, hLower, bLower = hu, bu, hl, bl
					}

					customDesigns = append(customDesigns, domain.OrderItemCustomDesign{
						ID:              uuid.New(),
						OrderItemID:     orderItem.ID,
						Version:         version,
						PhysicalSizeID:  physicalSizeID,
						PreviewURL:      previewURL,
						HeaderTextUpper: hUpper,
						BodyTextUpper:   bUpper,
						HeaderTextLower: hLower,
						BodyTextLower:   bLower,
						DesignSnapshot:  rawDesign,
						CreatedAt:       now,
					})
				}
			}
		}
	}

	for _, item := range orderItems {
		var itemID string
		if item.ProductID != nil {
			itemID = item.ProductID.String()
		} else {
			itemID = item.ID.String()
		}
		cItem := paymentgateway.ChargeItem{
			ID:       itemID,
			Name:     item.ProductName,
			Quantity: item.Quantity,
			Price:    item.UnitPrice,
		}

		chargeItems = append(chargeItems, cItem)
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

	payment := paymentDomain.Payment{
		ID:        uuid.New(),
		OrderID:   order.ID,
		MethodID:  input.PaymentMethodID,
		Amount:    order.Total,
		Status:    paymentDomain.PaymentStatusPending,
		CreatedAt: now,
	}

	chargeResp, err := u.paymentGateway.Charge(
		ctx,
		paymentgateway.ChargeRequest{
			PaymentID:     payment.ID,
			OrderID:       order.ID,
			Amount:        order.Total,
			PaymentType:   string(method.Code),
			ExpiresAt:     expiresAt,
			CustomerEmail: customerEmail,
			CustomerName:  customerName,
			CustomerPhone: customerPhone,
			Items:         chargeItems,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("payment gateway charge failed: %w", err)
	}

	// Process payments through the configured gateway for
	// methods that require external payment handling.
	payment.Provider = method.Provider
	payment.ProviderPaymentID = &chargeResp.GatewayTransactionID
	payment.ProviderOrderID = &chargeResp.GatewayOrderID
	if !chargeResp.ExpiresAt.IsZero() {
		payment.ExpiresAt = &chargeResp.ExpiresAt
	}

	var (
		instructionContent   *string
		instruction          *paymentDomain.PaymentInstruction
		channelData          *paymentDomain.PaymentChannelData
		paymentAccountResult *PaymentAccountResult
	)

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

	err = u.transactor.WithinTransaction(ctx, func(exec transaction.Executor) error {
		if err := u.orderRepo.Save(ctx, exec, order); err != nil {
			return fmt.Errorf("failed to save order: %w", err)
		}

		if err := u.invoiceRepo.Save(ctx, exec, invoice); err != nil {
			return fmt.Errorf("failed to save invoice: %w", err)
		}

		if err := u.orderItemRepo.SaveBulk(ctx, exec, orderItems); err != nil {
			return fmt.Errorf("failed to save order items: %w", err)
		}

		if len(customDesigns) > 0 && u.customDesignRepo != nil {
			if err := u.customDesignRepo.SaveBulk(ctx, exec, customDesigns); err != nil {
				return fmt.Errorf("failed to save order item custom designs: %w", err)
			}
		}

		if err := u.invoiceItemRepo.SaveBulk(ctx, exec, invoiceItems); err != nil {
			return fmt.Errorf("failed to save invoice items: %w", err)
		}

		if err := u.paymentRepo.Save(ctx, exec, payment); err != nil {
			return fmt.Errorf("failed to save payment: %w", err)
		}

		if err := u.paymentEventRepo.Create(ctx, exec, paymentEvent); err != nil {
			return fmt.Errorf("failed to save payment event: %w", err)
		}

		// Persist gateway channel data (QR string, VA number, deep link)
		// so it survives the initial checkout response and can be
		// retrieved on any subsequent request.
		if chargeResp != nil && len(chargeResp.Instructions) > 0 {
			inst := chargeResp.Instructions[0]

			var actionURL *string
			if inst.Value != "" {
				v := inst.Value
				actionURL = &v
			}

			cd := paymentDomain.PaymentChannelData{
				ID:          uuid.New(),
				PaymentID:   payment.ID,
				ChannelType: method.Type,
				DisplayName: inst.Label,
				ActionURL:   actionURL,
				ExpiresAt:   payment.ExpiresAt,
				CreatedAt:   now,
			}

			if err := u.paymentChannelDataRepo.Save(ctx, exec, cd); err != nil {
				return fmt.Errorf("failed to save payment channel data: %w", err)
			}

			channelData = &cd
		}

		for _, item := range orderItems {
			if item.ProductID == nil {
				continue
			}
			if err := u.inventoryRepo.Reserve(ctx, exec,
				*item.ProductID,
				item.ShopID,
				item.Quantity,
			); err != nil {
				return fmt.Errorf("failed to reserve inventory for product %s: %w", item.ProductID, err)
			}
		}

		cart, err := u.cartRepo.GetWithItemsByCustomerID(ctx, exec, input.CustomerID)
		if err != nil {
			return fmt.Errorf("failed to load cart with items: %w", err)
		}
		if cart != nil {
			for _, shop := range pricingResult.Shops {
				for _, item := range shop.Items {
					if item.CartItemID != nil {
						if cart.RemoveItemByID(*item.CartItemID) {
							continue
						}
					}
					if item.ProductID != nil {
						if !cart.RemoveItem(*item.ProductID, shop.ShopID) {
							cart.RemoveProduct(*item.ProductID)
						}
					}
				}
			}

			if err := u.cartRepo.Save(ctx, exec, cart); err != nil {
				return fmt.Errorf("failed to update cart: %w", err)
			}
		}

		ins, err := u.paymentInstructionRepo.GetByPaymentMethodID(ctx, u.executor,
			input.PaymentMethodID,
		)
		if err != nil {
			return fmt.Errorf("failed to retrieve payment instruction: %w", err)
		}

		instruction = ins

		return nil
	})
	if err != nil {
		if payment.ProviderOrderID != nil {
			// Best-effort rollback:
			// the payment transaction has already been created
			// at the gateway, but persisting the order/payment
			// in the system database failed.
			//
			// Attempt to cancel the gateway transaction to avoid
			// leaving an orphaned payable transaction.
			_ = u.paymentGateway.CancelTransaction(ctx, *payment.ProviderOrderID)
		}

		return nil, err
	}

	// Render payment instructions with transaction-specific values
	// such as invoice number, amount, expiration time, and account details
	if instruction != nil {
		var vaNumber string
		if chargeResp != nil && len(chargeResp.Instructions) > 0 {
			inst := chargeResp.Instructions[0]
			vaNumber = inst.Value
		}

		var effectiveExpiresAt time.Time
		if payment.ExpiresAt != nil {
			effectiveExpiresAt = *payment.ExpiresAt
		} else {
			effectiveExpiresAt = expiresAt
		}

		content, err := markdown.Render(
			instruction.Content,
			map[string]string{
				"invoice_number": invoice.Number,
				"amount":         strconv.FormatInt(order.Total, 10),
				"expired_at":     effectiveExpiresAt.Format(time.RFC3339),
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
		ChannelData:    channelData,
		Instruction:    instructionContent,
		Total:          order.Total,
	}

	return &result, nil
}

func (u *CreateOrderUsecase) validateAndCalculatePricing(
	ctx context.Context,
	input CreateOrderInput,
) (*paymentDomain.PaymentMethod, *repository.PricingResult, error) {
	// Ensure the selected payment method exists and can be used
	// before creating any order-related records
	method, err := u.paymentMethodRepo.GetByID(ctx, u.executor, input.PaymentMethodID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to retrieve payment method: %w", err)
	}
	if method == nil {
		return nil, nil, apperrors.NewNotFound("payment method not found")
	}
	if !u.paymentGateway.Supports(string(method.Code)) {
		return nil, nil, apperrors.NewBadRequest(fmt.Sprintf("payment method %q is not supported by the payment gateway", method.Code))
	}

	pricingInput := repository.PricingInput{
		CustomerID:      input.CustomerID,
		AddressID:       &input.AddressID,
		PaymentMethodID: &input.PaymentMethodID,
		Shops:           make([]repository.PricingShopInput, 0, len(input.Shops)),
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
					ProductID:    item.ProductID,
					CartItemID:   item.CartItemID,
					IsCustom:     item.IsCustom,
					CustomDesign: item.CustomDesign,
					Quantity:     item.Quantity,
				},
			)
		}

		pricingInput.Shops = append(pricingInput.Shops, shopInput)
	}

	// Calculate the final order pricing, including item subtotals,
	// shipping fees, payment fees, and the grand total
	pricingResult, err := u.pricingService.Calculate(ctx, u.executor, pricingInput)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to calculate order prices: %w", err)
	}

	return method, pricingResult, nil
}

func (u *CreateOrderUsecase) fetchCustomerDetails(
	ctx context.Context,
	userID uuid.UUID,
) (string, string, string, error) {
	user, err := u.userRepo.GetByID(ctx, u.executor, userID)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to retrieve user: %w", err)
	}
	if user == nil {
		return "", "", "", apperrors.NewNotFound("user not found")
	}

	account, err := u.accountRepo.GetByUserID(ctx, u.executor, user.ID)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to retrieve account: %w", err)
	}
	if account == nil {
		return "", "", "", apperrors.NewNotFound("account not found")
	}

	var customerPhone string
	if user.Phone != nil {
		customerPhone = *user.Phone
	}

	return user.Name, account.Email, customerPhone, nil
}
