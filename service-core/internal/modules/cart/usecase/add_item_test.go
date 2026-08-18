package usecase

import (
	"context"
	"testing"

	cartDomain "service-core/internal/modules/cart/domain"
	inventoryDomain "service-core/internal/modules/inventory/domain"
	inventoryRepo "service-core/internal/modules/inventory/repository"
	productDomain "service-core/internal/modules/product/domain"
	productRepo "service-core/internal/modules/product/repository"
	shopDomain "service-core/internal/modules/shop/domain"
	transaction "service-core/internal/shared/transaction"

	"github.com/google/uuid"
)

type mockAddItemInvRepo struct {
	inventoryRepo.InventoryRepository
	inv *inventoryDomain.Inventory
}

func (m *mockAddItemInvRepo) GetByProductIDAndShopID(
	ctx context.Context,
	exec transaction.Executor,
	productID, shopID uuid.UUID,
) (*inventoryDomain.Inventory, error) {
	return m.inv, nil
}

type mockAddItemProdRepo struct {
	productRepo.ProductRepository
	prod *productDomain.Product
}

func (m *mockAddItemProdRepo) GetByID(
	ctx context.Context,
	exec transaction.Executor,
	id uuid.UUID,
) (*productDomain.Product, error) {
	return m.prod, nil
}

func TestAddItem_Success(t *testing.T) {
	ctx := context.Background()
	customerID := uuid.New()
	shopID := uuid.New()
	productID := uuid.New()

	cartR := &mockSaveCartRepository{
		cart: &cartDomain.Cart{
			ID:         uuid.New(),
			CustomerID: customerID,
			Items:      []cartDomain.CartItem{},
		},
	}
	shopR := &mockCartShopRepo{
		shop: &shopDomain.Shop{
			ID:             shopID,
			Name:           "Operable Shop",
			IsActive:       true,
			ApprovalStatus: shopDomain.ShopApprovalStatusApproved,
		},
	}
	invR := &mockAddItemInvRepo{
		inv: &inventoryDomain.Inventory{
			ProductID:  productID,
			ShopID:     shopID,
			TotalStock: 10,
		},
	}
	prodR := &mockAddItemProdRepo{
		prod: &productDomain.Product{
			ID:     productID,
			Name:   "Rose Bouquet",
			Status: productDomain.ProductStatusActive,
		},
	}
	tx := &mockTransactor{}
	exec := &mockExecutor{}

	uc := NewAddItemUsecase(exec, tx, cartR, invR, prodR, shopR)

	err := uc.Execute(ctx, AddItemInput{
		CustomerID: customerID,
		ShopID:     shopID,
		ProductID:  productID,
		Quantity:   2,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !cartR.saved {
		t.Errorf("expected cart to be saved")
	}
}

func TestAddItem_RejectInactiveShop(t *testing.T) {
	ctx := context.Background()
	customerID := uuid.New()
	shopID := uuid.New()
	productID := uuid.New()

	cartR := &mockSaveCartRepository{}
	shopR := &mockCartShopRepo{
		shop: &shopDomain.Shop{
			ID:             shopID,
			Name:           "Inactive Shop",
			IsActive:       false,
			ApprovalStatus: shopDomain.ShopApprovalStatusApproved,
		},
	}
	invR := &mockAddItemInvRepo{}
	prodR := &mockAddItemProdRepo{}
	tx := &mockTransactor{}
	exec := &mockExecutor{}

	uc := NewAddItemUsecase(exec, tx, cartR, invR, prodR, shopR)

	err := uc.Execute(ctx, AddItemInput{
		CustomerID: customerID,
		ShopID:     shopID,
		ProductID:  productID,
		Quantity:   1,
	})

	if err == nil {
		t.Errorf("expected error for inactive shop, got nil")
	}
}

func TestAddItem_RejectPendingShop(t *testing.T) {
	ctx := context.Background()
	customerID := uuid.New()
	shopID := uuid.New()
	productID := uuid.New()

	cartR := &mockSaveCartRepository{}
	shopR := &mockCartShopRepo{
		shop: &shopDomain.Shop{
			ID:             shopID,
			Name:           "Pending Shop",
			IsActive:       true,
			ApprovalStatus: shopDomain.ShopApprovalStatusPending,
		},
	}
	invR := &mockAddItemInvRepo{}
	prodR := &mockAddItemProdRepo{}
	tx := &mockTransactor{}
	exec := &mockExecutor{}

	uc := NewAddItemUsecase(exec, tx, cartR, invR, prodR, shopR)

	err := uc.Execute(ctx, AddItemInput{
		CustomerID: customerID,
		ShopID:     shopID,
		ProductID:  productID,
		Quantity:   1,
	})

	if err == nil {
		t.Errorf("expected error for pending shop, got nil")
	}
}

func TestAddItem_RejectInactiveProduct(t *testing.T) {
	ctx := context.Background()
	customerID := uuid.New()
	shopID := uuid.New()
	productID := uuid.New()

	cartR := &mockSaveCartRepository{
		cart: &cartDomain.Cart{
			ID:         uuid.New(),
			CustomerID: customerID,
			Items:      []cartDomain.CartItem{},
		},
	}
	shopR := &mockCartShopRepo{
		shop: &shopDomain.Shop{
			ID:             shopID,
			Name:           "Active Shop",
			IsActive:       true,
			ApprovalStatus: shopDomain.ShopApprovalStatusApproved,
		},
	}
	invR := &mockAddItemInvRepo{
		inv: &inventoryDomain.Inventory{
			ProductID:  productID,
			ShopID:     shopID,
			TotalStock: 10,
		},
	}
	prodR := &mockAddItemProdRepo{
		prod: &productDomain.Product{
			ID:     productID,
			Name:   "Inactive Flower",
			Status: productDomain.ProductStatusInactive,
		},
	}
	tx := &mockTransactor{}
	exec := &mockExecutor{}

	uc := NewAddItemUsecase(exec, tx, cartR, invR, prodR, shopR)

	err := uc.Execute(ctx, AddItemInput{
		CustomerID: customerID,
		ShopID:     shopID,
		ProductID:  productID,
		Quantity:   1,
	})

	if err == nil {
		t.Errorf("expected error for inactive product, got nil")
	}
}

func TestAddItem_RejectArchivedProduct(t *testing.T) {
	ctx := context.Background()
	customerID := uuid.New()
	shopID := uuid.New()
	productID := uuid.New()

	cartR := &mockSaveCartRepository{
		cart: &cartDomain.Cart{
			ID:         uuid.New(),
			CustomerID: customerID,
			Items:      []cartDomain.CartItem{},
		},
	}
	shopR := &mockCartShopRepo{
		shop: &shopDomain.Shop{
			ID:             shopID,
			Name:           "Active Shop",
			IsActive:       true,
			ApprovalStatus: shopDomain.ShopApprovalStatusApproved,
		},
	}
	invR := &mockAddItemInvRepo{
		inv: &inventoryDomain.Inventory{
			ProductID:  productID,
			ShopID:     shopID,
			TotalStock: 10,
		},
	}
	prodR := &mockAddItemProdRepo{
		prod: &productDomain.Product{
			ID:     productID,
			Name:   "Archived Flower",
			Status: productDomain.ProductStatusArchived,
		},
	}
	tx := &mockTransactor{}
	exec := &mockExecutor{}

	uc := NewAddItemUsecase(exec, tx, cartR, invR, prodR, shopR)

	err := uc.Execute(ctx, AddItemInput{
		CustomerID: customerID,
		ShopID:     shopID,
		ProductID:  productID,
		Quantity:   1,
	})

	if err == nil {
		t.Errorf("expected error (not found) for archived product, got nil")
	}
}
