package domain_test

import (
	"encoding/json"
	"testing"

	"service-core/internal/modules/product/domain"
)

func TestCalculateCustomProductPrice(t *testing.T) {
	rawJSON := []byte(`{
		"layout": {
			"physicalSizeId": "medium",
			"border": { "style": "ornate", "colorHex": "#F5C842", "widthPx": 12 }
		},
		"sections": {
			"upper": {
				"bgColorHex": "#C0392B",
				"header": { "fontColorHex": "#FFD700" },
				"body": { "fontColorHex": "#FFFFFF" }
			},
			"lower": {
				"bgColorHex": "#1A3A5C",
				"header": { "fontColorHex": "#FFFFFF" },
				"body": { "fontColorHex": "#FFFFFF" }
			}
		},
		"decorations": {
			"topCrest": { "visible": true, "variantId": "classic", "primaryColorHex": "#E63946", "secondaryColorHex": "#F1FAEE" },
			"bottomCrest": { "visible": false, "variantId": "classic", "primaryColorHex": "#E63946", "secondaryColorHex": "#F1FAEE" }
		},
		"elements": [
			{ "id": "b1", "type": "brush", "brushType": "flower", "colorHex": "#E85D75" },
			{ "id": "b2", "type": "brush", "brushType": "rose", "colorHex": "#90BE6D" },
			{ "id": "img1", "type": "image", "src": "https://example.com/logo.png" }
		]
	}`)

	payload, err := domain.ParseCustomDesignPayload(json.RawMessage(rawJSON))
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}

	matrix := domain.DefaultCustomPricingMatrix()
	breakdown := domain.CalculateCustomProductPrice(*payload, matrix)

	if breakdown.BaseSizePrice != 200000 {
		t.Errorf("expected BaseSizePrice 200000, got %d", breakdown.BaseSizePrice)
	}
	if breakdown.BrushCount != 2 {
		t.Errorf("expected BrushCount 2, got %d", breakdown.BrushCount)
	}
	if breakdown.BrushTotalFee != 4000 {
		t.Errorf("expected BrushTotalFee 4000, got %d", breakdown.BrushTotalFee)
	}
	if breakdown.ImageCount != 1 {
		t.Errorf("expected ImageCount 1, got %d", breakdown.ImageCount)
	}
	if breakdown.ImageTotalFee != 20000 {
		t.Errorf("expected ImageTotalFee 20000, got %d", breakdown.ImageTotalFee)
	}
	if breakdown.BorderFee != 15000 {
		t.Errorf("expected BorderFee 15000, got %d", breakdown.BorderFee)
	}
	if breakdown.TopCrestFee != 25000 {
		t.Errorf("expected TopCrestFee 25000, got %d", breakdown.TopCrestFee)
	}
	if breakdown.TotalPrice <= 0 {
		t.Errorf("expected TotalPrice > 0, got %d", breakdown.TotalPrice)
	}
}
