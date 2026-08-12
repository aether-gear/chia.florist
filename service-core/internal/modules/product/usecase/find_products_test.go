package usecase

import (
	"context"
	"testing"

	storage "service-core/internal/infra/storage"
	inventoryDomain "service-core/internal/modules/inventory/domain"
	inventoryRepo "service-core/internal/modules/inventory/repository"
	"service-core/internal/modules/product/domain"
	"service-core/internal/modules/product/repository"
	shopDomain "service-core/internal/modules/shop/domain"
	shopRepo "service-core/internal/modules/shop/repository"
	transaction "service-core/internal/shared/transaction"

	"github.com/google/uuid"
)

type mockFindProductRepo struct {
	repository.ProductRepository
	capturedParams repository.FindProductParams
	products       []domain.ProductWithInventory
}

func (m *mockFindProductRepo) FindProductsWithInventory(
	ctx context.Context,
	exec transaction.Executor,
	params repository.FindProductParams,
) ([]domain.ProductWithInventory, int, error) {
	m.capturedParams = params
	return m.products, len(m.products), nil
}

type mockFindInventoryRepo struct {
	inventoryRepo.InventoryRepository
}

func (m *mockFindInventoryRepo) ListByProductIDs(
	ctx context.Context,
	exec transaction.Executor,
	productIDs []uuid.UUID,
) (map[uuid.UUID][]inventoryDomain.Inventory, error) {
	return map[uuid.UUID][]inventoryDomain.Inventory{}, nil
}

type mockFindImgRepo struct {
	repository.ProductImageRepository
}

func (m *mockFindImgRepo) ListByProductIDs(
	ctx context.Context,
	exec transaction.Executor,
	productIDs []uuid.UUID,
) (map[uuid.UUID][]domain.ProductImage, error) {
	return map[uuid.UUID][]domain.ProductImage{}, nil
}

type mockFindShopRepo struct {
	shopRepo.ShopRepository
}

func (m *mockFindShopRepo) FindByIDs(
	ctx context.Context,
	exec transaction.Executor,
	IDs []uuid.UUID,
) ([]shopDomain.Shop, error) {
	return []shopDomain.Shop{}, nil
}

type mockFileStore struct {
	storage.Provider
}

func (m *mockFileStore) PublicURL(key, bucket string) string {
	return "http://localhost/" + key
}

func TestFindProducts_WithShopFilter(t *testing.T) {
	ctx := context.Background()
	shopID := uuid.New()
	shopIDStr := shopID.String()
	shopSlug := "central-store"

	productRepo := &mockFindProductRepo{}
	invRepo := &mockFindInventoryRepo{}
	imgRepo := &mockFindImgRepo{}
	shopRepo := &mockFindShopRepo{}
	fileStore := &mockFileStore{}
	exec := &mockExecutor{}

	uc := NewFindProductsUsecase(productRepo, invRepo, imgRepo, shopRepo, fileStore, exec)

	// Case 1: Query with ShopID
	_, _, err := uc.Execute(ctx, FindProductsInput{
		ShopID: &shopIDStr,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if productRepo.capturedParams.ShopID == nil || *productRepo.capturedParams.ShopID != shopID {
		t.Errorf("expected ShopID %v, got %v", shopID, productRepo.capturedParams.ShopID)
	}

	// Case 2: Query with ShopSlug
	_, _, err = uc.Execute(ctx, FindProductsInput{
		ShopSlug: &shopSlug,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if productRepo.capturedParams.ShopSlug == nil || *productRepo.capturedParams.ShopSlug != shopSlug {
		t.Errorf("expected ShopSlug %v, got %v", shopSlug, productRepo.capturedParams.ShopSlug)
	}

	// Case 3: Query without Shop filter (legacy backward compatibility)
	_, _, err = uc.Execute(ctx, FindProductsInput{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if productRepo.capturedParams.ShopID != nil || productRepo.capturedParams.ShopSlug != nil {
		t.Errorf("expected nil shop filters for legacy query, got shopID: %v, shopSlug: %v",
			productRepo.capturedParams.ShopID, productRepo.capturedParams.ShopSlug)
	}
}
