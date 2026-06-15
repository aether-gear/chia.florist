package usecase

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	apperrors "service-core/internal/common/errors"
	addressDomain "service-core/internal/modules/address/domain"
	addressRepo "service-core/internal/modules/address/repository"
	authenDomain "service-core/internal/modules/authentication/domain"
	courierRepo "service-core/internal/modules/courier/repository"
	inventoryDomain "service-core/internal/modules/inventory/domain"
	inventoryRepo "service-core/internal/modules/inventory/repository"
	productDomain "service-core/internal/modules/product/domain"
	productRepo "service-core/internal/modules/product/repository"
	shipmentRepo "service-core/internal/modules/shipment/repository"
	shopRepository "service-core/internal/modules/shop/repository"
	transaction "service-core/internal/shared/transaction"

	"github.com/google/uuid"
)

type CheckoutUsecase struct {
	exec            transaction.Executor
	addressRepo     addressRepo.UserAddressRepository
	courierRepo     courierRepo.CourierRepository
	inventoryRepo   inventoryRepo.InventoryRepository
	productRepo     productRepo.ProductRepository
	shipmentRepo    shipmentRepo.ShippingCostProvider
	shopAddressRepo addressRepo.ShopAddressRepository
	shopRepo        shopRepository.ShopRepository
}

func NewCheckoutUsecase(
	exec transaction.Executor,
	addressRepo addressRepo.UserAddressRepository,
	courierRepo courierRepo.CourierRepository,
	inventoryRepo inventoryRepo.InventoryRepository,
	productRepo productRepo.ProductRepository,
	shipmentRepo shipmentRepo.ShippingCostProvider,
	shopAddressRepo addressRepo.ShopAddressRepository,
	shopRepo shopRepository.ShopRepository,
) *CheckoutUsecase {
	return &CheckoutUsecase{
		exec:            exec,
		addressRepo:     addressRepo,
		courierRepo:     courierRepo,
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

type CheckoutShopInput struct {
	ShopID uuid.UUID
	Items  []CheckoutItemInput
}

type CheckoutInput struct {
	AddressID *uuid.UUID
	Couriers  []string
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

type ShopResult struct {
	Items       []CheckoutItemResult
	ShippingFee []CheckoutCouriersResult
}

type CheckoutResult struct {
	Address         CheckoutAddressResult
	Shops           []ShopResult
	Subtotal        int64
	TotalEstimation int64
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
		GetDefaultByUserID(ctx, u.exec, authCtx.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve destination address: %w", err)
	}
	if defAddress == nil {
		return nil, apperrors.NewConflict(addressDomain.ErrNotFoundDefaultAddress.Error())
	}

	// Fallback to the user's default address,
	// if no ID was specified
	var destAddress addressDomain.Address
	if input.AddressID != nil {
		address, err := u.addressRepo.
			GetByID(ctx, u.exec, *input.AddressID)
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

	// Initialize couriers for cost estimation
	// It will get skip if input is empty
	var couriers []string
	if len(input.Couriers) > 0 {
		couriers = append(couriers, input.Couriers...)
	}

	var productIDs []uuid.UUID
	for _, shopGroup := range input.ShopInput {
		for _, item := range shopGroup.Items {
			productIDs = append(productIDs, item.ProductID)
		}
	}

	products, err := u.productRepo.
		FindByIDs(ctx, u.exec, productIDs)
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
		ListByProductIDs(ctx, u.exec, productIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to load inventory for cart products: %w", err)
	}

	// Calculate estimation, subtotal and weight
	// based on shop item (giant loop)
	var (
		totalSubtotal        int64
		shopsResult          []ShopResult
		totalSubtotalWithEst int64
		totalWeight          int = 0
	)
	for _, shopGroup := range input.ShopInput {
		// Map shop items, so products consist of
		// having inventory and are owned by a shop ID
		//
		// Include calculations of weight, total subtotal
		// and estimation total
		var shopItems []CheckoutItemResult
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
			totalSubtotalWithEst += itemSubtotal

			weight := defaultShippingWeightGrams
			if product.Weight != nil {
				weight = int(*product.Weight)
			}
			totalWeight += weight * shopItem.Quantity

			shopItems = append(shopItems, CheckoutItemResult{
				ProductID:   product.ID,
				ShopID:      shopGroup.ShopID,
				Name:        product.Name,
				Price:       product.Price,
				Quantity:    shopItem.Quantity,
				Subtotal:    itemSubtotal,
				TotalWeight: weight * shopItem.Quantity,
			})
		}

		// No default courier
		// Assuming no estimation cost
		if len(couriers) != 0 {
			originAddr, err := u.shopAddressRepo.
				GetDefaultByShopID(ctx, u.exec, shopGroup.ShopID)
			if err != nil {
				return nil, fmt.Errorf("failed to retrieve origin address: %w", err)
			}
			if originAddr == nil {
				return nil, apperrors.NewNotFound(errors.New("origin address not found").Error())
			}

			originDistrictID, err := strconv.Atoi(originAddr.Detail.DistrictID)
			if err != nil {
				return nil, fmt.Errorf("invalid origin district id: %w", err)
			}

			costOptions, err := u.shipmentRepo.
				CalculateCost(
					ctx,
					u.exec,
					shipmentRepo.CalculateCostInput{
						OriginID:      originDistrictID,
						DestinationID: destDistrictID,
						Weight:        totalWeight,
						Couriers:      couriers,
					},
				)
			if err != nil {
				return nil, fmt.Errorf("failed to calculate shipping cost: %w", err)
			}

			var leastFee int64
			var checkoutCouriers []CheckoutCouriersResult
			for i, option := range costOptions {
				checkoutCouriers = append(checkoutCouriers, CheckoutCouriersResult{
					Code:    option.Code,
					Service: option.Service,
					ETD:     option.Etd,
					Fee:     option.Cost,
				})

				if i == 0 || option.Cost < leastFee {
					leastFee = option.Cost
				}
			}

			shopsResult = append(shopsResult, ShopResult{
				Items:       shopItems,
				ShippingFee: checkoutCouriers,
			})

			// Calculate based on how many estimation per shop.
			// Even then, only take the least price
			if len(costOptions) > 0 {
				totalSubtotalWithEst += leastFee
			}
		} else {
			shopsResult = append(shopsResult, ShopResult{
				Items: shopItems,
			})
		}
	}

	address := CheckoutAddressResult{
		ID:            destAddress.ID,
		RecipientName: destAddress.ReceiverName,
		Phone:         destAddress.Phone,
		FullAddress:   destAddress.Detail.FullAddress,
	}

	result := CheckoutResult{
		Address:  address,
		Shops:    shopsResult,
		Subtotal: totalSubtotal,
	}
	if len(couriers) != 0 {
		result.TotalEstimation = totalSubtotalWithEst
	}

	return &result, nil
}
