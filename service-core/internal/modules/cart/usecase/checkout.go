package usecase

import (
	"context"
	"fmt"
	"strconv"

	apperrors "service-core/internal/common/errors"
	addressDomain "service-core/internal/modules/address/domain"
	addressRepo "service-core/internal/modules/address/repository"
	authenDomain "service-core/internal/modules/authentication/domain"
	"service-core/internal/modules/cart/domain"
	courierRepo "service-core/internal/modules/courier/repository"
	inventoryDomain "service-core/internal/modules/inventory/domain"
	inventoryRepo "service-core/internal/modules/inventory/repository"
	productDomain "service-core/internal/modules/product/domain"
	productRepo "service-core/internal/modules/product/repository"
	shipmentRepo "service-core/internal/modules/shipment/repository"
	shopDomain "service-core/internal/modules/shop/domain"
	shopRepo "service-core/internal/modules/shop/repository"
	transaction "service-core/internal/shared/transaction"

	"github.com/google/uuid"
)

type CheckoutUsecase struct {
	executor        transaction.Executor
	addressRepo     addressRepo.UserAddressRepository
	courierShopRepo courierRepo.ShopCourierRepository
	inventoryRepo   inventoryRepo.InventoryRepository
	productRepo     productRepo.ProductRepository
	shipmentRepo    shipmentRepo.ShippingCostProvider
	shopAddressRepo addressRepo.ShopAddressRepository
	shopRepo        shopRepo.ShopRepository
}

func NewCheckoutUsecase(
	executor transaction.Executor,
	addressRepo addressRepo.UserAddressRepository,
	courierShopRepo courierRepo.ShopCourierRepository,
	inventoryRepo inventoryRepo.InventoryRepository,
	productRepo productRepo.ProductRepository,
	shipmentRepo shipmentRepo.ShippingCostProvider,
	shopAddressRepo addressRepo.ShopAddressRepository,
	shopRepo shopRepo.ShopRepository,
) *CheckoutUsecase {
	return &CheckoutUsecase{
		executor:        executor,
		addressRepo:     addressRepo,
		courierShopRepo: courierShopRepo,
		inventoryRepo:   inventoryRepo,
		productRepo:     productRepo,
		shipmentRepo:    shipmentRepo,
		shopAddressRepo: shopAddressRepo,
		shopRepo:        shopRepo,
	}
}

type CheckoutItemInput struct {
	ProductID uuid.UUID
	Quantity  int
}

type SelectedCourierInput struct {
	Code    string
	Service string
}

type CheckoutShopInput struct {
	ShopID  uuid.UUID
	Items   []CheckoutItemInput
	Courier *SelectedCourierInput
}

type CheckoutInput struct {
	AddressID *uuid.UUID
	ShopInput []CheckoutShopInput
}

type CheckoutAddressResult struct {
	ID            uuid.UUID
	RecipientName string
	Phone         *string
	FullAddress   string
}

type CheckoutCouriersResult struct {
	Code    string
	Service string
	Name    string
	ETD     string
	Fee     int64
}

type CheckoutItemResult struct {
	ProductID   uuid.UUID
	ShopID      uuid.UUID
	Name        string
	Price       int64
	Quantity    int
	Subtotal    int64
	TotalWeight int
}

type SelectedCourierResult struct {
	Code    string
	Service string
	Fee     int64
}

type ShopResult struct {
	ShopID          uuid.UUID
	ShopSlug        string
	ShopName        string
	Items           []CheckoutItemResult
	SelectedCourier *SelectedCourierResult
	CostCouriers    []CheckoutCouriersResult
	Subtotal        int64
	Total           int64
}

type CheckoutResult struct {
	Address          CheckoutAddressResult
	Shops            []ShopResult
	Subtotal         int64
	TotalShippingFee int64
	TotalAll         int64
}

const defaultShippingWeightGrams = 1000

func (u *CheckoutUsecase) Execute(
	ctx context.Context,
	authCtx authenDomain.AuthContext,
	input CheckoutInput,
) (*CheckoutResult, error) {
	// Early fallback.
	// Make sure the default address exists
	defAddress, err := u.addressRepo.
		GetDefaultByUserID(ctx, u.executor, authCtx.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve destination address: %w", err)
	}
	if defAddress == nil {
		return nil, apperrors.NewConflict(addressDomain.ErrNotFoundDefaultAddress.Error())
	}

	// Collect IDs upfront to avoid per-item database queries.
	// Products are loaded for inventory checks,
	// and shops are loaded for origin address
	// and courier availability.
	var productIDs []uuid.UUID
	var shopIDs []uuid.UUID
	for _, shopGroup := range input.ShopInput {
		for _, item := range shopGroup.Items {
			productIDs = append(productIDs, item.ProductID)
		}

		shopIDs = append(shopIDs, shopGroup.ShopID)
	}

	products, err := u.productRepo.
		FindByIDs(ctx, u.executor, productIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to load products: %w", err)
	}

	productMap := make(
		map[uuid.UUID]productDomain.Product,
		len(products),
	)
	for _, p := range products {
		productMap[p.ID] = p
	}

	inventoryMap, err := u.inventoryRepo.
		ListByProductIDs(ctx, u.executor, productIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to load inventory for cart products: %w", err)
	}

	courierShopMap, err := u.courierShopRepo.
		ListsByShopIDs(ctx, u.executor, shopIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve shop courier: %w", err)
	}

	originAddresses, err := u.shopAddressRepo.
		GetDefaultsByShopIDs(ctx, u.executor, shopIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve origin addresses: %w", err)
	}

	shops, err := u.shopRepo.
		FindByIDs(ctx, u.executor, shopIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to load products: %w", err)
	}

	shopMap := make(
		map[uuid.UUID]shopDomain.Shop,
		len(shops),
	)
	for _, s := range shops {
		shopMap[s.ID] = s
	}

	// Checkout shipping destination defaults to the user's
	// primary address unless a specific address is requested.
	var destAddress addressDomain.Address
	if input.AddressID != nil {
		address, err := u.addressRepo.
			GetByID(ctx, u.executor, *input.AddressID)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve destination address: %w", err)
		}
		if address == nil {
			return nil, apperrors.NewNotFound(addressDomain.ErrAddressNotFound.Error())
		}

		destAddress = *address
	} else {
		destAddress = *defAddress
	}

	destDistrictID, err := strconv.Atoi(destAddress.Detail.DistrictID)
	if err != nil {
		return nil, fmt.Errorf("invalid destination district id: %w", err)
	}

	// Calculate estimation, subtotal and weight
	// based on shop item (giant loop)
	var (
		totalWeight      int = 0
		totalSubtotal    int64
		totalShippingFee int64
		shopsResult      []ShopResult
	)
	for _, shopGroup := range input.ShopInput {
		// Map shop items, so products consist of
		// having inventory and are owned by a shop ID
		//
		// Include calculations of weight, total subtotal
		// and estimation total
		var shopSubtotal int64
		var checkoutItems []CheckoutItemResult
		for _, shopItem := range shopGroup.Items {
			shopInventories := inventoryMap[shopItem.ProductID]

			var available int
			for _, inv := range shopInventories {
				if inv.ShopID == shopGroup.ShopID {
					available = inv.Available()
					break
				}
			}

			if shopItem.Quantity > available {
				return nil, apperrors.NewConflict(inventoryDomain.ErrInsufficientStock.Error())
			}

			product, ok := productMap[shopItem.ProductID]
			if !ok {
				return nil, apperrors.NewNotFound(productDomain.ErrProductNotFound.Error())
			}

			itemSubtotal := product.Price * int64(shopItem.Quantity)
			totalSubtotal += itemSubtotal
			shopSubtotal = itemSubtotal

			weight := defaultShippingWeightGrams
			if product.Weight != nil {
				weight = int(*product.Weight)
			}
			totalWeight += weight * shopItem.Quantity

			checkoutItems = append(checkoutItems, CheckoutItemResult{
				ProductID:   product.ID,
				ShopID:      shopGroup.ShopID,
				Name:        product.Name,
				Price:       product.Price,
				Quantity:    shopItem.Quantity,
				Subtotal:    itemSubtotal,
				TotalWeight: weight * shopItem.Quantity,
			})
		}

		// Shipping estimation requires at least one active courier
		// Reject checkout when the shop cannot ship orders
		couriers := courierShopMap[shopGroup.ShopID]
		if couriers == nil {
			return nil, apperrors.NewConflict(domain.ErrShopCouriersNotFound.Error())
		}

		// Calculate shipping for the selected courier,
		// or all available couriers when none is selected yet
		var codes []string
		if shopGroup.Courier != nil {
			codes = append(codes, shopGroup.Courier.Code)
		} else {
			for _, courier := range couriers {
				codes = append(codes, courier.Code)
			}
		}

		originAddr := originAddresses[shopGroup.ShopID]
		originDistrictID, err := strconv.Atoi(originAddr.Detail.DistrictID)
		if err != nil {
			return nil, fmt.Errorf("invalid origin district id: %w", err)
		}

		costOptions, err := u.shipmentRepo.
			CalculateCost(
				ctx,
				u.executor,
				shipmentRepo.CalculateCostInput{
					OriginID:      originDistrictID,
					DestinationID: destDistrictID,
					Weight:        totalWeight,
					Couriers:      codes,
				},
			)
		if err != nil {
			return nil, fmt.Errorf("failed to calculate shipping cost: %w", err)
		}

		var (
			leastFee         int64
			selectedFee      int64
			leastFeeCode     string
			leastFeeService  string
			hasSelected      bool
			checkoutCouriers = make([]CheckoutCouriersResult, 0, len(costOptions))
		)
		for i, option := range costOptions {
			checkoutCouriers = append(checkoutCouriers, CheckoutCouriersResult{
				Code:    option.Code,
				Service: option.Service,
				Name:    option.Name,
				ETD:     option.Etd,
				Fee:     option.Cost,
			})

			if i == 0 || option.Cost < leastFee {
				leastFee = option.Cost
				leastFeeCode = option.Code
				leastFeeService = option.Service
			}

			if shopGroup.Courier != nil &&
				option.Code == shopGroup.Courier.Code &&
				option.Service == shopGroup.Courier.Service {

				// Capture the fee of the customer-selected courier
				// so it can be used in the final checkout total.
				selectedFee = option.Cost

				hasSelected = true
			}
		}

		// Use the selected courier fee when specified
		// Otherwise, fall back to the cheapest available option
		shippingFee := leastFee
		if shopGroup.Courier != nil {
			if !hasSelected {
				return nil, apperrors.NewBadRequest("selected courier service is unavailable")
			}

			shippingFee = selectedFee
		}

		totalShippingFee += shippingFee

		shopResult := ShopResult{
			ShopID:   shopGroup.ShopID,
			ShopSlug: shopMap[shopGroup.ShopID].Slug,
			ShopName: shopMap[shopGroup.ShopID].Name,
			Items:    checkoutItems,
			Subtotal: shopSubtotal,
			Total:    shopSubtotal + shippingFee,
		}

		if shopGroup.Courier == nil {
			// No courier has been selected yet
			// Return available shipping options for customer selection
			shopResult.CostCouriers = checkoutCouriers
			shopResult.SelectedCourier = &SelectedCourierResult{
				Code:    leastFeeCode,
				Service: leastFeeService,
				Fee:     shippingFee,
			}
		} else {
			// Persist the selected courier and its resolved shipping fee
			shopResult.SelectedCourier = &SelectedCourierResult{
				Code:    shopGroup.Courier.Code,
				Service: shopGroup.Courier.Service,
				Fee:     shippingFee,
			}
		}

		shopsResult = append(shopsResult, shopResult)
	}

	address := CheckoutAddressResult{
		ID:            destAddress.ID,
		RecipientName: destAddress.ReceiverName,
		Phone:         destAddress.Phone,
		FullAddress:   destAddress.Detail.FullAddress,
	}

	result := CheckoutResult{
		Address:          address,
		Shops:            shopsResult,
		Subtotal:         totalSubtotal,
		TotalShippingFee: totalShippingFee,
		TotalAll:         totalSubtotal + totalShippingFee,
	}

	return &result, nil
}
