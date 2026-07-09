package usecase

import (
	"context"
	"errors"
	"testing"

	applogger "service-core/internal/common/logger"
	authenDomain "service-core/internal/modules/authentication/domain"
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

// UserDeletionService Mock
type mockUserDeletionSvc struct {
	deleteCalls int
	deleteErr   error
}

func (m *mockUserDeletionSvc) DeleteUserRecord(_ context.Context, _ transaction.Executor, _ uuid.UUID) error {
	m.deleteCalls++
	return m.deleteErr
}

// CustomerDeletionService Mock
type mockCustomerDeletionSvc struct {
	deleteCalls int
	deleteErr   error
}

func (m *mockCustomerDeletionSvc) DeleteCustomerRecord(_ context.Context, _ transaction.Executor, _ uuid.UUID) error {
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

	customerDeletionSvcMock := &mockCustomerDeletionSvc{}
	userDeletionSvcMock := &mockUserDeletionSvc{}

	uc := NewDeleteCustomerAccountUsecase(
		&delMockTransactor{},
		userDeletionSvcMock,
		customerDeletionSvcMock,
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

	if customerDeletionSvcMock.deleteCalls != 1 {
		t.Errorf("expected 1 customer deletion service call, got %d", customerDeletionSvcMock.deleteCalls)
	}
	if userDeletionSvcMock.deleteCalls != 1 {
		t.Errorf("expected 1 user deletion service call, got %d", userDeletionSvcMock.deleteCalls)
	}
}

func TestDeleteCustomerAccount_ForbiddenForNonCustomers(t *testing.T) {
	ctx := context.Background()

	uc := NewDeleteCustomerAccountUsecase(
		&delMockTransactor{},
		&mockUserDeletionSvc{},
		&mockCustomerDeletionSvc{},
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

	// Simulating error on customer deletion service
	customerDeletionSvcMock := &mockCustomerDeletionSvc{deleteErr: errors.New("db error")}

	uc := NewDeleteCustomerAccountUsecase(
		&delMockTransactor{},
		&mockUserDeletionSvc{},
		customerDeletionSvcMock,
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
