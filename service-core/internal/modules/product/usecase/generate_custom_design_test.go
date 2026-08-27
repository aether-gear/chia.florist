package usecase

import (
	"context"
	"errors"
	"testing"

	apperrors "service-core/internal/common/errors"
	"service-core/internal/infra/genai"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockAIProvider struct {
	response *genai.CustomDesignPayloadV3
	err      error
	called   bool
}

func (m *mockAIProvider) GenerateBoardDesign(ctx context.Context, prompt string, designCtx genai.DesignContext) (*genai.CustomDesignPayloadV3, error) {
	m.called = true
	if m.err != nil {
		return nil, m.err
	}
	return m.response, nil
}

func TestGenerateCustomDesignUsecase_Execute(t *testing.T) {
	custID := uuid.New()

	t.Run("Successful Generation", func(t *testing.T) {
		expectedPayload := &genai.CustomDesignPayloadV3{}
		expectedPayload.Metadata.Version = "3.0.0"

		mockProv := &mockAIProvider{
			response: expectedPayload,
		}
		uc := NewGenerateCustomDesignUsecase(mockProv, 5)

		res, err := uc.Execute(context.Background(), GenerateCustomDesignInput{
			CustomerID: custID,
			Prompt:     "Wedding board for John & Jane",
			Occasion:   "Wedding",
		})

		require.NoError(t, err)
		assert.Equal(t, "3.0.0", res.Metadata.Version)
		assert.True(t, mockProv.called)
	})

	t.Run("Unauthorized - Nil CustomerID", func(t *testing.T) {
		mockProv := &mockAIProvider{}
		uc := NewGenerateCustomDesignUsecase(mockProv, 5)

		res, err := uc.Execute(context.Background(), GenerateCustomDesignInput{
			CustomerID: uuid.Nil,
			Prompt:     "Wedding board",
		})

		assert.Error(t, err)
		assert.Nil(t, res)
		appErr, ok := err.(*apperrors.AppError)
		require.True(t, ok)
		assert.Equal(t, 401, appErr.StatusCode)
	})

	t.Run("Invalid Input - Short Prompt", func(t *testing.T) {
		mockProv := &mockAIProvider{}
		uc := NewGenerateCustomDesignUsecase(mockProv, 5)

		res, err := uc.Execute(context.Background(), GenerateCustomDesignInput{
			CustomerID: custID,
			Prompt:     "hi",
		})

		assert.Error(t, err)
		assert.Nil(t, res)
		appErr, ok := err.(*apperrors.AppError)
		require.True(t, ok)
		assert.Equal(t, 400, appErr.StatusCode)
	})

	t.Run("Rate Limit Exceeded", func(t *testing.T) {
		expectedPayload := &genai.CustomDesignPayloadV3{}
		mockProv := &mockAIProvider{
			response: expectedPayload,
		}
		uc := NewGenerateCustomDesignUsecase(mockProv, 2)

		// 1st request -> ok
		_, err := uc.Execute(context.Background(), GenerateCustomDesignInput{
			CustomerID: custID,
			Prompt:     "Valid prompt 1",
		})
		require.NoError(t, err)

		// 2nd request -> ok
		_, err = uc.Execute(context.Background(), GenerateCustomDesignInput{
			CustomerID: custID,
			Prompt:     "Valid prompt 2",
		})
		require.NoError(t, err)

		// 3rd request -> rate limit error 429
		_, err = uc.Execute(context.Background(), GenerateCustomDesignInput{
			CustomerID: custID,
			Prompt:     "Valid prompt 3",
		})
		require.Error(t, err)
		appErr, ok := err.(*apperrors.AppError)
		require.True(t, ok)
		assert.Equal(t, 429, appErr.StatusCode)
	})

	t.Run("AI Provider Error", func(t *testing.T) {
		mockProv := &mockAIProvider{
			err: errors.New("ai provider unreachable"),
		}
		uc := NewGenerateCustomDesignUsecase(mockProv, 5)

		res, err := uc.Execute(context.Background(), GenerateCustomDesignInput{
			CustomerID: custID,
			Prompt:     "Valid prompt",
		})
		assert.Error(t, err)
		assert.Nil(t, res)
	})
}

