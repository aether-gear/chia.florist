package genai

import (
	"context"
	"encoding/json"
)

// DesignContext provides optional context hints to enhance the AI prompt engineering.
type DesignContext struct {
	Occasion         string `json:"occasion,omitempty"`
	PreferredPalette string `json:"preferred_palette,omitempty"`
	Recipient        string `json:"recipient,omitempty"`
	Sender           string `json:"sender,omitempty"`
	PhysicalSizeID   string `json:"physical_size_id,omitempty"`
	Locale           string `json:"locale,omitempty"`
}

type TypographySpec struct {
	Text         *string `json:"text"`
	FontID       string  `json:"fontId"`
	FontSizePx   int     `json:"fontSizePx"`
	FontColorHex string  `json:"fontColorHex"`
	Alignment    string  `json:"alignment"`
}

type BoardSectionSpec struct {
	BgColorHex     string         `json:"bgColorHex"`
	CornerStyle    string         `json:"cornerStyle"`
	OpacityPercent int            `json:"opacityPercent"`
	Header         TypographySpec `json:"header"`
	Body           TypographySpec `json:"body"`
}

type BorderSpec struct {
	Style             string `json:"style"`
	ColorHex          string `json:"colorHex"`
	WidthPx           int    `json:"widthPx"`
	ShowCenterDivider bool   `json:"showCenterDivider"`
}

type CrestSpec struct {
	Visible           bool   `json:"visible"`
	VariantID         string `json:"variantId"`
	PrimaryColorHex   string `json:"primaryColorHex"`
	SecondaryColorHex string `json:"secondaryColorHex"`
	ScalePercent      int    `json:"scalePercent"`
}

type WatermarkSpec struct {
	Enabled        bool   `json:"enabled"`
	Text           string `json:"text"`
	OpacityPercent int    `json:"opacityPercent"`
}

type ElementTransform struct {
	XPercent     float64 `json:"xPercent"`
	YPercent     float64 `json:"yPercent"`
	ScalePercent float64 `json:"scalePercent"`
	RotationDeg  float64 `json:"rotationDeg"`
	ZIndex       int     `json:"zIndex"`
}

type DesignElement struct {
	ID         string           `json:"id"`
	Type       string           `json:"type"` // "image" | "brush"
	Transform  ElementTransform `json:"transform"`
	BrushType  *string          `json:"brushType,omitempty"`
	ColorHex   *string          `json:"colorHex,omitempty"`
	FrameStyle *string          `json:"frameStyle,omitempty"`
	Src        *string          `json:"src,omitempty"`
}

type CustomDesignPayloadV3 struct {
	Metadata struct {
		Version       string          `json:"version"`
		EditorVersion string          `json:"editorVersion"`
		Platform      string          `json:"platform"`
		Locale        string          `json:"locale"`
		CreatedAt     string          `json:"createdAt"`
		UpdatedAt     string          `json:"updatedAt"`
		Checksum      string          `json:"checksum"`
		FeatureFlags  map[string]bool `json:"featureFlags,omitempty"`
	} `json:"metadata"`
	Layout struct {
		PhysicalSizeID   string     `json:"physicalSizeId"`
		AspectRatioID    string     `json:"aspectRatioId"`
		UpperHeightRatio float64    `json:"upperHeightRatio"`
		Border           BorderSpec `json:"border"`
	} `json:"layout"`
	Sections struct {
		Upper BoardSectionSpec `json:"upper"`
		Lower BoardSectionSpec `json:"lower"`
	} `json:"sections"`
	Decorations struct {
		TopCrest    CrestSpec      `json:"topCrest"`
		BottomCrest CrestSpec      `json:"bottomCrest"`
		Watermark   *WatermarkSpec `json:"watermark,omitempty"`
	} `json:"decorations"`
	Elements []DesignElement `json:"elements"`
	Assets   struct {
		PreviewBase64   *string `json:"previewBase64"`
		PreviewAssetID  *string `json:"previewAssetId"`
		PreviewURL      *string `json:"previewUrl"`
		BucketPath      *string `json:"bucketPath"`
		StorageProvider *string `json:"storageProvider"`
	} `json:"assets"`
}

// RawDesignResponse represents the raw JSON structure returned by the AI provider.
type RawDesignResponse struct {
	Layout      *json.RawMessage `json:"layout"`
	Sections    *json.RawMessage `json:"sections"`
	Decorations *json.RawMessage `json:"decorations"`
	Elements    *json.RawMessage `json:"elements"`
}

// Provider defines the interface for interacting with Generative AI services.
type Provider interface {
	GenerateBoardDesign(ctx context.Context, prompt string, designCtx DesignContext) (*CustomDesignPayloadV3, error)
}
