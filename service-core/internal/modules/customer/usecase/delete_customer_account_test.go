package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	applogger "service-core/internal/common/logger"
	addressDomain "service-core/internal/modules/address/domain"
	authenDomain "service-core/internal/modules/authentication/domain"
	cartDomain "service-core/internal/modules/cart/domain"
	customerDomain "service-core/internal/modules/customer/domain"
	customerRepo "service-core/internal/modules/customer/repository"
	userDomain "service-core/internal/modules/user/domain"
	userRepo "service-core/internal/modules/user/repository"
	transaction "service-core/internal/shared/transaction"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// ===========================================================================
// Mocks for DeleteCustomerAccountUsecase Tests
// ===========================================================================

type delMockExecutor struct{}

func (m *delMockExecutor) Exec(_ context.Context, _ string, _ ...any) (pgconn.CommandTag, error) {
	return pgconn.NewCommandTag(""), nil
}
func (m *delMockExecutor) Query(_ context.Context, _ string, _ ...any) (pgx.Rows, error) {
	return nil, nil
}
func (m *delMockExecutor) QueryRow(_ context.Context, _ string, _ ...any) pgx.Row {
	return nil
}

type delMockTransactor struct{}

func (m *delMockTransactor) WithinTransaction(ctx context.Context, fn func(transaction.Executor) error) error {
	return fn(&delMockExecutor{})
}

type delMockAuditLogger struct{}

func (m *delMockAuditLogger) Log(_ context.Context, _ applogger.AuditEvent) {}

// CustomerAddressRepository Mock
type mockAddressRepo struct {
	deleteCalls int
}

func (m *mockAddressRepo) GetByID(_ context.Context, _ transaction.Executor, _ uuid.UUID) (*addressDomain.CustomerAddress, error) {
	return nil, nil
}
func (m *mockAddressRepo) GetDefaultByCustomerID(_ context.Context, _ transaction.Executor, _ uuid.UUID) (*addressDomain.CustomerAddress, error) {
	return nil, nil
}
func (m *mockAddressRepo) ListByCustomerID(_ context.Context, _ transaction.Executor, _ uuid.UUID) ([]addressDomain.CustomerAddress, error) {
	return nil, nil
}
func (m *mockAddressRepo) CountByCustomerID(_ context.Context, _ transaction.Executor, _ uuid.UUID) (*int, error) {
	return nil, nil
}
func (m *mockAddressRepo) UnsetDefaultByCustomerID(_ context.Context, _ transaction.Executor, _ uuid.UUID) error {
	return nil
}
func (m *mockAddressRepo) Save(_ context.Context, _ transaction.Executor, _ addressDomain.CustomerAddress) error {
	return nil
}
func (m *mockAddressRepo) Delete(_ context.Context, _ transaction.Executor, _ uuid.UUID) error {
	return nil
}
func (m *mockAddressRepo) DeleteByCustomerID(_ context.Context, _ transaction.Executor, _ uuid.UUID) error {
	m.deleteCalls++
	return nil
}

// CartRepository Mock
type mockCartRepo struct {
	deleteCalls int
}

func (m *mockCartRepo) GetWithItemsByCustomerID(_ context.Context, _ transaction.Executor, _ uuid.UUID) (*cartDomain.Cart, error) {
	return nil, nil
}
func (m *mockCartRepo) NewCart(_ context.Context, _ transaction.Executor, _ uuid.UUID) (*cartDomain.Cart, error) {
	return nil, nil
}
func (m *mockCartRepo) Save(_ context.Context, _ transaction.Executor, _ *cartDomain.Cart) error {
	return nil
}
func (m *mockCartRepo) DeleteByCustomerID(_ context.Context, _ transaction.Executor, _ uuid.UUID) error {
	m.deleteCalls++
	return nil
}

// CustomerRepository Mock
type mockCustomerRepo struct {
	deleteCalls int
	deleteErr   error
}

func (m *mockCustomerRepo) Create(_ context.Context, _ transaction.Executor, _ customerDomain.Customer) error {
	return nil
}
func (m *mockCustomerRepo) GetByID(_ context.Context, _ transaction.Executor, _ uuid.UUID) (*customerDomain.Customer, error) {
	return nil, nil
}
func (m *mockCustomerRepo) GetByUserID(_ context.Context, _ transaction.Executor, _ uuid.UUID) (*customerDomain.Customer, error) {
	return nil, nil
}
func (m *mockCustomerRepo) GetProfileByUserID(_ context.Context, _ transaction.Executor, _ uuid.UUID) (*customerDomain.CustomerProfile, error) {
	return nil, nil
}
func (m *mockCustomerRepo) FindCustomers(_ context.Context, _ transaction.Executor, _ customerRepo.FindCustomerParams) ([]customerDomain.CustomerProfile, int, error) {
	return nil, 0, nil
}
func (m *mockCustomerRepo) Delete(_ context.Context, _ transaction.Executor, _ uuid.UUID) error {
	m.deleteCalls++
	return m.deleteErr
}

// UserRepository Mock
type mockUserRepo struct {
	deleteCalls int
}

func (m *mockUserRepo) GetByID(_ context.Context, _ transaction.Executor, _ uuid.UUID) (*userDomain.User, error) {
	return nil, nil
}
func (m *mockUserRepo) GetByUsername(_ context.Context, _ transaction.Executor, _ string) (*userDomain.User, error) {
	return nil, nil
}
func (m *mockUserRepo) CreateUser(_ context.Context, _ transaction.Executor, _ userRepo.CreateUserProps) error {
	return nil
}
func (m *mockUserRepo) SaveProfile(_ context.Context, _ transaction.Executor, _ userRepo.SaveProfileProps) error {
	return nil
}
func (m *mockUserRepo) Delete(_ context.Context, _ transaction.Executor, _ uuid.UUID) error {
	m.deleteCalls++
	return nil
}

// AccountRepository Mock
type mockAccountRepo struct {
	deleteCalls int
}

func (m *mockAccountRepo) GetByEmail(_ context.Context, _ transaction.Executor, _ string) (*authenDomain.Account, error) {
	return nil, nil
}
func (m *mockAccountRepo) GetByID(_ context.Context, _ transaction.Executor, _ uuid.UUID) (*authenDomain.Account, error) {
	return nil, nil
}
func (m *mockAccountRepo) GetByUserID(_ context.Context, _ transaction.Executor, _ uuid.UUID) (*authenDomain.Account, error) {
	return nil, nil
}
func (m *mockAccountRepo) ActivateByUserID(_ context.Context, _ transaction.Executor, _ uuid.UUID) error {
	return nil
}
func (m *mockAccountRepo) UpdatePasswordByUserID(_ context.Context, _ transaction.Executor, _ uuid.UUID, _ string) error {
	return nil
}
func (m *mockAccountRepo) Create(_ context.Context, _ transaction.Executor, _ authenDomain.Account) error {
	return nil
}
func (m *mockAccountRepo) DeleteByUserID(_ context.Context, _ transaction.Executor, _ uuid.UUID) error {
	m.deleteCalls++
	return nil
}

// OAuthConnectionRepository Mock
type mockOAuthRepo struct {
	deleteCalls int
}

func (m *mockOAuthRepo) GetByProviderAndSubject(_ context.Context, _ transaction.Executor, _ authenDomain.OAuthProvider, _ string) (*authenDomain.OAuthConnection, error) {
	return nil, nil
}
func (m *mockOAuthRepo) GetByUserID(_ context.Context, _ transaction.Executor, _ uuid.UUID) (*authenDomain.OAuthConnection, error) {
	return nil, nil
}
func (m *mockOAuthRepo) Create(_ context.Context, _ transaction.Executor, _ authenDomain.OAuthConnection) error {
	return nil
}
func (m *mockOAuthRepo) UpdateLastLogin(_ context.Context, _ transaction.Executor, _ uuid.UUID, _ time.Time) error {
	return nil
}
func (m *mockOAuthRepo) DeleteByUserID(_ context.Context, _ transaction.Executor, _ uuid.UUID) error {
	m.deleteCalls++
	return nil
}

// SessionRepository Mock
type mockSessionRepo struct {
	revokeCalls int
}

func (m *mockSessionRepo) GetByID(_ context.Context, _ transaction.Executor, _ uuid.UUID) (*authenDomain.Session, error) {
	return nil, nil
}
func (m *mockSessionRepo) RevokeByID(_ context.Context, _ transaction.Executor, _ uuid.UUID) error {
	return nil
}
func (m *mockSessionRepo) RevokeAllByUserID(_ context.Context, _ transaction.Executor, _ uuid.UUID) error {
	m.revokeCalls++
	return nil
}
func (m *mockSessionRepo) UpdateLastActivityByID(_ context.Context, _ transaction.Executor, _ uuid.UUID) error {
	return nil
}
func (m *mockSessionRepo) Save(_ context.Context, _ transaction.Executor, _ authenDomain.Session) error {
	return nil
}

// ===========================================================================
// Tests
// ===========================================================================

func TestDeleteCustomerAccount_Success(t *testing.T) {
	ctx := context.Background()
	customerID := uuid.New()
	userID := uuid.New()

	customerRepoMock := &mockCustomerRepo{}
	userRepoMock := &mockUserRepo{}
	accountRepoMock := &mockAccountRepo{}
	oauthRepoMock := &mockOAuthRepo{}
	addressRepoMock := &mockAddressRepo{}
	cartRepoMock := &mockCartRepo{}
	sessionRepoMock := &mockSessionRepo{}

	uc := NewDeleteCustomerAccountUsecase(
		&delMockTransactor{},
		accountRepoMock,
		oauthRepoMock,
		sessionRepoMock,
		userRepoMock,
		customerRepoMock,
		addressRepoMock,
		cartRepoMock,
		&delMockAuditLogger{},
	)

	authCtx := authenDomain.AuthContext{
		UserID:          userID,
		CustomerID:      &customerID,
		IsAuthenticated: true,
	}

	err := uc.Execute(ctx, authCtx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if customerRepoMock.deleteCalls != 1 {
		t.Errorf("expected 1 customer delete call, got %d", customerRepoMock.deleteCalls)
	}
	if userRepoMock.deleteCalls != 1 {
		t.Errorf("expected 1 user delete call, got %d", userRepoMock.deleteCalls)
	}
	if accountRepoMock.deleteCalls != 1 {
		t.Errorf("expected 1 account delete call, got %d", accountRepoMock.deleteCalls)
	}
	if oauthRepoMock.deleteCalls != 1 {
		t.Errorf("expected 1 oauth delete call, got %d", oauthRepoMock.deleteCalls)
	}
	if addressRepoMock.deleteCalls != 1 {
		t.Errorf("expected 1 address delete call, got %d", addressRepoMock.deleteCalls)
	}
	if cartRepoMock.deleteCalls != 1 {
		t.Errorf("expected 1 cart delete call, got %d", cartRepoMock.deleteCalls)
	}
	if sessionRepoMock.revokeCalls != 1 {
		t.Errorf("expected 1 sessions revoke call, got %d", sessionRepoMock.revokeCalls)
	}
}

func TestDeleteCustomerAccount_ForbiddenForNonCustomers(t *testing.T) {
	ctx := context.Background()

	uc := NewDeleteCustomerAccountUsecase(
		&delMockTransactor{},
		&mockAccountRepo{},
		&mockOAuthRepo{},
		&mockSessionRepo{},
		&mockUserRepo{},
		&mockCustomerRepo{},
		&mockAddressRepo{},
		&mockCartRepo{},
		&delMockAuditLogger{},
	)

	// CustomerID is nil (e.g. staff member)
	authCtx := authenDomain.AuthContext{
		UserID:          uuid.New(),
		CustomerID:      nil,
		IsAuthenticated: true,
	}

	err := uc.Execute(ctx, authCtx)
	if err == nil {
		t.Fatal("expected error for non-customer deletion, got nil")
	}
}

func TestDeleteCustomerAccount_TransactionFailure(t *testing.T) {
	ctx := context.Background()
	customerID := uuid.New()
	userID := uuid.New()

	// Simulating error on customer delete
	customerRepoMock := &mockCustomerRepo{deleteErr: errors.New("db error")}

	uc := NewDeleteCustomerAccountUsecase(
		&delMockTransactor{},
		&mockAccountRepo{},
		&mockOAuthRepo{},
		&mockSessionRepo{},
		&mockUserRepo{},
		customerRepoMock,
		&mockAddressRepo{},
		&mockCartRepo{},
		&delMockAuditLogger{},
	)

	authCtx := authenDomain.AuthContext{
		UserID:          userID,
		CustomerID:      &customerID,
		IsAuthenticated: true,
	}

	err := uc.Execute(ctx, authCtx)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
