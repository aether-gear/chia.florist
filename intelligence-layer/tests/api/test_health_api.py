def test_health_check_endpoint(client):
    response = client.get("/api/v1/health")
    assert response.status_code == 200
    
    body = response.json()
    assert body["success"] is True
    assert "data" in body
    assert body["data"]["status"] == "healthy"
    assert "loaded_models" in body["data"]
