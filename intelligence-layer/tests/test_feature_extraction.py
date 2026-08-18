import os
import pytest
import numpy as np

from src.extractor import SyntheticDataExtractor
from src.feature_engineering import (
    TimeSeriesFeatureBuilder,
    InventoryStockoutFeatureBuilder,
    OperationalAnomalyFeatureBuilder,
    StandardScaler
)

def test_synthetic_data_extractor():
    extractor = SyntheticDataExtractor(seed=42, num_days=40, num_products=5, num_shops=2)
    data = extractor.extract_raw_data()

    assert "orders" in data
    assert "order_items" in data
    assert "products" in data
    assert "inventory" in data
    assert "payments" in data

    assert len(data["products"]) == 5
    assert len(data["inventory"]) == 10 # 5 products * 2 shops
    assert len(data["orders"]) > 0


def test_time_series_feature_builder():
    extractor = SyntheticDataExtractor(seed=42, num_days=45, num_products=5, num_shops=2)
    raw_data = extractor.extract_raw_data()

    builder = TimeSeriesFeatureBuilder(lags=[1, 7], rolling_windows=[7])
    headers, X, y = builder.build_features(raw_data)

    assert len(headers) > 0
    assert X.ndim == 2
    assert X.shape[1] == len(headers)
    assert y.ndim == 1
    assert len(X) == len(y)


def test_stockout_risk_feature_builder():
    extractor = SyntheticDataExtractor(seed=42, num_days=30, num_products=5, num_shops=2)
    raw_data = extractor.extract_raw_data()

    builder = InventoryStockoutFeatureBuilder()
    headers, X, y = builder.build_features(raw_data)

    assert len(headers) == 6
    assert X.shape[1] == 6
    assert len(X) == len(raw_data["inventory"])
    assert len(y) == len(X)


def test_anomaly_feature_builder():
    extractor = SyntheticDataExtractor(seed=42, num_days=30, num_products=5, num_shops=2)
    raw_data = extractor.extract_raw_data()

    builder = OperationalAnomalyFeatureBuilder()
    headers, X, y = builder.build_features(raw_data)

    assert len(headers) == 4
    assert X.shape[1] == 4
    assert len(X) == len(raw_data["payments"])
    assert len(y) == len(X)


def test_standard_scaler():
    X = np.array([[10.0, 200.0], [20.0, 400.0], [30.0, 600.0]], dtype=np.float32)
    scaler = StandardScaler()
    scaled_X = scaler.fit_transform(X)

    means = np.mean(scaled_X, axis=0)
    stds = np.std(scaled_X, axis=0)

    np.testing.assert_allclose(means, [0.0, 0.0], atol=1e-5)
    np.testing.assert_allclose(stds, [1.0, 1.0], atol=1e-5)
