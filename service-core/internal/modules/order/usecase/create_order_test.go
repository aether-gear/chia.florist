package usecase

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	paymentgateway "service-core/internal/infra/payment-gateway"
	authenDomain "service-core/internal/modules/authentication/domain"
	cartDomain "service-core/internal/modules/cart/domain"
	inventoryDomain "service-core/internal/modules/inventory/domain"
	orderDomain "service-core/internal/modules/order/domain"
	orderRepo "service-core/internal/modules/order/repository"
	paymentDomain "service-core/internal/modules/payment/domain"
	userDomain "service-core/internal/modules/user/domain"
	userRepo "service-core/internal/modules/user/repository"
	transaction "service-core/internal/shared/transaction"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// ===========================================================================
// Mock infrastructure (all in this package)
// ===========================================================================

// --- transaction mocks ---

type coMockExecutor struct{}

func (m *coMockExecutor) Exec(_ context.Context, _ string, _ ...any) (pgconn.CommandTag, error) {
	return pgconn.NewCommandTag(""), nil
}
func (m *coMockExecutor) Query(_ context.Context, _ string, _ ...any) (pgx.Rows, error) {
	return nil, nil
}
func (m *coMockExecutor) QueryRow(_ context.Context, _ string, _ ...any) pgx.Row {
	return nil
}

type coMockTransactor struct{}

func (m *coMockTransactor) WithinTransaction(ctx context.Context, fn func(transaction.Executor) error) error {
	return fn(&coMockExecutor{})
}

// coFailingTransactor always returns an error from WithinTransaction.
type coFailingTransactor struct{ err error }

func (t *coFailingTransactor) WithinTransaction(_ context.Context, _ func(transaction.Executor) error) error {
	return t.err
}

// --- pricing service ---

type coMockPricingService struct {
	result *orderRepo.PricingResult
	err    error
}

func (m *coMockPricingService) Calculate(_ context.Context, _ transaction.Executor, _ orderRepo.PricingInput) (*orderRepo.PricingResult, error) {
	return m.result, m.err
}

// --- user repo ---

type coMockUserRepo struct {
	user *userDomain.User
	err  error
}

func (m *coMockUserRepo) GetByID(_ context.Context, _ transaction.Executor, _ uuid.UUID) (*userDomain.User, error) {
	return m.user, m.err
}
func (m *coMockUserRepo) GetByUsername(_ context.Context, _ transaction.Executor, _ string) (*userDomain.User, error) {
	return nil, nil
}
func (m *coMockUserRepo) CreateUser(_ context.Context, _ transaction.Executor, _ userRepo.CreateUserProps) error {
	return nil
}
func (m *coMockUserRepo) SaveProfile(_ context.Context, _ transaction.Executor, _ userRepo.SaveProfileProps) error {
	return nil
}
func (m *coMockUserRepo) Delete(_ context.Context, _ transaction.Executor, _ uuid.UUID) error {
	return nil
}

// --- account repo ---

type coMockAccountRepo struct {
	account *authenDomain.Account
	err     error
}

func (m *coMockAccountRepo) GetByEmail(_ context.Context, _ transaction.Executor, _ string) (*authenDomain.Account, error) {
	return nil, nil
}
func (m *coMockAccountRepo) GetByID(_ context.Context, _ transaction.Executor, _ uuid.UUID) (*authenDomain.Account, error) {
	return nil, nil
}
func (m *coMockAccountRepo) GetByUserID(_ context.Context, _ transaction.Executor, _ uuid.UUID) (*authenDomain.Account, error) {
	return m.account, m.err
}
func (m *coMockAccountRepo) ActivateByUserID(_ context.Context, _ transaction.Executor, _ uuid.UUID) error {
	return nil
}
func (m *coMockAccountRepo) Create(_ context.Context, _ transaction.Executor, _ authenDomain.Account) error {
	return nil
}
func (m *coMockAccountRepo) UpdatePasswordByUserID(_ context.Context, _ transaction.Executor, _ uuid.UUID, _ string) error {
	return nil
}
func (m *coMockAccountRepo) DeleteByUserID(_ context.Context, _ transaction.Executor, _ uuid.UUID) error {
	return nil
}

// --- payment method repo ---

type coMockPaymentMethodRepo struct {
	method *paymentDomain.PaymentMethod
	err    error
}

func (m *coMockPaymentMethodRepo) Save(_ context.Context, _ transaction.Executor, _ paymentDomain.PaymentMethod) error {
	return nil
}
func (m *coMockPaymentMethodRepo) FindByName(_ context.Context, _ transaction.Executor, _ string) (*paymentDomain.PaymentMethod, error) {
	return nil, nil
}
func (m *coMockPaymentMethodRepo) GetByID(_ context.Context, _ transaction.Executor, _ uuid.UUID) (*paymentDomain.PaymentMethod, error) {
	return m.method, m.err
}
func (m *coMockPaymentMethodRepo) ListAll(_ context.Context, _ transaction.Executor) ([]paymentDomain.PaymentMethod, error) {
	return nil, nil
}

// --- payment repo ---

type coMockPaymentRepo struct {
	payments map[uuid.UUID]*paymentDomain.Payment
}

func (m *coMockPaymentRepo) GetByID(_ context.Context, _ transaction.Executor, id uuid.UUID) (*paymentDomain.Payment, error) {
	return m.payments[id], nil
}
func (m *coMockPaymentRepo) GetByOrderID(_ context.Context, _ transaction.Executor, orderID uuid.UUID) (*paymentDomain.Payment, error) {
	for _, p := range m.payments {
		if p.OrderID == orderID {
			return p, nil
		}
	}
	return nil, nil
}
func (m *coMockPaymentRepo) ListByOrderIDs(_ context.Context, _ transaction.Executor, _ []uuid.UUID) ([]paymentDomain.Payment, error) {
	return nil, nil
}
func (m *coMockPaymentRepo) UpdateStatus(_ context.Context, _ transaction.Executor, _ uuid.UUID, _ paymentDomain.PaymentStatus) error {
	return nil
}
func (m *coMockPaymentRepo) Save(_ context.Context, _ transaction.Executor, payment paymentDomain.Payment) error {
	m.payments[payment.ID] = &payment
	return nil
}

// --- payment account repo ---

type coMockPaymentAccountRepo struct {
	incremented        []uuid.UUID
	leastLoadedAccount *paymentDomain.PaymentAccount
}

func (m *coMockPaymentAccountRepo) Save(_ context.Context, _ transaction.Executor, _ paymentDomain.PaymentAccount) error {
	return nil
}
func (m *coMockPaymentAccountRepo) GetByID(_ context.Context, _ transaction.Executor, _ uuid.UUID) (*paymentDomain.PaymentAccount, error) {
	return nil, nil
}
func (m *coMockPaymentAccountRepo) RetrieveLeastLoaded(_ context.Context, _ transaction.Executor, _ uuid.UUID) (*paymentDomain.PaymentAccount, error) {
	return m.leastLoadedAccount, nil
}
func (m *coMockPaymentAccountRepo) IncrementLoad(_ context.Context, _ transaction.Executor, accountID uuid.UUID) error {
	m.incremented = append(m.incremented, accountID)
	return nil
}
func (m *coMockPaymentAccountRepo) DecrementLoad(_ context.Context, _ transaction.Executor, _ uuid.UUID) error {
	return nil
}
func (m *coMockPaymentAccountRepo) ListByMethodID(_ context.Context, _ transaction.Executor, _ uuid.UUID) ([]paymentDomain.PaymentAccount, error) {
	return nil, nil
}
func (m *coMockPaymentAccountRepo) ListAll(_ context.Context, _ transaction.Executor) ([]paymentDomain.PaymentAccount, error) {
	return nil, nil
}

// --- payment event repo ---

type coMockPaymentEventRepo struct {
	events []paymentDomain.PaymentEvent
}

func (m *coMockPaymentEventRepo) GetByID(_ context.Context, _ transaction.Executor, _ uuid.UUID) (*paymentDomain.PaymentEvent, error) {
	return nil, nil
}
func (m *coMockPaymentEventRepo) ListByPaymentID(_ context.Context, _ transaction.Executor, _ uuid.UUID) ([]paymentDomain.PaymentEvent, error) {
	return nil, nil
}
func (m *coMockPaymentEventRepo) Create(_ context.Context, _ transaction.Executor, event paymentDomain.PaymentEvent) error {
	m.events = append(m.events, event)
	return nil
}

// --- payment instruction repo ---

type coMockPaymentInstructionRepo struct {
	instruction *paymentDomain.PaymentInstruction
}

func (m *coMockPaymentInstructionRepo) GetByPaymentMethodID(_ context.Context, _ transaction.Executor, _ uuid.UUID) (*paymentDomain.PaymentInstruction, error) {
	return m.instruction, nil
}
func (m *coMockPaymentInstructionRepo) Save(_ context.Context, _ transaction.Executor, _ paymentDomain.PaymentInstruction) error {
	return nil
}

// --- order repo ---

type coMockOrderRepo struct {
	orders map[uuid.UUID]*orderDomain.Order
}

func (m *coMockOrderRepo) GetByID(_ context.Context, _ transaction.Executor, id uuid.UUID) (*orderDomain.Order, error) {
	return m.orders[id], nil
}
func (m *coMockOrderRepo) GetByNumber(_ context.Context, _ transaction.Executor, _ string) (*orderDomain.Order, error) {
	return nil, nil
}
func (m *coMockOrderRepo) UpdateStatus(_ context.Context, _ transaction.Executor, id uuid.UUID, status orderDomain.OrderStatus) error {
	if o, ok := m.orders[id]; ok {
		o.Status = status
	}
	return nil
}
func (m *coMockOrderRepo) Save(_ context.Context, _ transaction.Executor, order orderDomain.Order) error {
	m.orders[order.ID] = &order
	return nil
}
func (m *coMockOrderRepo) FindOrders(_ context.Context, _ transaction.Executor, _ orderRepo.FindOrderParams) ([]orderDomain.Order, int, error) {
	return nil, 0, nil
}

// --- order item repo ---

type coMockOrderItemRepo struct {
	items map[uuid.UUID][]orderDomain.OrderItem
}

func (m *coMockOrderItemRepo) ListByOrderID(_ context.Context, _ transaction.Executor, orderID uuid.UUID) ([]orderDomain.OrderItem, error) {
	return m.items[orderID], nil
}
func (m *coMockOrderItemRepo) ListByOrderIDs(_ context.Context, _ transaction.Executor, _ []uuid.UUID) ([]orderDomain.OrderItem, error) {
	return nil, nil
}
func (m *coMockOrderItemRepo) SaveBulk(_ context.Context, _ transaction.Executor, _ []orderDomain.OrderItem) error {
	return nil
}

// --- invoice repos ---

type coMockInvoiceRepo struct{}

func (m *coMockInvoiceRepo) GetByID(_ context.Context, _ transaction.Executor, _ uuid.UUID) (*orderDomain.Invoice, error) {
	return nil, nil
}
func (m *coMockInvoiceRepo) GetByOrderID(_ context.Context, _ transaction.Executor, _ uuid.UUID) (*orderDomain.Invoice, error) {
	return nil, nil
}
func (m *coMockInvoiceRepo) Save(_ context.Context, _ transaction.Executor, _ orderDomain.Invoice) error {
	return nil
}

type coMockInvoiceItemRepo struct{}

func (m *coMockInvoiceItemRepo) ListByInvoiceID(_ context.Context, _ transaction.Executor, _ uuid.UUID) ([]orderDomain.InvoiceItem, error) {
	return nil, nil
}
func (m *coMockInvoiceItemRepo) SaveBulk(_ context.Context, _ transaction.Executor, _ []orderDomain.InvoiceItem) error {
	return nil
}

// --- inventory repo ---

type coMockInventoryRepo struct {
	reserveCalls []string
	reserveErr   error
}

func (m *coMockInventoryRepo) GetByProductIDAndShopID(_ context.Context, _ transaction.Executor, _ uuid.UUID, _ uuid.UUID) (*inventoryDomain.Inventory, error) {
	return nil, nil
}
func (m *coMockInventoryRepo) ListByProductID(_ context.Context, _ transaction.Executor, _ uuid.UUID) ([]inventoryDomain.Inventory, error) {
	return nil, nil
}
func (m *coMockInventoryRepo) ListByProductIDs(_ context.Context, _ transaction.Executor, _ []uuid.UUID) (map[uuid.UUID][]inventoryDomain.Inventory, error) {
	return nil, nil
}
func (m *coMockInventoryRepo) ListByShopID(_ context.Context, _ transaction.Executor, _ uuid.UUID) ([]inventoryDomain.Inventory, error) {
	return nil, nil
}
func (m *coMockInventoryRepo) Create(_ context.Context, _ transaction.Executor, _ *inventoryDomain.Inventory) error {
	return nil
}
func (m *coMockInventoryRepo) Reserve(_ context.Context, _ transaction.Executor, productID uuid.UUID, shopID uuid.UUID, qty int) error {
	m.reserveCalls = append(m.reserveCalls, fmt.Sprintf("%s-%s-%d", productID, shopID, qty))
	return m.reserveErr
}
func (m *coMockInventoryRepo) Release(_ context.Context, _ transaction.Executor, _ uuid.UUID, _ uuid.UUID, _ int) error {
	return nil
}
func (m *coMockInventoryRepo) Commit(_ context.Context, _ transaction.Executor, _ uuid.UUID, _ uuid.UUID, _ int) error {
	return nil
}

// --- cart repo ---

type coMockCartRepo struct {
	cart  *cartDomain.Cart
	saved bool
}

func (m *coMockCartRepo) GetWithItemsByCustomerID(_ context.Context, _ transaction.Executor, _ uuid.UUID) (*cartDomain.Cart, error) {
	return m.cart, nil
}
func (m *coMockCartRepo) NewCart(_ context.Context, _ transaction.Executor, _ uuid.UUID) (*cartDomain.Cart, error) {
	return nil, nil
}
func (m *coMockCartRepo) Save(_ context.Context, _ transaction.Executor, _ *cartDomain.Cart) error {
	m.saved = true
	return nil
}
func (m *coMockCartRepo) DeleteByCustomerID(_ context.Context, _ transaction.Executor, _ uuid.UUID) error {
	return nil
}

// --- gateway mock ---

type coMockGateway struct {
	chargeResp    *paymentgateway.ChargeResponse
	chargeErr     error
	cancelCalled  bool
	cancelOrderID string
}

func (m *coMockGateway) Charge(_ context.Context, _ paymentgateway.ChargeRequest) (*paymentgateway.ChargeResponse, error) {
	return m.chargeResp, m.chargeErr
}
func (m *coMockGateway) ParseNotification(_ context.Context, _ paymentgateway.NotificationPayload) (*paymentgateway.NotificationResult, error) {
	return nil, nil
}
func (m *coMockGateway) CancelTransaction(_ context.Context, gatewayOrderID string) error {
	m.cancelCalled = true
	m.cancelOrderID = gatewayOrderID
	return nil
}

// capturingGateway captures the ChargeRequest for inspection.
type capturingGateway struct {
	resp     *paymentgateway.ChargeResponse
	captured *paymentgateway.ChargeRequest
}

func (c *capturingGateway) Charge(_ context.Context, req paymentgateway.ChargeRequest) (*paymentgateway.ChargeResponse, error) {
	*c.captured = req
	return c.resp, nil
}
func (c *capturingGateway) ParseNotification(_ context.Context, _ paymentgateway.NotificationPayload) (*paymentgateway.NotificationResult, error) {
	return nil, nil
}
func (c *capturingGateway) CancelTransaction(_ context.Context, _ string) error { return nil }

// ===========================================================================
// Test helpers
// ===========================================================================

func coDefaultUser() *userDomain.User {
	ph := "+628111111111"
	return &userDomain.User{ID: uuid.New(), Name: "Test User", Phone: &ph}
}

func coDefaultAccount(userID uuid.UUID) *authenDomain.Account {
	return &authenDomain.Account{ID: uuid.New(), UserID: userID, Email: "test@chia.florist"}
}

func coDefaultMethod(methodID uuid.UUID, name string) *paymentDomain.PaymentMethod {
	return &paymentDomain.PaymentMethod{ID: methodID, Name: name, Type: paymentDomain.TypeQRCode, IsActive: true}
}

func coDefaultPricing(productID, shopID uuid.UUID) *orderRepo.PricingResult {
	return &orderRepo.PricingResult{
		Subtotal:         100000,
		TotalShippingFee: 15000,
		GrandTotal:       115000,
		Shops: []orderRepo.PricingShopResult{
			{
				ShopID:   shopID,
				ShopName: "Florist Kage",
				Items: []orderRepo.PricingItemResult{
					{ProductID: productID, ProductName: "Bouquet A", Quantity: 2, UnitPrice: 50000, Subtotal: 100000},
				},
				SelectedCourier: orderRepo.SelectedCourierResult{Code: "jne", Service: "REG", Fee: 15000},
			},
		},
	}
}

func coDefaultChargeResp(providerOrderID string) *paymentgateway.ChargeResponse {
	return &paymentgateway.ChargeResponse{
		GatewayTransactionID: "midtrans-tx-001",
		GatewayOrderID:       providerOrderID,
		PaymentType:          "qris",
		GrossAmount:          115000,
		Status:               "pending",
		ExpiresAt:            time.Now().Add(24 * time.Hour),
	}
}

// buildUC assembles a CreateOrderUsecase from the provided dependencies.
func buildUC(
	transactor transaction.Transactor,
	gateway paymentgateway.Provider,
	pricing *coMockPricingService,
	paymentMethodRepo *coMockPaymentMethodRepo,
	paymentAccRepo *coMockPaymentAccountRepo,
	userRepo *coMockUserRepo,
	accountRepo *coMockAccountRepo,
	inventoryRepo *coMockInventoryRepo,
	paymentStore *coMockPaymentRepo,
	orderStore *coMockOrderRepo,
	cartRepo *coMockCartRepo,
) *CreateOrderUsecase {
	if transactor == nil {
		transactor = &coMockTransactor{}
	}
	if cartRepo == nil {
		cartRepo = &coMockCartRepo{}
	}
	return NewCreateOrderUsecase(
		&coMockExecutor{},
		transactor,
		accountRepo,
		orderStore,
		&coMockOrderItemRepo{items: map[uuid.UUID][]orderDomain.OrderItem{}},
		&coMockInvoiceRepo{},
		&coMockInvoiceItemRepo{},
		paymentStore,
		paymentMethodRepo,
		paymentAccRepo,
		&coMockPaymentEventRepo{},
		&coMockPaymentInstructionRepo{},
		inventoryRepo,
		cartRepo,
		userRepo,
		gateway,
		pricing,
	)
}

func coInput(customerID, paymentMethodID uuid.UUID, isManual bool, productID, shopID uuid.UUID) CreateOrderInput {
	return CreateOrderInput{
		UserID:          uuid.New(),
		CustomerID:      customerID,
		AddressID:       uuid.New(),
		PaymentMethodID: paymentMethodID,
		IsManual:        isManual,
		Shops: []OrderShopInput{
			{
				ShopID:   shopID,
				ShopName: "Florist Kage",
				Items:    []OrderItemInput{{ProductID: productID, ProductName: "Bouquet A", Quantity: 2}},
			},
		},
	}
}

// ===========================================================================
// Happy-path: gateway (non-manual) order
// ===========================================================================

func TestCreateOrder_Gateway_Success(t *testing.T) {
	ctx := context.Background()
	customerID := uuid.New()
	productID := uuid.New()
	shopID := uuid.New()
	methodID := uuid.New()

	user := coDefaultUser()
	acc := coDefaultAccount(user.ID)
	method := coDefaultMethod(methodID, "qris")
	pricing := coDefaultPricing(productID, shopID)
	providerOrderID := uuid.New().String()

	gateway := &coMockGateway{
		chargeResp: &paymentgateway.ChargeResponse{
			GatewayTransactionID: "midtrans-tx-001",
			GatewayOrderID:       providerOrderID,
			PaymentType:          "qris",
			GrossAmount:          115000,
			Status:               "pending",
			Instructions: []paymentgateway.PaymentInstruction{
				{Type: "qris", Label: "QRIS", Value: "qr-string-value"},
			},
			ExpiresAt: time.Now().Add(24 * time.Hour),
		},
	}

	paymentStore := &coMockPaymentRepo{payments: map[uuid.UUID]*paymentDomain.Payment{}}
	orderStore := &coMockOrderRepo{orders: map[uuid.UUID]*orderDomain.Order{}}

	uc := buildUC(nil, gateway,
		&coMockPricingService{result: pricing},
		&coMockPaymentMethodRepo{method: method},
		&coMockPaymentAccountRepo{},
		&coMockUserRepo{user: user},
		&coMockAccountRepo{account: acc},
		&coMockInventoryRepo{},
		paymentStore,
		orderStore,
		nil,
	)

	result, err := uc.Execute(ctx, coInput(customerID, methodID, false, productID, shopID))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.OrderID == uuid.Nil {
		t.Error("OrderID should not be nil")
	}
	if result.Total != 115000 {
		t.Errorf("Total = %v, want 115000", result.Total)
	}
	if result.PaymentAccount == nil {
		t.Fatal("PaymentAccount should not be nil")
	}
	if result.PaymentAccount.QRString == nil || *result.PaymentAccount.QRString != "qr-string-value" {
		t.Errorf("QRString mismatch: %v", result.PaymentAccount.QRString)
	}
}

// ===========================================================================
// Happy-path: manual order
// ===========================================================================

func TestCreateOrder_Manual_Success(t *testing.T) {
	ctx := context.Background()
	customerID := uuid.New()
	productID := uuid.New()
	shopID := uuid.New()
	methodID := uuid.New()
	accountID := uuid.New()
	accNumber := "1234567890"

	user := coDefaultUser()
	acc := coDefaultAccount(user.ID)
	method := &paymentDomain.PaymentMethod{ID: methodID, Name: "mandiri", Type: paymentDomain.TypeBankTransfer, IsActive: true}
	pricing := coDefaultPricing(productID, shopID)

	paRepo := &coMockPaymentAccountRepo{
		leastLoadedAccount: &paymentDomain.PaymentAccount{
			ID:            accountID,
			AccountName:   "Mandiri Kage",
			AccountNumber: &accNumber,
		},
	}

	paymentStore := &coMockPaymentRepo{payments: map[uuid.UUID]*paymentDomain.Payment{}}

	uc := buildUC(nil, &coMockGateway{},
		&coMockPricingService{result: pricing},
		&coMockPaymentMethodRepo{method: method},
		paRepo,
		&coMockUserRepo{user: user},
		&coMockAccountRepo{account: acc},
		&coMockInventoryRepo{},
		paymentStore,
		&coMockOrderRepo{orders: map[uuid.UUID]*orderDomain.Order{}},
		nil,
	)

	result, err := uc.Execute(ctx, coInput(customerID, methodID, true, productID, shopID))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.PaymentAccount == nil {
		t.Fatal("PaymentAccount should not be nil for manual payment")
	}
	if result.PaymentAccount.AccountName != "Mandiri Kage" {
		t.Errorf("AccountName = %v, want Mandiri Kage", result.PaymentAccount.AccountName)
	}
	if result.PaymentAccount.AccountNumber == nil || *result.PaymentAccount.AccountNumber != accNumber {
		t.Errorf("AccountNumber mismatch: %v", result.PaymentAccount.AccountNumber)
	}
	if len(paRepo.incremented) != 1 || paRepo.incremented[0] != accountID {
		t.Errorf("expected IncrementLoad for %v, got %v", accountID, paRepo.incremented)
	}
}

// ===========================================================================
// Charge response instruction type mapping
// ===========================================================================

func TestCreateOrder_ChargeResponse_BankTransfer_PopulatesAccountNumber(t *testing.T) {
	ctx := context.Background()
	customerID := uuid.New()
	productID := uuid.New()
	shopID := uuid.New()
	methodID := uuid.New()
	user := coDefaultUser()
	acc := coDefaultAccount(user.ID)
	method := &paymentDomain.PaymentMethod{ID: methodID, Name: "mandiri", Type: paymentDomain.TypeBankTransfer, IsActive: true}
	pricing := coDefaultPricing(productID, shopID)

	gateway := &coMockGateway{
		chargeResp: &paymentgateway.ChargeResponse{
			GatewayTransactionID: "tx-001",
			GatewayOrderID:       uuid.New().String(),
			GrossAmount:          115000,
			Instructions: []paymentgateway.PaymentInstruction{
				{Type: "bank_transfer", Label: "BCA Virtual Account", Value: "1234567890"},
			},
		},
	}

	uc := buildUC(nil, gateway,
		&coMockPricingService{result: pricing},
		&coMockPaymentMethodRepo{method: method},
		&coMockPaymentAccountRepo{},
		&coMockUserRepo{user: user},
		&coMockAccountRepo{account: acc},
		&coMockInventoryRepo{},
		&coMockPaymentRepo{payments: map[uuid.UUID]*paymentDomain.Payment{}},
		&coMockOrderRepo{orders: map[uuid.UUID]*orderDomain.Order{}},
		nil,
	)

	result, err := uc.Execute(ctx, coInput(customerID, methodID, false, productID, shopID))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.PaymentAccount == nil {
		t.Fatal("PaymentAccount should not be nil")
	}
	if result.PaymentAccount.AccountNumber == nil || *result.PaymentAccount.AccountNumber != "1234567890" {
		t.Errorf("AccountNumber = %v, want 1234567890", result.PaymentAccount.AccountNumber)
	}
	if result.PaymentAccount.QRString != nil {
		t.Error("QRString should be nil for bank_transfer")
	}
}

func TestCreateOrder_ChargeResponse_EWallet_PopulatesQRString(t *testing.T) {
	ctx := context.Background()
	customerID := uuid.New()
	productID := uuid.New()
	shopID := uuid.New()
	methodID := uuid.New()
	user := coDefaultUser()
	acc := coDefaultAccount(user.ID)
	method := &paymentDomain.PaymentMethod{ID: methodID, Name: "gopay", Type: paymentDomain.TypeEWallet, IsActive: true}
	pricing := coDefaultPricing(productID, shopID)

	gateway := &coMockGateway{
		chargeResp: &paymentgateway.ChargeResponse{
			GatewayTransactionID: "tx-002",
			GatewayOrderID:       uuid.New().String(),
			GrossAmount:          115000,
			Instructions: []paymentgateway.PaymentInstruction{
				{Type: "ewallet", Label: "deeplink-redirect", Value: "https://gopay.co/pay/link"},
			},
		},
	}

	uc := buildUC(nil, gateway,
		&coMockPricingService{result: pricing},
		&coMockPaymentMethodRepo{method: method},
		&coMockPaymentAccountRepo{},
		&coMockUserRepo{user: user},
		&coMockAccountRepo{account: acc},
		&coMockInventoryRepo{},
		&coMockPaymentRepo{payments: map[uuid.UUID]*paymentDomain.Payment{}},
		&coMockOrderRepo{orders: map[uuid.UUID]*orderDomain.Order{}},
		nil,
	)

	result, err := uc.Execute(ctx, coInput(customerID, methodID, false, productID, shopID))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.PaymentAccount == nil {
		t.Fatal("PaymentAccount should not be nil")
	}
	if result.PaymentAccount.QRString == nil || *result.PaymentAccount.QRString != "https://gopay.co/pay/link" {
		t.Errorf("QRString = %v, want deeplink URL", result.PaymentAccount.QRString)
	}
}

// ===========================================================================
// Aggregate invariant: adjustment item when item sum != order total
// ===========================================================================

func TestCreateOrder_InvariantAdjustmentItemAddedWhenSumDiffers(t *testing.T) {
	ctx := context.Background()
	customerID := uuid.New()
	productID := uuid.New()
	shopID := uuid.New()
	methodID := uuid.New()
	user := coDefaultUser()
	acc := coDefaultAccount(user.ID)
	method := coDefaultMethod(methodID, "qris")

	// GrandTotal=120000, but items+shipping=115000 → 5000 adjustment
	pricing := &orderRepo.PricingResult{
		Subtotal:         100000,
		TotalShippingFee: 15000,
		GrandTotal:       120000,
		Shops: []orderRepo.PricingShopResult{
			{
				ShopID:   shopID,
				ShopName: "Kage",
				Items: []orderRepo.PricingItemResult{
					{ProductID: productID, ProductName: "Bouquet", Quantity: 2, UnitPrice: 50000, Subtotal: 100000},
				},
				SelectedCourier: orderRepo.SelectedCourierResult{Code: "jne", Service: "REG", Fee: 15000},
			},
		},
	}

	var capturedReq paymentgateway.ChargeRequest
	gateway := &capturingGateway{
		resp:     coDefaultChargeResp(uuid.New().String()),
		captured: &capturedReq,
	}

	uc := buildUC(nil, gateway,
		&coMockPricingService{result: pricing},
		&coMockPaymentMethodRepo{method: method},
		&coMockPaymentAccountRepo{},
		&coMockUserRepo{user: user},
		&coMockAccountRepo{account: acc},
		&coMockInventoryRepo{},
		&coMockPaymentRepo{payments: map[uuid.UUID]*paymentDomain.Payment{}},
		&coMockOrderRepo{orders: map[uuid.UUID]*orderDomain.Order{}},
		nil,
	)

	_, err := uc.Execute(ctx, coInput(customerID, methodID, false, productID, shopID))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var hasAdjustment bool
	var itemSum int64
	for _, item := range capturedReq.Items {
		itemSum += item.Price * int64(item.Quantity)
		if item.ID == "adjustment" {
			hasAdjustment = true
		}
	}
	if !hasAdjustment {
		t.Error("expected an adjustment charge item when item sum differs from order total")
	}
	if itemSum != 120000 {
		t.Errorf("charge items sum = %d, want 120000", itemSum)
	}
}

// ===========================================================================
// Aggregate invariant: inventory reserved for each order item
// ===========================================================================

func TestCreateOrder_InvariantInventoryReservedPerItem(t *testing.T) {
	ctx := context.Background()
	customerID := uuid.New()
	productID := uuid.New()
	shopID := uuid.New()
	methodID := uuid.New()
	user := coDefaultUser()
	acc := coDefaultAccount(user.ID)
	method := coDefaultMethod(methodID, "qris")
	pricing := coDefaultPricing(productID, shopID)

	invRepo := &coMockInventoryRepo{}
	gateway := &coMockGateway{chargeResp: coDefaultChargeResp(uuid.New().String())}

	uc := buildUC(nil, gateway,
		&coMockPricingService{result: pricing},
		&coMockPaymentMethodRepo{method: method},
		&coMockPaymentAccountRepo{},
		&coMockUserRepo{user: user},
		&coMockAccountRepo{account: acc},
		invRepo,
		&coMockPaymentRepo{payments: map[uuid.UUID]*paymentDomain.Payment{}},
		&coMockOrderRepo{orders: map[uuid.UUID]*orderDomain.Order{}},
		nil,
	)

	_, err := uc.Execute(ctx, coInput(customerID, methodID, false, productID, shopID))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := fmt.Sprintf("%s-%s-2", productID, shopID)
	if len(invRepo.reserveCalls) != 1 || invRepo.reserveCalls[0] != expected {
		t.Errorf("expected inventory reserve %v, got %v", expected, invRepo.reserveCalls)
	}
}

// ===========================================================================
// Aggregate invariant: payment amount equals order total
// ===========================================================================

func TestCreateOrder_InvariantPaymentAmountEqualsOrderTotal(t *testing.T) {
	ctx := context.Background()
	customerID := uuid.New()
	productID := uuid.New()
	shopID := uuid.New()
	methodID := uuid.New()
	user := coDefaultUser()
	acc := coDefaultAccount(user.ID)
	method := coDefaultMethod(methodID, "qris")
	pricing := coDefaultPricing(productID, shopID) // GrandTotal = 115000

	paymentStore := &coMockPaymentRepo{payments: map[uuid.UUID]*paymentDomain.Payment{}}
	gateway := &coMockGateway{chargeResp: coDefaultChargeResp(uuid.New().String())}

	uc := buildUC(nil, gateway,
		&coMockPricingService{result: pricing},
		&coMockPaymentMethodRepo{method: method},
		&coMockPaymentAccountRepo{},
		&coMockUserRepo{user: user},
		&coMockAccountRepo{account: acc},
		&coMockInventoryRepo{},
		paymentStore,
		&coMockOrderRepo{orders: map[uuid.UUID]*orderDomain.Order{}},
		nil,
	)

	result, err := uc.Execute(ctx, coInput(customerID, methodID, false, productID, shopID))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Total != 115000 {
		t.Errorf("result.Total = %v, want 115000", result.Total)
	}
	for _, p := range paymentStore.payments {
		if p.Amount != 115000 {
			t.Errorf("payment.Amount = %v, want 115000", p.Amount)
		}
	}
}

// ===========================================================================
// Aggregate invariant: cart items removed after order
// ===========================================================================

func TestCreateOrder_InvariantCartItemsRemovedAfterOrder(t *testing.T) {
	ctx := context.Background()
	customerID := uuid.New()
	productID := uuid.New()
	shopID := uuid.New()
	methodID := uuid.New()
	user := coDefaultUser()
	acc := coDefaultAccount(user.ID)
	method := coDefaultMethod(methodID, "qris")
	pricing := coDefaultPricing(productID, shopID)

	cartRepo := &coMockCartRepo{
		cart: &cartDomain.Cart{
			ID:         uuid.New(),
			CustomerID: customerID,
			Items: []cartDomain.CartItem{
				{ID: uuid.New(), ProductID: productID, ShopID: shopID, Quantity: 2},
			},
		},
	}
	gateway := &coMockGateway{chargeResp: coDefaultChargeResp(uuid.New().String())}

	uc := buildUC(nil, gateway,
		&coMockPricingService{result: pricing},
		&coMockPaymentMethodRepo{method: method},
		&coMockPaymentAccountRepo{},
		&coMockUserRepo{user: user},
		&coMockAccountRepo{account: acc},
		&coMockInventoryRepo{},
		&coMockPaymentRepo{payments: map[uuid.UUID]*paymentDomain.Payment{}},
		&coMockOrderRepo{orders: map[uuid.UUID]*orderDomain.Order{}},
		cartRepo,
	)

	_, err := uc.Execute(ctx, coInput(customerID, methodID, false, productID, shopID))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !cartRepo.saved {
		t.Error("expected cart to be saved (items removed) after order creation")
	}
}

// ===========================================================================
// Failure recovery: gateway charge fails → no DB writes
// ===========================================================================

func TestCreateOrder_GatewayChargeFailure_NoDBWrites(t *testing.T) {
	ctx := context.Background()
	customerID := uuid.New()
	productID := uuid.New()
	shopID := uuid.New()
	methodID := uuid.New()
	user := coDefaultUser()
	acc := coDefaultAccount(user.ID)
	method := coDefaultMethod(methodID, "qris")
	pricing := coDefaultPricing(productID, shopID)

	gateway := &coMockGateway{chargeErr: errors.New("midtrans: gateway unavailable")}
	paymentStore := &coMockPaymentRepo{payments: map[uuid.UUID]*paymentDomain.Payment{}}

	uc := buildUC(nil, gateway,
		&coMockPricingService{result: pricing},
		&coMockPaymentMethodRepo{method: method},
		&coMockPaymentAccountRepo{},
		&coMockUserRepo{user: user},
		&coMockAccountRepo{account: acc},
		&coMockInventoryRepo{},
		paymentStore,
		&coMockOrderRepo{orders: map[uuid.UUID]*orderDomain.Order{}},
		nil,
	)

	_, err := uc.Execute(ctx, coInput(customerID, methodID, false, productID, shopID))
	if err == nil {
		t.Fatal("expected error when gateway charge fails, got nil")
	}
	if len(paymentStore.payments) != 0 {
		t.Errorf("expected no payments saved after charge failure, got %d", len(paymentStore.payments))
	}
}

// ===========================================================================
// Coordination protocol: DB fails after successful charge → cancel called
// ===========================================================================

func TestCreateOrder_DBWriteFailAfterCharge_BestEffortCancel(t *testing.T) {
	ctx := context.Background()
	customerID := uuid.New()
	productID := uuid.New()
	shopID := uuid.New()
	methodID := uuid.New()
	user := coDefaultUser()
	acc := coDefaultAccount(user.ID)
	method := coDefaultMethod(methodID, "qris")
	pricing := coDefaultPricing(productID, shopID)

	providerOrderID := "gateway-order-ref"
	gateway := &coMockGateway{
		chargeResp: &paymentgateway.ChargeResponse{
			GatewayTransactionID: "tx-cancel",
			GatewayOrderID:       providerOrderID,
			GrossAmount:          115000,
		},
	}

	uc := buildUC(
		&coFailingTransactor{err: errors.New("database: connection lost")},
		gateway,
		&coMockPricingService{result: pricing},
		&coMockPaymentMethodRepo{method: method},
		&coMockPaymentAccountRepo{},
		&coMockUserRepo{user: user},
		&coMockAccountRepo{account: acc},
		&coMockInventoryRepo{},
		&coMockPaymentRepo{payments: map[uuid.UUID]*paymentDomain.Payment{}},
		&coMockOrderRepo{orders: map[uuid.UUID]*orderDomain.Order{}},
		nil,
	)

	_, err := uc.Execute(ctx, coInput(customerID, methodID, false, productID, shopID))
	if err == nil {
		t.Fatal("expected error when DB write fails")
	}
	if !gateway.cancelCalled {
		t.Error("expected CancelTransaction to be called as best-effort rollback after DB failure")
	}
	if gateway.cancelOrderID != providerOrderID {
		t.Errorf("CancelTransaction called with %q, want %q", gateway.cancelOrderID, providerOrderID)
	}
}

// ===========================================================================
// Coordination protocol: charge never succeeded → cancel NOT called
// ===========================================================================

func TestCreateOrder_DBWriteFailNoCharge_NoBestEffortCancel(t *testing.T) {
	ctx := context.Background()
	customerID := uuid.New()
	productID := uuid.New()
	shopID := uuid.New()
	methodID := uuid.New()
	user := coDefaultUser()
	acc := coDefaultAccount(user.ID)
	method := coDefaultMethod(methodID, "qris")
	pricing := coDefaultPricing(productID, shopID)

	gateway := &coMockGateway{chargeErr: errors.New("gateway error")}

	uc := buildUC(
		&coFailingTransactor{err: errors.New("db error")},
		gateway,
		&coMockPricingService{result: pricing},
		&coMockPaymentMethodRepo{method: method},
		&coMockPaymentAccountRepo{},
		&coMockUserRepo{user: user},
		&coMockAccountRepo{account: acc},
		&coMockInventoryRepo{},
		&coMockPaymentRepo{payments: map[uuid.UUID]*paymentDomain.Payment{}},
		&coMockOrderRepo{orders: map[uuid.UUID]*orderDomain.Order{}},
		nil,
	)

	_, err := uc.Execute(ctx, coInput(customerID, methodID, false, productID, shopID))
	if err == nil {
		t.Fatal("expected error")
	}
	if gateway.cancelCalled {
		t.Error("CancelTransaction should NOT be called when no gateway charge was made")
	}
}

// ===========================================================================
// Failure: pricing service error → early exit
// ===========================================================================

func TestCreateOrder_PricingServiceFailure(t *testing.T) {
	ctx := context.Background()
	methodID := uuid.New()
	method := coDefaultMethod(methodID, "qris")
	user := coDefaultUser()
	acc := coDefaultAccount(user.ID)

	gateway := &coMockGateway{}

	uc := buildUC(nil, gateway,
		&coMockPricingService{err: errors.New("pricing: product not found")},
		&coMockPaymentMethodRepo{method: method},
		&coMockPaymentAccountRepo{},
		&coMockUserRepo{user: user},
		&coMockAccountRepo{account: acc},
		&coMockInventoryRepo{},
		&coMockPaymentRepo{payments: map[uuid.UUID]*paymentDomain.Payment{}},
		&coMockOrderRepo{orders: map[uuid.UUID]*orderDomain.Order{}},
		nil,
	)

	_, err := uc.Execute(ctx, coInput(uuid.New(), methodID, false, uuid.New(), uuid.New()))
	if err == nil {
		t.Fatal("expected error when pricing fails")
	}
	if gateway.cancelCalled {
		t.Error("CancelTransaction should not be called when pricing fails before charge")
	}
}

// ===========================================================================
// Failure: payment method not found
// ===========================================================================

func TestCreateOrder_PaymentMethodNotFound(t *testing.T) {
	ctx := context.Background()
	user := coDefaultUser()
	acc := coDefaultAccount(user.ID)

	uc := buildUC(nil, &coMockGateway{},
		&coMockPricingService{result: coDefaultPricing(uuid.New(), uuid.New())},
		&coMockPaymentMethodRepo{method: nil}, // not found
		&coMockPaymentAccountRepo{},
		&coMockUserRepo{user: user},
		&coMockAccountRepo{account: acc},
		&coMockInventoryRepo{},
		&coMockPaymentRepo{payments: map[uuid.UUID]*paymentDomain.Payment{}},
		&coMockOrderRepo{orders: map[uuid.UUID]*orderDomain.Order{}},
		nil,
	)

	_, err := uc.Execute(ctx, coInput(uuid.New(), uuid.New(), false, uuid.New(), uuid.New()))
	if err == nil {
		t.Fatal("expected error when payment method not found")
	}
}
