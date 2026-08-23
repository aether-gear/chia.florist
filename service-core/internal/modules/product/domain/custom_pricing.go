package domain

import (
	"encoding/json"
	"fmt"
	"math"
	"slices"
	"strings"
)

const (
	DEFAULT_PHYSICAL_SIZE_ID = "medium"

	DEFAULT_SIZE_PRICE_SMALL  int64 = 150000
	DEFAULT_SIZE_PRICE_MEDIUM int64 = 200000
	DEFAULT_SIZE_PRICE_LARGE  int64 = 280000

	DEFAULT_PRICE_PER_BRUSH_STROKE int64 = 2000
	DEFAULT_FREE_COLOR_THRESHOLD   int   = 3
	DEFAULT_PRICE_PER_EXTRA_COLOR  int64 = 10000
	DEFAULT_PREMIUM_BORDER_FEE     int64 = 15000

	DEFAULT_CREST_PRICE_CLASSIC int64 = 25000
	DEFAULT_CREST_PRICE_MODERN  int64 = 30000
	DEFAULT_CREST_PRICE_GRAND   int64 = 45000
	DEFAULT_CREST_FEE_FALLBACK  int64 = 25000

	DEFAULT_PRICE_PER_IMAGE_ELEMENT int64 = 20000

	DEFAULT_HEX_COLOR_FALLBACK = "#FFFFFF"
)

// CustomDesignPayload represents the standardized custom flower board
// design approach (supporting v3.0.0, v1.0.0, and normalized legacy payloads).
type CustomDesignPayload struct {
	Metadata    DesignMetadata    `json:"metadata"`
	Layout      DesignLayout      `json:"layout"`
	Sections    DesignSections    `json:"sections"`
	Decorations DesignDecorations `json:"decorations"`
	Elements    []DesignElement   `json:"elements"`
	Assets      DesignAssets      `json:"assets"`
}

type DesignMetadata struct {
	Version       string          `json:"version"`
	EditorVersion string          `json:"editorVersion"`
	Platform      string          `json:"platform"`
	Locale        string          `json:"locale"`
	CreatedAt     string          `json:"createdAt"`
	UpdatedAt     string          `json:"updatedAt"`
	Checksum      string          `json:"checksum"`
	FeatureFlags  map[string]bool `json:"featureFlags,omitempty"`
}

type DesignLayout struct {
	PhysicalSizeID   string     `json:"physicalSizeId"` // "small", "medium", "large"
	AspectRatioID    string     `json:"aspectRatioId,omitempty"`
	UpperHeightRatio float64    `json:"upperHeightRatio"`
	Border           BorderSpec `json:"border"`
}

type BorderSpec struct {
	Style             string `json:"style"`    // "none", "solid", "double", "dashed", "dotted", "groove", "ridge", "ornate"
	ColorHex          string `json:"colorHex"` // 6-digit hex string (#RRGGBB)
	WidthPx           int    `json:"widthPx"`
	ShowCenterDivider bool   `json:"showCenterDivider"`
}

type DesignSections struct {
	Upper SectionSpec `json:"upper"`
	Lower SectionSpec `json:"lower"`
}

type SectionSpec struct {
	BGColorHex        string         `json:"bgColorHex"`
	CornerStyle       string         `json:"cornerStyle"` // "none", "rounded", "cut", "ornate", "floral"
	OpacityPercent    *int           `json:"opacityPercent,omitempty"`
	Header            TypographySpec `json:"header"`
	Body              TypographySpec `json:"body"`
	HeaderBorder      *bool          `json:"headerBorder,omitempty"`
	HeaderBorderColor *string        `json:"headerBorderColor,omitempty"`
	HeaderBorderWidth *int           `json:"headerBorderWidth,omitempty"`
	BodyBorder        *bool          `json:"bodyBorder,omitempty"`
	BodyBorderColor   *string        `json:"bodyBorderColor,omitempty"`
	BodyBorderWidth   *int           `json:"bodyBorderWidth,omitempty"`
}

type TypographySpec struct {
	Text         *string `json:"text"`
	FontID       string  `json:"fontId"` // "inter", "playfair", "dancing", "bebas", "merriweather", "pacifico"
	FontSizePx   int     `json:"fontSizePx"`
	FontColorHex string  `json:"fontColorHex"`
	Alignment    string  `json:"alignment"` // "left", "center", "right"
}

type DesignDecorations struct {
	TopCrest    CrestSpec      `json:"topCrest"`
	BottomCrest CrestSpec      `json:"bottomCrest"`
	Watermark   *WatermarkSpec `json:"watermark,omitempty"`
}

type CrestSpec struct {
	Visible           bool   `json:"visible"`
	VariantID         string `json:"variantId"` // "classic", "modern", "grand"
	PrimaryColorHex   string `json:"primaryColorHex"`
	SecondaryColorHex string `json:"secondaryColorHex"`
	ScalePercent      int    `json:"scalePercent"`
}

type WatermarkSpec struct {
	Enabled        bool   `json:"enabled"`
	Text           string `json:"text"`
	OpacityPercent int    `json:"opacityPercent"`
}

type DesignElement struct {
	ID         string           `json:"id"`
	Type       string           `json:"type"`                 // "image" or "brush"
	Src        string           `json:"src,omitempty"`        // Image URL / Data URI
	FrameStyle string           `json:"frameStyle,omitempty"` // "none", "square", "circle"
	BrushType  string           `json:"brushType,omitempty"`  // "flower", "rose"
	ColorHex   string           `json:"colorHex,omitempty"`
	Crop       *ElementCrop     `json:"crop,omitempty"`
	Transform  ElementTransform `json:"transform"`
}

type ElementCrop struct {
	XPercent float64 `json:"xPercent"`
	YPercent float64 `json:"yPercent"`
	Zoom     float64 `json:"zoom"`
}

type ElementTransform struct {
	XPercent     float64 `json:"xPercent"`
	YPercent     float64 `json:"yPercent"`
	ScalePercent float64 `json:"scalePercent"`
	RotationDeg  float64 `json:"rotationDeg"`
	ZIndex       *int    `json:"zIndex,omitempty"`
}

type DesignAssets struct {
	PreviewBase64   *string `json:"previewBase64,omitempty"`
	PreviewAssetID  *string `json:"previewAssetId,omitempty"`
	PreviewURL      *string `json:"previewUrl,omitempty"`
	BucketPath      *string `json:"bucketPath,omitempty"`
	StorageProvider *string `json:"storageProvider,omitempty"`
}

// CustomPricingMatrix contains configurable pricing rules
// used by the pricing engine.
type CustomPricingMatrix struct {
	SizePrices           map[string]int64
	PricePerBrushStroke  int64
	FreeColorThreshold   int
	PricePerExtraColor   int64
	PremiumBorderStyles  []string
	PremiumBorderFee     int64
	CrestVariantPrices   map[string]int64
	PricePerImageElement int64
}

// DefaultCustomPricingMatrix returns the default pricing configuration.
func DefaultCustomPricingMatrix() CustomPricingMatrix {
	return CustomPricingMatrix{
		SizePrices: map[string]int64{
			"small":  DEFAULT_SIZE_PRICE_SMALL,
			"medium": DEFAULT_SIZE_PRICE_MEDIUM,
			"large":  DEFAULT_SIZE_PRICE_LARGE,
		},
		PricePerBrushStroke: DEFAULT_PRICE_PER_BRUSH_STROKE,
		FreeColorThreshold:  DEFAULT_FREE_COLOR_THRESHOLD,
		PricePerExtraColor:  DEFAULT_PRICE_PER_EXTRA_COLOR,
		PremiumBorderStyles: []string{"double", "groove", "ridge", "ornate"},
		PremiumBorderFee:    DEFAULT_PREMIUM_BORDER_FEE,
		CrestVariantPrices: map[string]int64{
			"classic": DEFAULT_CREST_PRICE_CLASSIC,
			"modern":  DEFAULT_CREST_PRICE_MODERN,
			"grand":   DEFAULT_CREST_PRICE_GRAND,
		},
		PricePerImageElement: DEFAULT_PRICE_PER_IMAGE_ELEMENT,
	}
}

// CustomPricingBreakdown contains the detailed fee
// calculation for a custom design.
type CustomPricingBreakdown struct {
	PhysicalSizeID    string `json:"physical_size_id"`
	BaseSizePrice     int64  `json:"base_size_price"`
	BrushCount        int    `json:"brush_count"`
	BrushTotalFee     int64  `json:"brush_total_fee"`
	UniqueColorsCount int    `json:"unique_colors_count"`
	ExtraColorsCount  int    `json:"extra_colors_count"`
	ColorTotalFee     int64  `json:"color_total_fee"`
	BorderStyle       string `json:"border_style"`
	BorderFee         int64  `json:"border_fee"`
	TopCrestFee       int64  `json:"top_crest_fee"`
	BottomCrestFee    int64  `json:"bottom_crest_fee"`
	CrestTotalFee     int64  `json:"crest_total_fee"`
	ImageCount        int    `json:"image_count"`
	ImageTotalFee     int64  `json:"image_total_fee"`
	TotalPrice        int64  `json:"total_price"`
}

// NormalizeHexColor converts a color value into canonical #RRGGBB
// uppercase form, falling back when invalid.
func NormalizeHexColor(color string, fallback string) string {
	if color == "" {
		return fallback
	}
	trimmed := strings.ToUpper(strings.TrimSpace(color))
	if len(trimmed) == 7 && trimmed[0] == '#' {
		return trimmed
	}
	if len(trimmed) == 4 && trimmed[0] == '#' {
		return fmt.Sprintf("#%c%c%c%c%c%c", trimmed[1], trimmed[1], trimmed[2], trimmed[2], trimmed[3], trimmed[3])
	}
	return fallback
}

// ParseCustomDesignPayload parses a raw JSON payload, seamlessly handling
// v3.0.0, v1.0.0, and legacy flat schemas.
func ParseCustomDesignPayload(raw json.RawMessage) (*CustomDesignPayload, error) {
	if len(raw) == 0 {
		return nil, fmt.Errorf("empty custom design payload")
	}

	var rawMap map[string]json.RawMessage
	if err := json.Unmarshal(raw, &rawMap); err != nil {
		return nil, fmt.Errorf("failed to parse custom design payload: %w", err)
	}

	// Check if this is standard structured payload (has "sections" or "layout")
	if _, hasSections := rawMap["sections"]; hasSections {
		var payload CustomDesignPayload
		if err := json.Unmarshal(raw, &payload); err != nil {
			return nil, fmt.Errorf("failed to unmarshal structured custom design: %w", err)
		}
		if payload.Metadata.Version == "" {
			payload.Metadata.Version = "3.0.0"
		}
		if payload.Layout.PhysicalSizeID == "" {
			payload.Layout.PhysicalSizeID = DEFAULT_PHYSICAL_SIZE_ID
		}
		return &payload, nil
	}

	// Fallback / legacy flat payload parser
	var payload CustomDesignPayload
	payload.Metadata.Version = "3.0.0"
	payload.Layout.PhysicalSizeID = DEFAULT_PHYSICAL_SIZE_ID

	if layoutRaw, ok := rawMap["layout"]; ok {
		_ = json.Unmarshal(layoutRaw, &payload.Layout)
	}
	if metaRaw, ok := rawMap["metadata"]; ok {
		_ = json.Unmarshal(metaRaw, &payload.Metadata)
	}
	if decRaw, ok := rawMap["decorations"]; ok {
		_ = json.Unmarshal(decRaw, &payload.Decorations)
	}
	if elemRaw, ok := rawMap["elements"]; ok {
		_ = json.Unmarshal(elemRaw, &payload.Elements)
	}
	if assetRaw, ok := rawMap["assets"]; ok {
		_ = json.Unmarshal(assetRaw, &payload.Assets)
	}

	// Check top-level physicalSizeId / physical_size_id
	if sizeRaw, ok := rawMap["physicalSizeId"]; ok {
		var sizeStr string
		if err := json.Unmarshal(sizeRaw, &sizeStr); err == nil && sizeStr != "" {
			payload.Layout.PhysicalSizeID = sizeStr
		}
	} else if sizeRaw, ok := rawMap["physical_size_id"]; ok {
		var sizeStr string
		if err := json.Unmarshal(sizeRaw, &sizeStr); err == nil && sizeStr != "" {
			payload.Layout.PhysicalSizeID = sizeStr
		}
	}

	// Check top-level heightRatio / upperHeightRatio
	if hrRaw, ok := rawMap["heightRatio"]; ok {
		var hr float64
		if err := json.Unmarshal(hrRaw, &hr); err == nil && hr > 0 {
			payload.Layout.UpperHeightRatio = hr
		}
	} else if hrRaw, ok := rawMap["upperHeightRatio"]; ok {
		var hr float64
		if err := json.Unmarshal(hrRaw, &hr); err == nil && hr > 0 {
			payload.Layout.UpperHeightRatio = hr
		}
	}

	// Parse flat topCrest / bottomCrest if present
	if tcRaw, ok := rawMap["topCrest"]; ok {
		var tc map[string]interface{}
		if err := json.Unmarshal(tcRaw, &tc); err == nil {
			if vis, ok := tc["enabled"].(bool); ok {
				payload.Decorations.TopCrest.Visible = vis
			} else if vis, ok := tc["visible"].(bool); ok {
				payload.Decorations.TopCrest.Visible = vis
			}
			if st, ok := tc["style"].(string); ok {
				payload.Decorations.TopCrest.VariantID = st
			} else if st, ok := tc["variantId"].(string); ok {
				payload.Decorations.TopCrest.VariantID = st
			}
			if p, ok := tc["primary"].(string); ok {
				payload.Decorations.TopCrest.PrimaryColorHex = p
			} else if p, ok := tc["primaryColorHex"].(string); ok {
				payload.Decorations.TopCrest.PrimaryColorHex = p
			}
			if s, ok := tc["secondary"].(string); ok {
				payload.Decorations.TopCrest.SecondaryColorHex = s
			} else if s, ok := tc["secondaryColorHex"].(string); ok {
				payload.Decorations.TopCrest.SecondaryColorHex = s
			}
		}
	}
	if bcRaw, ok := rawMap["bottomCrest"]; ok {
		var bc map[string]interface{}
		if err := json.Unmarshal(bcRaw, &bc); err == nil {
			if vis, ok := bc["enabled"].(bool); ok {
				payload.Decorations.BottomCrest.Visible = vis
			} else if vis, ok := bc["visible"].(bool); ok {
				payload.Decorations.BottomCrest.Visible = vis
			}
			if st, ok := bc["style"].(string); ok {
				payload.Decorations.BottomCrest.VariantID = st
			} else if st, ok := bc["variantId"].(string); ok {
				payload.Decorations.BottomCrest.VariantID = st
			}
			if p, ok := bc["primary"].(string); ok {
				payload.Decorations.BottomCrest.PrimaryColorHex = p
			} else if p, ok := bc["primaryColorHex"].(string); ok {
				payload.Decorations.BottomCrest.PrimaryColorHex = p
			}
			if s, ok := bc["secondary"].(string); ok {
				payload.Decorations.BottomCrest.SecondaryColorHex = s
			} else if s, ok := bc["secondaryColorHex"].(string); ok {
				payload.Decorations.BottomCrest.SecondaryColorHex = s
			}
		}
	}

	// Parse flat previewBase64 / previewUrl if present
	if pbRaw, ok := rawMap["previewBase64"]; ok {
		var pb string
		if err := json.Unmarshal(pbRaw, &pb); err == nil && pb != "" {
			payload.Assets.PreviewBase64 = &pb
		}
	}
	if puRaw, ok := rawMap["previewUrl"]; ok {
		var pu string
		if err := json.Unmarshal(puRaw, &pu); err == nil && pu != "" {
			payload.Assets.PreviewURL = &pu
		}
	}

	// Parse flat upper section if present
	if upperRaw, ok := rawMap["upper"]; ok {
		var flatUpper map[string]interface{}
		if err := json.Unmarshal(upperRaw, &flatUpper); err == nil {
			if bg, ok := flatUpper["bgColor"].(string); ok {
				payload.Sections.Upper.BGColorHex = bg
			} else if bgHex, ok := flatUpper["bgColorHex"].(string); ok {
				payload.Sections.Upper.BGColorHex = bgHex
			}
			if cs, ok := flatUpper["cornerStyle"].(string); ok {
				payload.Sections.Upper.CornerStyle = cs
			}
			if ht, ok := flatUpper["headerText"].(string); ok && ht != "" {
				payload.Sections.Upper.Header.Text = &ht
			}
			if hf, ok := flatUpper["headerFont"].(string); ok {
				payload.Sections.Upper.Header.FontID = hf
			}
			if hc, ok := flatUpper["headerColor"].(string); ok {
				payload.Sections.Upper.Header.FontColorHex = hc
			}
			if bt, ok := flatUpper["bodyText"].(string); ok && bt != "" {
				payload.Sections.Upper.Body.Text = &bt
			}
			if bf, ok := flatUpper["bodyFont"].(string); ok {
				payload.Sections.Upper.Body.FontID = bf
			}
			if bc, ok := flatUpper["bodyColor"].(string); ok {
				payload.Sections.Upper.Body.FontColorHex = bc
			}
		}
	}

	// Parse flat lower section if present
	if lowerRaw, ok := rawMap["lower"]; ok {
		var flatLower map[string]interface{}
		if err := json.Unmarshal(lowerRaw, &flatLower); err == nil {
			if bg, ok := flatLower["bgColor"].(string); ok {
				payload.Sections.Lower.BGColorHex = bg
			} else if bgHex, ok := flatLower["bgColorHex"].(string); ok {
				payload.Sections.Lower.BGColorHex = bgHex
			}
			if cs, ok := flatLower["cornerStyle"].(string); ok {
				payload.Sections.Lower.CornerStyle = cs
			}
			if ht, ok := flatLower["headerText"].(string); ok && ht != "" {
				payload.Sections.Lower.Header.Text = &ht
			}
			if hf, ok := flatLower["headerFont"].(string); ok {
				payload.Sections.Lower.Header.FontID = hf
			}
			if hc, ok := flatLower["headerColor"].(string); ok {
				payload.Sections.Lower.Header.FontColorHex = hc
			}
			if bt, ok := flatLower["bodyText"].(string); ok && bt != "" {
				payload.Sections.Lower.Body.Text = &bt
			}
			if bf, ok := flatLower["bodyFont"].(string); ok {
				payload.Sections.Lower.Body.FontID = bf
			}
			if bc, ok := flatLower["bodyColor"].(string); ok {
				payload.Sections.Lower.Body.FontColorHex = bc
			}
		}
	}

	// Parse flat border if present
	if borderRaw, ok := rawMap["border"]; ok {
		var flatBorder map[string]interface{}
		if err := json.Unmarshal(borderRaw, &flatBorder); err == nil {
			if s, ok := flatBorder["style"].(string); ok {
				payload.Layout.Border.Style = s
			}
			if c, ok := flatBorder["color"].(string); ok {
				payload.Layout.Border.ColorHex = c
			} else if cHex, ok := flatBorder["colorHex"].(string); ok {
				payload.Layout.Border.ColorHex = cHex
			}
			if w, ok := flatBorder["width"].(float64); ok {
				payload.Layout.Border.WidthPx = int(w)
			} else if wPx, ok := flatBorder["widthPx"].(float64); ok {
				payload.Layout.Border.WidthPx = int(wPx)
			}
			if ctr, ok := flatBorder["center"].(bool); ok {
				payload.Layout.Border.ShowCenterDivider = ctr
			} else if div, ok := flatBorder["showCenterDivider"].(bool); ok {
				payload.Layout.Border.ShowCenterDivider = div
			}
		}
	}

	return &payload, nil
}

// CalculateCustomProductPrice computes the per-unit price
// and returns a detailed fee breakdown.
func CalculateCustomProductPrice(design CustomDesignPayload, matrix CustomPricingMatrix) CustomPricingBreakdown {
	sizeID := strings.ToLower(strings.TrimSpace(design.Layout.PhysicalSizeID))
	if sizeID == "" {
		sizeID = DEFAULT_PHYSICAL_SIZE_ID
	}
	basePrice, ok := matrix.SizePrices[sizeID]
	if !ok {
		basePrice = DEFAULT_SIZE_PRICE_MEDIUM
	}

	var brushCount int
	var imageCount int
	colorSet := make(map[string]bool)

	addHexColor := func(c string) {
		if c != "" {
			colorSet[NormalizeHexColor(c, DEFAULT_HEX_COLOR_FALLBACK)] = true
		}
	}

	// Collect unique colors used throughout the design
	addHexColor(design.Sections.Upper.BGColorHex)
	addHexColor(design.Sections.Lower.BGColorHex)
	addHexColor(design.Sections.Upper.Header.FontColorHex)
	addHexColor(design.Sections.Upper.Body.FontColorHex)
	addHexColor(design.Sections.Lower.Header.FontColorHex)
	addHexColor(design.Sections.Lower.Body.FontColorHex)

	// Border color
	if design.Layout.Border.Style != "" &&
		design.Layout.Border.Style != "none" &&
		design.Layout.Border.WidthPx > 0 {

		addHexColor(design.Layout.Border.ColorHex)
	}

	// Apply crest-specific fees
	var topCrestFee int64
	if design.Decorations.TopCrest.Visible {
		addHexColor(design.Decorations.TopCrest.PrimaryColorHex)
		addHexColor(design.Decorations.TopCrest.SecondaryColorHex)
		variant := strings.ToLower(design.Decorations.TopCrest.VariantID)
		if fee, ok := matrix.CrestVariantPrices[variant]; ok {
			topCrestFee = fee
		} else {
			topCrestFee = DEFAULT_CREST_FEE_FALLBACK
		}
	}

	var bottomCrestFee int64
	if design.Decorations.BottomCrest.Visible {
		addHexColor(design.Decorations.BottomCrest.PrimaryColorHex)
		addHexColor(design.Decorations.BottomCrest.SecondaryColorHex)
		variant := strings.ToLower(design.Decorations.BottomCrest.VariantID)
		if fee, ok := matrix.CrestVariantPrices[variant]; ok {
			bottomCrestFee = fee
		} else {
			bottomCrestFee = DEFAULT_CREST_FEE_FALLBACK
		}
	}

	for _, el := range design.Elements {
		if el.Type == "brush" {
			brushCount++
			addHexColor(el.ColorHex)
		} else if el.Type == "image" {
			imageCount++
		}
	}

	// Calculate additional border charge
	var borderFee int64
	borderStyle := strings.ToLower(design.Layout.Border.Style)
	if borderStyle != "" && borderStyle != "none" && design.Layout.Border.WidthPx > 0 {
		if slices.Contains(matrix.PremiumBorderStyles, borderStyle) {
			borderFee = matrix.PremiumBorderFee
		}
	}

	// Build the final pricing breakdown
	brushTotalFee := int64(brushCount) * matrix.PricePerBrushStroke
	imageTotalFee := int64(imageCount) * matrix.PricePerImageElement
	crestTotalFee := topCrestFee + bottomCrestFee

	uniqueColorsCount := len(colorSet)
	extraColorsCount := int(math.Max(0, float64(uniqueColorsCount-matrix.FreeColorThreshold)))
	colorTotalFee := int64(extraColorsCount) * matrix.PricePerExtraColor

	totalPrice := basePrice + brushTotalFee + imageTotalFee + crestTotalFee + borderFee + colorTotalFee

	return CustomPricingBreakdown{
		PhysicalSizeID:    sizeID,
		BaseSizePrice:     basePrice,
		BrushCount:        brushCount,
		BrushTotalFee:     brushTotalFee,
		UniqueColorsCount: uniqueColorsCount,
		ExtraColorsCount:  extraColorsCount,
		ColorTotalFee:     colorTotalFee,
		BorderStyle:       borderStyle,
		BorderFee:         borderFee,
		TopCrestFee:       topCrestFee,
		BottomCrestFee:    bottomCrestFee,
		CrestTotalFee:     crestTotalFee,
		ImageCount:        imageCount,
		ImageTotalFee:     imageTotalFee,
		TotalPrice:        totalPrice,
	}
}

// ExtractDesignSummary returns key fields for quick order indexing and admin fulfillment sheets.
func ExtractDesignSummary(design CustomDesignPayload) (
	physicalSizeID string,
	previewURL *string,
	headerUpper *string,
	bodyUpper *string,
	headerLower *string,
	bodyLower *string,
) {
	physicalSizeID = design.Layout.PhysicalSizeID
	if physicalSizeID == "" {
		physicalSizeID = DEFAULT_PHYSICAL_SIZE_ID
	}

	if design.Assets.PreviewURL != nil && *design.Assets.PreviewURL != "" {
		previewURL = design.Assets.PreviewURL
	}

	if design.Sections.Upper.Header.Text != nil && *design.Sections.Upper.Header.Text != "" {
		headerUpper = design.Sections.Upper.Header.Text
	}
	if design.Sections.Upper.Body.Text != nil && *design.Sections.Upper.Body.Text != "" {
		bodyUpper = design.Sections.Upper.Body.Text
	}
	if design.Sections.Lower.Header.Text != nil && *design.Sections.Lower.Header.Text != "" {
		headerLower = design.Sections.Lower.Header.Text
	}
	if design.Sections.Lower.Body.Text != nil && *design.Sections.Lower.Body.Text != "" {
		bodyLower = design.Sections.Lower.Body.Text
	}

	return physicalSizeID, previewURL, headerUpper, bodyUpper, headerLower, bodyLower
}

