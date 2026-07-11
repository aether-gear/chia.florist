package usecase

import (
	"context"
	"errors"
	"testing"

	"service-core/internal/modules/threat_intel/domain"
)

func TestGetGeoIPUsecase_Execute_Success(t *testing.T) {
	ctx := context.Background()
	ip := "8.8.8.8"
	expectedGeo := map[string]any{
		"city_name": "Mountain View",
	}

	provider := &mockThreatIntelProvider{
		geolocationData: expectedGeo,
	}
	u := NewGetGeoIPUsecase(provider)

	res, err := u.Execute(ctx, ip)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.IP != ip {
		t.Errorf("expected IP %s, got %s", ip, res.IP)
	}

	if res.RawReport["city_name"] != "Mountain View" {
		t.Errorf("expected city_name key value to be 'Mountain View', got %v", res.RawReport["city_name"])
	}

	if provider.calledIP != ip {
		t.Errorf("expected provider called IP %s, got %s", ip, provider.calledIP)
	}
}

func TestGetGeoIPUsecase_Execute_EmptyIP(t *testing.T) {
	ctx := context.Background()
	provider := &mockThreatIntelProvider{}
	u := NewGetGeoIPUsecase(provider)

	_, err := u.Execute(ctx, "")
	if !errors.Is(err, domain.ErrInvalidIP) {
		t.Fatalf("expected ErrInvalidIP, got %v", err)
	}
}

func TestGetGeoIPUsecase_Execute_ProviderError(t *testing.T) {
	ctx := context.Background()
	providerErr := errors.New("provider connection failed")
	provider := &mockThreatIntelProvider{
		geolocationErr: providerErr,
	}
	u := NewGetGeoIPUsecase(provider)

	_, err := u.Execute(ctx, "8.8.8.8")
	if err == nil || err.Error() != "provider connection failed" {
		t.Fatalf("expected provider error, got %v", err)
	}
}
