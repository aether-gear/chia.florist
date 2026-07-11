package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"service-core/internal/modules/threat_intel/domain"
	"service-core/internal/modules/threat_intel/repository"
	config "service-core/internal/shared/config"
)

type threatIntelProviderImpl struct {
	client *http.Client
	cfg    config.WAFConfig
}

func NewThreatIntelProvider(cfg config.WAFConfig) repository.ThreatIntelProvider {
	return &threatIntelProviderImpl{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		cfg: cfg,
	}
}

func (p *threatIntelProviderImpl) GetReputation(
	ctx context.Context,
	ip string,
	apiKey string,
) (map[string]any, error) {
	if apiKey == "" {
		apiKey = p.cfg.VirusTotalAPIKey
	}
	if apiKey == "" {
		return nil, domain.ErrAPIKeyRequired
	}

	endpoint := "https://www.virustotal.com/api/v3/ip_addresses"
	url := fmt.Sprintf("%s/%s", endpoint, ip)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("create VT request: %w", err)
	}

	req.Header.Set("x-apikey", apiKey)

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do VT request: %w", domain.ErrProviderUnavailable)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("VT api returned status %d: %w", resp.StatusCode, domain.ErrProviderUnavailable)
	}

	var result map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode VT response: %w", err)
	}

	return result, nil
}

func (p *threatIntelProviderImpl) GetGeolocation(ctx context.Context, ip string) (map[string]any, error) {
	if p.cfg.IP2LocationAPIKey != "" {
		url := fmt.Sprintf("https://api.ip2location.io/?ip=%s&key=%s", ip, p.cfg.IP2LocationAPIKey)
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return nil, fmt.Errorf("create ip2location request: %w", err)
		}

		resp, err := p.client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("do ip2location request: %w", domain.ErrProviderUnavailable)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("ip2location api returned status %d: %w", resp.StatusCode, domain.ErrProviderUnavailable)
		}

		var result map[string]any
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return nil, fmt.Errorf("decode ip2location response: %w", err)
		}

		return result, nil
	}

	endpoint := "http://ip-api.com/json"
	url := fmt.Sprintf("%s/%s", endpoint, ip)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("create GeoIP request: %w", err)
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do GeoIP request: %w", domain.ErrProviderUnavailable)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GeoIP api returned status %d: %w", resp.StatusCode, domain.ErrProviderUnavailable)
	}

	var rawResult map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&rawResult); err != nil {
		return nil, fmt.Errorf("decode GeoIP response: %w", err)
	}

	result := make(map[string]any)
	for k, v := range rawResult {
		result[k] = v
	}

	if lat, ok := rawResult["lat"]; ok {
		result["latitude"] = lat
	}
	if lon, ok := rawResult["lon"]; ok {
		result["longitude"] = lon
	}
	if city, ok := rawResult["city"]; ok {
		result["city_name"] = city
	}
	if country, ok := rawResult["country"]; ok {
		result["country_name"] = country
	}
	if countryCode, ok := rawResult["countryCode"]; ok {
		result["country_code"] = countryCode
	}
	if regionName, ok := rawResult["regionName"]; ok {
		result["region_name"] = regionName
	}

	return result, nil
}
