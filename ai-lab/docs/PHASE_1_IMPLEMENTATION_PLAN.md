# 🛠️ Phase 1 Implementation Plan: Feature Extraction & Data Processing Pipeline

> **Module**: `ai-lab`  
> **Target Service**: `service-core` database & analytics data  
> **Date**: 2026-08-10  
> **Reference Document**: [MODEL_TRAINING_REPORT.md](file:///d:/__Projects/kage/chia.florist/ai-lab/docs/MODEL_TRAINING_REPORT.md)  
> **Status**: Proposed / Pending Review  

---

## Executive Summary

Phase 1 establishes the automated **Feature Extraction & Data Processing Pipeline** for `ai-lab`. This pipeline bridges the raw transactional data in `service-core` (or synthetic mock generation when DB is offline) and the model training workflows. 

The pipeline extracts raw tables (`orders`, `order_items`, `product_stock_history`, `inventory`, `product_performance`, `payments`, `shipments`, `audit_logs`), performs feature engineering (lags, rolling averages, duration z-scores, stock velocity ratios), normalizes features, and exports clean datasets to `ai-lab/data/processed/`.

---

## 1. Architectural Overview

```
 ┌─────────────────────────────────────────────────────────────────────────────┐
 │                         1. DATA EXTRACTION LAYER                            │
 │                                                                             │
 │  ┌─────────────────────────────────┐   ┌─────────────────────────────────┐  │
 │  │ service-core PostgreSQL Database │   │ Synthetic Extractor Fallback    │  │
 │  │ (orders, stock, payments, logs) │   │ (Generates synthetic benchmark) │  │
 │  └────────────────┬────────────────┘   └────────────────┬────────────────┘  │
 └───────────────────┼─────────────────────────────────────┼───────────────────┘
                     └──────────────────┬──────────────────┘
                                        │ Raw Tables / Dataframes
                                        ▼
 ┌─────────────────────────────────────────────────────────────────────────────┐
 │                       2. FEATURE ENGINEERING LAYER                          │
 │                          (`src/feature_engineering.py`)                     │
 │                                                                             │
 │  ┌─────────────────────────┐ ┌────────────────────────┐ ┌────────────────┐ │
 │  │ Time-Series Aggregations│ │ Inventory Risk Metrics │ │ Anomaly Matrix │ │
 │  │ • Lags [1, 7, 14, 30]d  │ │ • Stock Velocity Ratio │ │ • Latency Zs   │ │
 │  │ • Rolling Mean & Std    │ │ • Reorder Buffer Days  │ │ • Event Spikes │ │
 │  │ • Holiday & Day-of-Week │ │ • Margin-Weighted Risk │ │ • Failure Pct  │ │
 │  └────────────┬────────────┘ └───────────┬────────────┘ └───────┬────────┘ │
 └───────────────┼──────────────────────────┼──────────────────────┼──────────┘
                 └────────────────────┬─────┴──────────────────────┘
                                      │ Scaled & Imputed Matrices
                                      ▼
 ┌─────────────────────────────────────────────────────────────────────────────┐
 │                         3. STORAGE & PIPELINE OUTPUT                        │
 │                          (`ai-lab/data/processed/`)                         │
 │                                                                             │
 │  • demand_forecasting_features.csv  / .pt                                  │
 │  • stockout_risk_features.csv       / .pt                                  │
 │  • anomaly_detection_features.csv   / .pt                                  │
 └─────────────────────────────────────────────────────────────────────────────┘
```

---

## 2. Planned Components & Code Structure

| File Path | Type | Role & Purpose |
|---|---|---|
| [configs/feature_extraction.yaml](file:///d:/__Projects/kage/chia.florist/ai-lab/configs/feature_extraction.yaml) | `[NEW]` | YAML config for DB URIs, feature window sizes, scaling parameters, and target output paths. |
| [src/extractor.py](file:///d:/__Projects/kage/chia.florist/ai-lab/src/extractor.py) | `[NEW]` | SQL extraction engine for `service-core` PostgreSQL + Synthetic fallback generator. |
| [src/feature_engineering.py](file:///d:/__Projects/kage/chia.florist/ai-lab/src/feature_engineering.py) | `[NEW]` | Core transformers: Time-series lags, stockout indicators, anomaly latency Z-scores, scaler modules. |
| [src/data_loader.py](file:///d:/__Projects/kage/chia.florist/ai-lab/src/data_loader.py) | `[MODIFY]` | Updates PyTorch `Dataset` loaders to load saved processed matrices from `data/processed/`. |
| [extract.py](file:///d:/__Projects/kage/chia.florist/ai-lab/extract.py) | `[NEW]` | CLI tool to execute extraction pipeline (`python extract.py --config configs/feature_extraction.yaml`). |
| [tests/test_feature_extraction.py](file:///d:/__Projects/kage/chia.florist/ai-lab/tests/test_feature_extraction.py) | `[NEW]` | PyTest suite covering dataset extraction, lag generation, scaling, and file serialization. |

---

## 3. Dataset Specifications & Feature Definitions

### A. Demand Forecasting Feature Matrix
* **Target Output**: `data/processed/demand_forecasting_features.csv`
* **Features**:
  - `date`, `shop_id`, `product_id`, `category_id`
  - `units_sold_lag1`, `units_sold_lag7`, `units_sold_lag14`, `units_sold_lag30`
  - `rolling_7d_mean`, `rolling_7d_std`, `rolling_30d_mean`
  - `view_count`, `sales_velocity_7d`, `sales_velocity_30d`, `gross_margin_pct`
  - `day_of_week`, `month`, `is_weekend`, `is_holiday`
  - **Target Label**: `units_sold_next_7d` (numeric float)

### B. Stockout Risk Feature Matrix
* **Target Output**: `data/processed/stockout_risk_features.csv`
* **Features**:
  - `product_id`, `shop_id`, `available_stock`, `reserved_stock`
  - `stock_burn_rate_7d` (units sold per day)
  - `supplier_lead_time_days`
  - `estimated_days_until_stockout` = `available_stock / (stock_burn_rate_7d + 1e-5)`
  - `reorder_urgency_ratio` = `supplier_lead_time_days / (estimated_days_until_stockout + 1e-5)`
  - **Target Label**: `stockout_within_lead_time` (binary 0 or 1)

### C. Operational Anomaly Feature Matrix
* **Target Output**: `data/processed/anomaly_detection_features.csv`
* **Features**:
  - `timestamp_hour`, `payment_method_id`, `courier_code`
  - `payment_creation_to_paid_seconds_avg`
  - `payment_failure_rate_1h`
  - `order_fulfillment_duration_seconds`
  - `status_transition_delay_zscore`
  - `audit_log_event_count_1h`
  - **Target Label**: `is_anomaly_sample` (unsupervised scoring metric / synthetic validation flag)

---

## 4. Verification & Testing Strategy

1. **Automated Unit Tests**:
   - Run `pytest tests/test_feature_extraction.py -v` to verify math transforms, date encodings, lag computations, and synthetic dataset generation.
2. **CLI Dry Run Validation**:
   - Run `python extract.py --config configs/feature_extraction.yaml --dry-run` to test pipeline initialization without writing disk artifacts.
3. **End-to-End Extraction**:
   - Run `python extract.py --config configs/feature_extraction.yaml` and verify dataset files under `ai-lab/data/processed/`.

---

## 5. Next Steps

Upon approval of this implementation plan, execution will proceed with creating configuration files, building `src/extractor.py` and `src/feature_engineering.py`, implementing `extract.py`, writing unit tests, and verifying generated dataset files.
