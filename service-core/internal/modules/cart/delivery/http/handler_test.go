package http

import (
	"encoding/json"
	"testing"

	"github.com/google/uuid"
)

func TestParseCheckoutInput_CustomCartItem(t *testing.T) {
	h := &CartHandler{}

	shopID := uuid.New().String()
	cartItemID := uuid.New().String()

	req := checkoutCalculateRequest{
		Shops: []checkoutShopRequest{
			{
				ShopID: shopID,
				Items: []checkoutItemRequest{
					{
						CartItemID:   &cartItemID,
						ProductID:    nil,
						Quantity:     1,
						CustomDesign: json.RawMessage(`{"layout":{"physicalSizeId":"medium"}}`),
					},
				},
			},
		},
	}

	input, err := h.parseCheckoutInput(req)
	if err != nil {
		t.Fatalf("expected parseCheckoutInput to succeed for custom cart item, got error: %v", err)
	}

	if len(input.ShopInput) != 1 {
		t.Fatalf("expected 1 shop input, got %d", len(input.ShopInput))
	}

	item := input.ShopInput[0].Items[0]
	if !item.IsCustom {
		t.Errorf("expected IsCustom to be true")
	}
	if item.CartItemID == nil || item.CartItemID.String() != cartItemID {
		t.Errorf("expected CartItemID %s, got %v", cartItemID, item.CartItemID)
	}
}

func TestParseCheckoutInput_CustomDirectItemWithoutProductID(t *testing.T) {
	h := &CartHandler{}

	shopID := uuid.New().String()
	customStr := "custom"

	req := checkoutCalculateRequest{
		Shops: []checkoutShopRequest{
			{
				ShopID: shopID,
				Items: []checkoutItemRequest{
					{
						ProductID:          &customStr,
						ProductVariantType: "custom",
						Quantity:           1,
						CustomDesign:       json.RawMessage(`{"layout":{"physicalSizeId":"medium"}}`),
					},
				},
			},
		},
	}

	input, err := h.parseCheckoutInput(req)
	if err != nil {
		t.Fatalf("expected parseCheckoutInput to succeed for custom direct item, got error: %v", err)
	}

	item := input.ShopInput[0].Items[0]
	if !item.IsCustom {
		t.Errorf("expected IsCustom to be true")
	}
	if item.ProductID != nil {
		t.Errorf("expected ProductID to be nil for custom product, got %v", item.ProductID)
	}
}
