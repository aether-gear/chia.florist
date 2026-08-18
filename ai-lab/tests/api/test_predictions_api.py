def test_demand_prediction_api(client):
    payload = {
        "product_id": "TEST-SKU-001",
        "gross_margin_pct": 0.4,
        "view_count": 800,
        "day_of_week": 2,
        "month": 8,
        "is_weekend": 0.0,
        "units_sold_lag_1": 10.0,
        "units_sold_lag_7": 12.0,
        "units_sold_lag_14": 8.0,
        "units_sold_lag_30": 11.0,
        "rolling_7d_mean": 9.5,
        "rolling_7d_std": 1.2,
        "rolling_30d_mean": 8.8,
        "rolling_30d_std": 1.5
    }
    response = client.post("/api/v1/predict/demand", json=payload)
    assert response.status_code == 200
    body = response.json()
    assert body["success"] is True
    assert body["data"]["product_id"] == "TEST-SKU-001"
    assert "predicted_units_sold_7d" in body["data"]

def test_stockout_risk_api(client):
    payload = {
        "product_id": "TEST-SKU-002",
        "stock": 2.0,
        "reserved_stock": 1.0,
        "stock_burn_rate_7d": 4.0,
        "supplier_lead_time_days": 7.0
    }
    response = client.post("/api/v1/predict/stockout-risk", json=payload)
    assert response.status_code == 200
    body = response.json()
    assert body["success"] is True
    assert body["data"]["product_id"] == "TEST-SKU-002"
    assert "stockout_probability" in body["data"]
    assert "will_stockout" in body["data"]

def test_courier_sla_api(client):
    payload = {
        "courier_code": "jne",
        "shipping_cost": 25000.0,
        "dispatch_day_of_week": 1,
        "dispatch_hour": 14,
        "dispatch_is_weekend": 0.0
    }
    response = client.post("/api/v1/predict/courier-sla", json=payload)
    assert response.status_code == 200
    body = response.json()
    assert body["success"] is True
    assert body["data"]["courier_code"] == "jne"
    assert "estimated_duration_hours" in body["data"]

def test_anomaly_detection_api(client):
    payload = {
        "amount": 5000000.0,
        "time_to_pay_sec": 3600.0,
        "is_manual_transfer": 1.0,
        "is_failed_status": 0.0
    }
    response = client.post("/api/v1/predict/anomaly", json=payload)
    assert response.status_code == 200
    body = response.json()
    assert body["success"] is True
    assert "is_anomaly" in body["data"]
    assert "anomaly_score" in body["data"]


