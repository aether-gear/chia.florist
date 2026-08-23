package intelligencelayer

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestClient_HealthCheck(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/health" {
			t.Fatalf("expected path /health, got %s", r.URL.Path)
		}

		resp := APIResponse[HealthData]{
			Success: true,
			Data: &HealthData{
				Status:       "healthy",
				Version:      "1.0.0",
				LoadedModels: []string{"demand", "stockout", "courier", "anomaly"},
			},
		}
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL, 1*time.Second, nil)
	health, err := client.HealthCheck(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if health.Status != "healthy" {
		t.Errorf("expected healthy, got %s", health.Status)
	}
	if len(health.LoadedModels) != 4 {
		t.Errorf("expected 4 models, got %d", len(health.LoadedModels))
	}
}

func TestClient_PredictDemand(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/predict/demand" {
			t.Fatalf("expected path /predict/demand, got %s", r.URL.Path)
		}

		var req DemandForecastRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode request body: %v", err)
		}

		if req.ProductID != "SKU-001" {
			t.Errorf("expected SKU-001, got %s", req.ProductID)
		}

		resp := APIResponse[DemandForecastResponse]{
			Success: true,
			Data: &DemandForecastResponse{
				ProductID:            req.ProductID,
				PredictedUnitsSold7d: 98.5,
				ConfidenceTier:       "high",
			},
		}
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL, 1*time.Second, nil)
	res, err := client.PredictDemand(context.Background(), DemandForecastRequest{
		ProductID:      "SKU-001",
		GrossMarginPct: 0.45,
		ViewCount:      100,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.PredictedUnitsSold7d != 98.5 {
		t.Errorf("expected 98.5, got %f", res.PredictedUnitsSold7d)
	}
	if res.ConfidenceTier != "high" {
		t.Errorf("expected high, got %s", res.ConfidenceTier)
	}
}

func TestClient_PredictStockoutRisk(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/predict/stockout-risk" {
			t.Fatalf("expected path /predict/stockout-risk, got %s", r.URL.Path)
		}

		resp := APIResponse[StockoutRiskResponse]{
			Success: true,
			Data: &StockoutRiskResponse{
				ProductID:           "SKU-LILY-01",
				StockoutProbability: 0.9421,
				WillStockout:        true,
				RiskLevel:           "CRITICAL",
			},
		}
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL, 1*time.Second, nil)
	res, err := client.PredictStockoutRisk(context.Background(), StockoutRiskRequest{
		ProductID:            "SKU-LILY-01",
		Stock:                5,
		ReservedStock:        2,
		StockBurnRate7d:      3.5,
		SupplierLeadTimeDays: 7,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.RiskLevel != "CRITICAL" {
		t.Errorf("expected CRITICAL, got %s", res.RiskLevel)
	}
	if !res.WillStockout {
		t.Errorf("expected will_stockout=true")
	}
}

func TestClient_PredictCourierSLA(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/predict/courier-sla" {
			t.Fatalf("expected path /predict/courier-sla, got %s", r.URL.Path)
		}

		resp := APIResponse[CourierSLAResponse]{
			Success: true,
			Data: &CourierSLAResponse{
				CourierCode:            "jne",
				EstimatedDurationHours: 26.4,
				SLAConfidenceScore:     0.95,
				DeliveryStatus:         "ON_TRACK",
			},
		}
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL, 1*time.Second, nil)
	res, err := client.PredictCourierSLA(context.Background(), CourierSLARequest{
		CourierCode:        "jne",
		ShippingCost:       22000,
		DispatchDayOfWeek:  1,
		DispatchHour:       14,
		DispatchIsWeekend:  0.0,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.EstimatedDurationHours != 26.4 {
		t.Errorf("expected 26.4, got %f", res.EstimatedDurationHours)
	}
	if res.DeliveryStatus != "ON_TRACK" {
		t.Errorf("expected ON_TRACK, got %s", res.DeliveryStatus)
	}
}

func TestClient_DetectAnomaly(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/predict/anomaly" {
			t.Fatalf("expected path /predict/anomaly, got %s", r.URL.Path)
		}

		resp := APIResponse[AnomalyCheckResponse]{
			Success: true,
			Data: &AnomalyCheckResponse{
				IsAnomaly:    true,
				AnomalyScore: -0.1425,
				Severity:     "HIGH",
				Reasons: []string{
					"Excessive payment completion delay (7200s)",
				},
			},
		}
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL, 1*time.Second, nil)
	res, err := client.DetectAnomaly(context.Background(), AnomalyCheckRequest{
		Amount:           4500000,
		TimeToPaySec:     7200,
		IsFailedStatus:   1.0,
		IsManualTransfer: 1.0,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !res.IsAnomaly {
		t.Errorf("expected is_anomaly=true")
	}
	if res.Severity != "HIGH" {
		t.Errorf("expected severity=HIGH, got %s", res.Severity)
	}
}

func TestClient_ServerError_503(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
		_, _ = w.Write([]byte(`{"success":false,"data":null,"error":"Model 'stockout' is not available"}`))
	}))
	defer server.Close()

	client := NewClient(server.URL, 1*time.Second, nil)
	_, err := client.PredictStockoutRisk(context.Background(), StockoutRiskRequest{ProductID: "SKU-1"})
	if err == nil {
		t.Fatalf("expected error on 503, got nil")
	}
}
