package usecase

import (
	"context"
	"testing"

	"service-core/internal/modules/courier/domain"
	shopDomain "service-core/internal/modules/shop/domain"
	shopRepo "service-core/internal/modules/shop/repository"
	transaction "service-core/internal/shared/transaction"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type mockExecutor struct{}

func (m *mockExecutor) Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error) {
	return pgconn.NewCommandTag("UPDATE 1"), nil
}
func (m *mockExecutor) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	return nil, nil
}
func (m *mockExecutor) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	return nil
}

type mockTransactor struct{}

func (m *mockTransactor) WithinTransaction(ctx context.Context, fn func(tx transaction.Executor) error) error {
	return fn(&mockExecutor{})
}

type mockCourierRepo struct {
	validCodes []string
}

func (m *mockCourierRepo) ListAll(ctx context.Context, exec transaction.Executor) ([]string, error) {
	return m.validCodes, nil
}
func (m *mockCourierRepo) GetActiveCodes(ctx context.Context, exec transaction.Executor, codes []string) ([]string, error) {
	return m.validCodes, nil
}
func (m *mockCourierRepo) ValidateCouriers(ctx context.Context, exec transaction.Executor, codes []string) ([]string, error) {
	return codes, nil
}

type mockShopCourierRepo struct {
	couriers map[string]*domain.ShopCourier
}

func (m *mockShopCourierRepo) ListByShopID(ctx context.Context, exec transaction.Executor, shopID uuid.UUID) ([]domain.ShopCourier, error) {
	var list []domain.ShopCourier
	for _, c := range m.couriers {
		list = append(list, *c)
	}
	return list, nil
}
func (m *mockShopCourierRepo) GetByShopIDAndCode(ctx context.Context, exec transaction.Executor, shopID uuid.UUID, code string) (*domain.ShopCourier, error) {
	if c, ok := m.couriers[code]; ok {
		return c, nil
	}
	return nil, nil
}
func (m *mockShopCourierRepo) ListsByShopIDs(ctx context.Context, exec transaction.Executor, shopIDs []uuid.UUID) (map[uuid.UUID][]domain.ShopCourier, error) {
	return nil, nil
}
func (m *mockShopCourierRepo) SaveShopCourier(ctx context.Context, exec transaction.Executor, sc domain.ShopCourier) error {
	if m.couriers == nil {
		m.couriers = make(map[string]*domain.ShopCourier)
	}
	copySc := sc
	m.couriers[sc.Code] = &copySc
	return nil
}
func (m *mockShopCourierRepo) VerifyShopCourier(ctx context.Context, exec transaction.Executor, shopID uuid.UUID, code string, status domain.CourierVerificationStatus, active bool, verifiedBy uuid.UUID, rejectionReason *string) error {
	if c, ok := m.couriers[code]; ok {
		c.VerificationStatus = status
		c.Active = active
		c.VerifiedBy = &verifiedBy
		c.RejectionReason = rejectionReason
	}
	return nil
}

type mockShopRepo struct {
	shopRepo.ShopRepository
	shop *shopDomain.Shop
}

func (m *mockShopRepo) GetByID(ctx context.Context, exec transaction.Executor, id uuid.UUID) (*shopDomain.Shop, error) {
	return m.shop, nil
}

func TestConfigureShopCourier_StaffActivateRequiresNameAndAddress(t *testing.T) {
	ctx := context.Background()
	shopID := uuid.New()
	shop := &shopDomain.Shop{ID: shopID, Name: "Test Shop"}

	sRepo := &mockShopRepo{shop: shop}
	cRepo := &mockCourierRepo{validCodes: []string{"jne"}}
	scRepo := &mockShopCourierRepo{couriers: make(map[string]*domain.ShopCourier)}
	exec := &mockExecutor{}
	tx := &mockTransactor{}

	uc := NewConfigureShopCourierUsecase(exec, tx, cRepo, scRepo, sRepo)

	// Test: Active toggle ON with missing name
	_, err := uc.UpdateSingle(ctx, UpdateSingleShopCourierInput{
		ShopID:          shopID,
		Code:            "jne",
		Name:            nil,
		LocationAddress: ptr("Jl. Melati 123"),
		Active:          true,
		IsAdmin:         false,
	})
	if err == nil || err.Error() != domain.ErrCourierNameRequired.Error() {
		t.Fatalf("expected ErrCourierNameRequired, got %v", err)
	}

	// Test: Active toggle ON with missing address
	_, err = uc.UpdateSingle(ctx, UpdateSingleShopCourierInput{
		ShopID:          shopID,
		Code:            "jne",
		Name:            ptr("JNE Pickup"),
		LocationAddress: nil,
		Active:          true,
		IsAdmin:         false,
	})
	if err == nil || err.Error() != domain.ErrCourierLocationRequired.Error() {
		t.Fatalf("expected ErrCourierLocationRequired, got %v", err)
	}

	// Test: Active toggle ON with valid name and address by staff -> status should be PENDING and active FALSE
	saved, err := uc.UpdateSingle(ctx, UpdateSingleShopCourierInput{
		ShopID:          shopID,
		Code:            "jne",
		Name:            ptr("JNE Pickup"),
		LocationAddress: ptr("Jl. Melati 123"),
		Active:          true,
		IsAdmin:         false,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if saved.VerificationStatus != domain.CourierVerificationPending {
		t.Errorf("expected status 'pending', got '%s'", saved.VerificationStatus)
	}
	if saved.Active != false {
		t.Errorf("expected active=false for staff activation before admin verification, got %v", saved.Active)
	}
}

func TestConfigureShopCourier_AdminDirectActivate(t *testing.T) {
	ctx := context.Background()
	shopID := uuid.New()
	adminID := uuid.New()
	shop := &shopDomain.Shop{ID: shopID, Name: "Test Shop"}

	sRepo := &mockShopRepo{shop: shop}
	cRepo := &mockCourierRepo{validCodes: []string{"jne"}}
	scRepo := &mockShopCourierRepo{couriers: make(map[string]*domain.ShopCourier)}
	exec := &mockExecutor{}
	tx := &mockTransactor{}

	uc := NewConfigureShopCourierUsecase(exec, tx, cRepo, scRepo, sRepo)

	// Admin direct activation -> status should be VERIFIED and active TRUE
	saved, err := uc.UpdateSingle(ctx, UpdateSingleShopCourierInput{
		ShopID:          shopID,
		Code:            "jne",
		Name:            ptr("JNE Admin Configured"),
		LocationAddress: ptr("Jl. Gatot Subroto 45"),
		Active:          true,
		IsAdmin:         true,
		AdminStaffID:    &adminID,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if saved.VerificationStatus != domain.CourierVerificationVerified {
		t.Errorf("expected status 'verified', got '%s'", saved.VerificationStatus)
	}
	if saved.Active != true {
		t.Errorf("expected active=true for admin activation, got %v", saved.Active)
	}
}

func TestConfigureShopCourier_Execute_BulkActivationRequiresDetails(t *testing.T) {
	ctx := context.Background()
	shopID := uuid.New()
	shop := &shopDomain.Shop{ID: shopID, Name: "Test Shop"}

	sRepo := &mockShopRepo{shop: shop}
	cRepo := &mockCourierRepo{validCodes: []string{"jne"}}
	scRepo := &mockShopCourierRepo{couriers: make(map[string]*domain.ShopCourier)}
	exec := &mockExecutor{}
	tx := &mockTransactor{}

	uc := NewConfigureShopCourierUsecase(exec, tx, cRepo, scRepo, sRepo)

	// Attempt bulk activate without existing name/address -> must fail
	err := uc.Execute(ctx, shopID, []ConfigureShopCourierInput{
		{Code: "jne", Active: true},
	})
	if err == nil {
		t.Fatal("expected error when bulk activating courier without branch name/address, got nil")
	}

	// Deactivate via bulk -> should succeed even if empty
	err = uc.Execute(ctx, shopID, []ConfigureShopCourierInput{
		{Code: "jne", Active: false},
	})
	if err != nil {
		t.Fatalf("unexpected error when bulk deactivating: %v", err)
	}
}

func TestConfigureShopCourier_DeactivateAllowsEmptyDetails(t *testing.T) {
	ctx := context.Background()
	shopID := uuid.New()
	shop := &shopDomain.Shop{ID: shopID, Name: "Test Shop"}

	sRepo := &mockShopRepo{shop: shop}
	cRepo := &mockCourierRepo{validCodes: []string{"jne"}}
	scRepo := &mockShopCourierRepo{couriers: make(map[string]*domain.ShopCourier)}
	exec := &mockExecutor{}
	tx := &mockTransactor{}

	uc := NewConfigureShopCourierUsecase(exec, tx, cRepo, scRepo, sRepo)

	// Deactivate with nil name and address -> should succeed
	saved, err := uc.UpdateSingle(ctx, UpdateSingleShopCourierInput{
		ShopID:          shopID,
		Code:            "jne",
		Name:            nil,
		LocationAddress: nil,
		Active:          false,
		IsAdmin:         false,
	})
	if err != nil {
		t.Fatalf("unexpected error deactivating courier: %v", err)
	}
	if saved.Active != false {
		t.Errorf("expected active=false, got %v", saved.Active)
	}
}

func ptr(s string) *string {
	return &s
}
