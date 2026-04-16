package usecase

import (
	cartD "service-core/internal/modules/cart/domain"
	cartR "service-core/internal/modules/cart/repository"
	productR "service-core/internal/modules/product/repository"

	"github.com/google/uuid"
)

type GetCartResult struct {
	Cart     *cartD.Cart
	Products map[uuid.UUID]productR.ProductWithInventory
}

type GetCartUsecase struct {
	cartRepo    cartR.CartRepository
	productRepo productR.ProductRepository
}

func NewGetCartUsecase(cR cartR.CartRepository, pR productR.ProductRepository) *GetCartUsecase {
	return &GetCartUsecase{
		cartRepo:    cR,
		productRepo: pR,
	}
}

func (u *GetCartUsecase) Execute(userID uuid.UUID) (*GetCartResult, error) {
	cart, err := u.cartRepo.GetWithItemsByUserID(userID)
	if err != nil {
		return nil, err
	}

	if cart == nil {
		cart, err = u.cartRepo.NewCart(userID)
		if err != nil {
			return nil, err
		}
	}

	if len(cart.Items) == 0 {
		return &GetCartResult{
			Cart:     cart,
			Products: map[uuid.UUID]productR.ProductWithInventory{},
		}, nil
	}

	productIDs := make([]uuid.UUID, 0, len(cart.Items))
	for _, item := range cart.Items {
		productIDs = append(productIDs, item.ProductID)
	}

	products, err := u.productRepo.FindByIDs(productIDs)
	if err != nil {
		return nil, err
	}

	productMap := make(map[uuid.UUID]productR.ProductWithInventory)
	for _, p := range products {
		productMap[p.Product.ID] = p
	}

	return &GetCartResult{
		Cart:     cart,
		Products: productMap,
	}, nil
}
