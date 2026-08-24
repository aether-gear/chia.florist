package domain_test

import (
	"encoding/json"
	"testing"

	"service-core/internal/modules/product/domain"
)

func TestCalculateCustomProductPrice_V3(t *testing.T) {
	rawJSON := []byte(`{
		"metadata": {
			"version": "3.0.0",
			"editorVersion": "3.0.0",
			"platform": "web",
			"locale": "id-ID",
			"featureFlags": { "standaloneModule": true }
		},
		"layout": {
			"physicalSizeId": "medium",
			"aspectRatioId": "portrait-3-4",
			"upperHeightRatio": 0.58,
			"border": { "style": "ornate", "colorHex": "#F5C842", "widthPx": 12, "showCenterDivider": true }
		},
		"sections": {
			"upper": {
				"bgColorHex": "#C0392B",
				"cornerStyle": "rounded",
				"opacityPercent": 100,
				"header": { "text": "Selamat & Sukses", "fontId": "playfair", "fontSizePx": 36, "fontColorHex": "#FFD700", "alignment": "center" },
				"body": { "text": "Jane Doe", "fontId": "inter", "fontSizePx": 20, "fontColorHex": "#FFFFFF", "alignment": "center" }
			},
			"lower": {
				"bgColorHex": "#1A3A5C",
				"cornerStyle": "none",
				"opacityPercent": 100,
				"header": { "text": null, "fontId": "bebas", "fontSizePx": 26, "fontColorHex": "#FFFFFF", "alignment": "center" },
				"body": { "text": "PT. Tech Nusantara", "fontId": "inter", "fontSizePx": 22, "fontColorHex": "#FFFFFF", "alignment": "center" }
			}
		},
		"decorations": {
			"topCrest": { "visible": true, "variantId": "grand", "primaryColorHex": "#E63946", "secondaryColorHex": "#F1FAEE", "scalePercent": 40 },
			"bottomCrest": { "visible": true, "variantId": "modern", "primaryColorHex": "#10B981", "secondaryColorHex": "#ECFDF5", "scalePercent": 40 },
			"watermark": { "enabled": false, "text": "Chia Florist", "opacityPercent": 20 }
		},
		"elements": [
			{ "id": "b1", "type": "brush", "brushType": "flower", "colorHex": "#E85D75", "transform": { "xPercent": 50, "yPercent": 50, "scalePercent": 48, "rotationDeg": 0, "zIndex": 10 } },
			{ "id": "b2", "type": "brush", "brushType": "rose", "colorHex": "#90BE6D", "transform": { "xPercent": 20, "yPercent": 20, "scalePercent": 30, "rotationDeg": 15, "zIndex": 11 } },
			{ "id": "img1", "type": "image", "src": "https://example.com/logo.png", "frameStyle": "square", "transform": { "xPercent": 15, "yPercent": 15, "scalePercent": 22, "rotationDeg": 0, "zIndex": 12 } }
		],
		"assets": {
			"previewUrl": "https://example.com/preview.png"
		}
	}`)

	payload, err := domain.ParseCustomDesignPayload(json.RawMessage(rawJSON))
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}

	if payload.Metadata.Version != "3.0.0" {
		t.Errorf("expected version 3.0.0, got %s", payload.Metadata.Version)
	}
	if payload.Layout.AspectRatioID != "portrait-3-4" {
		t.Errorf("expected aspectRatioId portrait-3-4, got %s", payload.Layout.AspectRatioID)
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
	if breakdown.TopCrestFee != 45000 {
		t.Errorf("expected TopCrestFee 45000 (grand), got %d", breakdown.TopCrestFee)
	}
	if breakdown.BottomCrestFee != 30000 {
		t.Errorf("expected BottomCrestFee 30000 (modern), got %d", breakdown.BottomCrestFee)
	}
	if breakdown.CrestTotalFee != 75000 {
		t.Errorf("expected CrestTotalFee 75000, got %d", breakdown.CrestTotalFee)
	}
	if breakdown.TotalPrice <= 0 {
		t.Errorf("expected TotalPrice > 0, got %d", breakdown.TotalPrice)
	}

	// Test ExtractDesignSummary
	sizeID, prevURL, hUpper, bUpper, hLower, bLower := domain.ExtractDesignSummary(*payload)
	if sizeID != "medium" {
		t.Errorf("expected size medium, got %s", sizeID)
	}
	if prevURL == nil || *prevURL != "https://example.com/preview.png" {
		t.Errorf("expected previewUrl https://example.com/preview.png, got %v", prevURL)
	}
	if hUpper == nil || *hUpper != "Selamat & Sukses" {
		t.Errorf("expected header upper 'Selamat & Sukses', got %v", hUpper)
	}
	if bUpper == nil || *bUpper != "Jane Doe" {
		t.Errorf("expected body upper 'Jane Doe', got %v", bUpper)
	}
	if hLower != nil {
		t.Errorf("expected header lower nil, got %v", hLower)
	}
	if bLower == nil || *bLower != "PT. Tech Nusantara" {
		t.Errorf("expected body lower 'PT. Tech Nusantara', got %v", bLower)
	}
}

func TestCalculateCustomProductPrice_FlatLegacy(t *testing.T) {
	flatJSON := []byte(`{
		"upper": {
			"bgColor": "#C0392B",
			"headerText": "Happy Wedding",
			"headerColor": "#FFD700",
			"bodyText": "Romeo & Juliet",
			"bodyColor": "#FFFFFF"
		},
		"lower": {
			"bgColor": "#1A3A5C",
			"bodyText": "From Friends",
			"bodyColor": "#FFFFFF"
		},
		"border": {
			"style": "solid",
			"color": "#F5C842",
			"width": 8
		},
		"physicalSizeId": "large",
		"elements": [
			{ "id": "b1", "type": "brush", "brushType": "flower", "colorHex": "#E85D75" }
		]
	}`)

	payload, err := domain.ParseCustomDesignPayload(json.RawMessage(flatJSON))
	if err != nil {
		t.Fatalf("unexpected parse error on flat payload: %v", err)
	}

	matrix := domain.DefaultCustomPricingMatrix()
	breakdown := domain.CalculateCustomProductPrice(*payload, matrix)

	if breakdown.BaseSizePrice != 280000 {
		t.Errorf("expected Large BaseSizePrice 280000, got %d", breakdown.BaseSizePrice)
	}
	if breakdown.BrushCount != 1 {
		t.Errorf("expected BrushCount 1, got %d", breakdown.BrushCount)
	}
	if breakdown.BorderFee != 0 {
		t.Errorf("expected standard solid border fee 0, got %d", breakdown.BorderFee)
	}
}

