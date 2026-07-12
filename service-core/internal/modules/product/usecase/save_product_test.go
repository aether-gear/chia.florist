package usecase

import (
	"context"
	"errors"
	"testing"

	"service-core/internal/modules/product/domain"
	"service-core/internal/modules/product/repository"
	transaction "service-core/internal/shared/transaction"

	"github.com/google/uuid"
)

// ===========================================================================
// Mocks for SaveProductUsecase Tests
// ===========================================================================

type mockProductRepository struct {
	repository.ProductRepository
	savedProduct *domain.Product
	saveErr      error
}

func (m *mockProductRepository) SaveProduct(
	ctx context.Context,
	exec transaction.Executor,
	product *domain.Product,
) error {
	m.savedProduct = product
	return m.saveErr
}

type mockSlugGenerator struct{}

func (m *mockSlugGenerator) Generate(input string) string {
	return "mocked-" + input
}

type mockExecutor struct {
	transaction.Executor
}

// ===========================================================================
// Tests
// ===========================================================================

func TestSaveProduct_Success_NewProduct(t *testing.T) {
	ctx := context.Background()
	repo := &mockProductRepository{}
	slugGen := &mockSlugGenerator{}
	exec := &mockExecutor{}

	uc := NewSaveProductUsecase(repo, slugGen, exec)

	desc := "New product description"
	weight := 1.5

	input := SaveProductInput{
		ID:          nil,
		SKU:         "PROD-001",
		Name:        "Test Product",
		Description: &desc,
		Status:      "active",
		Price:       10000,
		Weight:      &weight,
	}

	err := uc.Execute(ctx, input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if repo.savedProduct == nil {
		t.Fatal("expected product to be saved, got nil")
	}

	if repo.savedProduct.ID == uuid.Nil {
		t.Error("expected new UUID to be generated, got Nil")
	}

	if repo.savedProduct.SKU != "PROD-001" {
		t.Errorf("expected SKU 'PROD-001', got '%s'", repo.savedProduct.SKU)
	}

	if repo.savedProduct.Name != "Test Product" {
		t.Errorf("expected Name 'Test Product', got '%s'", repo.savedProduct.Name)
	}

	if repo.savedProduct.Slug != "mocked-Test Product" {
		t.Errorf("expected Slug 'mocked-Test Product', got '%s'", repo.savedProduct.Slug)
	}

	if *repo.savedProduct.Description != desc {
		t.Errorf("expected Description '%s', got '%s'", desc, *repo.savedProduct.Description)
	}

	if repo.savedProduct.Status != domain.ProductStatusActive {
		t.Errorf("expected Status 'active', got '%s'", repo.savedProduct.Status)
	}

	if repo.savedProduct.Price != 10000 {
		t.Errorf("expected Price 10000, got %d", repo.savedProduct.Price)
	}

	if *repo.savedProduct.Weight != weight {
		t.Errorf("expected Weight %v, got %v", weight, *repo.savedProduct.Weight)
	}
}

func TestSaveProduct_Success_UpdateProduct(t *testing.T) {
	ctx := context.Background()
	repo := &mockProductRepository{}
	slugGen := &mockSlugGenerator{}
	exec := &mockExecutor{}

	uc := NewSaveProductUsecase(repo, slugGen, exec)

	productID := uuid.New()
	desc := "Updated description"
	weight := 2.0

	input := SaveProductInput{
		ID:          &productID,
		SKU:         "PROD-001-UPD",
		Name:        "Updated Product",
		Description: &desc,
		Status:      "inactive",
		Price:       15000,
		Weight:      &weight,
	}

	err := uc.Execute(ctx, input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if repo.savedProduct == nil {
		t.Fatal("expected product to be saved, got nil")
	}

	if repo.savedProduct.ID != productID {
		t.Errorf("expected product ID %v, got %v", productID, repo.savedProduct.ID)
	}

	if repo.savedProduct.SKU != "PROD-001-UPD" {
		t.Errorf("expected SKU 'PROD-001-UPD', got '%s'", repo.savedProduct.SKU)
	}

	if repo.savedProduct.Name != "Updated Product" {
		t.Errorf("expected Name 'Updated Product', got '%s'", repo.savedProduct.Name)
	}

	if repo.savedProduct.Slug != "mocked-Updated Product" {
		t.Errorf("expected Slug 'mocked-Updated Product', got '%s'", repo.savedProduct.Slug)
	}

	if *repo.savedProduct.Description != desc {
		t.Errorf("expected Description '%s', got '%s'", desc, *repo.savedProduct.Description)
	}

	if repo.savedProduct.Status != domain.ProductStatusInactive {
		t.Errorf("expected Status 'inactive', got '%s'", repo.savedProduct.Status)
	}

	if repo.savedProduct.Price != 15000 {
		t.Errorf("expected Price 15000, got %d", repo.savedProduct.Price)
	}

	if *repo.savedProduct.Weight != weight {
		t.Errorf("expected Weight %v, got %v", weight, *repo.savedProduct.Weight)
	}
}

func TestSaveProduct_InvalidInput_EmptyName(t *testing.T) {
	ctx := context.Background()
	repo := &mockProductRepository{}
	slugGen := &mockSlugGenerator{}
	exec := &mockExecutor{}

	uc := NewSaveProductUsecase(repo, slugGen, exec)

	input := SaveProductInput{
		ID:     nil,
		SKU:    "PROD-001",
		Name:   "",
		Status: "active",
		Price:  10000,
	}

	err := uc.Execute(ctx, input)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestSaveProduct_InvalidInput_NegativePrice(t *testing.T) {
	ctx := context.Background()
	repo := &mockProductRepository{}
	slugGen := &mockSlugGenerator{}
	exec := &mockExecutor{}

	uc := NewSaveProductUsecase(repo, slugGen, exec)

	input := SaveProductInput{
		ID:     nil,
		SKU:    "PROD-001",
		Name:   "Test Product",
		Status: "active",
		Price:  -50,
	}

	err := uc.Execute(ctx, input)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestSaveProduct_RepoError(t *testing.T) {
	ctx := context.Background()
	expectedErr := errors.New("database connection failed")
	repo := &mockProductRepository{saveErr: expectedErr}
	slugGen := &mockSlugGenerator{}
	exec := &mockExecutor{}

	uc := NewSaveProductUsecase(repo, slugGen, exec)

	input := SaveProductInput{
		ID:     nil,
		SKU:    "PROD-001",
		Name:   "Test Product",
		Status: "active",
		Price:  10000,
	}

	err := uc.Execute(ctx, input)
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected error '%v', got '%v'", expectedErr, err)
	}
}
