package usecase

import (
	"context"
	"testing"

	"service-core/internal/modules/courier/domain"
	shopDomain "service-core/internal/modules/shop/domain"

	"github.com/google/uuid"
)

func TestVerifyShopCourier_ApproveAndReject(t *testing.T) {
	ctx := context.Background()
	shopID := uuid.New()
	adminID := uuid.New()
	shop := &shopDomain.Shop{ID: shopID, Name: "Test Shop"}

	sRepo := &mockShopRepo{shop: shop}
	cRepo := &mockCourierRepo{validCodes: []string{"sicepat"}}
	scRepo := &mockShopCourierRepo{couriers: make(map[string]*domain.ShopCourier)}
	exec := &mockExecutor{}
	tx := &mockTransactor{}

	// Seed a pending courier
	courierName := "SiCepat Hub 1"
	courierAddr := "Jl. Sudirman No. 1"
	scRepo.couriers["sicepat"] = &domain.ShopCourier{
		ShopID:             shopID,
		Code:               "sicepat",
		Name:               &courierName,
		LocationAddress:    &courierAddr,
		Active:             false,
		VerificationStatus: domain.CourierVerificationPending,
	}

	uc := NewVerifyShopCourierUsecase(exec, tx, cRepo, scRepo, sRepo)

	// Test: Admin Approves/Verifies
	verified, err := uc.Execute(ctx, VerifyShopCourierInput{
		ShopID:       shopID,
		Code:         "sicepat",
		Action:       "verify",
		AdminStaffID: adminID,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if verified.VerificationStatus != domain.CourierVerificationVerified {
		t.Errorf("expected status 'verified', got '%s'", verified.VerificationStatus)
	}
	if verified.Active != true {
		t.Errorf("expected active=true after verify, got %v", verified.Active)
	}

	// Test: Admin Rejects with reason
	rejectReason := "Address could not be verified by logistic partner"
	rejected, err := uc.Execute(ctx, VerifyShopCourierInput{
		ShopID:          shopID,
		Code:            "sicepat",
		Action:          "reject",
		RejectionReason: &rejectReason,
		AdminStaffID:    adminID,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rejected.VerificationStatus != domain.CourierVerificationRejected {
		t.Errorf("expected status 'rejected', got '%s'", rejected.VerificationStatus)
	}
	if rejected.Active != false {
		t.Errorf("expected active=false after reject, got %v", rejected.Active)
	}
	if rejected.RejectionReason == nil || *rejected.RejectionReason != rejectReason {
		t.Errorf("expected rejection reason '%s', got %v", rejectReason, rejected.RejectionReason)
	}
}

func TestVerifyShopCourier_RequiresBranchNameAndLocationAddress(t *testing.T) {
	ctx := context.Background()
	shopID := uuid.New()
	adminID := uuid.New()
	shop := &shopDomain.Shop{ID: shopID, Name: "Test Shop"}

	sRepo := &mockShopRepo{shop: shop}
	cRepo := &mockCourierRepo{validCodes: []string{"sicepat"}}
	scRepo := &mockShopCourierRepo{couriers: make(map[string]*domain.ShopCourier)}
	exec := &mockExecutor{}
	tx := &mockTransactor{}

	// Seed courier missing branch name
	courierAddr := "Jl. Sudirman No. 1"
	scRepo.couriers["sicepat"] = &domain.ShopCourier{
		ShopID:             shopID,
		Code:               "sicepat",
		Name:               nil,
		LocationAddress:    &courierAddr,
		Active:             false,
		VerificationStatus: domain.CourierVerificationPending,
	}

	uc := NewVerifyShopCourierUsecase(exec, tx, cRepo, scRepo, sRepo)

	// Verifying courier with missing branch name must fail
	_, err := uc.Execute(ctx, VerifyShopCourierInput{
		ShopID:       shopID,
		Code:         "sicepat",
		Action:       "verify",
		AdminStaffID: adminID,
	})
	if err == nil {
		t.Fatal("expected error when verifying courier without branch name, got nil")
	}

	// Now seed courier with missing address
	courierName := "SiCepat Hub"
	scRepo.couriers["sicepat"] = &domain.ShopCourier{
		ShopID:             shopID,
		Code:               "sicepat",
		Name:               &courierName,
		LocationAddress:    nil,
		Active:             false,
		VerificationStatus: domain.CourierVerificationPending,
	}

	// Verifying courier with missing location address must fail
	_, err = uc.Execute(ctx, VerifyShopCourierInput{
		ShopID:       shopID,
		Code:         "sicepat",
		Action:       "verify",
		AdminStaffID: adminID,
	})
	if err == nil {
		t.Fatal("expected error when verifying courier without location address, got nil")
	}
}

