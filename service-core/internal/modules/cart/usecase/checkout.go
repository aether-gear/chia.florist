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
	ProductID uuid.UUID
	ShopID    uuid.UUID
	Name      string
	Price     int64
	Quantity  int
	Subtotal  int64
}

type ShopResult struct {
	Items       []CheckoutItemResult
	ShippingFee []CheckoutCouriersResult
}

type CheckoutResult struct {
	Address  CheckoutAddressResult
	Shops    []ShopResult
	Subtotal int64
	Total    int64
}

const defaultShippingWeightGrams = 1000

func (u *CheckoutUsecase) Execute(
	ctx context.Context,
	authCtx authenDomain.AuthContext,
	input CheckoutInput,
) (*CheckoutResult, error) {
	// Early fallback.
	// Make sure default address is always exist - Deuterrr
	var destAddress addressDomain.Address
	defAddress, err := u.addressRepo.
		GetDefaultByUserID(ctx, u.exec, authCtx.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve destination address: %w", err)
	}
	if defAddress == nil {
		return nil, apperrors.NewConflict(addressDomain.ErrNotFoundDefaultAddress.Error())
	}

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

	couriers, err := u.courierRepo.ListAll(ctx, u.exec)
	if err != nil {
		return nil, fmt.Errorf("failed to load couriers: %w", err)
	}
	if len(couriers) == 0 {
		return nil, apperrors.NewInternal(errors.New("no courier service available at the moment"))
	}

	destDistrictID, err := strconv.Atoi(destAddress.Detail.DistrictID)
	if err != nil {
		return nil, fmt.Errorf("invalid destination district id: %w", err)
	}

	var shopsResult []ShopResult
	var totalSubtotal int64
	for _, shopGroup := range input.ShopInput {
		var shopItems []CheckoutItemResult
		for _, shopItem := range shopGroup.Items {
			product, ok := productMap[shopItem.ProductID]
			if !ok {
				return nil, apperrors.NewNotFound(productDomain.ErrProductNotFound.Error())
			}

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

			itemSubtotal := product.Price * int64(shopItem.Quantity)
			shopItems = append(shopItems, CheckoutItemResult{
				ProductID: product.ID,
				ShopID:    shopGroup.ShopID,
				Name:      product.Name,
				Price:     product.Price,
				Quantity:  shopItem.Quantity,
				Subtotal:  itemSubtotal,
			})

			totalSubtotal += itemSubtotal
		}

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

		totalWeight := 0
		for _, item := range shopGroup.Items {
			product := productMap[item.ProductID]

			weight := defaultShippingWeightGrams
			if product.Weight != nil {
				weight = int(*product.Weight)
			}

			totalWeight += weight * item.Quantity
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

		var checkoutCouriers []CheckoutCouriersResult
		for _, option := range costOptions {
			checkoutCouriers = append(checkoutCouriers, CheckoutCouriersResult{
				Code:    option.Code,
				Service: option.Service,
				ETD:     option.Etd,
				Fee:     option.Cost,
			})
		}

		shopsResult = append(shopsResult, ShopResult{
			Items:       shopItems,
			ShippingFee: checkoutCouriers,
		})
	}

	address := CheckoutAddressResult{
		ID:            destAddress.ID,
		RecipientName: destAddress.ReceiverName,
		Phone:         destAddress.Phone,
		FullAddress:   destAddress.Detail.FullAddress,
	}

	return &CheckoutResult{
		Address:  address,
		Shops:    shopsResult,
		Subtotal: totalSubtotal,
	}, nil
}
