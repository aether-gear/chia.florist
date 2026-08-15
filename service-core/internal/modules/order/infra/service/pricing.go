package service

import (
	"context"
	"fmt"
	"strconv"

	apperrors "service-core/internal/common/errors"
	shipping "service-core/internal/infra/shipping"
	addressDomain "service-core/internal/modules/address/domain"
	addressRepo "service-core/internal/modules/address/repository"
	cartDomain "service-core/internal/modules/cart/domain"
	cartRepo "service-core/internal/modules/cart/repository"
	courierRepo "service-core/internal/modules/courier/repository"
	inventoryDomain "service-core/internal/modules/inventory/domain"
	inventoryRepo "service-core/internal/modules/inventory/repository"
	orderRepo "service-core/internal/modules/order/repository"
	paymentRepo "service-core/internal/modules/payment/repository"
	productDomain "service-core/internal/modules/product/domain"
	productRepo "service-core/internal/modules/product/repository"
	shopDomain "service-core/internal/modules/shop/domain"
	shopRepo "service-core/internal/modules/shop/repository"
	transaction "service-core/internal/shared/transaction"

	"github.com/google/uuid"
)

type pricingServiceImpl struct {
	addressRepo       addressRepo.CustomerAddressRepository
	cartRepo          cartRepo.CartRepository
	courierShopRepo   courierRepo.ShopCourierRepository
	inventoryRepo     inventoryRepo.InventoryRepository
	paymentMethodRepo paymentRepo.PaymentMethodRepository
	productRepo       productRepo.ProductRepository
	shippingProvider  shipping.ShippingProvider
	shopAddressRepo   addressRepo.ShopAddressRepository
	shopRepo          shopRepo.ShopRepository
}

func NewPricingService(
	addressRepo addressRepo.CustomerAddressRepository,
	cartRepo cartRepo.CartRepository,
	courierShopRepo courierRepo.ShopCourierRepository,
	inventoryRepo inventoryRepo.InventoryRepository,
	paymentMethodRepo paymentRepo.PaymentMethodRepository,
	productRepo productRepo.ProductRepository,
	shippingProvider shipping.ShippingProvider,
	shopAddressRepo addressRepo.ShopAddressRepository,
	shopRepo shopRepo.ShopRepository,
) orderRepo.PricingService {
	return &pricingServiceImpl{
		addressRepo:       addressRepo,
		cartRepo:          cartRepo,
		courierShopRepo:   courierShopRepo,
		inventoryRepo:     inventoryRepo,
		paymentMethodRepo: paymentMethodRepo,
		productRepo:       productRepo,
		shippingProvider:  shippingProvider,
		shopAddressRepo:   shopAddressRepo,
		shopRepo:          shopRepo,
	}
}

const defaultShippingWeightGrams = 1000

func (s *pricingServiceImpl) Calculate(
	ctx context.Context,
	exec transaction.Executor,
	input orderRepo.PricingInput,
) (*orderRepo.PricingResult, error) {
	// Checkout shipping destination defaults to the user's
	// primary address unless a specific address is requested
	var destAddress *addressDomain.CustomerAddress
	if input.AddressID != nil {
		addr, err := s.addressRepo.GetByID(ctx, exec,
			*input.AddressID,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve destination address: %w", err)
		}
		if addr == nil {
			return nil, apperrors.NewNotFound(addressDomain.ErrAddressNotFound.Error())
		}

		destAddress = addr
	} else {
		addr, err := s.addressRepo.GetDefaultByCustomerID(ctx, exec,
			input.CustomerID,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve default destination address: %w", err)
		}
		if addr == nil {
			return nil, apperrors.NewConflict(addressDomain.ErrNotFoundDefaultAddress.Error())
		}

		destAddress = addr
	}

	destDistrictID, err := strconv.Atoi(destAddress.Detail.DistrictID)
	if err != nil {
		return nil, fmt.Errorf("invalid destination district id: %w", err)
	}

	// Hydrate stored cart item metadata (e.g. custom_design, product_id, is_custom)
	// from the user's cart in PostgreSQL when cart_item_id is provided.
	var cartItemMap map[uuid.UUID]cartDomain.CartItem
	if input.CustomerID != uuid.Nil && s.cartRepo != nil {
		if userCart, err := s.cartRepo.GetWithItemsByCustomerID(ctx, exec, input.CustomerID); err == nil && userCart != nil {
			cartItemMap = make(map[uuid.UUID]cartDomain.CartItem, len(userCart.Items))
			for _, ci := range userCart.Items {
				cartItemMap[ci.ID] = ci
			}
		}
	}

	if cartItemMap != nil {
		for i := range input.Shops {
			for j := range input.Shops[i].Items {
				item := &input.Shops[i].Items[j]
				if item.CartItemID != nil {
					if ci, ok := cartItemMap[*item.CartItemID]; ok {
						if ci.ProductVariantType == cartDomain.ProductVariantTypeCustom || ci.ProductID == nil {
							item.IsCustom = true
							if len(item.CustomDesign) == 0 {
								item.CustomDesign = ci.CustomDesign
							}
						} else if item.ProductID == nil {
							item.ProductID = ci.ProductID
						}
					}
				}
			}
		}
	}

	// Collect IDs upfront to avoid per-item database queries
	//
	// Products are loaded for inventory checks,
	// and shops are loaded for origin address
	// and courier availability.
	var productIDs []uuid.UUID
	var shopIDs []uuid.UUID
	for _, shopGroup := range input.Shops {
		for _, item := range shopGroup.Items {
			if !item.IsCustom && item.ProductID != nil {
				productIDs = append(productIDs, *item.ProductID)
			}
		}
		shopIDs = append(shopIDs, shopGroup.ShopID)
	}

	var products []productDomain.Product
	if len(productIDs) > 0 {
		var err error
		products, err = s.productRepo.FindByIDs(ctx, exec, productIDs)
		if err != nil {
			return nil, fmt.Errorf("failed to load products: %w", err)
		}
	}

	productMap := make(map[uuid.UUID]productDomain.Product, len(products))
	for _, p := range products {
		productMap[p.ID] = p
	}

	var inventoryMap map[uuid.UUID][]inventoryDomain.Inventory
	if len(productIDs) > 0 {
		var err error
		inventoryMap, err = s.inventoryRepo.ListByProductIDs(ctx, exec, productIDs)
		if err != nil {
			return nil, fmt.Errorf("failed to load inventory for products: %w", err)
		}
	}

	courierShopMap, err := s.courierShopRepo.ListsByShopIDs(ctx, exec,
		shopIDs,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve shop courier configurations: %w", err)
	}

	originAddresses, err := s.shopAddressRepo.GetDefaultsByShopIDs(ctx, exec,
		shopIDs,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve origin addresses: %w", err)
	}

	shops, err := s.shopRepo.FindByIDs(ctx, exec,
		shopIDs,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load shops: %w", err)
	}

	shopMap := make(map[uuid.UUID]shopDomain.Shop, len(shops))
	for _, sh := range shops {
		shopMap[sh.ID] = sh
	}

	for _, shopID := range shopIDs {
		sh, ok := shopMap[shopID]
		if !ok || !sh.IsOperable() {
			shopName := "Unknown Shop"
			if ok {
				shopName = sh.Name
			}
			return nil, apperrors.NewConflict(fmt.Sprintf("shop '%s' is inactive or not approved for transactions", shopName))
		}
	}

	// Calculate estimation, subtotal and weight
	// based on shop item (giant loop)
	var (
		totalSubtotal    int64
		totalShippingFee int64
		shopsResult      []orderRepo.PricingShopResult
	)

	for _, shopGroup := range input.Shops {
		// Map shop items, so products consist of
		// having inventory and are owned by a shop ID
		//
		// Include calculations of weight, total subtotal
		// and estimation total
		var (
			shopSubtotal    int64
			shopItemsWeight int
			pricingItems    []orderRepo.PricingItemResult
		)

		for _, shopItem := range shopGroup.Items {
			if shopItem.IsCustom || shopItem.ProductID == nil {
				var unitPrice int64 = productDomain.DEFAULT_SIZE_PRICE_MEDIUM
				weight := 2500 // default medium board weight

				if len(shopItem.CustomDesign) > 0 {
					if parsedDesign, err := productDomain.ParseCustomDesignPayload(shopItem.CustomDesign); err == nil {
						breakdown := productDomain.CalculateCustomProductPrice(
							*parsedDesign,
							productDomain.DefaultCustomPricingMatrix(),
						)

						unitPrice = breakdown.TotalPrice
						if breakdown.PhysicalSizeID == "small" {
							weight = 1500
						} else if breakdown.PhysicalSizeID == "large" {
							weight = 4000
						}
					}
				}

				itemSubtotal := unitPrice * int64(shopItem.Quantity)
				shopSubtotal += itemSubtotal
				totalSubtotal += itemSubtotal
				itemWeight := weight * shopItem.Quantity
				shopItemsWeight += itemWeight

				pricingItems = append(pricingItems, orderRepo.PricingItemResult{
					ProductID:   shopItem.ProductID,
					CartItemID:  shopItem.CartItemID,
					IsCustom:    true,
					ProductName: "(Custom Flower Board)",
					Quantity:    shopItem.Quantity,
					UnitPrice:   unitPrice,
					Subtotal:    itemSubtotal,
					WeightGrams: itemWeight,
				})
				continue
			}

			pid := *shopItem.ProductID
			product, ok := productMap[pid]
			if !ok || product.Status == productDomain.ProductStatusArchived {
				return nil, apperrors.NewNotFound(productDomain.ErrProductNotFound.Error())
			}
			if product.Status != productDomain.ProductStatusActive {
				return nil, apperrors.NewConflict(fmt.Sprintf("product '%s' is currently not available for purchase", product.Name))
			}

			shopInventories := inventoryMap[pid]
			var available int
			for _, inv := range shopInventories {
				if inv.ShopID == shopGroup.ShopID {
					available = inv.Available()
					break
				}
			}

			if shopItem.Quantity > available {
				return nil, apperrors.NewConflict(fmt.Sprintf("%s is out of stock", product.Name))
			}


			itemSubtotal := product.Price * int64(shopItem.Quantity)
			shopSubtotal += itemSubtotal
			totalSubtotal += itemSubtotal

			weight := defaultShippingWeightGrams
			if product.Weight != nil {
				weight = int(*product.Weight)
			}

			itemWeight := weight * shopItem.Quantity
			shopItemsWeight += itemWeight

			pID := product.ID
			pricingItems = append(pricingItems, orderRepo.PricingItemResult{
				ProductID:   &pID,
				CartItemID:  shopItem.CartItemID,
				IsCustom:    false,
				ProductName: product.Name,
				Quantity:    shopItem.Quantity,
				UnitPrice:   product.Price,
				Subtotal:    itemSubtotal,
				WeightGrams: itemWeight,
			})
		}

		// Shipping estimation requires
		// at least one active courier
		//
		// Reject checkout when the shop
		// cannot ship orders
		couriers := courierShopMap[shopGroup.ShopID]
		if couriers == nil {
			return nil, apperrors.NewConflict(cartDomain.ErrShopCouriersNotFound.Error())
		}

		// Calculate shipping for the selected courier,
		// or all available couriers when none is selected yet
		var codes []string
		if shopGroup.CourierCode != nil && *shopGroup.CourierCode != "" {
			codes = append(codes, *shopGroup.CourierCode)
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

		costOptions, err := s.shippingProvider.CalculateRates(
			ctx,
			shipping.CalculateRatesInput{
				OriginID:      originDistrictID,
				DestinationID: destDistrictID,
				Weight:        shopItemsWeight,
				Couriers:      codes,
			},
		)
		if err != nil {
			return nil, fmt.Errorf("failed to calculate shipping cost: %w", err)
		}

		var (
			leastFee        int64
			selectedFee     int64
			leastFeeCode    string
			leastFeeService string
			hasSelected     bool
			courierOptions  = make([]orderRepo.CourierOption, 0, len(costOptions))
		)

		for i, option := range costOptions {
			courierOptions = append(courierOptions, orderRepo.CourierOption{
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

			if shopGroup.CourierCode != nil &&
				shopGroup.CourierService != nil &&
				option.Code == *shopGroup.CourierCode &&
				option.Service == *shopGroup.CourierService {

				// Capture the fee of the customer-selected courier
				// so it can be used in the final checkout total
				selectedFee = option.Cost
				hasSelected = true
			}
		}

		// Use the selected courier fee when specified
		//
		// Otherwise, fall back to the cheapest available option
		shippingFee := leastFee
		resolvedCourierCode := leastFeeCode
		resolvedCourierService := leastFeeService

		if shopGroup.CourierCode != nil &&
			shopGroup.CourierService != nil &&
			*shopGroup.CourierCode != "" {

			if !hasSelected {
				return nil, apperrors.NewBadRequest("selected courier service is unavailable")
			}

			shippingFee = selectedFee
			resolvedCourierCode = *shopGroup.CourierCode
			resolvedCourierService = *shopGroup.CourierService
		}

		totalShippingFee += shippingFee

		shopDetails := shopMap[shopGroup.ShopID]
		shopResult := orderRepo.PricingShopResult{
			ShopID:   shopGroup.ShopID,
			ShopName: shopDetails.Name,
			ShopSlug: shopDetails.Slug,
			Items:    pricingItems,
			SelectedCourier: orderRepo.SelectedCourierResult{
				Code:    resolvedCourierCode,
				Service: resolvedCourierService,
				Fee:     shippingFee,
			},
			Subtotal: shopSubtotal,
			Total:    shopSubtotal + shippingFee,
		}

		// No courier has been selected yet
		//
		// Return available shipping options for customer selection
		if shopGroup.CourierCode == nil ||
			*shopGroup.CourierCode == "" {
			shopResult.CourierOptions = courierOptions
		}

		shopsResult = append(shopsResult, shopResult)
	}

	var (
		paymentMethods        []orderRepo.PaymentMethodPricingResult
		selectedPaymentMethod *orderRepo.PaymentMethodPricingResult
		totalAll              = totalSubtotal + totalShippingFee
	)

	if input.PaymentMethodID != nil {
		pm, err := s.paymentMethodRepo.GetByID(ctx, exec,
			*input.PaymentMethodID,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve payment method: %w", err)
		}
		if pm == nil {
			return nil, apperrors.NewNotFound("selected payment method is not available")
		}

		fee := pm.CalculateFee(totalAll)
		selectedPaymentMethod = &orderRepo.PaymentMethodPricingResult{
			PaymentMethodID: pm.ID,
			Name:            pm.Name,
			Type:            string(pm.Type),
			Description:     pm.Description,
			Fee:             fee,
			Subtotal:        totalAll,
			Total:           totalAll + fee,
		}
	} else {
		pms, err := s.paymentMethodRepo.ListAll(ctx, exec,
			nil,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to load payment methods: %w", err)
		}
		if len(pms) == 0 {
			return nil, apperrors.NewNotFound("payment method is not available")
		}

		paymentMethods = make([]orderRepo.PaymentMethodPricingResult, 0, len(pms))
		for _, pm := range pms {
			fee := pm.CalculateFee(totalAll)
			paymentMethods = append(paymentMethods, orderRepo.PaymentMethodPricingResult{
				PaymentMethodID: pm.ID,
				Name:            pm.Name,
				Type:            string(pm.Type),
				Description:     pm.Description,
				Fee:             fee,
				Subtotal:        totalAll,
				Total:           totalAll + fee,
			})
		}
	}

	// When a payment method is already selected,
	// use its fee directly
	//
	// Otherwise, use the cheapest available option
	// to estimate the grand total
	var leastFeePayMethod int64
	if selectedPaymentMethod != nil {
		leastFeePayMethod = selectedPaymentMethod.Fee
	} else {
		for i, method := range paymentMethods {
			if i == 0 || method.Fee < leastFeePayMethod {
				leastFeePayMethod = method.Fee
			}
		}
	}

	result := &orderRepo.PricingResult{
		Address: orderRepo.PricingAddressResult{
			ID:            destAddress.ID,
			RecipientName: destAddress.ReceiverName,
			Phone:         destAddress.Phone,
			FullAddress:   destAddress.Detail.FullAddress,
		},
		Shops:                 shopsResult,
		Subtotal:              totalSubtotal,
		TotalShippingFee:      totalShippingFee,
		GrandTotal:            totalAll + leastFeePayMethod,
		PaymentMethods:        paymentMethods,
		SelectedPaymentMethod: selectedPaymentMethod,
	}

	return result, nil
}
