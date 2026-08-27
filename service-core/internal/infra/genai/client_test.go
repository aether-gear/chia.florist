package genai

import (
	"context"
	"testing"

	"service-core/internal/shared/config"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildUserPrompt(t *testing.T) {
	prompt := BuildUserPrompt("Papan bunga pernikahan untuk Dimas & Sarah", DesignContext{
		Occasion:         "Pernikahan",
		PreferredPalette: "Pastel Gold",
		Recipient:        "Dimas & Sarah",
		Sender:           "PT Maju Jaya",
		PhysicalSizeID:   "large",
	})

	assert.Contains(t, prompt, "Papan bunga pernikahan untuk Dimas & Sarah")
	assert.Contains(t, prompt, "Occasion / Event: Pernikahan")
	assert.Contains(t, prompt, "Preferred Color Palette: Pastel Gold")
	assert.Contains(t, prompt, "Recipient / Honoree: Dimas & Sarah")
	assert.Contains(t, prompt, "Sender / Organization: PT Maju Jaya")
	assert.Contains(t, prompt, "Physical Size: large")
}

func TestClient_GenerateBoardDesign_DeterministicFallback(t *testing.T) {
	cfg := config.GenAIConfig{
		Enabled:   true,
		APIKey:    "", // Triggers deterministic fallback
		TimeoutMS: 5000,
	}
	client := NewClient(cfg)

	t.Run("Wedding Occasion", func(t *testing.T) {
		res, err := client.GenerateBoardDesign(context.Background(), "Buatkan bunga papan wedding", DesignContext{
			Occasion: "Wedding",
		})
		require.NoError(t, err)
		require.NotNil(t, res)

		assert.Equal(t, "3.0.0", res.Metadata.Version)
		assert.Equal(t, "HAPPY WEDDING", *res.Sections.Upper.Header.Text)
		assert.Equal(t, "dancing", res.Sections.Upper.Header.FontID)
		assert.NotEmpty(t, res.Elements)
	})

	t.Run("Condolence Occasion", func(t *testing.T) {
		res, err := client.GenerateBoardDesign(context.Background(), "Bunga duka cita", DesignContext{
			Occasion: "Duka Cita",
		})
		require.NoError(t, err)
		require.NotNil(t, res)

		assert.Equal(t, "TURUT BERDUKA CITA", *res.Sections.Upper.Header.Text)
	})

	t.Run("Empty Prompt Error", func(t *testing.T) {
		res, err := client.GenerateBoardDesign(context.Background(), "   ", DesignContext{})
		assert.Error(t, err)
		assert.Nil(t, res)
	})
}
