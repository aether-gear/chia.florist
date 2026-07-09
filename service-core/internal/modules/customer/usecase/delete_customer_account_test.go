package usecase

import (
	"context"
	"errors"
	"testing"

	applogger "service-core/internal/common/logger"
	addressDomain "service-core/internal/modules/address/domain"
	authenDomain "service-core/internal/modules/authentication/domain"
	cartDomain "service-core/internal/modules/cart/domain"
	"service-core/internal/modules/customer/domain"
	"service-core/internal/modules/customer/repository"
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

func (m *mockCustomerRepo) Create(_ context.Context, _ transaction.Executor, _ domain.Customer) error {
	return nil
}
func (m *mockCustomerRepo) GetByID(_ context.Context, _ transaction.Executor, _ uuid.UUID) (*domain.Customer, error) {
	return nil, nil
}
func (m *mockCustomerRepo) GetByUserID(_ context.Context, _ transaction.Executor, _ uuid.UUID) (*domain.Customer, error) {
	return nil, nil
}
func (m *mockCustomerRepo) GetProfileByUserID(_ context.Context, _ transaction.Executor, _ uuid.UUID) (*domain.CustomerProfile, error) {
	return nil, nil
}
func (m *mockCustomerRepo) FindCustomers(_ context.Context, _ transaction.Executor, _ repository.FindCustomerParams) ([]domain.CustomerProfile, int, error) {
	return nil, 0, nil
}
func (m *mockCustomerRepo) Delete(_ context.Context, _ transaction.Executor, _ uuid.UUID) error {
	m.deleteCalls++
	return m.deleteErr
}

// UserDeletionService Mock
type mockUserDeletionSvc struct {
	deleteCalls int
	deleteErr   error
}

func (m *mockUserDeletionSvc) DeleteUserRecord(_ context.Context, _ transaction.Executor, _ uuid.UUID) error {
	m.deleteCalls++
	return m.deleteErr
}

// ===========================================================================
// Tests
// ===========================================================================

func TestDeleteCustomerAccount_Success(t *testing.T) {
	ctx := context.Background()
	customerID := uuid.New()
	userID := uuid.New()

	customerRepoMock := &mockCustomerRepo{}
	userDeletionSvcMock := &mockUserDeletionSvc{}
	addressRepoMock := &mockAddressRepo{}
	cartRepoMock := &mockCartRepo{}

	uc := NewDeleteCustomerAccountUsecase(
		&delMockTransactor{},
		userDeletionSvcMock,
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
	if userDeletionSvcMock.deleteCalls != 1 {
		t.Errorf("expected 1 user deletion service call, got %d", userDeletionSvcMock.deleteCalls)
	}
	if addressRepoMock.deleteCalls != 1 {
		t.Errorf("expected 1 address delete call, got %d", addressRepoMock.deleteCalls)
	}
	if cartRepoMock.deleteCalls != 1 {
		t.Errorf("expected 1 cart delete call, got %d", cartRepoMock.deleteCalls)
	}
}

func TestDeleteCustomerAccount_ForbiddenForNonCustomers(t *testing.T) {
	ctx := context.Background()

	uc := NewDeleteCustomerAccountUsecase(
		&delMockTransactor{},
		&mockUserDeletionSvc{},
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
		&mockUserDeletionSvc{},
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
