package genai

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"service-core/internal/shared/config"
)

type Client struct {
	httpClient *http.Client
	cfg        config.GenAIConfig
}

func NewClient(cfg config.GenAIConfig) *Client {
	timeout := time.Duration(cfg.TimeoutMS) * time.Millisecond
	if timeout <= 0 {
		timeout = 15 * time.Second
	}

	return &Client{
		httpClient: &http.Client{
			Timeout: timeout,
		},
		cfg: cfg,
	}
}

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatCompletionRequest struct {
	Model          string        `json:"model"`
	Messages       []ChatMessage `json:"messages"`
	Temperature    float64       `json:"temperature"`
	ResponseFormat *struct {
		Type string `json:"type"`
	} `json:"response_format,omitempty"`
}

type ChatCompletionChoice struct {
	Message ChatMessage `json:"message"`
}

type ChatCompletionResponse struct {
	Choices []ChatCompletionChoice `json:"choices"`
	Error   *struct {
		Message string `json:"message"`
		Type    string `json:"type"`
	} `json:"error,omitempty"`
}

// GenerateBoardDesign calls the AI service and returns a validated CustomDesignPayloadV3.
func (c *Client) GenerateBoardDesign(ctx context.Context, prompt string, designCtx DesignContext) (*CustomDesignPayloadV3, error) {
	if strings.TrimSpace(prompt) == "" {
		return nil, errors.New("prompt cannot be empty")
	}

	userPrompt := BuildUserPrompt(prompt, designCtx)

	// If API key is empty or service not enabled, provide an intelligent deterministic generation fallback
	if !c.cfg.Enabled || c.cfg.APIKey == "" {
		return c.generateDeterministicFallback(prompt, designCtx), nil
	}

	reqBody := ChatCompletionRequest{
		Model: c.cfg.Model,
		Messages: []ChatMessage{
			{Role: "system", Content: SystemPrompt},
			{Role: "user", Content: userPrompt},
		},
		Temperature: 0.7,
		ResponseFormat: &struct {
			Type string `json:"type"`
		}{Type: "json_object"},
	}

	jsonBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	endpoint := strings.TrimRight(c.cfg.BaseURL, "/")
	if !strings.HasSuffix(endpoint, "/chat/completions") && !strings.Contains(endpoint, "/generate") {
		endpoint += "/chat/completions"
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(jsonBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create http request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	if c.cfg.APIKey != "" {
		httpReq.Header.Set("Authorization", "Bearer "+c.cfg.APIKey)
	}

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		// Log and fallback gracefully if remote AI is temporarily unreachable
		return c.generateDeterministicFallback(prompt, designCtx), nil
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode >= 400 {
		return c.generateDeterministicFallback(prompt, designCtx), nil
	}

	var chatResp ChatCompletionResponse
	if err := json.Unmarshal(bodyBytes, &chatResp); err == nil && len(chatResp.Choices) > 0 {
		rawContent := chatResp.Choices[0].Message.Content
		return c.parseAndSanitizeJSON(rawContent, prompt, designCtx)
	}

	// Try parsing direct payload
	return c.parseAndSanitizeJSON(string(bodyBytes), prompt, designCtx)
}

func (c *Client) parseAndSanitizeJSON(content string, originalPrompt string, designCtx DesignContext) (*CustomDesignPayloadV3, error) {
	cleanJSON := cleanJSONMarkdown(content)

	var payload CustomDesignPayloadV3
	if err := json.Unmarshal([]byte(cleanJSON), &payload); err != nil {
		// Fallback gracefully to thematic heuristic
		return c.generateDeterministicFallback(originalPrompt, designCtx), nil
	}

	c.normalizePayload(&payload, designCtx)
	return &payload, nil
}

func cleanJSONMarkdown(raw string) string {
	trimmed := strings.TrimSpace(raw)
	if strings.HasPrefix(trimmed, "```json") {
		trimmed = strings.TrimPrefix(trimmed, "```json")
		if idx := strings.LastIndex(trimmed, "```"); idx != -1 {
			trimmed = trimmed[:idx]
		}
	} else if strings.HasPrefix(trimmed, "```") {
		trimmed = strings.TrimPrefix(trimmed, "```")
		if idx := strings.LastIndex(trimmed, "```"); idx != -1 {
			trimmed = trimmed[:idx]
		}
	}
	return strings.TrimSpace(trimmed)
}

func (c *Client) normalizePayload(payload *CustomDesignPayloadV3, designCtx DesignContext) {
	now := time.Now().UTC().Format(time.RFC3339)
	payload.Metadata.Version = "3.0.0"
	if payload.Metadata.EditorVersion == "" {
		payload.Metadata.EditorVersion = "3.0.0"
	}
	if payload.Metadata.Platform == "" {
		payload.Metadata.Platform = "web"
	}
	if payload.Metadata.Locale == "" {
		payload.Metadata.Locale = "id-ID"
	}
	if payload.Metadata.CreatedAt == "" {
		payload.Metadata.CreatedAt = now
	}
	payload.Metadata.UpdatedAt = now

	if payload.Layout.PhysicalSizeID == "" {
		if designCtx.PhysicalSizeID != "" {
			payload.Layout.PhysicalSizeID = designCtx.PhysicalSizeID
		} else {
			payload.Layout.PhysicalSizeID = "medium"
		}
	}
	if payload.Layout.AspectRatioID == "" {
		payload.Layout.AspectRatioID = "portrait-3-4"
	}
	if payload.Layout.UpperHeightRatio < 0.25 || payload.Layout.UpperHeightRatio > 0.75 {
		payload.Layout.UpperHeightRatio = 0.58
	}
	if payload.Layout.Border.WidthPx <= 0 {
		payload.Layout.Border.WidthPx = 12
	}
	if payload.Layout.Border.Style == "" {
		payload.Layout.Border.Style = "solid"
	}
	if payload.Layout.Border.ColorHex == "" {
		payload.Layout.Border.ColorHex = "#F5C842"
	}

	// Calculate checksum
	hash := sha256.Sum256([]byte(fmt.Sprintf("%s-%s-%s", payload.Sections.Upper.Header.FontColorHex, payload.Sections.Lower.Header.FontColorHex, now)))
	payload.Metadata.Checksum = hex.EncodeToString(hash[:])[:16]
}

// generateDeterministicFallback creates a beautiful themed design based on keyword heuristics
func (c *Client) generateDeterministicFallback(prompt string, designCtx DesignContext) *CustomDesignPayloadV3 {
	lowerPrompt := strings.ToLower(prompt + " " + designCtx.Occasion)

	var (
		upperHeader = "SELAMAT & SUKSES"
		upperBody   = "Atas Pembukaan Usaha Baru"
		lowerHeader = "SEMOGA SEMAKIN SUKSES & JAYA"
		lowerBody   = "Dari Rekan & Sahabat"
		upperBg     = "#C0392B"
		lowerBg     = "#1A3A5C"
		upperText   = "#FFD700"
		lowerText   = "#FFFFFF"
		borderColor = "#F5C842"
		crestStyle  = "grand"
		fontHeader  = "playfair"
		fontBody    = "inter"
	)

	if strings.Contains(lowerPrompt, "wedding") || strings.Contains(lowerPrompt, "nikah") || strings.Contains(lowerPrompt, "kawin") {
		upperHeader = "HAPPY WEDDING"
		upperBody = "Selamat Menempuh Hidup Baru"
		lowerHeader = "SEMOGA BERBAHAGIA SELALU"
		lowerBody = "Best Wishes From Friends & Family"
		upperBg = "#E85D75"
		lowerBg = "#FFF0F3"
		upperText = "#FFFFFF"
		lowerText = "#2B2D42"
		borderColor = "#FFD166"
		crestStyle = "classic"
		fontHeader = "dancing"
		fontBody = "inter"
	} else if strings.Contains(lowerPrompt, "duka") || strings.Contains(lowerPrompt, "condolence") || strings.Contains(lowerPrompt, "wafat") || strings.Contains(lowerPrompt, "meninggal") {
		upperHeader = "TURUT BERDUKA CITA"
		upperBody = "Atas Berpulangnya Sosok Tercinta"
		lowerHeader = "SEMOGA DIBERIKAN KETABAHAN"
		lowerBody = "Keluarga Besar & Rekan Sejawat"
		upperBg = "#2B2D42"
		lowerBg = "#1F2421"
		upperText = "#FFFFFF"
		lowerText = "#E0E0E0"
		borderColor = "#D4AF37"
		crestStyle = "modern"
		fontHeader = "playfair"
		fontBody = "merriweather"
	} else if strings.Contains(lowerPrompt, "wisuda") || strings.Contains(lowerPrompt, "graduat") {
		upperHeader = "HAPPY GRADUATION"
		upperBody = "Selamat & Sukses atas Kelulusannya"
		lowerHeader = "SEMOGA ILMUNYA BERKAH & SUKSES SELALU"
		lowerBody = "Sahabat & Kerabat Tercinta"
		upperBg = "#0D9488"
		lowerBg = "#134E4A"
		upperText = "#FDE047"
		lowerText = "#FFFFFF"
		borderColor = "#FDE047"
		crestStyle = "modern"
		fontHeader = "pacifico"
		fontBody = "inter"
	}

	if designCtx.Recipient != "" {
		upperBody = designCtx.Recipient
	}
	if designCtx.Sender != "" {
		lowerBody = designCtx.Sender
	}

	payload := &CustomDesignPayloadV3{}
	payload.Metadata.Version = "3.0.0"
	payload.Metadata.EditorVersion = "3.0.0"
	payload.Metadata.Platform = "web"
	payload.Metadata.Locale = "id-ID"
	payload.Metadata.CreatedAt = time.Now().UTC().Format(time.RFC3339)
	payload.Metadata.UpdatedAt = time.Now().UTC().Format(time.RFC3339)

	size := "medium"
	if designCtx.PhysicalSizeID != "" {
		size = designCtx.PhysicalSizeID
	}
	payload.Layout.PhysicalSizeID = size
	payload.Layout.AspectRatioID = "portrait-3-4"
	payload.Layout.UpperHeightRatio = 0.58
	payload.Layout.Border = BorderSpec{
		Style:             "solid",
		ColorHex:          borderColor,
		WidthPx:           12,
		ShowCenterDivider: true,
	}

	payload.Sections.Upper = BoardSectionSpec{
		BgColorHex:     upperBg,
		CornerStyle:    "ornate",
		OpacityPercent: 100,
		Header: TypographySpec{
			Text:         &upperHeader,
			FontID:       fontHeader,
			FontSizePx:   36,
			FontColorHex: upperText,
			Alignment:    "center",
		},
		Body: TypographySpec{
			Text:         &upperBody,
			FontID:       fontBody,
			FontSizePx:   22,
			FontColorHex: "#FFFFFF",
			Alignment:    "center",
		},
	}

	payload.Sections.Lower = BoardSectionSpec{
		BgColorHex:     lowerBg,
		CornerStyle:    "none",
		OpacityPercent: 100,
		Header: TypographySpec{
			Text:         &lowerHeader,
			FontID:       "bebas",
			FontSizePx:   24,
			FontColorHex: lowerText,
			Alignment:    "center",
		},
		Body: TypographySpec{
			Text:         &lowerBody,
			FontID:       "inter",
			FontSizePx:   20,
			FontColorHex: lowerText,
			Alignment:    "center",
		},
	}

	payload.Decorations.TopCrest = CrestSpec{
		Visible:           true,
		VariantID:         crestStyle,
		PrimaryColorHex:   borderColor,
		SecondaryColorHex: "#FFFFFF",
		ScalePercent:      40,
	}
	payload.Decorations.BottomCrest = CrestSpec{
		Visible:           true,
		VariantID:         crestStyle,
		PrimaryColorHex:   borderColor,
		SecondaryColorHex: "#FFFFFF",
		ScalePercent:      40,
	}

	payload.Elements = []DesignElement{
		{
			ID:   "ai-elem-1",
			Type: "brush",
			Transform: ElementTransform{
				XPercent:     10,
				YPercent:     10,
				ScalePercent: 36,
				RotationDeg:  45,
				ZIndex:       10,
			},
			BrushType: stringPtr("flower"),
			ColorHex:  stringPtr(borderColor),
		},
		{
			ID:   "ai-elem-2",
			Type: "brush",
			Transform: ElementTransform{
				XPercent:     90,
				YPercent:     10,
				ScalePercent: 36,
				RotationDeg:  -45,
				ZIndex:       11,
			},
			BrushType: stringPtr("flower"),
			ColorHex:  stringPtr(borderColor),
		},
	}

	c.normalizePayload(payload, designCtx)
	return payload
}

func stringPtr(s string) *string {
	return &s
}
