package usecase

import (
	"context"
	"errors"
	"testing"

	"service-core/internal/modules/threat_intel/domain"
)

type mockThreatIntelProvider struct {
	reputationData  map[string]any
	geolocationData map[string]any
	reputationErr   error
	geolocationErr  error
	calledIP        string
	calledAPIKey    string
}

func (m *mockThreatIntelProvider) GetReputation(ctx context.Context, ip string, apiKey string) (map[string]any, error) {
	m.calledIP = ip
	m.calledAPIKey = apiKey
	if m.reputationErr != nil {
		return nil, m.reputationErr
	}
	return m.reputationData, nil
}

func (m *mockThreatIntelProvider) GetGeolocation(ctx context.Context, ip string) (map[string]any, error) {
	m.calledIP = ip
	if m.geolocationErr != nil {
		return nil, m.geolocationErr
	}
	return m.geolocationData, nil
}

func TestAnalyzeIPUsecase_Execute_Success(t *testing.T) {
	ctx := context.Background()
	ip := "8.8.8.8"
	apiKey := "test_api_key"
	expectedReputation := map[string]any{
		"harmless": 10,
	}

	provider := &mockThreatIntelProvider{
		reputationData: expectedReputation,
	}
	u := NewAnalyzeIPUsecase(provider)

	res, err := u.Execute(ctx, ip, apiKey)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.IP != ip {
		t.Errorf("expected IP %s, got %s", ip, res.IP)
	}

	if res.RawReport["harmless"] != 10 {
		t.Errorf("expected harmless key value to be 10, got %v", res.RawReport["harmless"])
	}

	if provider.calledIP != ip {
		t.Errorf("expected provider called IP %s, got %s", ip, provider.calledIP)
	}

	if provider.calledAPIKey != apiKey {
		t.Errorf("expected provider called apiKey %s, got %s", apiKey, provider.calledAPIKey)
	}
}

func TestAnalyzeIPUsecase_Execute_EmptyIP(t *testing.T) {
	ctx := context.Background()
	provider := &mockThreatIntelProvider{}
	u := NewAnalyzeIPUsecase(provider)

	_, err := u.Execute(ctx, "", "key")
	if !errors.Is(err, domain.ErrInvalidIP) {
		t.Fatalf("expected ErrInvalidIP, got %v", err)
	}
}

func TestAnalyzeIPUsecase_Execute_ProviderError(t *testing.T) {
	ctx := context.Background()
	providerErr := errors.New("provider connection failed")
	provider := &mockThreatIntelProvider{
		reputationErr: providerErr,
	}
	u := NewAnalyzeIPUsecase(provider)

	_, err := u.Execute(ctx, "8.8.8.8", "key")
	if err == nil || err.Error() != "provider connection failed" {
		t.Fatalf("expected provider error, got %v", err)
	}
}

func TestAnalyzeIPUsecase_Execute_APIKeyRequiredError(t *testing.T) {
	ctx := context.Background()
	provider := &mockThreatIntelProvider{
		reputationErr: domain.ErrAPIKeyRequired,
	}
	u := NewAnalyzeIPUsecase(provider)

	_, err := u.Execute(ctx, "8.8.8.8", "")
	if !errors.Is(err, domain.ErrAPIKeyRequired) {
		t.Fatalf("expected ErrAPIKeyRequired, got %v", err)
	}
}
